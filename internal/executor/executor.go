package executor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yi-nology/git-platform-sdk/gitbackend"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// RunManager handles sync run lifecycle operations.
type RunManager interface {
	CreateRun(task *model.SyncTask, trigger string, webhookEventID *uint) (*model.SyncRun, error)
	CreateRunStep(step *model.SyncRunStep) error
	UpdateRunStep(step *model.SyncRunStep) error
	CompleteRun(run *model.SyncRun) error
	UpdateTaskLastRun(task *model.SyncTask, run *model.SyncRun) error
}

// RepoProvider provides repository lookup by key.
type RepoProvider interface {
	GetRepoByKey(key string) (*model.Repo, error)
}

// PlatformProvider provides platform lookup by ID.
type PlatformProvider interface {
	GetPlatformByID(id uint) (*model.Platform, error)
}

// Service is the interface the executor depends on.
// It exposes only the operations the executor needs — no DAO leakage.
type Service interface {
	GetTempDir(taskKey string) string
	GetConfig() *model.Config
	RunManager
	RepoProvider
	PlatformProvider
}

type Executor struct {
	service Service
	backend gitbackend.GitBackend
}

func NewExecutor(svc Service) (*Executor, error) {
	// git.backend 配置接线:空值由 SDK 自动选择(native 优先,回退 gogit)
	backendType := ""
	if cfg := svc.GetConfig(); cfg != nil {
		backendType = cfg.Git.Backend
	}
	backend, err := gitbackend.NewGitBackend(gitbackend.Options{Type: backendType})
	if err != nil {
		return nil, fmt.Errorf("init git backend failed: %w", err)
	}
	return &Executor{
		service: svc,
		backend: backend,
	}, nil
}

