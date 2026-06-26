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

type RunWriter interface {
	Create(run *model.SyncRun) error
	Update(run *model.SyncRun) error
}

type TaskUpdater interface {
	Update(task *model.SyncTask) error
}

type RepoReader interface {
	FindByKey(key string) (*model.Repo, error)
}

type ConfigProvider interface {
	GetTempDir(taskKey string) string
	GetConfig() *model.Config
}

type DAOProvider interface {
	RunDAO() RunWriter
	TaskDAO() TaskUpdater
	RepoDAO() RepoReader
}

type Service interface {
	ConfigProvider
	DAOProvider
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

func (e *Executor) Execute(ctx context.Context, task *model.SyncTask, trigger string) (*model.SyncRun, error) {
	run := &model.SyncRun{
		TaskKey:       task.Key,
		TriggerSource: trigger,
		Status:        "running",
		StartTime:     time.Now(),
	}
	if err := e.service.RunDAO().Create(run); err != nil {
		return nil, err
	}

	var details strings.Builder
	fmt.Fprintf(&details, "=== Sync Task: %s ===\n", task.Name)
	fmt.Fprintf(&details, "Trigger: %s\n", trigger)
	fmt.Fprintf(&details, "Time: %s\n\n", time.Now().Format(time.RFC3339))

	defer func() {
		run.EndTime = timePtr(time.Now())
		run.Details = details.String()
		if err := e.service.RunDAO().Update(run); err != nil {
			slog.Error("failed to update sync run", "error", err)
		}

		task.LastRunAt = run.EndTime
		task.LastStatus = run.Status
		if err := e.service.TaskDAO().Update(task); err != nil {
			slog.Error("failed to update task status", "error", err)
		}
	}()

	sourceRepo, err := e.service.RepoDAO().FindByKey(task.SourceRepoKey)
	if err != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("source repo not found: %v", err)
		return run, err
	}

	targetRepo, err := e.service.RepoDAO().FindByKey(task.TargetRepoKey)
	if err != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("target repo not found: %v", err)
		return run, err
	}

	workDir := e.service.GetTempDir(task.Key)
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		run.Status = "failed"
		run.ErrorMessage = fmt.Sprintf("create work dir failed: %v", err)
		return run, err
	}

	repoDir := filepath.Join(workDir, "repo")

	timeout := e.service.GetConfig().Sync.DefaultTimeout
	if timeout <= 0 {
		timeout = 300
	}
	execCtx, execCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer execCancel()

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		details.WriteString("Step 1: Initial clone of source repo...\n")
		if err := e.cloneRepo(execCtx, repoDir, sourceRepo, task, &details); err != nil {
			run.Status = "failed"
			run.ErrorMessage = fmt.Sprintf("clone failed: %v", err)
			return run, err
		}
	} else {
		details.WriteString("Step 1: Fetch updates from source repo...\n")
		if err := e.fetchRepo(execCtx, repoDir, task, &details); err != nil {
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
	maxRetries := e.service.GetConfig().Sync.RetryCount
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var pushErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			fmt.Fprintf(&details, "\nRetry attempt %d/%d...\n", attempt, maxRetries)
			time.Sleep(time.Duration(attempt*500) * time.Millisecond)
		}

		pushErr = e.push(execCtx, repoDir, task, sourceRepo, &details)
		if pushErr == nil {
			break
		}

		if attempt < maxRetries {
			details.WriteString("Push failed, retrying fetch...\n")
			if err := e.fetchRepo(execCtx, repoDir, task, &details); err != nil {
				fmt.Fprintf(&details, "Retry fetch failed: %v\n", err)
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

func (e *Executor) authConfig(repo *model.Repo) gitbackend.AuthConfig {
	if repo.AccessToken != "" {
		return gitbackend.AuthConfig{
			Type:  gitbackend.AuthHTTPToken,
			Token: repo.AccessToken,
		}
	}
	return gitbackend.AuthConfig{Type: gitbackend.AuthNone}
}

func (e *Executor) cloneRepo(ctx context.Context, dir string, repo *model.Repo, task *model.SyncTask, details *strings.Builder) error {
	err := e.backend.Clone(ctx, gitbackend.CloneOptions{
		URL:          repo.CloneURL,
		Path:         dir,
		Branch:       task.SourceBranch,
		Depth:        1,
		SingleBranch: true,
		Auth:         e.authConfig(repo),
	})
	if err != nil {
		fmt.Fprintf(details, "clone error: %v\n", err)
	}
	return err
}

func (e *Executor) fetchRepo(ctx context.Context, dir string, task *model.SyncTask, details *strings.Builder) error {
	_, err := e.backend.Fetch(ctx, gitbackend.FetchOptions{
		RepoPath: dir,
		Remote:   "origin",
		Branches: []string{task.SourceBranch},
		Tags:     task.GitTags,
		Prune:    task.GitPrune,
		Auth:     gitbackend.AuthConfig{Type: gitbackend.AuthNone},
	})
	if err != nil {
		fmt.Fprintf(details, "fetch error: %v\n", err)
		return err
	}

	if err := e.backend.Checkout(ctx, dir, task.SourceBranch); err != nil {
		fmt.Fprintf(details, "checkout error: %v\n", err)
		return err
	}

	return nil
}

func (e *Executor) ensureRemote(ctx context.Context, dir string, repo *model.Repo, details *strings.Builder) error {
	remotes, err := e.backend.ListRemoteBranches(ctx, dir, "target")
	if err != nil || len(remotes) == 0 {
		return e.backend.AddRemote(ctx, dir, "target", repo.CloneURL)
	}
	return nil
}

func (e *Executor) push(ctx context.Context, dir string, task *model.SyncTask, repo *model.Repo, details *strings.Builder) error {
	refSpec := fmt.Sprintf("%s:%s", task.SourceBranch, task.TargetBranch)

	_, err := e.backend.Push(ctx, gitbackend.PushOptions{
		RepoPath: dir,
		Remote:   "target",
		RefSpecs: []string{refSpec},
		Force:    task.GitForce,
		Auth:     e.authConfig(repo),
	})
	if err != nil {
		fmt.Fprintf(details, "push error: %v\n", err)
	}
	return err
}

func timePtr(t time.Time) *time.Time {
	return &t
}

type SyncPreview struct {
	CanSync       bool
	SourceExists  bool
	TargetExists  bool
	CommitsBehind int
	CommitsAhead  int
	LatestCommit  string
	Message       string
}

func (e *Executor) Preview(ctx context.Context, task *model.SyncTask) (*SyncPreview, error) {
	preview := &SyncPreview{}

	sourceRepo, err := e.service.RepoDAO().FindByKey(task.SourceRepoKey)
	if err != nil {
		preview.Message = fmt.Sprintf("source repo error: %v", err)
		return preview, nil
	}
	preview.SourceExists = sourceRepo != nil
	if !preview.SourceExists {
		preview.Message = "source repo not found"
		return preview, nil
	}

	targetRepo, err := e.service.RepoDAO().FindByKey(task.TargetRepoKey)
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
