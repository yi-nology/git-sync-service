package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/yi-nology/git-sync-service/internal/lock"
)

// TestLocalGuard_TaskMutexAndConcurrency 验证进程内 guard:同 taskKey 互斥 + 全局并发上限。
func TestLocalGuard_TaskMutexAndConcurrency(t *testing.T) {
	g := newLocalGuard(2) // max 2 并发
	ctx := context.Background()

	// 同 taskKey 第二次获取应返回 ErrTaskRunning
	rel, err := g.Acquire(ctx, "taskA")
	if err != nil {
		t.Fatalf("acquire taskA: %v", err)
	}
	if _, err := g.Acquire(ctx, "taskA"); !errors.Is(err, ErrTaskRunning) {
		t.Fatalf("second acquire taskA should be ErrTaskRunning, got %v", err)
	}

	// 占满并发(taskB),再取 taskC 应 ErrTooManyConcurrent
	relB, err := g.Acquire(ctx, "taskB")
	if err != nil {
		t.Fatalf("acquire taskB: %v", err)
	}
	if _, err := g.Acquire(ctx, "taskC"); !errors.Is(err, ErrTooManyConcurrent) {
		t.Fatalf("acquire taskC should be ErrTooManyConcurrent, got %v", err)
	}

	// 释放 taskA 后空出 1 槽,taskC 可以拿到(占住)
	rel()
	relC, err := g.Acquire(ctx, "taskC")
	if err != nil {
		t.Fatalf("taskC should acquire after taskA released: %v", err)
	}
	// 现在 taskB + taskC 占满,taskD 应被拦
	if _, err := g.Acquire(ctx, "taskD"); !errors.Is(err, ErrTooManyConcurrent) {
		t.Fatalf("taskD should be ErrTooManyConcurrent, got %v", err)
	}
	relB()
	relC()
}

// TestLocalGuard_ConcurrentStress 并发场景下不应超卖并发槽。
func TestLocalGuard_ConcurrentStress(t *testing.T) {
	g := newLocalGuard(3)
	ctx := context.Background()
	var wg sync.WaitGroup
	var acquired, running, maxRunning int32
	var mu sync.Mutex
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "task"
			// 大部分用同一 key(被互斥拦),少部分用独立 key(吃并发槽)
			if i%5 == 0 {
				key = "task-" + itoa(i)
			}
			rel, err := g.Acquire(ctx, key)
			if err != nil {
				return
			}
			defer rel()
			mu.Lock()
			acquired++
			running++
			if running > maxRunning {
				maxRunning = running
			}
			mu.Unlock()
			// 模拟一点工作
			// hold briefly
			mu.Lock()
			running--
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	_ = acquired
	if maxRunning > 3 {
		t.Fatalf("concurrency exceeded: maxRunning=%d > 3", maxRunning)
	}
}

// TestRedisGuard_CrossInstanceMutex 是多实例核心测试:
// 两个 guard(模拟两个实例)共享同一 redis,同 taskKey 必须跨实例互斥。
func TestRedisGuard_CrossInstanceMutex(t *testing.T) {
	mr := miniredis.RunT(t)
	g1, err := newRedisGuard(mr.Addr(), "", 0, 2, lock.RedisPoolOptions{})
	if err != nil {
		t.Fatalf("newRedisGuard g1: %v", err)
	}
	defer func() { _ = g1.Close() }()
	g2, err := newRedisGuard(mr.Addr(), "", 0, 2, lock.RedisPoolOptions{})
	if err != nil {
		t.Fatalf("newRedisGuard g2: %v", err)
	}
	defer func() { _ = g2.Close() }()
	ctx := context.Background()

	// g1 拿到 taskA
	rel1, err := g1.Acquire(ctx, "taskA")
	if err != nil {
		t.Fatalf("g1 acquire taskA: %v", err)
	}
	// g2(另一个实例)再拿 taskA 必须失败 —— 跨实例互斥
	if _, err := g2.Acquire(ctx, "taskA"); !errors.Is(err, ErrTaskRunning) {
		t.Fatalf("g2 acquire taskA should be ErrTaskRunning (cross-instance), got %v", err)
	}
	// g1 释放后 g2 能拿到
	rel1()
	rel2, err := g2.Acquire(ctx, "taskA")
	if err != nil {
		t.Fatalf("g2 acquire taskA after release: %v", err)
	}
	rel2()
}

// TestRedisGuard_CrossInstanceConcurrencyCap 跨实例并发上限:
// g1 占满槽,g2 拿不到(即便不同 taskKey)。
func TestRedisGuard_CrossInstanceConcurrencyCap(t *testing.T) {
	mr := miniredis.RunT(t)
	g1, err := newRedisGuard(mr.Addr(), "", 0, 1, lock.RedisPoolOptions{}) // max 1
	if err != nil {
		t.Fatalf("g1: %v", err)
	}
	defer func() { _ = g1.Close() }()
	g2, err := newRedisGuard(mr.Addr(), "", 0, 1, lock.RedisPoolOptions{})
	if err != nil {
		t.Fatalf("g2: %v", err)
	}
	defer func() { _ = g2.Close() }()
	ctx := context.Background()

	rel1, err := g1.Acquire(ctx, "taskA")
	if err != nil {
		t.Fatalf("g1 acquire: %v", err)
	}
	// 不同 taskKey,但全局槽已满 → ErrTooManyConcurrent
	if _, err := g2.Acquire(ctx, "taskB"); !errors.Is(err, ErrTooManyConcurrent) {
		t.Fatalf("g2 should be ErrTooManyConcurrent, got %v", err)
	}
	rel1()
	// 释放后 g2 能拿
	rel2, err := g2.Acquire(ctx, "taskB")
	if err != nil {
		t.Fatalf("g2 after release: %v", err)
	}
	rel2()
}

// itoa 轻量整数转字符串,避免引入 strconv 仅为此。
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}
