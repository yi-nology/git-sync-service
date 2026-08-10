package lock

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRedisLock_TryLock(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-lock-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	// Clean up any existing lock
	lock.Unlock(ctx, key)

	// Should acquire lock
	ok, err := lock.TryLock(ctx, key)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	// Should fail to acquire same lock
	ok, err = lock.TryLock(ctx, key)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if ok {
		t.Fatal("Expected lock to be already held")
	}

	// Cleanup
	lock.Unlock(ctx, key)
}

func TestRedisLock_TryLockWithTTL(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-lock-ttl-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	// Clean up any existing lock
	lock.Unlock(ctx, key)

	// Acquire lock with short TTL
	ok, err := lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	// Should fail to acquire same lock
	ok, err = lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if ok {
		t.Fatal("Expected lock to be already held")
	}

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Should acquire lock after TTL expired
	ok, err = lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock after TTL expired")
	}

	// Cleanup
	lock.Unlock(ctx, key)
}

func TestRedisLock_Unlock(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-unlock-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	// Clean up any existing lock
	lock.Unlock(ctx, key)

	// Acquire lock
	ok, err := lock.TryLock(ctx, key)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	// Unlock
	err = lock.Unlock(ctx, key)
	if err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}

	// Should be able to acquire again
	ok, err = lock.TryLock(ctx, key)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock after unlock")
	}

	// Cleanup
	lock.Unlock(ctx, key)
}

func TestRedisLock_UnlockWithValue(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-unlock-value-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	// Clean up any existing lock
	lock.Unlock(ctx, key)

	// Acquire lock with value
	ok, value, err := lock.LockWithTTL(ctx, key, 30*time.Second)
	if err != nil {
		t.Fatalf("LockWithTTL failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}
	if value == "" {
		t.Fatal("Expected non-empty lock value")
	}

	// Unlock with value
	err = lock.UnlockWithValue(ctx, key, value)
	if err != nil {
		t.Fatalf("UnlockWithValue failed: %v", err)
	}

	// Should be able to acquire again
	ok, err = lock.TryLock(ctx, key)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock after unlock")
	}

	// Cleanup
	lock.Unlock(ctx, key)
}

func TestRedisLock_ExtendLock(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-extend-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	// Clean up any existing lock
	lock.Unlock(ctx, key)

	// Acquire lock with short TTL
	ok, err := lock.TryLockWithTTL(ctx, key, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	// Extend the lock
	err = lock.ExtendLock(ctx, key, 1*time.Second)
	if err != nil {
		t.Fatalf("ExtendLock failed: %v", err)
	}

	// Wait for original TTL
	time.Sleep(150 * time.Millisecond)

	// Lock should still be held (extended)
	ok, err = lock.TryLock(ctx, key)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if ok {
		t.Fatal("Expected lock to still be held after extend")
	}

	// Cleanup
	lock.Unlock(ctx, key)
}

func TestRedisLock_Ping(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	ctx := context.Background()
	err := lock.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

func TestRedisLock_Concurrent(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-concurrent-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	// Clean up any existing lock
	lock.Unlock(ctx, key)

	acquired := make(chan bool, 2)

	go func() {
		ok, _ := lock.TryLock(ctx, key)
		acquired <- ok
		if ok {
			time.Sleep(50 * time.Millisecond)
			lock.Unlock(ctx, key)
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

	if !results[0] && !results[1] {
		t.Fatal("At least one goroutine should have acquired the lock")
	}
}

func TestRedisLock_LockWithContext(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-lock-ctx-%d", time.Now().UnixNano())
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	// Clean up any existing lock
	lock.Unlock(ctx, key)

	// Acquire lock in first goroutine
	ok, err := lock.TryLock(ctx, key)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	// Try to acquire with short timeout context
	ctx2, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err = lock.Lock(ctx2, key)
	if err == nil {
		t.Fatal("expected error due to context timeout")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}

	// Cleanup
	lock.Unlock(ctx, key)
}

func TestRedisLock_Close(t *testing.T) {
	lock := NewRedisLock("localhost:6379", "", 15)

	err := lock.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestRedisLock_Client(t *testing.T) {
	lock := NewRedisLock("localhost:6379", "", 15)
	defer lock.Close()

	client := lock.Client()
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewRedisLock(t *testing.T) {
	lock := NewRedisLock("localhost:6379", "password", 0)
	if lock == nil {
		t.Fatal("expected non-nil lock")
	}
	if lock.client == nil {
		t.Fatal("expected non-nil client")
	}
	if lock.lockValues == nil {
		t.Fatal("expected non-nil lockValues map")
	}
}