func (e *Executor) Execute(ctx context.Context, task *model.SyncTask, trigger string, webhookEventID *uint) (*model.SyncRun, error) {
	run, err := e.service.CreateRun(task, trigger, webhookEventID)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()

	// Build a summary details string for backward compatibility
	var details strings.Builder
	fmt.Fprintf(&details, "=== Sync Task: %s ===\n", task.Name)
	fmt.Fprintf(&details, "Trigger: %s\n", trigger)
	fmt.Fprintf(&details, "Time: %s\n\n", startTime.Format(time.RFC3339))

	defer func() {
		run.EndTime = timePtr(time.Now())
		run.DurationMs = run.EndTime.Sub(startTime).Milliseconds()
		run.Details = details.String()
		if err := e.service.CompleteRun(run); err != nil {
			slog.Error("failed to complete sync run", "error", err)
		}
		if err := e.service.UpdateTaskLastRun(task, run); err != nil {
			slog.Error("failed to update task status", "error", err)
		}
	}()

	sourceRepo, err := e.service.GetRepoByKey(task.SourceRepoKey)
	if err != nil {
		return failRun(run, fmt.Errorf("query source repo failed: %v", err))
	}
	if sourceRepo == nil {
		return failRun(run, fmt.Errorf("source repo not found: %s", task.SourceRepoKey))
	}

	targetRepo, err := e.service.GetRepoByKey(task.TargetRepoKey)
	if err != nil {
		return failRun(run, fmt.Errorf("query target repo failed: %v", err))
	}
	if targetRepo == nil {
		return failRun(run, fmt.Errorf("target repo not found: %s", task.TargetRepoKey))
	}

	workDir := e.service.GetTempDir(task.Key)
	if err := os.MkdirAll(workDir, 0o750); err != nil {
		return failRun(run, fmt.Errorf("create work dir failed: %v", err))
	}

	defer func() {
		// 成功:保留 workdir,下次走增量 fetch(避免每次全量 clone);
		// 非成功(失败/panic):清理,避免损坏的 .git / index.lock 影响下次执行。
		// 注:同 taskKey 的并发执行已由 Service.runningTasks 互斥,不存在并发写同一 workdir。
		if run.Status != model.StatusSuccess {
			if rmErr := os.RemoveAll(workDir); rmErr != nil {
				slog.Error("failed to cleanup temp dir", "error", rmErr, "dir", workDir)
			}
		}
	}()

	repoDir := filepath.Join(workDir, RepoDir)

	timeout := e.service.GetConfig().Sync.DefaultTimeout
	if timeout <= 0 {
		timeout = 300
	}
	execCtx, execCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer execCancel()

	// Step 1: Clone or Fetch
	var step1Name string
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		step1Name = model.StepClone
	} else {
		step1Name = model.StepFetch
	}

	step1 := e.beginStep(run.ID, step1Name)
	if step1Name == model.StepClone {
		details.WriteString("Step 1: Initial clone of source repo...\n")
		if err := e.cloneRepo(execCtx, repoDir, sourceRepo, task); err != nil {
			e.failStep(step1, err)
			fmt.Fprintf(&details, "clone error: %v\n", err)
			return failRun(run, fmt.Errorf("clone failed: %v", err))
		}
	} else {
		details.WriteString("Step 1: Fetch updates from source repo...\n")
		if err := e.fetchRepo(execCtx, repoDir, task, sourceRepo); err != nil {
			e.failStep(step1, err)
			fmt.Fprintf(&details, "fetch error: %v\n", err)
			return failRun(run, fmt.Errorf("fetch failed: %v", err))
		}
	}
	e.completeStep(step1, "")
	details.WriteString("Step 1: completed\n")

	// Step 2: Ensure target remote
	step2 := e.beginStep(run.ID, model.StepEnsureRemote)
	details.WriteString("\nStep 2: Ensure target remote exists...\n")
	if err := e.ensureRemote(execCtx, repoDir, targetRepo); err != nil {
		e.failStep(step2, err)
		fmt.Fprintf(&details, "add remote error: %v\n", err)
		return failRun(run, fmt.Errorf("add remote failed: %v", err))
	}
	e.completeStep(step2, "")
	details.WriteString("Step 2: completed\n")

	// Step 3: Push with retry
	step3 := e.beginStep(run.ID, model.StepPush)
	details.WriteString("\nStep 3: Push to target...\n")
	maxRetries := e.service.GetConfig().Sync.RetryCount
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var pushErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			step3.RetryCount = attempt - 1
			fmt.Fprintf(&details, "\nRetry attempt %d/%d...\n", attempt, maxRetries)
			// 退避期间监听 context,超时/取消时立即中止,不再干等
			backoff := time.Duration(attempt*500) * time.Millisecond
			select {
			case <-execCtx.Done():
				pushErr = execCtx.Err()
				fmt.Fprintf(&details, "retry aborted (context done): %v\n", pushErr)
			case <-time.After(backoff):
			}
			if pushErr != nil {
				break
			}
		}

		pushErr = e.push(execCtx, repoDir, task, targetRepo)
		if pushErr == nil {
			break
		}

		if attempt < maxRetries {
			details.WriteString("Push failed, retrying fetch...\n")
			if err := e.fetchRepo(execCtx, repoDir, task, sourceRepo); err != nil {
				fmt.Fprintf(&details, "Retry fetch failed: %v\n", err)
			}
		}
	}

	run.RetryTotal = max(0, step3.RetryCount)

	if pushErr != nil {
		e.failStep(step3, pushErr)
		fmt.Fprintf(&details, "push error: %v\n", pushErr)
		return failRun(run, fmt.Errorf("push failed after %d attempts: %v", maxRetries, pushErr))
	}
	e.completeStep(step3, "")
	details.WriteString("Step 3: completed\n")

	run.Status = model.StatusSuccess
	details.WriteString("\n=== Sync completed successfully ===")
	return run, nil
}

// failRun 标记 run 失败并按错误分类填充字段,统一 Execute 的失败出口。
func failRun(run *model.SyncRun, err error) (*model.SyncRun, error) {
	run.Status = model.StatusFailed
	run.ErrorMessage = err.Error()
	run.ErrorType = ClassifyError(err)
	return run, err
}

func (e *Executor) authConfig(repo *model.Repo) gitbackend.AuthConfig {
	// 平台记录只查一次:token 仅作回退,SkipTLSVerify 则无论 token 来源
	// 是否为仓库自身都要生效(自签名证书平台的 git clone/fetch/push 依赖它)。
	var skipTLS bool
	var platformToken string
	if repo.PlatformID > 0 {
		if platform, err := e.service.GetPlatformByID(repo.PlatformID); err == nil && platform != nil {
			skipTLS = platform.SkipTLSVerify
			platformToken = platform.AccessToken
		}
	}
	token := repo.AccessToken
	if token == "" {
		token = platformToken
	}
	if token != "" {
		slog.Debug("using access token", "repo", repo.Key, "fromPlatform", repo.AccessToken == "")
		return gitbackend.AuthConfig{
			Type:             gitbackend.AuthHTTPToken,
			Token:            token,
			InsecureSkipTLS:  skipTLS,
		}
	}
	slog.Debug("no auth configured", "repo", repo.Key)
	return gitbackend.AuthConfig{Type: gitbackend.AuthNone, InsecureSkipTLS: skipTLS}
}

