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
		ctx2 := context.WithValue(ctx, lockValueKey{}, value)
		*lockValueFromContext(ctx2) = value
	}
	return ok, nil
}

type lockValueKey struct{}

func lockValueFromContext(ctx context.Context) *string {
	val := ctx.Value(lockValueKey{})
	if val == nil {
		return new(string)
	}
	return val.(*string)
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
		ok, _, err := l.LockWithTTL(ctx, key, defaultLockTTL)
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
	value := lockValueFromContext(ctx)
	if *value == "" {
		_, err := l.client.Del(ctx, lockKey).Result()
		if err != nil {
			return fmt.Errorf("redis del failed: %w", err)
		}
		return nil
	}
	return l.UnlockWithValue(ctx, key, *value)
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
