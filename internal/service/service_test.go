package service

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yi-nology/git-sync-service/internal/lock"
	"github.com/yi-nology/git-sync-service/sync/model"
)

func TestServiceGetTempDir(t *testing.T) {
	cfg := &model.Config{
		Git: model.GitConfig{
			TempDir: "/tmp/test",
		},
	}

	svc := &Service{
		config: cfg,
	}

	result := svc.GetTempDir("my-task")
	expected := "/tmp/test/my-task"
	assert.Equal(t, expected, result, "Expected %s, got %s", expected, result)
}

func TestServiceGetConfig(t *testing.T) {
	cfg := &model.Config{
		Server: model.ServerConfig{
			Port: 8080,
		},
	}

	svc := &Service{
		config: cfg,
	}

	result := svc.GetConfig()
	assert.Same(t, cfg, result, "Expected same config pointer")
}

func TestServiceGetAPIKey(t *testing.T) {
	cfg := &model.Config{
		Server: model.ServerConfig{
			APIKey: "test-api-key",
		},
	}

	svc := &Service{
		config: cfg,
	}

	result := svc.GetAPIKey()
	assert.Equal(t, "test-api-key", result, "Expected 'test-api-key'")
}

func TestServiceHealthCheckNoRedis(t *testing.T) {
	// Skip this test as it requires a database connection
	t.Skip("Skipping HealthCheck test - requires database")
}

func TestServiceHealthCheckWithRedis(t *testing.T) {
	// Skip this test as it requires a database connection
	t.Skip("Skipping HealthCheck test - requires database")
}

func TestConcurrencyGuardLocal(t *testing.T) {
	guard := newLocalGuard(2)
	defer func() { _ = guard.Close() }()

	ctx := context.Background()

	// First acquire should succeed
	release1, err := guard.Acquire(ctx, "task1")
	require.NoError(t, err, "First acquire failed")

	// Second acquire of same task should fail
	_, err = guard.Acquire(ctx, "task1")
	assert.Equal(t, ErrTaskRunning, err, "Expected ErrTaskRunning")

	// Acquire different task should succeed
	release2, err := guard.Acquire(ctx, "task2")
	require.NoError(t, err, "Second task acquire failed")

	// Third acquire should fail (max concurrent reached)
	_, err = guard.Acquire(ctx, "task3")
	assert.Equal(t, ErrTooManyConcurrent, err, "Expected ErrTooManyConcurrent")

	// Release first task
	release1()

	// Now we should be able to acquire task3
	release3, err := guard.Acquire(ctx, "task3")
	require.NoError(t, err, "Third task acquire failed")

	// Cleanup
	release2()
	release3()
}

func TestConcurrencyGuardLocalDefaultMax(t *testing.T) {
	// Test that default max concurrent is 5
	guard := newLocalGuard(0)
	defer func() { _ = guard.Close() }()

	ctx := context.Background()

	// Should be able to acquire 5 tasks
	releases := make([]func(), 5)
	for i := 0; i < 5; i++ {
		release, err := guard.Acquire(ctx, "task"+string(rune('a'+i)))
		require.NoError(t, err, "Failed to acquire task %d", i)
		releases[i] = release
	}

	// 6th should fail
	_, err := guard.Acquire(ctx, "task6")
	assert.Equal(t, ErrTooManyConcurrent, err, "Expected ErrTooManyConcurrent")

	// Cleanup
	for _, release := range releases {
		release()
	}
}

func TestNewGuardLocal(t *testing.T) {
	guard, err := newGuard("", "", 0, 5, lock.RedisPoolOptions{})
	require.NoError(t, err, "Failed to create local guard")
	defer func() { _ = guard.Close() }()

	// Should be a localGuard
	_, ok := guard.(*localGuard)
	assert.True(t, ok, "Expected localGuard type")
}

func TestNewGuardRedis(t *testing.T) {
	// Use miniredis for testing
	mr := miniredis.RunT(t)
	defer mr.Close()

	guard, err := newGuard(mr.Addr(), "", 0, 5, lock.RedisPoolOptions{})
	require.NoError(t, err, "Failed to create Redis guard")
	defer func() { _ = guard.Close() }()

	// Should be a redisGuard
	_, ok := guard.(*redisGuard)
	assert.True(t, ok, "Expected redisGuard type")
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	assert.Equal(t, 30*time.Second, taskLockTTL, "Expected taskLockTTL to be 30s")
	assert.Equal(t, "task:", taskLockKeyPfx, "Expected taskLockKeyPfx to be 'task:'")
	assert.Equal(t, "sync-tasks", semSlotKey, "Expected semSlotKey to be 'sync-tasks'")
	assert.Equal(t, 30*time.Second, semSlotTTL, "Expected semSlotTTL to be 30s")
	assert.Equal(t, 10*time.Second, renewInterval, "Expected renewInterval to be 10s")
}
