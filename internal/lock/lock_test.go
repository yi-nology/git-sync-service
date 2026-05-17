package lock

import (
	"context"
	"testing"
	"time"
)

func TestLocalLockBasic(t *testing.T) {
	l := &LocalLock{}
	ctx := context.Background()

	ok, err := l.TryLock(ctx, "test-key")
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if !ok {
		t.Fatal("TryLock should succeed on first attempt")
	}

	ok, err = l.TryLock(ctx, "test-key")
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	if ok {
		t.Fatal("TryLock should fail when already locked")
	}

	if err := l.Unlock(ctx, "test-key"); err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}

	ok, err = l.TryLock(ctx, "test-key")
	if err != nil {
		t.Fatalf("TryLock after unlock failed: %v", err)
	}
	if !ok {
		t.Fatal("TryLock should succeed after unlock")
	}
	l.Unlock(ctx, "test-key")
}

func TestLocalLockTTLExpiry(t *testing.T) {
	l := &LocalLock{}
	ctx := context.Background()

	ok, err := l.TryLockWithTTL(ctx, "ttl-key", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLockWithTTL failed: %v", err)
	}
	if !ok {
		t.Fatal("TryLockWithTTL should succeed")
	}

	time.Sleep(150 * time.Millisecond)

	ok, err = l.TryLock(ctx, "ttl-key")
	if err != nil {
		t.Fatalf("TryLock after TTL failed: %v", err)
	}
	if !ok {
		t.Fatal("TryLock should succeed after TTL expired")
	}
	l.Unlock(ctx, "ttl-key")
}

func TestLocalLockDifferentKeys(t *testing.T) {
	l := &LocalLock{}
	ctx := context.Background()

	ok1, _ := l.TryLock(ctx, "key-a")
	ok2, _ := l.TryLock(ctx, "key-b")

	if !ok1 || !ok2 {
		t.Fatal("Different keys should be independently lockable")
	}

	l.Unlock(ctx, "key-a")
	l.Unlock(ctx, "key-b")
}

func TestLocalLockBlocking(t *testing.T) {
	l := &LocalLock{}
	ctx := context.Background()

	l.TryLock(ctx, "block-key")

	done := make(chan struct{})
	go func() {
		l.Unlock(ctx, "block-key")
		close(done)
	}()

	err := l.Lock(ctx, "block-key")
	if err != nil {
		t.Fatalf("Lock should succeed after unlock: %v", err)
	}
	<-done
	l.Unlock(ctx, "block-key")
}

func TestLocalLockCanceledContext(t *testing.T) {
	l := &LocalLock{}
	ctx, cancel := context.WithCancel(context.Background())

	l.TryLock(ctx, "cancel-key")

	cancel()

	err := l.Lock(ctx, "cancel-key")
	if err == nil {
		t.Fatal("Lock should fail with canceled context")
	}
	l.Unlock(ctx, "cancel-key")
}
