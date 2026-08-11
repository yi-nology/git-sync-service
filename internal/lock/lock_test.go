package lock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestLocalLock_TryLock(t *testing.T) {
	lock := NewLocalLock()
	ctx := context.Background()

	ok, err := lock.TryLock(ctx, "test-key")
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	ok, err = lock.TryLock(ctx, "test-key")
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if ok {
		t.Fatal("Expected lock to be already held")
	}
}

func TestLocalLock_TryLockWithTTL(t *testing.T) {
	lock := NewLocalLock()
	ctx := context.Background()

	ok, err := lock.TryLockWithTTL(ctx, "test-key", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	ok, err = lock.TryLockWithTTL(ctx, "test-key", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if ok {
		t.Fatal("Expected lock to be already held")
	}

	time.Sleep(150 * time.Millisecond)

	ok, err = lock.TryLockWithTTL(ctx, "test-key", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock after TTL expired")
	}
}

func TestLocalLock_Unlock(t *testing.T) {
	lock := NewLocalLock()
	ctx := context.Background()

	ok, err := lock.TryLock(ctx, "test-key")
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock")
	}

	err = lock.Unlock(ctx, "test-key")
	if err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}

	ok, err = lock.TryLock(ctx, "test-key")
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire lock after unlock")
	}
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

	if !results[0] && !results[1] {
		t.Fatal("At least one goroutine should have acquired the lock")
	}
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
	if err := client.Close(); err != nil {
		t.Errorf("Failed to close client: %v", err)
	}
}

func TestSemaphore_Acquire(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-%d", time.Now().UnixNano())
	sem := NewSemaphore(client, key, 2)
	defer client.Del(ctx, sem.key)

	// Should acquire first two slots
	ok1, err := sem.Acquire(ctx, "worker-1")
	if err != nil {
		t.Fatalf("Acquire worker-1 failed: %v", err)
	}
	if !ok1 {
		t.Fatal("Expected worker-1 to acquire semaphore")
	}

	ok2, err := sem.Acquire(ctx, "worker-2")
	if err != nil {
		t.Fatalf("Acquire worker-2 failed: %v", err)
	}
	if !ok2 {
		t.Fatal("Expected worker-2 to acquire semaphore")
	}

	// Third acquire should fail (max is 2)
	ok3, err := sem.Acquire(ctx, "worker-3")
	if err != nil {
		t.Fatalf("Acquire worker-3 failed: %v", err)
	}
	if ok3 {
		t.Fatal("Expected worker-3 to fail acquiring semaphore")
	}
}

func TestSemaphore_Release(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-release-%d", time.Now().UnixNano())
	sem := NewSemaphore(client, key, 1)
	defer client.Del(ctx, sem.key)

	// Acquire the single slot
	ok, err := sem.Acquire(ctx, "worker-1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire semaphore")
	}

	// Release the slot
	err = sem.Release(ctx, "worker-1")
	if err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	// Should be able to acquire again
	ok, err = sem.Acquire(ctx, "worker-2")
	if err != nil {
		t.Fatalf("Acquire after release failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire semaphore after release")
	}
}

func TestSemaphore_Concurrent(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-concurrent-%d", time.Now().UnixNano())
	maxSlots := 3
	sem := NewSemaphore(client, key, maxSlots)
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
		if errors[i] != nil {
			t.Errorf("Worker %d got error: %v", i, errors[i])
		}
		if results[i] {
			acquired++
		}
	}

	// Exactly maxSlots should have acquired
	if acquired != maxSlots {
		t.Errorf("Expected %d workers to acquire, got %d", maxSlots, acquired)
	}
}

func TestSemaphore_Cleanup(t *testing.T) {
	client := getTestRedisClient(t)
	defer closeClient(t, client)

	ctx := context.Background()
	key := fmt.Sprintf("test-semaphore-cleanup-%d", time.Now().UnixNano())
	sem := NewSemaphore(client, key, 1)
	defer client.Del(ctx, sem.key)

	// Acquire a slot
	ok, err := sem.Acquire(ctx, "old-worker")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected to acquire semaphore")
	}

	// Cleanup with long duration should remove nothing (entries are fresh)
	err = sem.Cleanup(ctx, 1*time.Hour)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Old worker should still be in the semaphore
	ok, err = sem.Acquire(ctx, "new-worker")
	if err != nil {
		t.Fatalf("Acquire new-worker failed: %v", err)
	}
	if ok {
		t.Fatal("Expected new-worker to fail (semaphore full)")
	}

	// Cleanup with very short duration should remove old entries
	err = sem.Cleanup(ctx, 0)
	if err != nil {
		t.Fatalf("Cleanup with 0 duration failed: %v", err)
	}

	// Now new-worker should succeed
	ok, err = sem.Acquire(ctx, "new-worker")
	if err != nil {
		t.Fatalf("Acquire new-worker after cleanup failed: %v", err)
	}
	if !ok {
		t.Fatal("Expected new-worker to acquire after cleanup")
	}
}
