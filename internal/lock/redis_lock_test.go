package lock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRedisLock_TryLock(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-lock-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	// Clean up any existing lock
	_ = lock.Unlock(ctx, key)

	// Should acquire lock
	ok, err := lock.TryLock(ctx, key)
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock")

	// Should fail to acquire same lock
	ok, err = lock.TryLock(ctx, key)
	require.NoError(t, err, "TryLock failed")
	require.False(t, ok, "Expected lock to be already held")

	// Cleanup
	_ = lock.Unlock(ctx, key)
}

func TestRedisLock_TryLockWithTTL(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-lock-ttl-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	// Clean up any existing lock
	_ = lock.Unlock(ctx, key)

	// Acquire lock with short TTL
	ok, err := lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	require.NoError(t, err, "TryLockWithTTL failed")
	require.True(t, ok, "Expected to acquire lock")

	// Should fail to acquire same lock
	ok, err = lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	require.NoError(t, err, "TryLockWithTTL failed")
	require.False(t, ok, "Expected lock to be already held")

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Should acquire lock after TTL expired
	ok, err = lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	require.NoError(t, err, "TryLockWithTTL failed")
	require.True(t, ok, "Expected to acquire lock after TTL expired")

	// Cleanup
	_ = lock.Unlock(ctx, key)
}

func TestRedisLock_Unlock(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-unlock-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	// Clean up any existing lock
	_ = lock.Unlock(ctx, key)

	// Acquire lock
	ok, err := lock.TryLock(ctx, key)
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock")

	// Unlock
	err = lock.Unlock(ctx, key)
	require.NoError(t, err, "Unlock failed")

	// Should be able to acquire again
	ok, err = lock.TryLock(ctx, key)
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock after unlock")

	// Cleanup
	_ = lock.Unlock(ctx, key)
}

func TestRedisLock_UnlockWithValue(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-unlock-value-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	// Clean up any existing lock
	_ = lock.Unlock(ctx, key)

	// Acquire lock with value
	ok, value, err := lock.LockWithTTL(ctx, key, 30*time.Second)
	require.NoError(t, err, "LockWithTTL failed")
	require.True(t, ok, "Expected to acquire lock")
	require.NotEmpty(t, value, "Expected non-empty lock value")

	// Unlock with value
	err = lock.UnlockWithValue(ctx, key, value)
	require.NoError(t, err, "UnlockWithValue failed")

	// Should be able to acquire again
	ok, err = lock.TryLock(ctx, key)
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock after unlock")

	// Cleanup
	_ = lock.Unlock(ctx, key)
}

func TestRedisLock_ExtendLock(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-extend-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	// Clean up any existing lock
	_ = lock.Unlock(ctx, key)

	// Acquire lock with short TTL
	ok, err := lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	require.NoError(t, err, "TryLockWithTTL failed")
	require.True(t, ok, "Expected to acquire lock")

	// Extend the lock
	err = lock.ExtendLock(ctx, key, 1*time.Second)
	require.NoError(t, err, "ExtendLock failed")

	// Wait for original TTL
	time.Sleep(150 * time.Millisecond)

	// Lock should still be held (extended)
	ok, err = lock.TryLock(ctx, key)
	require.NoError(t, err, "TryLock failed")
	require.False(t, ok, "Expected lock to still be held after extend")

	// Cleanup
	_ = lock.Unlock(ctx, key)
}

func TestRedisLock_Ping(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	ctx := context.Background()
	err := lock.Ping(ctx)
	require.NoError(t, err, "Ping failed")
}

func TestRedisLock_Concurrent(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-concurrent-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	// Clean up any existing lock
	_ = lock.Unlock(ctx, key)

	acquired := make(chan bool, 2)

	go func() {
		ok, _ := lock.TryLock(ctx, key)
		acquired <- ok
		if ok {
			time.Sleep(50 * time.Millisecond)
			_ = lock.Unlock(ctx, key)
		}
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ok, _ := lock.TryLock(ctx, key)
		acquired <- ok
	}()

	results := make([]bool, 2)
	for i := 0; i < 2; i++ {
		results[i] = <-acquired
	}

	require.True(t, results[0] || results[1], "At least one goroutine should have acquired the lock")
}

func TestRedisLock_LockWithContext(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-lock-ctx-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	// Clean up any existing lock
	_ = lock.Unlock(ctx, key)

	// Acquire lock in first goroutine
	ok, err := lock.TryLock(ctx, key)
	require.NoError(t, err, "TryLock failed")
	require.True(t, ok, "Expected to acquire lock")

	// Try to acquire with short timeout context
	ctx2, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err = lock.Lock(ctx2, key)
	require.Error(t, err, "expected error due to context timeout")
	require.Equal(t, context.DeadlineExceeded, err, "expected DeadlineExceeded")

	// Cleanup
	_ = lock.Unlock(ctx, key)
}

func TestRedisLock_Close(t *testing.T) {
	lock := NewRedisLock("localhost:6379", "", 15)

	err := lock.Close()
	require.NoError(t, err, "Close failed")
}

func TestRedisLock_Client(t *testing.T) {
	lock := NewRedisLock("localhost:6379", "", 15)
	defer func() { _ = lock.Close() }()

	client := lock.Client()
	require.NotNil(t, client, "expected non-nil client")
}

func TestNewRedisLock(t *testing.T) {
	lock := NewRedisLock("localhost:6379", "password", 0)
	require.NotNil(t, lock, "expected non-nil lock")
	require.NotNil(t, lock.client, "expected non-nil client")
	require.NotNil(t, lock.lockValues, "expected non-nil lockValues map")
}
