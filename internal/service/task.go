package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
)

func (s *Service) ListTasks(ctx context.Context, repoKey string, offset, limit int) ([]*model.SyncTask, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	if repoKey != "" {
		return s.taskDAO.FindByRepoKey(repoKey, page)
	}
	return s.taskDAO.FindAll(page)
}

func (s *Service) GetTask(ctx context.Context, key string) (*model.SyncTask, error) {
	return s.taskDAO.FindByKey(key)
}

func (s *Service) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.SyncTask, error) {
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

	if err := s.taskDAO.Create(task); err != nil {
		return nil, err
	}

	if task.Cron != "" {
		if err := s.addCronJob(task); err != nil {
			return nil, fmt.Errorf("create cron job failed: %w", err)
		}
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, req *model.UpdateTaskRequest) (*model.SyncTask, error) {
	task, err := s.taskDAO.FindByKey(req.Key)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
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

	if err := s.taskDAO.Update(task); err != nil {
		return nil, err
	}

	if task.Cron != "" && task.Enabled {
		if err := s.addCronJob(task); err != nil {
			return nil, fmt.Errorf("update cron job failed: %w", err)
		}
	} else {
		s.removeCronJob(task.Key)
	}

	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, key string) error {
	s.removeCronJob(key)
	return s.taskDAO.Delete(key)
}

func (s *Service) RunTask(ctx context.Context, taskKey string) error {
	return s.RunTaskWithTrigger(ctx, taskKey, "manual")
}

func (s *Service) RunTaskWithTrigger(ctx context.Context, taskKey, trigger string) error {
	task, err := s.taskDAO.FindByKey(taskKey)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	_, err = s.executor.Execute(ctx, task, trigger)
	return err
}

func (s *Service) PreviewSync(ctx context.Context, req *model.PreviewSyncRequest) (*model.PreviewSyncResult, error) {
	sourceRepo, err := s.repoDAO.FindByKey(req.SourceRepoKey)
	if err != nil {
		return nil, err
	}
	targetRepo, err := s.repoDAO.FindByKey(req.TargetRepoKey)
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

func (s *Service) ListHistory(ctx context.Context, taskKey string, offset, limit int) ([]*model.SyncRun, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return s.runDAO.FindByTaskKey(taskKey, page)
}

func (s *Service) DeleteHistory(ctx context.Context, id uint) error {
	return s.runDAO.Delete(id)
}
