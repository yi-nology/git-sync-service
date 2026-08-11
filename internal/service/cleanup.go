package service

import (
	"log/slog"
	"time"
)

func (s *Service) CleanupOldData(maxAge time.Duration) (events, runs, steps int64, err error) {
	events, err = s.webhooks.CleanupOldEvents(maxAge)
	if err != nil {
		return 0, 0, 0, err
	}
	runs, err = s.tasks.CleanupOldRuns(maxAge)
	if err != nil {
		return events, 0, 0, err
	}
	steps, err = s.tasks.CleanupOldRunSteps(maxAge)
	if err != nil {
		return events, runs, 0, err
	}
	slog.Info("data cleanup completed", "events_deleted", events, "runs_deleted", runs, "steps_deleted", steps)
	return events, runs, steps, nil
}

func (s *Service) cleanupTriggerTimes() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.lastTriggerTime.Range(func(key, value interface{}) bool {
				if t, ok := value.(time.Time); ok && now.Sub(t) > 1*time.Hour {
					s.lastTriggerTime.Delete(key)
				}
				return true
			})
		case <-s.cleanupDone:
			return
		}
	}
}
