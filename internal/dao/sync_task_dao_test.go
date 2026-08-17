package dao

import (
	"testing"
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupSyncTaskTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.SyncTask{}, &model.SyncRun{}, &model.SyncRunStep{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return db
}

func TestSyncTaskDAO_CreateAndFindByKey(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create a task
	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		SourceBranch:  "main",
		TargetRepoKey: "target-repo",
		TargetBranch:  "main",
		SyncMode:      "mirror",
		Cron:          "0 * * * *",
		Enabled:       true,
	}

	if err := d.Create(task); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if task.ID == 0 {
		t.Error("expected task ID to be set after create")
	}

	// Find by key
	got, err := d.FindByKey("test-task")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Key != "test-task" {
		t.Errorf("expected key 'test-task', got '%s'", got.Key)
	}

	if got.Name != "Test Task" {
		t.Errorf("expected name 'Test Task', got '%s'", got.Name)
	}

	if got.SourceRepoKey != "source-repo" {
		t.Errorf("expected source repo key 'source-repo', got '%s'", got.SourceRepoKey)
	}

	if got.TargetRepoKey != "target-repo" {
		t.Errorf("expected target repo key 'target-repo', got '%s'", got.TargetRepoKey)
	}

	if got.SyncMode != "mirror" {
		t.Errorf("expected sync mode 'mirror', got '%s'", got.SyncMode)
	}

	if got.Cron != "0 * * * *" {
		t.Errorf("expected cron '0 * * * *', got '%s'", got.Cron)
	}

	if !got.Enabled {
		t.Error("expected enabled to be true")
	}
}

func TestSyncTaskDAO_FindAll(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create multiple tasks with unique webhook tokens
	tasks := []*model.SyncTask{
		{Key: "task1", Name: "Task 1", Enabled: true, WebhookToken: "token1"},
		{Key: "task2", Name: "Task 2", Enabled: true, WebhookToken: "token2"},
		{Key: "task3", Name: "Task 3", Enabled: false, WebhookToken: "token3"},
	}

	for _, task := range tasks {
		if err := d.Create(task); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find all
	got, total, err := d.FindAll(DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find all failed: %v", err)
	}

	if total != 3 {
		t.Fatalf("expected 3 tasks, got %d", total)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(got))
	}
}

func TestSyncTaskDAO_Update(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create a task
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
		Cron: "0 * * * *",
	}

	if err := d.Create(task); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update the task
	task.Name = "Updated Task"
	task.Cron = "*/5 * * * *"

	if err := d.Update(task); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Get the updated task
	got, err := d.FindByKey("test-task")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Name != "Updated Task" {
		t.Errorf("expected name 'Updated Task', got '%s'", got.Name)
	}

	if got.Cron != "*/5 * * * *" {
		t.Errorf("expected cron '*/5 * * * *', got '%s'", got.Cron)
	}
}

func TestSyncTaskDAO_Delete(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create a task
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := d.Create(task); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete the task
	if err := d.Delete("test-task"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Try to get the deleted task - should return nil (soft delete)
	got, err := d.FindByKey("test-task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting deleted task")
	}
}

func TestSyncTaskDAO_FindByKey_NotFound(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Try to get a non-existent task
	got, err := d.FindByKey("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent task")
	}
}

