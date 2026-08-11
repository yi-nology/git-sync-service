package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/internal/executor"
	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
	"gorm.io/gorm"
)

// Compile-time check: Service satisfies executor.Service interface.
var _ executor.Service = (*Service)(nil)

type Config = model.Config

type Service struct {
	config    *Config
	db        *gorm.DB
	repos     *RepoService
	tasks     *TaskService
	webhooks  *WebhookService
	platforms *PlatformService
	cron            *cron.Cron
	cronEntryIDs    map[string]cron.EntryID
	cronMu          sync.RWMutex
	executor        *executor.Executor
	lastTriggerTime sync.Map
	cleanupDone     chan struct{}
	bgCtx           context.Context
	bgCancel        context.CancelFunc
	wg              sync.WaitGroup
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

	repoDAO, err := dao.NewRepoDAO(db)
	if err != nil {
		return nil, fmt.Errorf("init repo DAO failed: %w", err)
	}

	providerMgr := sdkprov.NewManager(30 * time.Minute)
	taskDAO := dao.NewSyncTaskDAO(db)
	runDAO := dao.NewSyncRunDAO(db)
	runStepDAO := dao.NewSyncRunStepDAO(db)
	ruleDAO := dao.NewWebhookRuleDAO(db)
	eventDAO := dao.NewWebhookEventDAO(db)
	platformDAO := dao.NewPlatformDAO(db)

	repoService := NewRepoService(repoDAO, providerMgr)
	taskService := NewTaskService(taskDAO, runDAO, runStepDAO, repoDAO)
	webhookService := NewWebhookService(ruleDAO, eventDAO, repoDAO)
	platformService := NewPlatformService(platformDAO, repoDAO, providerMgr)

	bgCtx, bgCancel := context.WithCancel(context.Background())

	svc := &Service{
		config:         cfg,
		db:             db,
		repos:          repoService,
		tasks:          taskService,
		webhooks:       webhookService,
		platforms:      platformService,
		cron:           cron.New(cron.WithSeconds()),
		cronEntryIDs:   make(map[string]cron.EntryID),
		cleanupDone:    make(chan struct{}),
		bgCtx:          bgCtx,
		bgCancel:       bgCancel,
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
	// Cancel background context to signal all goroutines
	s.bgCancel()
	close(s.cleanupDone)
	s.stopCronJobs()

	// Wait for background goroutines to finish (with timeout)
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		slog.Info("all background goroutines stopped")
	case <-time.After(10 * time.Second):
		slog.Warn("timeout waiting for background goroutines to stop")
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

// CreateRun creates a new sync run record. Satisfies executor.RunManager.
func (s *Service) CreateRun(task *model.SyncTask, trigger string, webhookEventID *uint) (*model.SyncRun, error) {
	return s.tasks.CreateRun(task, trigger, webhookEventID)
}

// CreateRunStep creates a new sync run step record. Satisfies executor.RunManager.
func (s *Service) CreateRunStep(step *model.SyncRunStep) error {
	return s.tasks.CreateRunStep(step)
}

// UpdateRunStep updates an existing sync run step record. Satisfies executor.RunManager.
func (s *Service) UpdateRunStep(step *model.SyncRunStep) error {
	return s.tasks.UpdateRunStep(step)
}

// CompleteRun updates a sync run with final status and details. Satisfies executor.RunManager.
func (s *Service) CompleteRun(run *model.SyncRun) error {
	return s.tasks.CompleteRun(run)
}

// UpdateTaskLastRun updates the task's last run status. Satisfies executor.RunManager.
func (s *Service) UpdateTaskLastRun(task *model.SyncTask, run *model.SyncRun) error {
	return s.tasks.UpdateTaskLastRun(task, run)
}

// GetRepoByKey returns a repository by key. Satisfies executor.RepoProvider.
func (s *Service) GetRepoByKey(key string) (*model.Repo, error) {
	return s.repos.GetRepoByKey(key)
}

// HealthCheck checks the health of all dependencies.
// Returns a map of component name to "ok" or error message.
func (s *Service) HealthCheck() map[string]string {
	status := map[string]string{
		"database": "ok",
		"redis":    "ok",
		"service":  "ok",
	}

	// Check database connectivity
	sqlDB, err := s.db.DB()
	if err != nil {
		status["database"] = err.Error()
	} else if err := sqlDB.Ping(); err != nil {
		status["database"] = err.Error()
	}

	// Check Redis connectivity (if configured)
	if s.config.Redis.Addr != "" {
		status["redis"] = "not checked (lock service removed)"
	} else {
		status["redis"] = "not configured"
	}

	return status
}

// Platform related methods

// CreatePlatform 创建平台
func (s *Service) CreatePlatform(ctx context.Context, platform *model.Platform) error {
	return s.platforms.CreatePlatform(ctx, platform)
}

// GetPlatform 获取平台
func (s *Service) GetPlatform(ctx context.Context, key string) (*model.Platform, error) {
	return s.platforms.GetPlatform(ctx, key)
}

// ListPlatforms 列出所有平台
func (s *Service) ListPlatforms(ctx context.Context) ([]*model.Platform, error) {
	return s.platforms.ListPlatforms(ctx)
}

// UpdatePlatform 更新平台
func (s *Service) UpdatePlatform(ctx context.Context, platform *model.Platform) error {
	return s.platforms.UpdatePlatform(ctx, platform)
}

// DeletePlatform 删除平台
func (s *Service) DeletePlatform(ctx context.Context, key string) error {
	return s.platforms.DeletePlatform(ctx, key)
}

// SetDefaultPlatform 设置默认平台
func (s *Service) SetDefaultPlatform(ctx context.Context, key string) error {
	return s.platforms.SetDefaultPlatform(ctx, key)
}

// UpdatePlatformStatus 更新平台状态
func (s *Service) UpdatePlatformStatus(ctx context.Context, key, status, testResult string) error {
	return s.platforms.UpdatePlatformStatus(ctx, key, status, testResult)
}

// TestPlatformConnection 测试平台连接
func (s *Service) TestPlatformConnection(ctx context.Context, key string) (*sdkprov.TestConnectionResult, error) {
	return s.platforms.TestPlatformConnection(ctx, key)
}

// ListPlatformRepos 列出平台上的仓库
func (s *Service) ListPlatformRepos(ctx context.Context, key, page, perPage string) ([]*sdkprov.PlatformRepo, error) {
	return s.platforms.ListPlatformRepos(ctx, key, page, perPage)
}

// SyncPlatformRepos 同步平台仓库到本地
func (s *Service) SyncPlatformRepos(ctx context.Context, key string) (int, error) {
	return s.platforms.SyncPlatformRepos(ctx, key)
}

// ListReposByPlatform 列出平台下的仓库
func (s *Service) ListReposByPlatform(ctx context.Context, platformKey string) ([]*model.Repo, error) {
	return s.platforms.ListReposByPlatform(ctx, platformKey)
}
