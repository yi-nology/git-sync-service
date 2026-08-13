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
	opLogs    *OperationLogService
	cron            *cron.Cron
	cronEntryIDs    map[string]cron.EntryID
	cronMu          sync.RWMutex
	executor        *executor.Executor
	lastTriggerTime sync.Map
	// guard 统一封装“同 taskKey 互斥 + 全局并发上限”;配 redis 时为分布式,否则进程内。
	guard concurrencyGuard
	cleanupDone chan struct{}
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
	opLogDAO := dao.NewOperationLogDAO(db)

	repoService := NewRepoService(repoDAO, platformDAO, providerMgr)
	taskService := NewTaskService(taskDAO, runDAO, runStepDAO, repoDAO)
	webhookService := NewWebhookService(ruleDAO, eventDAO, repoDAO)
	platformService := NewPlatformService(platformDAO, repoDAO, providerMgr)
	opLogService := NewOperationLogService(opLogDAO)

	bgCtx, bgCancel := context.WithCancel(context.Background())

	svc := &Service{
		config:       cfg,
		db:           db,
		repos:        repoService,
		tasks:        taskService,
		webhooks:     webhookService,
		platforms:    platformService,
		opLogs:       opLogService,
		cron:         cron.New(cron.WithSeconds()),
		cronEntryIDs: make(map[string]cron.EntryID),
		cleanupDone:  make(chan struct{}),
		bgCtx:        bgCtx,
		bgCancel:     bgCancel,
	}

	// 并发控制:配了 redis 用分布式(多实例安全),否则进程内(单实例)
	guard, err := newGuard(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, cfg.Sync.MaxConcurrent)
	if err != nil {
		return nil, fmt.Errorf("init concurrency guard failed: %w", err)
	}
	svc.guard = guard

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

	// 释放并发控制器(关闭 redis 连接等)
	if s.guard != nil {
		if err := s.guard.Close(); err != nil {
			slog.Error("failed to close concurrency guard", "error", err)
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

// GetPlatformByID returns a platform by ID. Satisfies executor.PlatformProvider.
func (s *Service) GetPlatformByID(id uint) (*model.Platform, error) {
	return s.platforms.GetPlatformByID(context.Background(), id)
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

// Operation log (audit) related methods

// RecordOperation 记录一条审计日志。
func (s *Service) RecordOperation(ctx context.Context, entry *model.OperationLog) error {
	return s.opLogs.Record(ctx, entry)
}

// ListOperations 按过滤条件分页返回审计日志。
func (s *Service) ListOperations(ctx context.Context, offset, limit int, filter *dao.OperationLogFilter) ([]*model.OperationLog, int64, error) {
	return s.opLogs.List(ctx, offset, limit, filter)
}

// OperationStats 返回今日、本周、总操作数。
func (s *Service) OperationStats(ctx context.Context) (today, week, total int64, err error) {
	return s.opLogs.Stats(ctx)
}
