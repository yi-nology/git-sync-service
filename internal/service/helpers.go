package service

import (
	"context"

	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
)

func (s *Service) ListAllTasks(ctx context.Context) ([]*model.SyncTask, error) {
	page := dao.DefaultPagination(0, 200)
	tasks, _, err := s.taskDAO.FindAll(page)
	return tasks, err
}

func (s *Service) ListRecentHistory(ctx context.Context, limit int) ([]*model.SyncRun, error) {
	page := dao.DefaultPagination(0, limit)
	return s.runDAO.FindRecent(page)
}

func (s *Service) ListHistoryByTask(ctx context.Context, taskKey string, limit int) ([]*model.SyncRun, error) {
	page := dao.DefaultPagination(0, limit)
	runs, _, err := s.runDAO.FindByTaskKey(taskKey, page)
	return runs, err
}
