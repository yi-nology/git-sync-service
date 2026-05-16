package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/sync/model"
)

func (s *Service) ListTasks(repoKey string) ([]*model.SyncTask, error) {
	if repoKey != "" {
		return s.taskDAO.FindByRepoKey(repoKey)
	}
	tasks, err := s.taskDAO.FindAllEnabled()
	if err != nil {
		return nil, err
	}
	var allTasks []*model.SyncTask
	for _, t := range tasks {
		allTasks = append(allTasks, t)
	}
	return allTasks, nil
}

func (s *Service) GetTask(key string) (*model.SyncTask, error) {
	return s.taskDAO.FindByKey(key)
}

func (s *Service) CreateTask(req *model.CreateTaskRequest) (*model.SyncTask, error) {
	task := &model.SyncTask{
		Key:           uuid.New().String(),
		Name:          req.Name,
		SourceRepoKey: req.SourceRepoKey,
		SourceBranch:  req.SourceBranch,
		TargetRepoKey: req.TargetRepoKey,
		TargetBranch:  req.TargetBranch,
		SyncMode:      req.SyncMode,
		Cron:          req.Cron,
		WebhookToken:  uuid.New().String()[:16],
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

func (s *Service) UpdateTask(req *model.UpdateTaskRequest) (*model.SyncTask, error) {
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

	if task.Cron != "" {
		if err := s.addCronJob(task); err != nil {
			return nil, fmt.Errorf("update cron job failed: %w", err)
		}
	}

	return task, nil
}

func (s *Service) DeleteTask(key string) error {
	if entryID, ok := s.cronEntryIDs[key]; ok {
		s.cron.Remove(entryID)
		delete(s.cronEntryIDs, key)
	}
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

func (s *Service) PreviewSync(req *model.PreviewSyncRequest) (*model.PreviewSyncResult, error) {
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

func (s *Service) ListHistory(taskKey string, limit int) ([]*model.SyncRun, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.runDAO.FindByTaskKey(taskKey, limit)
}

func (s *Service) DeleteHistory(id uint) error {
	return s.runDAO.Delete(id)
}
