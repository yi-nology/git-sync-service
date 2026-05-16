package sync

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/yi-nology/git-sync-service/internal/sync/model"
)

type Executor struct {
	service *Service
}

func NewExecutor(svc *Service) *Executor {
	return &Executor{service: svc}
}

func (e *Executor) Execute(ctx context.Context, task *model.SyncTask, trigger string) (*model.SyncRun, error) {
	lockCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	locked, err := e.service.tryAcquireTaskLock(lockCtx, task.Key)
	if err != nil {
		return nil, fmt.Errorf("acquire lock failed: %w", err)
	}
	if !locked {
		return nil, fmt.Errorf("task is already running")
	}
	defer e.service.releaseTaskLock(ctx, task.Key)

	semCtx, cancelSem := context.WithTimeout(ctx, 30*time.Second)
	defer cancelSem()

	semAcquired, err := e.service.acquireSemaphore(semCtx, task.Key)
	if err != nil {
		return nil, fmt.Errorf("acquire semaphore failed: %w", err)
	}
	if !semAcquired {
		return nil, fmt.Errorf("too many concurrent tasks")
	}
	defer e.service.releaseSemaphore(ctx, task.Key)

	run := &model.SyncRun{
		TaskKey:       task.Key,
		TriggerSource: trigger,
		Status:        "running",
		StartTime:     time.Now(),
	}
	if err := e.service.runDAO.Create(run); err != nil {
		return nil, err
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf("=== Sync Task: %s ===\n", task.Name))
	details.WriteString(fmt.Sprintf("Trigger: %s\n", trigger))
	details.WriteString(fmt.Sprintf("Time: %s\n\n", time.Now().Format(time.RFC3339)))

	defer func() {
		run.EndTime = timePtr(time.Now())
		run.Details = details.String()
		_ = e.service.runDAO.Update(run)

		task.LastRunAt = run.EndTime
		task.LastStatus = run.Status
		_ = e.service.taskDAO.Update(task)
	}()

	sourceRepo, err := e.service.repoDAO.FindByKey(task.SourceRepoKey)
	if err != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("source repo not found: %v", err)
		return run, err
	}

	targetRepo, err := e.service.repoDAO.FindByKey(task.TargetRepoKey)
	if err != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("target repo not found: %v", err)
		return run, err
	}

	workDir := e.service.getTempDir(task.Key)
	if err := os.MkdirAll(workDir, 0755); err != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("create work dir failed: %w", err)
		return run, err
	}
	defer os.RemoveAll(workDir)

	repoDir := filepath.Join(workDir, "repo")

	timeout := e.service.config.Sync.DefaultTimeout
	if timeout <= 0 {
		timeout = 300
	}
	execCtx, execCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer execCancel()

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		details.WriteString("Step 1: Initial clone of source repo...\n")
		if err := e.cloneRepo(execCtx, repoDir, sourceRepo, task.SourceBranch, &details); err != nil {
			run.Status = "failed"
			run.ErrorMessage = fmt.Sprintf("clone failed: %v", err)
			return run, err
		}
	} else {
		details.WriteString("Step 1: Fetch updates from source repo...\n")
		if err := e.fetchRepo(execCtx, repoDir, task.SourceBranch, &details); err != nil {
			run.Status = "failed"
			run.ErrorMessage = fmt.Sprintf("fetch failed: %v", err)
			return run, err
		}
	}

	details.WriteString("\nStep 2: Ensure target remote exists...\n")
	if err := e.ensureRemote(execCtx, repoDir, targetRepo, &details); err != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("add remote failed: %v", err)
		return run, err
	}

	details.WriteString("\nStep 3: Push to target...\n")
	maxRetries := e.service.config.Sync.RetryCount
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var pushErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			details.WriteString(fmt.Sprintf("\nRetry attempt %d/%d...\n", attempt, maxRetries))
			time.Sleep(time.Duration(attempt*500) * time.Millisecond)
		}

		pushErr = e.push(execCtx, repoDir, task, &details)
		if pushErr == nil {
			break
		}

		if attempt < maxRetries {
			details.WriteString(fmt.Sprintf("Push failed, retrying fetch...\n"))
			if err := e.fetchRepo(execCtx, repoDir, task.SourceBranch, &details); err != nil {
				details.WriteString(fmt.Sprintf("Retry fetch failed: %v\n", err))
			}
		}
	}

	if pushErr != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("push failed after %d attempts: %v", maxRetries, pushErr)
		return run, pushErr
	}

	run.Status = "success"
	details.WriteString("\n=== Sync completed successfully ===")
	return run, nil
}

