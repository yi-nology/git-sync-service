package service

import (
	"context"

	"github.com/yi-nology/git-sync-service/sync/model"
)

// ListAllTasks 获取所有任务
func (s *Service) ListAllTasks(ctx context.Context) ([]*model.SyncTask, error) {
	return s.taskDAO.FindAll()
}

// ListRecentHistory 获取最近的执行历史
func (s *Service) ListRecentHistory(ctx context.Context, limit int) ([]*model.SyncRun, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.runDAO.FindRecent(limit)
}

// ListHistoryByTask 获取指定任务的执行历史
func (s *Service) ListHistoryByTask(ctx context.Context, taskKey string, limit int) ([]*model.SyncRun, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.runDAO.FindByTaskKey(taskKey, limit)
}
