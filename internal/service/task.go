package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/yi-nology/git-sync-service/sync/model"
)

// ListTasks returns a paginated list of sync tasks.
func (s *Service) ListTasks(ctx context.Context, repoKey string, offset, limit int) ([]*model.SyncTask, int64, error) {
	return s.tasks.ListTasks(ctx, repoKey, offset, limit)
}

// GetTask returns a sync task by key.
func (s *Service) GetTask(ctx context.Context, key string) (*model.SyncTask, error) {
	return s.tasks.GetTask(ctx, key)
}

func (s *Service) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.SyncTask, error) {
	task, err := s.tasks.CreateTask(ctx, req)
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
	task, err := s.tasks.UpdateTask(ctx, req)
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
	if err := s.tasks.DeleteTask(ctx, key); err != nil {
		return err
	}
	// 清理该任务残留的 workdir(成功执行后保留的增量 fetch 工作区),避免磁盘泄漏
	if err := os.RemoveAll(s.GetTempDir(key)); err != nil {
		slog.Error("failed to cleanup workdir on task delete", "error", err, "taskKey", key)
	}
	return nil
}

func (s *Service) RunTask(ctx context.Context, taskKey string) error {
	return s.RunTaskWithTrigger(ctx, taskKey, model.TriggerManual, nil)
}

// CountTasksByStatus returns task counts grouped by last_status (key "total" = overall).
func (s *Service) CountTasksByStatus() (map[string]int64, error) {
	return s.tasks.CountTasksByStatus()
}

func (s *Service) RunTaskWithTrigger(ctx context.Context, taskKey, trigger string, webhookEventID *uint) error {
	task, err := s.tasks.FindTaskByKey(taskKey)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	// 已禁用的任务不执行(cron/webhook/手动均拦截),需要运行请先启用
	if !task.Enabled {
		return ErrTaskDisabled
	}

	// 全局并发上限(非阻塞):达到 MaxConcurrent 直接拒绝,保护 git 平台与内存
	select {
	case s.concurrencySem <- struct{}{}:
		defer func() { <-s.concurrencySem }()
	default:
		return ErrTooManyConcurrent
	}

	// 同任务互斥(非阻塞):已有执行流在跑则跳过,避免并发写同一 worktree
	// 注意:只有成功 store 后才注册 Delete,loaded 返回时不能误删别人的标记
	if _, loaded := s.runningTasks.LoadOrStore(taskKey, struct{}{}); loaded {
		return ErrTaskRunning
	}
	defer s.runningTasks.Delete(taskKey)

	_, err = s.executor.Execute(ctx, task, trigger, webhookEventID)
	return err
}

// PreviewSync previews a sync operation.
func (s *Service) PreviewSync(ctx context.Context, req *model.PreviewSyncRequest) (*model.PreviewSyncResult, error) {
	return s.tasks.PreviewSync(ctx, req)
}

// ListHistory returns a paginated list of sync run history.
func (s *Service) ListHistory(ctx context.Context, taskKey string, offset, limit int) ([]*model.SyncRun, int64, error) {
	return s.tasks.ListHistory(ctx, taskKey, offset, limit)
}

// DeleteHistory deletes a sync run by ID.
func (s *Service) DeleteHistory(ctx context.Context, id uint) error {
	return s.tasks.DeleteHistory(ctx, id)
}