func (e *Executor) cloneRepo(ctx context.Context, dir string, repo *model.Repo, branch string, details *strings.Builder) error {
	authURL := e.buildAuthURL(repo.CloneURL, repo.AccessToken)

	args := []string{"clone", "--branch", branch, "--single-branch", "--depth", "1", authURL, dir}
	cmd := exec.CommandContext(ctx, "git", args...)

	output, err := cmd.CombinedOutput()
	details.Write(output)

	return err
}

func (e *Executor) fetchRepo(ctx context.Context, dir string, branch string, details *strings.Builder) error {
	args := []string{"fetch", "origin", branch}
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	details.Write(output)
	if err != nil {
		return err
	}

	args = []string{"checkout", branch}
	cmd = exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	output, err = cmd.CombinedOutput()
	details.Write(output)
	if err != nil {
		return err
	}

	args = []string{"pull", "origin", branch}
	cmd = exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	output, err = cmd.CombinedOutput()
	details.Write(output)

	return err
}

func (e *Executor) ensureRemote(ctx context.Context, dir string, repo *model.Repo, details *strings.Builder) error {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "target")
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return e.addRemote(ctx, dir, repo, details)
	}

	authURL := e.buildAuthURL(repo.CloneURL, repo.AccessToken)
	cmd = exec.CommandContext(ctx, "git", "remote", "set-url", "target", authURL)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	details.Write(output)

	return err
}

func (e *Executor) addRemote(ctx context.Context, dir string, repo *model.Repo, details *strings.Builder) error {
	authURL := e.buildAuthURL(repo.CloneURL, repo.AccessToken)

	cmd := exec.CommandContext(ctx, "git", "remote", "add", "target", authURL)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	details.Write(output)

	return err
}

func (e *Executor) push(ctx context.Context, dir string, task *model.SyncTask, details *strings.Builder) error {
	args := []string{"push"}
	if task.GitForce {
		args = append(args, "--force")
	}
	if task.GitPrune {
		args = append(args, "--prune")
	}
	if task.GitNoVerify {
		args = append(args, "--no-verify")
	}
	if task.GitTags {
		args = append(args, "--tags")
	}

	refSpec := fmt.Sprintf("%s:%s", task.SourceBranch, task.TargetBranch)
	args = append(args, "target", refSpec)

	if task.PushOptions != "" {
		for _, opt := range strings.Split(task.PushOptions, ",") {
			args = append(args, "-o", strings.TrimSpace(opt))
		}
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	details.Write(output)

	return err
}

func (e *Executor) buildAuthURL(cloneURL, token string) string {
	if token == "" {
		return cloneURL
	}

	u, err := url.Parse(cloneURL)
	if err != nil {
		return cloneURL
	}

	u.User = url.UserPassword("oauth2", token)
	return u.String()
}

func timePtr(t time.Time) *time.Time {
	return &t
}

type SyncPreview struct {
	CanSync       bool
	SourceExists  bool
	TargetExists  bool
	CommitsBehind  int
	CommitsAhead int
	LatestCommit  string
	Message       string
}

func (e *Executor) Preview(ctx context.Context, task *model.SyncTask) (*SyncPreview, error) {
	preview := &SyncPreview{}

	sourceRepo, err := e.service.repoDAO.FindByKey(task.SourceRepoKey)
	if err != nil {
		preview.Message = fmt.Sprintf("source repo error: %v", err)
		return preview, nil
	}
	preview.SourceExists = sourceRepo != nil
	if !preview.SourceExists {
		preview.Message = "source repo not found"
		return preview, nil
	}

	targetRepo, err := e.service.repoDAO.FindByKey(task.TargetRepoKey)
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
