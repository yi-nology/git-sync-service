package service

import (
	"log/slog"
	"time"
)

func (s *Service) CleanupOldData(maxAge time.Duration) (events, runs, steps int64, err error) {
	// 三张表独立清理,互不依赖,并行执行缩短耗时。
	type result struct {
		count int64
		err   error
	}
	ch := make(chan result, 3)

	go func() {
		c, e := s.webhooks.CleanupOldEvents(maxAge)
		ch <- result{c, e}
	}()
	go func() {
		c, e := s.tasks.CleanupOldRuns(maxAge)
		ch <- result{c, e}
	}()
	go func() {
		c, e := s.tasks.CleanupOldRunSteps(maxAge)
		ch <- result{c, e}
	}()

	r1, r2, r3 := <-ch, <-ch, <-ch
	events, runs, steps = r1.count, r2.count, r3.count

	for _, r := range []result{r1, r2, r3} { // 任一失败都报告,但不中断其余结果
		if r.err != nil {
			err = r.err
		}
	}
	slog.Info("data cleanup completed", "events_deleted", events, "runs_deleted", runs, "steps_deleted", steps)
	return events, runs, steps, err
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
