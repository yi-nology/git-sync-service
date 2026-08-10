package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// TaskService handles sync task-related operations.
type TaskService struct {
	taskDAO *dao.SyncTaskDAO
	runDAO  *dao.SyncRunDAO
	repoDAO *dao.RepoDAO
}

// NewTaskService creates a new TaskService instance.
func NewTaskService(taskDAO *dao.SyncTaskDAO, runDAO *dao.SyncRunDAO, repoDAO *dao.RepoDAO) *TaskService {
	return &TaskService{
		taskDAO: taskDAO,
		runDAO:  runDAO,
		repoDAO: repoDAO,
	}
}

// ListTasks returns a paginated list of sync tasks.
func (ts *TaskService) ListTasks(ctx context.Context, repoKey string, offset, limit int) ([]*model.SyncTask, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	if repoKey != "" {
		return ts.taskDAO.FindByRepoKey(repoKey, page)
	}
	return ts.taskDAO.FindAll(page)
}

// GetTask returns a sync task by key.
func (ts *TaskService) GetTask(ctx context.Context, key string) (*model.SyncTask, error) {
	return ts.taskDAO.FindByKey(key)
}

// CreateTask creates a new sync task.
func (ts *TaskService) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.SyncTask, error) {
	task := &model.SyncTask{
		Key:           uuid.New().String(),
		Name:          req.Name,
		SourceRepoKey: req.SourceRepoKey,
		SourceBranch:  req.SourceBranch,
		TargetRepoKey: req.TargetRepoKey,
		TargetBranch:  req.TargetBranch,
		SyncMode:      req.SyncMode,
		Cron:          req.Cron,
		WebhookToken:  uuid.New().String(),
		Enabled:       true,
		GitTags:       req.GitTags,
		GitForce:      req.GitForce,
		GitPrune:      req.GitPrune,
		GitNoVerify:   req.GitNoVerify,
		PushOptions:   req.PushOptions,
	}

	if err := ts.taskDAO.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask updates an existing sync task.
func (ts *TaskService) UpdateTask(ctx context.Context, req *model.UpdateTaskRequest) (*model.SyncTask, error) {
	task, err := ts.taskDAO.FindByKey(req.Key)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.SourceBranch != "" {
		task.SourceBranch = req.SourceBranch
	}
	if req.TargetBranch != "" {
		task.TargetBranch = req.TargetBranch
	}
	if req.SyncMode != "" {
		task.SyncMode = req.SyncMode
	}
	if req.Cron != "" {
		task.Cron = req.Cron
	}
	task.Enabled = req.Enabled
	task.GitTags = req.GitTags
	task.GitForce = req.GitForce
	task.GitPrune = req.GitPrune
	task.GitNoVerify = req.GitNoVerify
	task.PushOptions = req.PushOptions

	if err := ts.taskDAO.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask deletes a sync task by key.
func (ts *TaskService) DeleteTask(ctx context.Context, key string) error {
	return ts.taskDAO.Delete(key)
}

// PreviewSync previews a sync operation.
func (ts *TaskService) PreviewSync(ctx context.Context, req *model.PreviewSyncRequest) (*model.PreviewSyncResult, error) {
	sourceRepo, err := ts.repoDAO.FindByKey(req.SourceRepoKey)
	if err != nil {
		return nil, err
	}
	targetRepo, err := ts.repoDAO.FindByKey(req.TargetRepoKey)
	if err != nil {
		return nil, err
	}

	result := &model.PreviewSyncResult{
		SourceExists: sourceRepo != nil,
		TargetExists: targetRepo != nil,
	}

	if sourceRepo == nil || targetRepo == nil {
		result.CanSync = false
		result.Message = "source or target repo not found"
		return result, nil
	}

	result.CanSync = true
	result.Message = "sync preview ready"
	return result, nil
}

// ListHistory returns a paginated list of sync run history.
func (ts *TaskService) ListHistory(ctx context.Context, taskKey string, offset, limit int) ([]*model.SyncRun, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return ts.runDAO.FindByTaskKey(taskKey, page)
}

// DeleteHistory deletes a sync run by ID.
func (ts *TaskService) DeleteHistory(ctx context.Context, id uint) error {
	return ts.runDAO.Delete(id)
}

// FindAllEnabledTasks returns all enabled sync tasks.
func (ts *TaskService) FindAllEnabledTasks() ([]*model.SyncTask, error) {
	return ts.taskDAO.FindAllEnabled()
}

// FindTaskByKey returns a sync task by key (internal use).
func (ts *TaskService) FindTaskByKey(key string) (*model.SyncTask, error) {
	return ts.taskDAO.FindByKey(key)
}

// UpdateTaskStatus updates the task's last run status.
func (ts *TaskService) UpdateTaskStatus(task *model.SyncTask) error {
	return ts.taskDAO.Update(task)
}

// CreateRun creates a new sync run record.
func (ts *TaskService) CreateRun(run *model.SyncRun) error {
	return ts.runDAO.Create(run)
}

// UpdateRun updates a sync run record.
func (ts *TaskService) UpdateRun(run *model.SyncRun) error {
	return ts.runDAO.Update(run)
}

// CleanupOldRuns removes sync runs older than the specified duration.
func (ts *TaskService) CleanupOldRuns(maxAge time.Duration) (int64, error) {
	return ts.runDAO.CleanupOlderThan(maxAge)
}

// RunTaskWithTrigger runs a sync task with the specified trigger source.
// This method requires access to the executor, so it takes a callback function.
func (ts *TaskService) RunTaskWithTrigger(ctx context.Context, taskKey, trigger string, executeFn func(ctx context.Context, task *model.SyncTask, trigger string) (*model.SyncRun, error)) error {
	task, err := ts.taskDAO.FindByKey(taskKey)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	_, err = executeFn(ctx, task, trigger)
	return err
}

// FormatCronJobKey formats a cron job key for a task.
func FormatCronJobKey(taskKey string) string {
	return fmt.Sprintf("cron:%s", taskKey)
}
