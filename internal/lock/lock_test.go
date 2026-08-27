package lock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalLock_TryLock(t *testing.T) {
	lock := NewLocalLock()
	ctx := context.Background()

	ok, err := lock.TryLock(ctx, "test-key")
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock")

	ok, err = lock.TryLock(ctx, "test-key")
	require.NoError(t, err, "TryLock failed")
	require.False(t, ok, "Expected lock to be already held")
}

func TestLocalLock_TryLockWithTTL(t *testing.T) {
	lock := NewLocalLock()
	ctx := context.Background()

	ok, err := lock.TryLockWithTTL(ctx, "test-key", 100*time.Millisecond)
	require.NoError(t, err, "TryLockWithTTL failed")
	require.True(t, ok, "Expected to acquire lock")

	ok, err = lock.TryLockWithTTL(ctx, "test-key", 100*time.Millisecond)
	require.NoError(t, err, "TryLockWithTTL failed")
	require.False(t, ok, "Expected lock to be already held")

	time.Sleep(150 * time.Millisecond)

	ok, err = lock.TryLockWithTTL(ctx, "test-key", 100*time.Millisecond)
	require.NoError(t, err, "TryLockWithTTL failed")
	require.True(t, ok, "Expected to acquire lock after TTL expired")
}

func TestLocalLock_Unlock(t *testing.T) {
	lock := NewLocalLock()
	ctx := context.Background()

	ok, err := lock.TryLock(ctx, "test-key")
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock")

	err = lock.Unlock(ctx, "test-key")
	require.NoError(t, err, "Unlock failed")

	ok, err = lock.TryLock(ctx, "test-key")
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock after unlock")
}

func TestLocalLock_Concurrent(t *testing.T) {
	lock := NewLocalLock()
	ctx := context.Background()

	acquired := make(chan bool, 2)

	go func() {
		ok, _ := lock.TryLock(ctx, "concurrent-key")
		acquired <- ok
		if ok {
			time.Sleep(50 * time.Millisecond)
			_ = lock.Unlock(ctx, "concurrent-key")
		}
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ok, _ := lock.TryLock(ctx, "concurrent-key")
		acquired <- ok
	}()

	results := make([]bool, 2)
	for i := 0; i < 2; i++ {
		results[i] = <-acquired
	}

	require.True(t, results[0] || results[1], "At least one goroutine should have acquired the lock")
}

// Helper function to create a test Redis client
func getTestRedisClient(t *testing.T) *redis.Client {
	t.Helper()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use a different DB for tests
	})
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping test")
	}
	return client
}

func closeClient(t *testing.T, client *redis.Client) {
	t.Helper()
	err := client.Close()
	assert.NoError(t, err, "Failed to close client")
}

func TestSemaphore_Acquire(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-%d", time.Now().UnixNano())
	sem := NewSemaphore(client, key, 2, 10*time.Second)
	defer client.Del(ctx, sem.key)

	// Should acquire first two slots
	ok1, err := sem.Acquire(ctx, "worker-1")
	require.NoError(t, err, "Acquire worker-1 failed")
	require.True(t, ok1, "Expected worker-1 to acquire semaphore")

	ok2, err := sem.Acquire(ctx, "worker-2")
	require.NoError(t, err, "Acquire worker-2 failed")
	require.True(t, ok2, "Expected worker-2 to acquire semaphore")

	// Third acquire should fail (max is 2)
	ok3, err := sem.Acquire(ctx, "worker-3")
	require.NoError(t, err, "Acquire worker-3 failed")
	require.False(t, ok3, "Expected worker-3 to fail acquiring semaphore")
}

func TestSemaphore_Release(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-release-%d", time.Now().UnixNano())
	sem := NewSemaphore(client, key, 1, 10*time.Second)
	defer client.Del(ctx, sem.key)

	// Acquire the single slot
	ok, err := sem.Acquire(ctx, "worker-1")
	require.NoError(t, err, "Acquire failed")
	require.True(t, ok, "Expected to acquire semaphore")

	// Release the slot
	err = sem.Release(ctx, "worker-1")
	require.NoError(t, err, "Release failed")

	// Should be able to acquire again
	ok, err = sem.Acquire(ctx, "worker-2")
	require.NoError(t, err, "Acquire after release failed")
	require.True(t, ok, "Expected to acquire semaphore after release")
}

func TestSemaphore_Concurrent(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-concurrent-%d", time.Now().UnixNano())
	maxSlots := 3
	sem := NewSemaphore(client, key, maxSlots, 10*time.Second)
	defer client.Del(ctx, sem.key)

	numWorkers := 10
	var wg sync.WaitGroup
	results := make([]bool, numWorkers)
	errors := make([]error, numWorkers)

	// Launch concurrent workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			workerID := fmt.Sprintf("worker-%d", id)
			ok, err := sem.Acquire(ctx, workerID)
			results[id] = ok
			errors[id] = err
		}(i)
	}

	wg.Wait()

	// Count successful acquisitions
	acquired := 0
	for i := 0; i < numWorkers; i++ {
		assert.NoError(t, errors[i], "Worker %d got error", i)
		if results[i] {
			acquired++
		}
	}

	// Exactly maxSlots should have acquired
	assert.Equal(t, maxSlots, acquired, "Expected %d workers to acquire, got %d", maxSlots, acquired)
}

func TestSemaphore_Cleanup(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-cleanup-%d", time.Now().UnixNano())
	sem := NewSemaphore(client, key, 3, 10*time.Second)
	defer client.Del(ctx, sem.key)

	// 直接构造槽位:score 为过期毫秒时间戳(Acquire 打分语义)。
	// stale-2h/stale-1h 已过期,live 尚未过期。
	now := time.Now()
	seed := map[string]time.Time{
		"stale-2h": now.Add(-2 * time.Hour),
		"stale-1h": now.Add(-1 * time.Hour),
		"live":     now.Add(10 * time.Second),
	}
	for member, expireAt := range seed {
		err := client.ZAdd(ctx, sem.key, redis.Z{
			Score:  float64(expireAt.UnixMilli()),
			Member: member,
		}).Err()
		require.NoError(t, err, "seed %s failed", member)
	}

	// Cleanup(90min): 清理"过期超过 90 分钟"的条目 → 只清 stale-2h
	err := sem.Cleanup(ctx, 90*time.Minute)
	require.NoError(t, err, "Cleanup failed")

	card, err := client.ZCard(ctx, sem.key).Result()
	require.NoError(t, err, "ZCard failed")
	require.Equal(t, int64(2), card, "Expected 2 members after Cleanup(90min)")

	// Cleanup(0): 清理所有已过期条目 → 只剩 live
	err = sem.Cleanup(ctx, 0)
	require.NoError(t, err, "Cleanup with 0 duration failed")

	members, _, err := client.ZScan(ctx, sem.key, 0, "", 0).Result()
	require.NoError(t, err, "ZScan failed")
	// ZScan 返回 [member, score] 交错列表
	require.Len(t, members, 2, "Expected only live member after Cleanup(0)")
	require.Equal(t, "live", members[0], "Expected only live member after Cleanup(0)")

	// live 的 score 不应被改动(仍为未来的过期时间)
	score, err := client.ZScore(ctx, sem.key, "live").Result()
	require.NoError(t, err, "ZScore failed")
	require.Greater(t, int64(score), now.UnixMilli(), "live score should stay in the future")
}
