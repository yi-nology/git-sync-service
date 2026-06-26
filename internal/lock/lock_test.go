package lock

import (
	"context"
	"testing"
	"time"
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
