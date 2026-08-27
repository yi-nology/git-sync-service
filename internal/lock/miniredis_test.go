package lock

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
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
	defer func() { _ = client.Close() }()
	rl := &RedisLock{client: client, lockValues: make(map[string]string)}
	ctx := context.Background()

	ok, value, err := rl.LockWithTTL(ctx, "k1", 2*time.Second)
	require.NoError(t, err, "LockWithTTL failed")
	require.True(t, ok, "LockWithTTL should acquire lock")

	got, _ := rl.ExtendLockWithValue(ctx, "k1", value, 5*time.Second)
	require.True(t, got, "owner extend should succeed")

	got, _ = rl.ExtendLockWithValue(ctx, "k1", "wrong-value", 5*time.Second)
	require.False(t, got, "non-owner extend should fail")

	err = rl.UnlockWithValue(ctx, "k1", value)
	require.NoError(t, err, "UnlockWithValue")

	got, _ = rl.ExtendLockWithValue(ctx, "k1", value, 5*time.Second)
	require.False(t, got, "extend after unlock should fail")
}

// TestSemaphore_RespectsMax_miniredis 验证不超过 max,release 后释放槽。
func TestSemaphore_RespectsMax_miniredis(t *testing.T) {
	client := newMiniredis(t)
	defer func() { _ = client.Close() }()
	ctx := context.Background()
	sem := NewSemaphore(client, "sem-max", 2, 5*time.Second)

	ok, _ := sem.Acquire(ctx, "a")
	require.True(t, ok, "a should acquire")

	ok, _ = sem.Acquire(ctx, "b")
	require.True(t, ok, "b should acquire")

	ok, _ = sem.Acquire(ctx, "c")
	require.False(t, ok, "c should fail (max=2)")

	err := sem.Release(ctx, "a")
	require.NoError(t, err, "release a")

	ok, _ = sem.Acquire(ctx, "c")
	require.True(t, ok, "c should acquire after a released")
}

// TestSemaphore_RenewKeepsSlot_miniredis 验证 Renew 让槽活过原始 TTL;
// 不续期则槽到期被清理。
func TestSemaphore_RenewKeepsSlot_miniredis(t *testing.T) {
	client := newMiniredis(t)
	defer func() { _ = client.Close() }()
	ctx := context.Background()
	sem := NewSemaphore(client, "sem-renew", 1, 300*time.Millisecond)

	ok, _ := sem.Acquire(ctx, "a")
	require.True(t, ok, "a should acquire")

	ok, _ = sem.Acquire(ctx, "b")
	require.False(t, ok, "b should fail while a holds the only slot")

	// 续期 a 几次,撑过原始 TTL
	for i := 0; i < 3; i++ {
		time.Sleep(150 * time.Millisecond)
		err := sem.Renew(ctx, "a")
		require.NoError(t, err, "renew a")
	}
	// 此时已过 450ms > 300ms 原始 TTL,但续期后 a 仍存活
	ok, _ = sem.Acquire(ctx, "b")
	require.False(t, ok, "b should still fail while a is renewed")

	// 不再续期,等过期后 b 应能拿到
	time.Sleep(500 * time.Millisecond)
	ok, _ = sem.Acquire(ctx, "b")
	require.True(t, ok, "b should acquire after a's slot expired")
}
