package lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultLockTTL = 30 * time.Second
)

var unlockScript = redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

var semaphoreAcquireScript = redis.NewScript(`
local key = KEYS[1]
local member = ARGV[1]
local max = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 清理过期的成员
redis.call('ZREMRANGEBYSCORE', key, '-inf', now - 1)

-- 检查当前数量
local current = redis.call('ZCARD', key)
if current >= max then
    return 0
end

-- 添加新成员
redis.call('ZADD', key, now, member)
return 1
`)

type DistLock interface {
	TryLock(ctx context.Context, key string) (bool, error)
	TryLockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string) error
}

type RedisLock struct {
	client     *redis.Client
	mu         sync.Mutex
	lockValues map[string]string
}

func NewRedisLock(addr, password string, db int) *RedisLock {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisLock{
		client:     client,
		lockValues: make(map[string]string),
	}
}

func generateLockValue() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (l *RedisLock) TryLock(ctx context.Context, key string) (bool, error) {
	return l.TryLockWithTTL(ctx, key, defaultLockTTL)
}

func (l *RedisLock) TryLockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := "git-sync:lock:" + key
	value := generateLockValue()
	ok, err := l.client.SetNX(ctx, lockKey, value, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx failed: %w", err)
	}
	if ok {
		l.mu.Lock()
		l.lockValues[key] = value
		l.mu.Unlock()
	}
	return ok, nil
}


func (l *RedisLock) LockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, string, error) {
	lockKey := "git-sync:lock:" + key
	value := generateLockValue()
	ok, err := l.client.SetNX(ctx, lockKey, value, ttl).Result()
	if err != nil {
		return false, "", fmt.Errorf("redis setnx failed: %w", err)
	}
	return ok, value, nil
}

func (l *RedisLock) UnlockWithValue(ctx context.Context, key, value string) error {
	lockKey := "git-sync:lock:" + key
	_, err := unlockScript.Run(ctx, l.client, []string{lockKey}, value).Result()
	if err != nil {
		return fmt.Errorf("redis unlock failed: %w", err)
	}
	return nil
}

func (l *RedisLock) Lock(ctx context.Context, key string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ok, _, err := l.LockWithTTL(ctx, key, defaultLockTTL)
		if err != nil {
			// If context was cancelled during the operation, return context error
			if ctx.Err() != nil {
				return ctx.Err()
			}
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
	l.mu.Lock()
	value, exists := l.lockValues[key]
	if exists {
		delete(l.lockValues, key)
	}
	l.mu.Unlock()

	if !exists || value == "" {
		lockKey := "git-sync:lock:" + key
		_, err := l.client.Del(ctx, lockKey).Result()
		if err != nil {
			return fmt.Errorf("redis del failed: %w", err)
		}
		return nil
	}
	return l.UnlockWithValue(ctx, key, value)
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
	now := time.Now().Unix()
	result, err := semaphoreAcquireScript.Run(ctx, s.client, []string{s.key}, identifier, s.max, now).Int()
	if err != nil {
		return false, fmt.Errorf("semaphore acquire failed: %w", err)
	}
	return result == 1, nil
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

type localLockEntry struct {
	expiresAt time.Time
}

type LocalLock struct {
	mu      sync.Mutex
	entries map[string]*localLockEntry
}

func NewLocalLock() *LocalLock {
	return &LocalLock{
		entries: make(map[string]*localLockEntry),
	}
}

func (l *LocalLock) TryLock(ctx context.Context, key string) (bool, error) {
	return l.TryLockWithTTL(ctx, key, defaultLockTTL)
}

func (l *LocalLock) TryLockWithTTL(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry, exists := l.entries[key]
	if !exists || now.After(entry.expiresAt) {
		l.entries[key] = &localLockEntry{expiresAt: now.Add(ttl)}
		return true, nil
	}
	return false, nil
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
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.entries, key)
	return nil
}