func TestSyncTaskDAO_CountByStatus(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create tasks with different statuses and unique webhook tokens
	tasks := []*model.SyncTask{
		{Key: "task1", Name: "Task 1", Enabled: true, LastStatus: "success", WebhookToken: "token1"},
		{Key: "task2", Name: "Task 2", Enabled: true, LastStatus: "success", WebhookToken: "token2"},
		{Key: "task3", Name: "Task 3", Enabled: true, LastStatus: "failed", WebhookToken: "token3"},
		{Key: "task4", Name: "Task 4", Enabled: false, WebhookToken: "token4"},
	}

	for _, task := range tasks {
		if err := d.Create(task); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Count by status
	counts, err := d.CountByStatus()
	if err != nil {
		t.Fatalf("count by status failed: %v", err)
	}

	if counts["total"] != 4 {
		t.Errorf("expected total 4, got %d", counts["total"])
	}

	if counts["success"] != 2 {
		t.Errorf("expected success 2, got %d", counts["success"])
	}

	if counts["failed"] != 1 {
		t.Errorf("expected failed 1, got %d", counts["failed"])
	}
}

func TestSyncTaskDAO_FindAllEnabled(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create tasks with different enabled status and unique webhook tokens
	// Note: FindAllEnabled requires cron to be set and non-empty
	// Use db.Create directly to avoid GORM default value issues with bool fields
	task1 := &model.SyncTask{Key: "task1", Name: "Task 1", Cron: "0 * * * *", WebhookToken: "token1"}
	task2 := &model.SyncTask{Key: "task2", Name: "Task 2", Cron: "*/5 * * * *", WebhookToken: "token2"}
	task3 := &model.SyncTask{Key: "task3", Name: "Task 3", Cron: "0 * * * *", WebhookToken: "token3"}
	task4 := &model.SyncTask{Key: "task4", Name: "Task 4", Cron: "", WebhookToken: "token4"}

	// Create tasks and explicitly set enabled status
	if err := db.Create(task1).Error; err != nil {
		t.Fatalf("create task1 failed: %v", err)
	}
	db.Model(task1).Update("enabled", true)

	if err := db.Create(task2).Error; err != nil {
		t.Fatalf("create task2 failed: %v", err)
	}
	db.Model(task2).Update("enabled", true)

	if err := db.Create(task3).Error; err != nil {
		t.Fatalf("create task3 failed: %v", err)
	}
	db.Model(task3).Update("enabled", false)

	if err := db.Create(task4).Error; err != nil {
		t.Fatalf("create task4 failed: %v", err)
	}
	db.Model(task4).Update("enabled", true)

	// Find all enabled tasks
	got, err := d.FindAllEnabled()
	if err != nil {
		t.Fatalf("find all enabled failed: %v", err)
	}

	// Should return task1 and task2 (enabled with cron)
	// Should NOT return task3 (disabled) or task4 (no cron)
	if len(got) != 2 {
		t.Fatalf("expected 2 enabled tasks, got %d", len(got))
	}
}

func TestSyncTaskDAO_FindByRepoKey(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create tasks with different repo keys and unique webhook tokens
	tasks := []*model.SyncTask{
		{Key: "task1", Name: "Task 1", SourceRepoKey: "repo1", WebhookToken: "token1"},
		{Key: "task2", Name: "Task 2", SourceRepoKey: "repo1", WebhookToken: "token2"},
		{Key: "task3", Name: "Task 3", SourceRepoKey: "repo2", WebhookToken: "token3"},
	}

	for _, task := range tasks {
		if err := d.Create(task); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find by repo key
	got, total, err := d.FindByRepoKey("repo1", DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find by repo key failed: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected 2 tasks, got %d", total)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(got))
	}
}

func TestSyncTaskDAO_Fields(t *testing.T) {
	now := time.Now()
	task := model.SyncTask{
		ID:            1,
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		SourceBranch:  "main",
		TargetRepoKey: "target-repo",
		TargetBranch:  "main",
		SyncMode:      "mirror",
		Cron:          "0 * * * *",
		WebhookToken:  "token123",
		Enabled:       true,
		GitTags:       true,
		GitForce:      false,
		GitPrune:      true,
		LastRunAt:     &now,
		LastStatus:    "success",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Key != "test-task" {
		t.Errorf("expected key 'test-task', got '%s'", task.Key)
	}

	if task.Name != "Test Task" {
		t.Errorf("expected name 'Test Task', got '%s'", task.Name)
	}

	if task.SourceRepoKey != "source-repo" {
		t.Errorf("expected source repo key 'source-repo', got '%s'", task.SourceRepoKey)
	}

	if task.TargetRepoKey != "target-repo" {
		t.Errorf("expected target repo key 'target-repo', got '%s'", task.TargetRepoKey)
	}

	if task.SyncMode != "mirror" {
		t.Errorf("expected sync mode 'mirror', got '%s'", task.SyncMode)
	}

	if task.Cron != "0 * * * *" {
		t.Errorf("expected cron '0 * * * *', got '%s'", task.Cron)
	}

	if task.WebhookToken != "token123" {
		t.Errorf("expected webhook token 'token123', got '%s'", task.WebhookToken)
	}

	if !task.Enabled {
		t.Error("expected enabled to be true")
	}

	if !task.GitTags {
		t.Error("expected git tags to be true")
	}

	if task.GitForce {
		t.Error("expected git force to be false")
	}

	if !task.GitPrune {
		t.Error("expected git prune to be true")
	}

	if task.LastStatus != "success" {
		t.Errorf("expected last status 'success', got '%s'", task.LastStatus)
	}
}
