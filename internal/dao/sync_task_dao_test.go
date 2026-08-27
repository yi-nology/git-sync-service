package dao

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func setupSyncTaskTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open test db")

	err = db.AutoMigrate(&model.SyncTask{}, &model.SyncRun{}, &model.SyncRunStep{})
	require.NoError(t, err, "failed to migrate test db")

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

	err := d.Create(task)
	require.NoError(t, err, "create failed")

	assert.NotZero(t, task.ID, "expected task ID to be set after create")

	// Find by key
	got, err := d.FindByKey("test-task")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "test-task", got.Key, "expected key 'test-task'")
	assert.Equal(t, "Test Task", got.Name, "expected name 'Test Task'")
	assert.Equal(t, "source-repo", got.SourceRepoKey, "expected source repo key 'source-repo'")
	assert.Equal(t, "target-repo", got.TargetRepoKey, "expected target repo key 'target-repo'")
	assert.Equal(t, "mirror", got.SyncMode, "expected sync mode 'mirror'")
	assert.Equal(t, "0 * * * *", got.Cron, "expected cron '0 * * * *'")
	assert.True(t, got.Enabled, "expected enabled to be true")
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
		err := d.Create(task)
		require.NoError(t, err, "create failed")
	}

	// Find all
	got, total, err := d.FindAll(DefaultPagination(0, 50))
	require.NoError(t, err, "find all failed")

	require.Equal(t, int64(3), total, "expected 3 tasks")
	require.Len(t, got, 3, "expected 3 tasks")
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

	err := d.Create(task)
	require.NoError(t, err, "create failed")

	// Update the task
	task.Name = "Updated Task"
	task.Cron = "*/5 * * * *"

	err = d.Update(task)
	require.NoError(t, err, "update failed")

	// Get the updated task
	got, err := d.FindByKey("test-task")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "Updated Task", got.Name, "expected name 'Updated Task'")
	assert.Equal(t, "*/5 * * * *", got.Cron, "expected cron '*/5 * * * *'")
}

func TestSyncTaskDAO_Delete(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Create a task
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	err := d.Create(task)
	require.NoError(t, err, "create failed")

	// Delete the task
	err = d.Delete("test-task")
	require.NoError(t, err, "delete failed")

	// Try to get the deleted task - should return nil (soft delete)
	got, err := d.FindByKey("test-task")
	require.NoError(t, err, "unexpected error")

	assert.Nil(t, got, "expected nil when getting deleted task")
}

func TestSyncTaskDAO_FindByKey_NotFound(t *testing.T) {
	db := setupSyncTaskTestDB(t)
	d := NewSyncTaskDAO(db)

	// Try to get a non-existent task
	got, err := d.FindByKey("nonexistent")
	require.NoError(t, err, "unexpected error")

	assert.Nil(t, got, "expected nil when getting non-existent task")
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
		err := d.Create(task)
		require.NoError(t, err, "create failed")
	}

	// Count by status
	counts, err := d.CountByStatus()
	require.NoError(t, err, "count by status failed")

	assert.Equal(t, int64(4), counts["total"], "expected total 4")
	assert.Equal(t, int64(2), counts["success"], "expected success 2")
	assert.Equal(t, int64(1), counts["failed"], "expected failed 1")
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
	err := db.Create(task1).Error
	require.NoError(t, err, "create task1 failed")
	db.Model(task1).Update("enabled", true)

	err = db.Create(task2).Error
	require.NoError(t, err, "create task2 failed")
	db.Model(task2).Update("enabled", true)

	err = db.Create(task3).Error
	require.NoError(t, err, "create task3 failed")
	db.Model(task3).Update("enabled", false)

	err = db.Create(task4).Error
	require.NoError(t, err, "create task4 failed")
	db.Model(task4).Update("enabled", true)

	// Find all enabled tasks
	got, err := d.FindAllEnabled()
	require.NoError(t, err, "find all enabled failed")

	// Should return task1 and task2 (enabled with cron)
	// Should NOT return task3 (disabled) or task4 (no cron)
	require.Len(t, got, 2, "expected 2 enabled tasks")
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
		err := d.Create(task)
		require.NoError(t, err, "create failed")
	}

	// Find by repo key
	got, total, err := d.FindByRepoKey("repo1", DefaultPagination(0, 50))
	require.NoError(t, err, "find by repo key failed")

	require.Equal(t, int64(2), total, "expected 2 tasks")
	require.Len(t, got, 2, "expected 2 tasks")
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

	assert.Equal(t, uint(1), task.ID, "expected ID 1")
	assert.Equal(t, "test-task", task.Key, "expected key 'test-task'")
	assert.Equal(t, "Test Task", task.Name, "expected name 'Test Task'")
	assert.Equal(t, "source-repo", task.SourceRepoKey, "expected source repo key 'source-repo'")
	assert.Equal(t, "target-repo", task.TargetRepoKey, "expected target repo key 'target-repo'")
	assert.Equal(t, "mirror", task.SyncMode, "expected sync mode 'mirror'")
	assert.Equal(t, "0 * * * *", task.Cron, "expected cron '0 * * * *'")
	assert.Equal(t, "token123", task.WebhookToken, "expected webhook token 'token123'")
	assert.True(t, task.Enabled, "expected enabled to be true")
	assert.True(t, task.GitTags, "expected git tags to be true")
	assert.False(t, task.GitForce, "expected git force to be false")
	assert.True(t, task.GitPrune, "expected git prune to be true")
	assert.Equal(t, "success", task.LastStatus, "expected last status 'success'")
}
