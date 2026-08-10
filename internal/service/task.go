package service

import (
	"context"
	"fmt"

	"github.com/yi-nology/git-sync-service/sync/model"
)

func (s *Service) ListTasks(ctx context.Context, repoKey string, offset, limit int) ([]*model.SyncTask, int64, error) {
	return s.taskService.ListTasks(ctx, repoKey, offset, limit)
}

func (s *Service) GetTask(ctx context.Context, key string) (*model.SyncTask, error) {
	return s.taskService.GetTask(ctx, key)
}

func (s *Service) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.SyncTask, error) {
	task, err := s.taskService.CreateTask(ctx, req)
	if err != nil {
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
	task, err := s.taskService.UpdateTask(ctx, req)
	if err != nil {
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
	return s.taskService.DeleteTask(ctx, key)
}

func (s *Service) RunTask(ctx context.Context, taskKey string) error {
	return s.RunTaskWithTrigger(ctx, taskKey, "manual")
}

func (s *Service) RunTaskWithTrigger(ctx context.Context, taskKey, trigger string) error {
	task, err := s.taskService.FindTaskByKey(taskKey)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	_, err = s.executor.Execute(ctx, task, trigger)
	return err
}

func (s *Service) PreviewSync(ctx context.Context, req *model.PreviewSyncRequest) (*model.PreviewSyncResult, error) {
	return s.taskService.PreviewSync(ctx, req)
}

func (s *Service) ListHistory(ctx context.Context, taskKey string, offset, limit int) ([]*model.SyncRun, int64, error) {
	return s.taskService.ListHistory(ctx, taskKey, offset, limit)
}

func (s *Service) DeleteHistory(ctx context.Context, id uint) error {
	return s.taskService.DeleteHistory(ctx, id)
}
