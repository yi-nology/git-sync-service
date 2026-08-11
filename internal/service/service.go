package service

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/internal/executor"
	"github.com/yi-nology/git-sync-service/internal/lock"
	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
	"gorm.io/gorm"
)

type Config = model.Config

type Service struct {
	config          *Config
	db              *gorm.DB
	repoService     *RepoService
	taskService     *TaskService
	webhookService  *WebhookService
	cron            *cron.Cron
	cronEntryIDs    map[string]cron.EntryID
	cronMu          sync.RWMutex
	lock            lock.DistLock
	semaphore       *lock.Semaphore
	semaphoreID     string
	executor        *executor.Executor
	lastTriggerTime sync.Map
	cleanupDone     chan struct{}
}

func NewService(cfg *Config) (*Service, error) {
	db, err := model.InitDB(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("init db failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	if err := os.MkdirAll(cfg.Git.TempDir, 0o755); err != nil {
		return nil, fmt.Errorf("create temp dir failed: %w", err)
	}

	var distLock lock.DistLock
	var sem *lock.Semaphore

	if cfg.Redis.Addr != "" {
		distLock = lock.NewRedisLock(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
		redisClient := distLock.(*lock.RedisLock).Client()
		maxConcurrent := cfg.Sync.MaxConcurrent
		if maxConcurrent <= 0 {
			maxConcurrent = 5
		}
		sem = lock.NewSemaphore(redisClient, "sync-tasks", maxConcurrent)
	} else {
		distLock = lock.NewLocalLock()
	}

	repoDAO, err := dao.NewRepoDAO(db)
	if err != nil {
		return nil, fmt.Errorf("init repo DAO failed: %w", err)
	}

	providerMgr := sdkprov.NewManager(30 * time.Minute)
	taskDAO := dao.NewSyncTaskDAO(db)
	runDAO := dao.NewSyncRunDAO(db)
	ruleDAO := dao.NewWebhookRuleDAO(db)
	eventDAO := dao.NewWebhookEventDAO(db)

	repoService := NewRepoService(repoDAO, providerMgr)
	taskService := NewTaskService(taskDAO, runDAO, repoDAO)
	webhookService := NewWebhookService(ruleDAO, eventDAO, repoDAO)

	svc := &Service{
		config:         cfg,
		db:             db,
		repoService:    repoService,
		taskService:    taskService,
		webhookService: webhookService,
		cron:           cron.New(cron.WithSeconds()),
		cronEntryIDs:   make(map[string]cron.EntryID),
		lock:           distLock,
		semaphore:      sem,
		semaphoreID:    uuid.New().String(),
		cleanupDone:    make(chan struct{}),
	}

	exec, err := executor.NewExecutor(svc)
	if err != nil {
		return nil, fmt.Errorf("init executor failed: %w", err)
	}
	svc.executor = exec

	go svc.cleanupTriggerTimes()

	return svc, nil
}

func (s *Service) Start() error {
	return s.startCronJobs()
}

func (s *Service) Stop() {
	close(s.cleanupDone)
	s.stopCronJobs()
	if closer, ok := s.lock.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			slog.Error("failed to close lock", "error", err)
		}
	}
}

func (s *Service) GetTempDir(taskKey string) string {
	return filepath.Join(s.config.Git.TempDir, taskKey)
}

func (s *Service) GetConfig() *model.Config {
	return s.config
}

func (s *Service) GetAPIKey() string {
	return s.config.Server.APIKey
}

func (s *Service) RunDAO() executor.RunWriter {
	return s.taskService.runDAO
}

func (s *Service) TaskDAO() executor.TaskUpdater {
	return s.taskService.taskDAO
}

func (s *Service) RepoDAO() executor.RepoReader {
	return s.repoService.repoDAO
}
