package service

import (
	"context"
	"log/slog"

	"github.com/yi-nology/git-sync-service/sync/model"
)

func (s *Service) addCronJob(task *model.SyncTask) error {
	s.cronMu.Lock()
	defer s.cronMu.Unlock()

	if entryID, ok := s.cronEntryIDs[task.Key]; ok {
		s.cron.Remove(entryID)
	}

	entryID, err := s.cron.AddFunc(task.Cron, func() {
		ctx := context.Background()
			if err := s.RunTaskWithTrigger(ctx, task.Key, model.TriggerCron, nil); err != nil {
			slog.Error("cron task failed", "taskKey", task.Key, "error", err)
		}
	})
	if err != nil {
		return err
	}

	s.cronEntryIDs[task.Key] = entryID
	return nil
}

func (s *Service) removeCronJob(taskKey string) {
	s.cronMu.Lock()
	defer s.cronMu.Unlock()

	if entryID, ok := s.cronEntryIDs[taskKey]; ok {
		s.cron.Remove(entryID)
		delete(s.cronEntryIDs, taskKey)
	}
}

func (s *Service) startCronJobs() error {
	tasks, err := s.tasks.FindAllEnabledTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Cron != "" {
			if err := s.addCronJob(task); err != nil {
				slog.Error("add cron job failed", "taskKey", task.Key, "error", err)
			}
		}
	}

	s.cron.Start()
	return nil
}

func (s *Service) stopCronJobs() {
	s.cron.Stop()
}
