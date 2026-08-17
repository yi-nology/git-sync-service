package dao

import (
	"testing"
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupSyncRunTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.SyncRun{}, &model.SyncRunStep{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return db
}

func TestSyncRunDAO_CreateAndFindByTaskKey(t *testing.T) {
	db := setupSyncRunTestDB(t)
	d := NewSyncRunDAO(db)

	// Create a run
	run := &model.SyncRun{
		TaskKey:       "test-task",
		TriggerSource: "manual",
		Status:        "running",
		StartTime:     time.Now(),
	}

	if err := d.Create(run); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if run.ID == 0 {
		t.Error("expected run ID to be set after create")
	}

	// Find by task key
	got, total, err := d.FindByTaskKey("test-task", DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find by task key failed: %v", err)
	}

	if total != 1 {
		t.Fatalf("expected 1 run, got %d", total)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 run, got %d", len(got))
	}

	if got[0].TaskKey != "test-task" {
		t.Errorf("expected task key 'test-task', got '%s'", got[0].TaskKey)
	}

	if got[0].TriggerSource != "manual" {
		t.Errorf("expected trigger source 'manual', got '%s'", got[0].TriggerSource)
	}

	if got[0].Status != "running" {
		t.Errorf("expected status 'running', got '%s'", got[0].Status)
	}
}

func TestSyncRunDAO_FindRecent(t *testing.T) {
	db := setupSyncRunTestDB(t)
	d := NewSyncRunDAO(db)

	// Create multiple runs
	runs := []*model.SyncRun{
		{TaskKey: "task1", TriggerSource: "manual", Status: "success"},
		{TaskKey: "task2", TriggerSource: "cron", Status: "failed"},
		{TaskKey: "task3", TriggerSource: "webhook", Status: "running"},
	}

	for _, run := range runs {
		if err := d.Create(run); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find recent
	got, err := d.FindRecent(DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find recent failed: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 runs, got %d", len(got))
	}
}

func TestSyncRunDAO_Update(t *testing.T) {
	db := setupSyncRunTestDB(t)
	d := NewSyncRunDAO(db)

	// Create a run
	run := &model.SyncRun{
		TaskKey:       "test-task",
		TriggerSource: "manual",
		Status:        "running",
		StartTime:     time.Now(),
	}

	if err := d.Create(run); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update the run
	now := time.Now()
	run.Status = "success"
	run.EndTime = &now
	run.Details = "Sync completed"
	run.DurationMs = 1000

	if err := d.Update(run); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Find by task key to verify update
	got, total, err := d.FindByTaskKey("test-task", DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find by task key failed: %v", err)
	}

	if total != 1 {
		t.Fatalf("expected 1 run, got %d", total)
	}

	if got[0].Status != "success" {
		t.Errorf("expected status 'success', got '%s'", got[0].Status)
	}

	if got[0].Details != "Sync completed" {
		t.Errorf("expected details 'Sync completed', got '%s'", got[0].Details)
	}

	if got[0].DurationMs != 1000 {
		t.Errorf("expected duration ms 1000, got %d", got[0].DurationMs)
	}
}

func TestSyncRunDAO_Delete(t *testing.T) {
	db := setupSyncRunTestDB(t)
	d := NewSyncRunDAO(db)

	// Create a run
	run := &model.SyncRun{
		TaskKey:       "test-task",
		TriggerSource: "manual",
		Status:        "running",
	}

	if err := d.Create(run); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete the run
	if err := d.Delete(run.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Find by task key should return empty
	got, total, err := d.FindByTaskKey("test-task", DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find by task key failed: %v", err)
	}

	if total != 0 {
		t.Fatalf("expected 0 runs, got %d", total)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 runs, got %d", len(got))
	}
}

func TestSyncRunDAO_CountByTaskKey(t *testing.T) {
	db := setupSyncRunTestDB(t)
	d := NewSyncRunDAO(db)

	// Create runs with different task keys
	runs := []*model.SyncRun{
		{TaskKey: "task1", Status: "success"},
		{TaskKey: "task1", Status: "failed"},
		{TaskKey: "task2", Status: "running"},
	}

	for _, run := range runs {
		if err := d.Create(run); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Count by task key
	count, err := d.CountByTaskKey("task1")
	if err != nil {
		t.Fatalf("count by task key failed: %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestSyncRunDAO_CleanupOlderThan(t *testing.T) {
	db := setupSyncRunTestDB(t)
	d := NewSyncRunDAO(db)

	// Create runs with different ages
	oldRun := &model.SyncRun{
		TaskKey:   "task1",
		Status:    "success",
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}

	newRun := &model.SyncRun{
		TaskKey:   "task2",
		Status:    "success",
		CreatedAt: time.Now(),
	}

	if err := db.Create(oldRun).Error; err != nil {
		t.Fatalf("create old run failed: %v", err)
	}

	if err := db.Create(newRun).Error; err != nil {
		t.Fatalf("create new run failed: %v", err)
	}

	// Cleanup older than 24 hours
	count, err := d.CleanupOlderThan(24 * time.Hour)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Verify only new run remains
	got, err := d.FindRecent(DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find recent failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 run, got %d", len(got))
	}

	if got[0].TaskKey != "task2" {
		t.Errorf("expected task key 'task2', got '%s'", got[0].TaskKey)
	}
}

func TestSyncRunDAO_Fields(t *testing.T) {
	now := time.Now()
	run := model.SyncRun{
		ID:             1,
		TaskKey:        "test-task",
		TriggerSource:  "manual",
		Status:         "success",
		StartTime:      now,
		EndTime:        &now,
		CommitRange:    "abc123..def456",
		Details:        "Sync completed",
		ErrorMessage:   "",
		ErrorType:      "",
		DurationMs:     1000,
		RetryTotal:     0,
		WebhookEventID: nil,
		CreatedAt:      now,
	}

	if run.ID != 1 {
		t.Errorf("expected ID 1, got %d", run.ID)
	}

	if run.TaskKey != "test-task" {
		t.Errorf("expected task key 'test-task', got '%s'", run.TaskKey)
	}

	if run.TriggerSource != "manual" {
		t.Errorf("expected trigger source 'manual', got '%s'", run.TriggerSource)
	}

	if run.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", run.Status)
	}

	if run.CommitRange != "abc123..def456" {
		t.Errorf("expected commit range 'abc123..def456', got '%s'", run.CommitRange)
	}

	if run.Details != "Sync completed" {
		t.Errorf("expected details 'Sync completed', got '%s'", run.Details)
	}

	if run.ErrorMessage != "" {
		t.Errorf("expected empty error message, got '%s'", run.ErrorMessage)
	}

	if run.ErrorType != "" {
		t.Errorf("expected empty error type, got '%s'", run.ErrorType)
	}

	if run.DurationMs != 1000 {
		t.Errorf("expected duration ms 1000, got %d", run.DurationMs)
	}

	if run.RetryTotal != 0 {
		t.Errorf("expected retry total 0, got %d", run.RetryTotal)
	}

	if run.WebhookEventID != nil {
		t.Error("expected webhook event ID to be nil")
	}
}
