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

// extendScript 仅当锁的 value 匹配(自己是持有者)才续期,避免误续别人的锁。
var extendScript = redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

var semaphoreAcquireScript = redis.NewScript(`
local key = KEYS[1]
local member = ARGV[1]
local max = tonumber(ARGV[2])
local now = tonumber(ARGV[3])      -- 当前时间(毫秒)
local expire = tonumber(ARGV[4])   -- 槽过期时间(毫秒)

-- 清理已过期的成员(score = 过期毫秒时间戳 <= now 即已过期)
redis.call('ZREMRANGEBYSCORE', key, '-inf', now)

local current = redis.call('ZCARD', key)
if current >= max then
    return 0
end

-- 占用一个槽,score 设为过期时间戳,到点未续期则被清理(进程崩溃也不会永久占槽)
redis.call('ZADD', key, expire, member)
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

// ExtendLockWithValue 续期:仅当 value 匹配(自己持有)才延长 TTL,返回是否成功续期。
// 供 watchdog 在长任务期间周期性续期,避免锁过期被别人抢走。
func (l *RedisLock) ExtendLockWithValue(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	lockKey := "git-sync:lock:" + key
	res, err := extendScript.Run(ctx, l.client, []string{lockKey}, value, ttl.Milliseconds()).Int()
	if err != nil {
		return false, fmt.Errorf("redis extend failed: %w", err)
	}
	return res == 1, nil
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
	client  *redis.Client
	key     string
	max     int
	slotTTL time.Duration // 单个槽的存活时长,持有期间需 Renew 续期,否则到期自动释放
}

// NewSemaphore 创建一个基于 redis ZSET 的分布式信号量。
// slotTTL 为单个槽的存活时长:持有者需周期性 Renew,进程崩溃后槽在 slotTTL 后自动释放。
func NewSemaphore(client *redis.Client, key string, max int, slotTTL time.Duration) *Semaphore {
	return &Semaphore{
		client:  client,
		key:     "git-sync:semaphore:" + key,
		max:     max,
		slotTTL: slotTTL,
	}
}

func (s *Semaphore) Acquire(ctx context.Context, identifier string) (bool, error) {
	now := time.Now().UnixMilli()
	expire := time.Now().Add(s.slotTTL).UnixMilli()
	result, err := semaphoreAcquireScript.Run(ctx, s.client, []string{s.key}, identifier, s.max, now, expire).Int()
	if err != nil {
		return false, fmt.Errorf("semaphore acquire failed: %w", err)
	}
	return result == 1, nil
}

// Renew 续期持有的槽(把 score 推到 now+slotTTL),供 watchdog 调用。
func (s *Semaphore) Renew(ctx context.Context, identifier string) error {
	expire := time.Now().Add(s.slotTTL).UnixMilli()
	_, err := s.client.ZAdd(ctx, s.key, redis.Z{Score: float64(expire), Member: identifier}).Result()
	return err
}

func (s *Semaphore) Release(ctx context.Context, identifier string) error {
	_, err := s.client.ZRem(ctx, s.key, identifier).Result()
	return err
}

func (s *Semaphore) Cleanup(ctx context.Context, olderThan time.Duration) error {
	// 注意:槽位 score 是毫秒时间戳(Acquire 用 UnixMilli),阈值必须同为毫秒,
	// 否则 ZRemRangeByScore 永远删不到任何条目(此前用 Unix() 秒导致差 1000 倍)
	cutoff := float64(time.Now().Add(-olderThan).UnixMilli())
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