// beginStep creates a new step record in "running" state.
func (e *Executor) beginStep(runID uint, stepName string) *model.SyncRunStep {
	step := &model.SyncRunStep{
		RunID:     runID,
		StepName:  stepName,
		Status:    model.StatusRunning,
		StartTime: time.Now(),
	}
	if err := e.service.CreateRunStep(step); err != nil {
		slog.Error("failed to create run step", "step", stepName, "error", err)
	}
	return step
}

// completeStep marks a step as successfully completed.
func (e *Executor) completeStep(step *model.SyncRunStep, output string) {
	now := time.Now()
	step.EndTime = &now
	step.DurationMs = now.Sub(step.StartTime).Milliseconds()
	step.Status = model.StatusSuccess
	step.Output = output
	if err := e.service.UpdateRunStep(step); err != nil {
		slog.Error("failed to update run step", "step", step.StepName, "error", err)
	}
}

// failStep marks a step as failed with error classification.
func (e *Executor) failStep(step *model.SyncRunStep, err error) {
	now := time.Now()
	step.EndTime = &now
	step.DurationMs = now.Sub(step.StartTime).Milliseconds()
	step.Status = model.StatusFailed
	step.ErrorMsg = err.Error()
	step.ErrorType = ClassifyError(err)
	if updateErr := e.service.UpdateRunStep(step); updateErr != nil {
		slog.Error("failed to update run step", "step", step.StepName, "error", updateErr)
	}
}

func (e *Executor) cloneRepo(ctx context.Context, dir string, repo *model.Repo, task *model.SyncTask) error {
	// Use shallow clone by default for efficiency
	depth := 1
	// If force push is enabled, use full clone to avoid "shallow update not allowed"
	if task.GitForce {
		depth = 0
	}
	return e.backend.Clone(ctx, gitbackend.CloneOptions{
		URL:          repo.CloneURL,
		Path:         dir,
		Branch:       task.SourceBranch,
		Depth:        depth,
		SingleBranch: true,
		Auth:         e.authConfig(repo),
	})
}

func (e *Executor) fetchRepo(ctx context.Context, dir string, task *model.SyncTask, repo *model.Repo) error {
	_, err := e.backend.Fetch(ctx, gitbackend.FetchOptions{
		RepoPath: dir,
		Remote:   RemoteOrigin,
		Branches: []string{task.SourceBranch},
		Tags:     task.GitTags,
		Prune:    task.GitPrune,
		Auth:     e.authConfig(repo),
	})
	if err != nil {
		return err
	}

	return e.backend.Checkout(ctx, dir, task.SourceBranch)
}

func (e *Executor) ensureRemote(ctx context.Context, dir string, repo *model.Repo) error {
	// 按配置的 remote 名称判断存在性(GetRemotes);此前用 ListRemoteBranches
	// 依赖 remote-tracking refs,从未 fetch 过的既有 remote 会被误判为不存在,
	// 导致 AddRemote 报 "remote target already exists"。
	remotes, err := e.backend.GetRemotes(ctx, dir)
	if err != nil {
		return err
	}
	for _, name := range remotes {
		if name == RemoteTarget {
			return nil
		}
	}
	return e.backend.AddRemote(ctx, dir, RemoteTarget, repo.CloneURL)
}

func (e *Executor) push(ctx context.Context, dir string, task *model.SyncTask, repo *model.Repo) error {
	refSpec := fmt.Sprintf("%s:%s", task.SourceBranch, task.TargetBranch)

	_, err := e.backend.Push(ctx, gitbackend.PushOptions{
		RepoPath: dir,
		Remote:   RemoteTarget,
		RefSpecs: []string{refSpec},
		Force:    task.GitForce,
		Auth:     e.authConfig(repo),
	})
	return err
}

func timePtr(t time.Time) *time.Time {
	return &t
}

