package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"git.enjoye.top/enjoydream/ekit/pkg/uidgen"
	"github.com/robfig/cron/v3"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/internal/executor"
	"github.com/yi-nology/git-sync-service/internal/lock"
	"github.com/yi-nology/git-sync-service/internal/provider"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

type Config = model.Config

type Service struct {
	config        *Config
	db            *gorm.DB
	repoDAO       *dao.RepoDAO
	taskDAO       *dao.SyncTaskDAO
	runDAO        *dao.SyncRunDAO
	ruleDAO       *dao.WebhookRuleDAO
	eventDAO      *dao.WebhookEventDAO
	providerMgr   *provider.ProviderManager
	cron          *cron.Cron
	cronEntryIDs  map[string]cron.EntryID
	cronMu        sync.RWMutex
	lock          lock.DistLock
	semaphore     *lock.Semaphore
	semaphoreID   string
	executor      *executor.Executor
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

	if err := os.MkdirAll(cfg.Git.TempDir, 0755); err != nil {
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
		distLock = &lock.LocalLock{}
	}

	svc := &Service{
		config:       cfg,
		db:           db,
		repoDAO:      dao.NewRepoDAO(db),
		taskDAO:      dao.NewSyncTaskDAO(db),
		runDAO:       dao.NewSyncRunDAO(db),
		ruleDAO:      dao.NewWebhookRuleDAO(db),
		eventDAO:     dao.NewWebhookEventDAO(db),
		providerMgr:  provider.NewProviderManager(),
		cron:         cron.New(cron.WithSeconds()),
		cronEntryIDs: make(map[string]cron.EntryID),
		lock:         distLock,
		semaphore:    sem,
		semaphoreID:  uidgen.UUID(),
	}

	svc.executor = executor.NewExecutor(svc)

	return svc, nil
}

func (s *Service) Start() error {
	tasks, err := s.taskDAO.FindAllEnabled()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Cron != "" {
			if err := s.addCronJob(task); err != nil {
				slog.Error("add cron job failed", "taskKey", task.Key, "error", err)
			}
		}
	}

	s.cron.Start()
	return nil
}

func (s *Service) Stop() {
	s.cron.Stop()
	if closer, ok := s.lock.(interface{ Close() error }); ok {
		closer.Close()
	}
}

func (s *Service) addCronJob(task *model.SyncTask) error {
	s.cronMu.Lock()
	defer s.cronMu.Unlock()

	if entryID, ok := s.cronEntryIDs[task.Key]; ok {
		s.cron.Remove(entryID)
	}

	entryID, err := s.cron.AddFunc(task.Cron, func() {
		ctx := context.Background()
		_ = s.RunTaskWithTrigger(ctx, task.Key, "cron")
	})
	if err != nil {
		return err
	}

	s.cronEntryIDs[task.Key] = entryID
	return nil
}

func (s *Service) removeCronJob(taskKey string) {
	s.cronMu.Lock()
	defer s.cronMu.Unlock()

	if entryID, ok := s.cronEntryIDs[taskKey]; ok {
		s.cron.Remove(entryID)
		delete(s.cronEntryIDs, taskKey)
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

func (s *Service) TryAcquireTaskLock(ctx context.Context, taskKey string) (bool, error) {
	lockKey := fmt.Sprintf("task:%s", taskKey)
	ttl := time.Duration(s.config.Sync.DefaultTimeout) * time.Second
	if ttl <= 0 {
		ttl = 300 * time.Second
	}
	return s.lock.TryLockWithTTL(ctx, lockKey, ttl)
}

func (s *Service) ReleaseTaskLock(ctx context.Context, taskKey string) error {
	lockKey := fmt.Sprintf("task:%s", taskKey)
	return s.lock.Unlock(ctx, lockKey)
}

func (s *Service) AcquireSemaphore(ctx context.Context, taskKey string) (bool, error) {
	if s.semaphore == nil {
		return true, nil
	}
	return s.semaphore.Acquire(ctx, s.semaphoreID+":"+taskKey)
}

func (s *Service) ReleaseSemaphore(ctx context.Context, taskKey string) {
	if s.semaphore != nil {
		s.semaphore.Release(ctx, s.semaphoreID+":"+taskKey)
	}
}

func (s *Service) RunDAO() interface {
	Create(run *model.SyncRun) error
	Update(run *model.SyncRun) error
} {
	return s.runDAO
}

func (s *Service) TaskDAO() interface {
	Update(task *model.SyncTask) error
} {
	return s.taskDAO
}

func (s *Service) RepoDAO() interface {
	FindByKey(key string) (*model.Repo, error)
} {
	return s.repoDAO
}

func (s *Service) CleanupOldData(maxAge time.Duration) (events int64, runs int64, err error) {
	events, err = s.eventDAO.CleanupOlderThan(maxAge)
	if err != nil {
		return 0, 0, err
	}
	runs, err = s.runDAO.CleanupOlderThan(maxAge)
	if err != nil {
		return events, 0, err
	}
	slog.Info("data cleanup completed", "events_deleted", events, "runs_deleted", runs)
	return events, runs, nil
}
