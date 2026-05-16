package lock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultLockTTL = 30 * time.Second
)

type DistLock interface {
	TryLock(ctx context.Context, key string) (bool, error)
	TryLockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string) error
}

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock(addr, password string, db int) *RedisLock {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisLock{client: client}
}

func (l *RedisLock) TryLock(ctx context.Context, key string) (bool, error) {
	return l.TryLockWithTTL(ctx, key, defaultLockTTL)
}

func (l *RedisLock) TryLockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := "git-sync:lock:" + key
	ok, err := l.client.SetNX(ctx, lockKey, "1", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx failed: %w", err)
	}
	return ok, nil
}

func (l *RedisLock) Lock(ctx context.Context, key string) error {
	for {
		ok, err := l.TryLock(ctx, key)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (l *RedisLock) Unlock(ctx context.Context, key string) error {
	lockKey := "git-sync:lock:" + key
	_, err := l.client.Del(ctx, lockKey).Result()
	if err != nil {
		return fmt.Errorf("redis del failed: %w", err)
	}
	return nil
}

func (l *RedisLock) ExtendLock(ctx context.Context, key string, ttl time.Duration) error {
	lockKey := "git-sync:lock:" + key
	_, err := l.client.Expire(ctx, lockKey, ttl).Result()
	if err != nil {
		return fmt.Errorf("redis expire failed: %w", err)
	}
	return nil
}

func (l *RedisLock) Ping(ctx context.Context) error {
	return l.client.Ping(ctx).Err()
}

func (l *RedisLock) Close() error {
	return l.client.Close()
}

func (l *RedisLock) Client() *redis.Client {
	return l.client
}

type Semaphore struct {
	client *redis.Client
	key    string
	max    int
}

func NewSemaphore(client *redis.Client, key string, max int) *Semaphore {
	return &Semaphore{
		client: client,
		key:    "git-sync:semaphore:" + key,
		max:    max,
	}
}

func (s *Semaphore) Acquire(ctx context.Context, identifier string) (bool, error) {
	now := float64(time.Now().Unix())
	_, err := s.client.ZAdd(ctx, s.key, redis.Z{Score: now, Member: identifier}).Result()
	if err != nil {
		return false, err
	}

	rank, err := s.client.ZRank(ctx, s.key, identifier).Result()
	if err != nil {
		return false, err
	}

	if int(rank) < s.max {
		return true, nil
	}

	s.client.ZRem(ctx, s.key, identifier)
	return false, nil
}

func (s *Semaphore) Release(ctx context.Context, identifier string) error {
	_, err := s.client.ZRem(ctx, s.key, identifier).Result()
	return err
}

func (s *Semaphore) Cleanup(ctx context.Context, olderThan time.Duration) error {
	cutoff := float64(time.Now().Add(-olderThan).Unix())
	_, err := s.client.ZRemRangeByScore(ctx, s.key, "-inf", fmt.Sprintf("%f", cutoff)).Result()
	return err
}

type LocalLock struct {
	mu sync.Map
}

func (l *LocalLock) TryLock(ctx context.Context, key string) (bool, error) {
	return l.TryLockWithTTL(ctx, key, defaultLockTTL)
}

func (l *LocalLock) TryLockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	_, loaded := l.mu.LoadOrStore(key, time.Now().Add(ttl))
	return !loaded, nil
}

func (l *LocalLock) Lock(ctx context.Context, key string) error {
	for {
		ok, _ := l.TryLock(ctx, key)
		if ok {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (l *LocalLock) Unlock(ctx context.Context, key string) error {
	l.mu.Delete(key)
	return nil
}
