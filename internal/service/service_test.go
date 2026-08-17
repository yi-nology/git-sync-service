package service

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
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
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
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
	if result != cfg {
		t.Error("Expected same config pointer")
	}
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
	if result != "test-api-key" {
		t.Errorf("Expected 'test-api-key', got '%s'", result)
	}
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
	if err != nil {
		t.Fatalf("First acquire failed: %v", err)
	}

	// Second acquire of same task should fail
	_, err = guard.Acquire(ctx, "task1")
	if err != ErrTaskRunning {
		t.Errorf("Expected ErrTaskRunning, got %v", err)
	}

	// Acquire different task should succeed
	release2, err := guard.Acquire(ctx, "task2")
	if err != nil {
		t.Fatalf("Second task acquire failed: %v", err)
	}

	// Third acquire should fail (max concurrent reached)
	_, err = guard.Acquire(ctx, "task3")
	if err != ErrTooManyConcurrent {
		t.Errorf("Expected ErrTooManyConcurrent, got %v", err)
	}

	// Release first task
	release1()

	// Now we should be able to acquire task3
	release3, err := guard.Acquire(ctx, "task3")
	if err != nil {
		t.Fatalf("Third task acquire failed: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Failed to acquire task %d: %v", i, err)
		}
		releases[i] = release
	}

	// 6th should fail
	_, err := guard.Acquire(ctx, "task6")
	if err != ErrTooManyConcurrent {
		t.Errorf("Expected ErrTooManyConcurrent, got %v", err)
	}

	// Cleanup
	for _, release := range releases {
		release()
	}
}

func TestNewGuardLocal(t *testing.T) {
	guard, err := newGuard("", "", 0, 5)
	if err != nil {
		t.Fatalf("Failed to create local guard: %v", err)
	}
	defer func() { _ = guard.Close() }()

	// Should be a localGuard
	_, ok := guard.(*localGuard)
	if !ok {
		t.Error("Expected localGuard type")
	}
}

func TestNewGuardRedis(t *testing.T) {
	// Use miniredis for testing
	mr := miniredis.RunT(t)
	defer mr.Close()

	guard, err := newGuard(mr.Addr(), "", 0, 5)
	if err != nil {
		t.Fatalf("Failed to create Redis guard: %v", err)
	}
	defer func() { _ = guard.Close() }()

	// Should be a redisGuard
	_, ok := guard.(*redisGuard)
	if !ok {
		t.Error("Expected redisGuard type")
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	if taskLockTTL != 30*time.Second {
		t.Errorf("Expected taskLockTTL to be 30s, got %v", taskLockTTL)
	}

	if taskLockKeyPfx != "task:" {
		t.Errorf("Expected taskLockKeyPfx to be 'task:', got '%s'", taskLockKeyPfx)
	}

	if semSlotKey != "sync-tasks" {
		t.Errorf("Expected semSlotKey to be 'sync-tasks', got '%s'", semSlotKey)
	}

	if semSlotTTL != 30*time.Second {
		t.Errorf("Expected semSlotTTL to be 30s, got %v", semSlotTTL)
	}

	if renewInterval != 10*time.Second {
		t.Errorf("Expected renewInterval to be 10s, got %v", renewInterval)
	}
}
