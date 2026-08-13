package lock

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// newMiniredis 启动一个进程内 redis,测试始终可运行(不依赖外部 redis)。
func newMiniredis(t *testing.T) *redis.Client {
	t.Helper()
	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

// TestRedisLock_ExtendLockWithValue_miniredis 验证 CAS 续期:
// 持有者能续期,非持有者/已释放后续期失败。
func TestRedisLock_ExtendLockWithValue_miniredis(t *testing.T) {
	client := newMiniredis(t)
	defer client.Close()
	rl := &RedisLock{client: client}
	ctx := context.Background()

	ok, value, err := rl.LockWithTTL(ctx, "k1", 2*time.Second)
	if err != nil || !ok {
		t.Fatalf("LockWithTTL failed: ok=%v err=%v", ok, err)
	}

	if got, _ := rl.ExtendLockWithValue(ctx, "k1", value, 5*time.Second); !got {
		t.Fatal("owner extend should succeed")
	}
	if got, _ := rl.ExtendLockWithValue(ctx, "k1", "wrong-value", 5*time.Second); got {
		t.Fatal("non-owner extend should fail")
	}

	if err := rl.UnlockWithValue(ctx, "k1", value); err != nil {
		t.Fatalf("UnlockWithValue: %v", err)
	}
	if got, _ := rl.ExtendLockWithValue(ctx, "k1", value, 5*time.Second); got {
		t.Fatal("extend after unlock should fail")
	}
}

// TestSemaphore_RespectsMax_miniredis 验证不超过 max,release 后释放槽。
func TestSemaphore_RespectsMax_miniredis(t *testing.T) {
	client := newMiniredis(t)
	defer client.Close()
	ctx := context.Background()
	sem := NewSemaphore(client, "sem-max", 2, 5*time.Second)

	if ok, _ := sem.Acquire(ctx, "a"); !ok {
		t.Fatal("a should acquire")
	}
	if ok, _ := sem.Acquire(ctx, "b"); !ok {
		t.Fatal("b should acquire")
	}
	if ok, _ := sem.Acquire(ctx, "c"); ok {
		t.Fatal("c should fail (max=2)")
	}
	if err := sem.Release(ctx, "a"); err != nil {
		t.Fatalf("release a: %v", err)
	}
	if ok, _ := sem.Acquire(ctx, "c"); !ok {
		t.Fatal("c should acquire after a released")
	}
}

// TestSemaphore_RenewKeepsSlot_miniredis 验证 Renew 让槽活过原始 TTL;
// 不续期则槽到期被清理。
func TestSemaphore_RenewKeepsSlot_miniredis(t *testing.T) {
	client := newMiniredis(t)
	defer client.Close()
	ctx := context.Background()
	sem := NewSemaphore(client, "sem-renew", 1, 300*time.Millisecond)

	if ok, _ := sem.Acquire(ctx, "a"); !ok {
		t.Fatal("a should acquire")
	}
	if ok, _ := sem.Acquire(ctx, "b"); ok {
		t.Fatal("b should fail while a holds the only slot")
	}

	// 续期 a 几次,撑过原始 TTL
	for i := 0; i < 3; i++ {
		time.Sleep(150 * time.Millisecond)
		if err := sem.Renew(ctx, "a"); err != nil {
			t.Fatalf("renew a: %v", err)
		}
	}
	// 此时已过 450ms > 300ms 原始 TTL,但续期后 a 仍存活
	if ok, _ := sem.Acquire(ctx, "b"); ok {
		t.Fatal("b should still fail while a is renewed")
	}

	// 不再续期,等过期后 b 应能拿到
	time.Sleep(500 * time.Millisecond)
	if ok, _ := sem.Acquire(ctx, "b"); !ok {
		t.Fatal("b should acquire after a's slot expired")
	}
}
