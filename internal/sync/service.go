package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/internal/sync/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Service struct {
	config       *Config
	db           *gorm.DB
	repoDAO      *RepoDAO
	taskDAO      *SyncTaskDAO
	runDAO       *SyncRunDAO
	ruleDAO      *WebhookRuleDAO
	eventDAO     *WebhookEventDAO
	providerMgr  *ProviderManager
	cron         *cron.Cron
	cronEntryIDs map[string]cron.EntryID
	lock         DistLock
	semaphore    *Semaphore
	semaphoreID  string
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

	var lock DistLock
	var sem *Semaphore

	if cfg.Redis.Addr != "" {
		lock = NewRedisLock(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
		redisClient := lock.(*RedisLock).client
		maxConcurrent := cfg.Sync.MaxConcurrent
		if maxConcurrent <= 0 {
			maxConcurrent = 5
		}
		sem = NewSemaphore(redisClient, "sync-tasks", maxConcurrent)
	} else {
		lock = &LocalLock{}
	}

	return &Service{
		config:       cfg,
		db:           db,
		repoDAO:      NewRepoDAO(db),
		taskDAO:      NewSyncTaskDAO(db),
		runDAO:       NewSyncRunDAO(db),
		ruleDAO:      NewWebhookRuleDAO(db),
		eventDAO:     NewWebhookEventDAO(db),
		providerMgr:  NewProviderManager(),
		cron:         cron.New(cron.WithSeconds()),
		cronEntryIDs: make(map[string]cron.EntryID),
		lock:         lock,
		semaphore:    sem,
		semaphoreID:  uuid.New().String(),
	}, nil
}

func (s *Service) Start() error {
	tasks, err := s.taskDAO.FindAllEnabled()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Cron != "" {
			if err := s.addCronJob(task); err != nil {
				fmt.Printf("add cron job for task %s failed: %v\n", task.Key, err)
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

func (s *Service) getTempDir(taskKey string) string {
	return filepath.Join(s.config.Git.TempDir, taskKey)
}

func (s *Service) tryAcquireTaskLock(ctx context.Context, taskKey string) (bool, error) {
	lockKey := fmt.Sprintf("task:%s", taskKey)
	ttl := time.Duration(s.config.Sync.DefaultTimeout) * time.Second
	if ttl <= 0 {
		ttl = 300 * time.Second
	}
	return s.lock.TryLockWithTTL(ctx, lockKey, ttl)
}

func (s *Service) releaseTaskLock(ctx context.Context, taskKey string) error {
	lockKey := fmt.Sprintf("task:%s", taskKey)
	return s.lock.Unlock(ctx, lockKey)
}

func (s *Service) acquireSemaphore(ctx context.Context, taskKey string) (bool, error) {
	if s.semaphore == nil {
		return true, nil
	}
	return s.semaphore.Acquire(ctx, s.semaphoreID+":"+taskKey)
}

func (s *Service) releaseSemaphore(ctx context.Context, taskKey string) {
	if s.semaphore != nil {
		s.semaphore.Release(ctx, s.semaphoreID+":"+taskKey)
	}
}

type LocalLock struct {
	mu sync.Map
}

func (l *LocalLock) TryLock(ctx context.Context, key string) (bool, error) {
	return l.TryLockWithTTL(ctx, key, defaultLockTTL)
}

func (l *LocalLock) TryLockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	_, loaded := l.mu.LoadOrStore(key, time.Now().Add(ttl))
	return !loaded, nil
}

func (l *LocalLock) Lock(ctx context.Context, key string) error {
	for {
		ok, _ := l.TryLock(ctx, key)
		if ok {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (l *LocalLock) Unlock(ctx context.Context, key string) error {
	l.mu.Delete(key)
	return nil
}

var _ DistLock = (*LocalLock)(nil)
