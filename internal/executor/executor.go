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
	backend, err := gitbackend.NewGitBackend(gitbackend.Options{})
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
		run.Status = model.StatusFailed
		run.ErrorMessage = fmt.Sprintf("query source repo failed: %v", err)
		run.ErrorType = ClassifyError(err)
		return run, err
	}
	if sourceRepo == nil {
		run.Status = model.StatusFailed
		run.ErrorMessage = fmt.Sprintf("source repo not found: %s", task.SourceRepoKey)
		run.ErrorType = model.ErrorConfig
		return run, fmt.Errorf("source repo not found: %s", task.SourceRepoKey)
	}

	targetRepo, err := e.service.GetRepoByKey(task.TargetRepoKey)
	if err != nil {
		run.Status = model.StatusFailed
		run.ErrorMessage = fmt.Sprintf("query target repo failed: %v", err)
		run.ErrorType = ClassifyError(err)
		return run, err
	}
	if targetRepo == nil {
		run.Status = model.StatusFailed
		run.ErrorMessage = fmt.Sprintf("target repo not found: %s", task.TargetRepoKey)
		run.ErrorType = model.ErrorConfig
		return run, fmt.Errorf("target repo not found: %s", task.TargetRepoKey)
	}

	workDir := e.service.GetTempDir(task.Key)
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		run.Status = model.StatusFailed
		run.ErrorMessage = fmt.Sprintf("create work dir failed: %v", err)
		run.ErrorType = ClassifyError(err)
		return run, err
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
			run.Status = model.StatusFailed
			run.ErrorMessage = fmt.Sprintf("clone failed: %v", err)
			run.ErrorType = ClassifyError(err)
			return run, err
		}
	} else {
		details.WriteString("Step 1: Fetch updates from source repo...\n")
		if err := e.fetchRepo(execCtx, repoDir, task, sourceRepo); err != nil {
			e.failStep(step1, err)
			fmt.Fprintf(&details, "fetch error: %v\n", err)
			run.Status = model.StatusFailed
			run.ErrorMessage = fmt.Sprintf("fetch failed: %v", err)
			run.ErrorType = ClassifyError(err)
			return run, err
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
		run.Status = model.StatusFailed
		run.ErrorMessage = fmt.Sprintf("add remote failed: %v", err)
		run.ErrorType = ClassifyError(err)
		return run, err
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
		run.Status = model.StatusFailed
		run.ErrorMessage = fmt.Sprintf("push failed after %d attempts: %v", maxRetries, pushErr)
		run.ErrorType = ClassifyError(pushErr)
		return run, pushErr
	}
	e.completeStep(step3, "")
	details.WriteString("Step 3: completed\n")

	run.Status = model.StatusSuccess
	details.WriteString("\n=== Sync completed successfully ===")
	return run, nil
}

func (e *Executor) authConfig(repo *model.Repo) gitbackend.AuthConfig {
	// Use repo's access token first
	if repo.AccessToken != "" {
		slog.Debug("using repo access token", "repo", repo.Key)
		return gitbackend.AuthConfig{
			Type:  gitbackend.AuthHTTPToken,
			Token: repo.AccessToken,
		}
	}
	// Fall back to platform's token
	if repo.PlatformID > 0 {
		platform, err := e.service.GetPlatformByID(repo.PlatformID)
		if err == nil && platform != nil && platform.AccessToken != "" {
			slog.Debug("using platform access token", "repo", repo.Key, "platform", platform.Name)
			return gitbackend.AuthConfig{
				Type:  gitbackend.AuthHTTPToken,
				Token: platform.AccessToken,
			}
		}
		slog.Debug("no platform token found", "repo", repo.Key, "platformID", repo.PlatformID)
	}
	slog.Debug("no auth configured", "repo", repo.Key)
	return gitbackend.AuthConfig{Type: gitbackend.AuthNone}
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
	remotes, err := e.backend.ListRemoteBranches(ctx, dir, RemoteTarget)
	if err != nil || len(remotes) == 0 {
		return e.backend.AddRemote(ctx, dir, RemoteTarget, repo.CloneURL)
	}
	return nil
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

// SyncPreview contains the result of a sync preview operation.
type SyncPreview struct {
	CanSync       bool
	SourceExists  bool
	TargetExists  bool
	CommitsBehind int
	CommitsAhead  int
	LatestCommit  string
	Message       string
}

// Preview checks if a sync operation can be performed.
func (e *Executor) Preview(ctx context.Context, task *model.SyncTask) (*SyncPreview, error) {
	preview := &SyncPreview{}

	sourceRepo, err := e.service.GetRepoByKey(task.SourceRepoKey)
	if err != nil {
		preview.Message = fmt.Sprintf("source repo error: %v", err)
		return preview, nil
	}
	preview.SourceExists = sourceRepo != nil
	if !preview.SourceExists {
		preview.Message = "source repo not found"
		return preview, nil
	}

	targetRepo, err := e.service.GetRepoByKey(task.TargetRepoKey)
	if err != nil {
		preview.Message = fmt.Sprintf("target repo error: %v", err)
		return preview, nil
	}
	preview.TargetExists = targetRepo != nil
	if !preview.TargetExists {
		preview.Message = "target repo not found"
		return preview, nil
	}

	preview.CanSync = true
	preview.Message = "sync can be performed"
	return preview, nil
}
