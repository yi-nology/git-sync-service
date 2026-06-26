package service

import (
	"context"
	"fmt"
	"time"
)

func (s *Service) TryAcquireTaskLock(ctx context.Context, taskKey string) (bool, error) {
	lockKey := fmt.Sprintf("task:%s", taskKey)
	ttl := time.Duration(s.config.Sync.DefaultTimeout) * time.Second
	if ttl <= 0 {
		ttl = 300 * time.Second
	}
	return s.lock.TryLockWithTTL(ctx, lockKey, ttl)
}

func (s *Service) ReleaseTaskLock(ctx context.Context, taskKey string) error {
	lockKey := fmt.Sprintf("task:%s", taskKey)
	return s.lock.Unlock(ctx, lockKey)
}

func (s *Service) AcquireSemaphore(ctx context.Context, taskKey string) (bool, error) {
	if s.semaphore == nil {
		return true, nil
	}
	return s.semaphore.Acquire(ctx, s.semaphoreID+":"+taskKey)
}

func (s *Service) ReleaseSemaphore(ctx context.Context, taskKey string) error {
	if s.semaphore != nil {
		return s.semaphore.Release(ctx, s.semaphoreID+":"+taskKey)
	}
	return nil
}
