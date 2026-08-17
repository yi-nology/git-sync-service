package service

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func setupTaskServiceTestDB(t *testing.T) (*gorm.DB, *TaskService) {
	t.Helper()
	// Set up encryption key for tests
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.SyncTask{}, &model.SyncRun{}, &model.SyncRunStep{}, &model.Repo{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	taskDAO := dao.NewSyncTaskDAO(db)
	runDAO := dao.NewSyncRunDAO(db)
	runStepDAO := dao.NewSyncRunStepDAO(db)
	repoDAO, err := dao.NewRepoDAO(db)
	if err != nil {
		t.Fatalf("failed to create RepoDAO: %v", err)
	}

	svc := NewTaskService(taskDAO, runDAO, runStepDAO, repoDAO)

	return db, svc
}

func TestTaskService_ListTasks(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)
	ctx := context.Background()

	// Create tasks directly in database
	tasks := []*model.SyncTask{
		{Key: "task1", Name: "Task 1", WebhookToken: "token1"},
		{Key: "task2", Name: "Task 2", WebhookToken: "token2"},
		{Key: "task3", Name: "Task 3", WebhookToken: "token3"},
	}

	for _, task := range tasks {
		if err := db.Create(task).Error; err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List all
	got, total, err := svc.ListTasks(ctx, "", 0, 50)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if total != 3 {
		t.Fatalf("expected 3 tasks, got %d", total)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(got))
	}
}

func TestTaskService_GetTask(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)
	ctx := context.Background()

	// Create task directly in database
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Get by key
	got, err := svc.GetTask(ctx, "test-task")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.Key != "test-task" {
		t.Errorf("expected key 'test-task', got '%s'", got.Key)
	}

	if got.Name != "Test Task" {
		t.Errorf("expected name 'Test Task', got '%s'", got.Name)
	}
}

func TestTaskService_GetTask_NotFound(t *testing.T) {
	_, svc := setupTaskServiceTestDB(t)
	ctx := context.Background()

	// Try to get a non-existent task
	got, err := svc.GetTask(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent task")
	}
}

func TestTaskService_CountTasksByStatus(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create tasks with different statuses
	tasks := []*model.SyncTask{
		{Key: "task1", Name: "Task 1", Enabled: true, LastStatus: "success", WebhookToken: "token1"},
		{Key: "task2", Name: "Task 2", Enabled: true, LastStatus: "success", WebhookToken: "token2"},
		{Key: "task3", Name: "Task 3", Enabled: true, LastStatus: "failed", WebhookToken: "token3"},
		{Key: "task4", Name: "Task 4", Enabled: false, WebhookToken: "token4"},
	}

	for _, task := range tasks {
		if err := db.Create(task).Error; err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Count by status
	counts, err := svc.CountTasksByStatus()
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

func TestTaskService_CreateRun(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create task directly in database
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Create run
	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("create run failed: %v", err)
	}

	if run == nil {
		t.Fatal("expected non-nil run")
	}

	if run.TaskKey != "test-task" {
		t.Errorf("expected task key 'test-task', got '%s'", run.TaskKey)
	}

	if run.TriggerSource != "manual" {
		t.Errorf("expected trigger source 'manual', got '%s'", run.TriggerSource)
	}

	if run.Status != "running" {
		t.Errorf("expected status 'running', got '%s'", run.Status)
	}
}

func TestTaskService_CompleteRun(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create task directly in database
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Create run
	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("create run failed: %v", err)
	}

	// Complete run
	run.Status = "success"
	run.Details = "Sync completed"

	if err := svc.CompleteRun(run); err != nil {
		t.Fatalf("complete run failed: %v", err)
	}

	// Verify run was completed
	var got model.SyncRun
	if err := db.First(&got, run.ID).Error; err != nil {
		t.Fatalf("get run failed: %v", err)
	}

	if got.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", got.Status)
	}

	if got.Details != "Sync completed" {
		t.Errorf("expected details 'Sync completed', got '%s'", got.Details)
	}
}

func TestTaskService_CreateRunStep(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create task and run directly in database
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("create run failed: %v", err)
	}

	// Create run step
	step := &model.SyncRunStep{
		RunID:    run.ID,
		StepName: "fetch",
		Status:   "running",
	}

	if err := svc.CreateRunStep(step); err != nil {
		t.Fatalf("create run step failed: %v", err)
	}

	if step.ID == 0 {
		t.Error("expected step ID to be set after create")
	}
}

func TestTaskService_UpdateRunStep(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create task and run directly in database
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("create run failed: %v", err)
	}

	// Create run step
	step := &model.SyncRunStep{
		RunID:    run.ID,
		StepName: "fetch",
		Status:   "running",
	}

	if err := svc.CreateRunStep(step); err != nil {
		t.Fatalf("create run step failed: %v", err)
	}

	// Update run step
	step.Status = "success"
	step.DurationMs = 500

	if err := svc.UpdateRunStep(step); err != nil {
		t.Fatalf("update run step failed: %v", err)
	}

	// Verify step was updated
	var got model.SyncRunStep
	if err := db.First(&got, step.ID).Error; err != nil {
		t.Fatalf("get run step failed: %v", err)
	}

	if got.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", got.Status)
	}

	if got.DurationMs != 500 {
		t.Errorf("expected duration ms 500, got %d", got.DurationMs)
	}
}

func TestTaskService_UpdateTaskLastRun(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create task directly in database
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Create run
	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("create run failed: %v", err)
	}

	// Complete run with end time
	now := time.Now()
	run.Status = "success"
	run.EndTime = &now

	if err := svc.CompleteRun(run); err != nil {
		t.Fatalf("complete run failed: %v", err)
	}

	// Update task last run
	if err := svc.UpdateTaskLastRun(task, run); err != nil {
		t.Fatalf("update task last run failed: %v", err)
	}

	// Verify task was updated
	var got model.SyncTask
	if err := db.First(&got, task.ID).Error; err != nil {
		t.Fatalf("get task failed: %v", err)
	}

	if got.LastStatus != "success" {
		t.Errorf("expected last status 'success', got '%s'", got.LastStatus)
	}

	if got.LastRunAt == nil {
		t.Error("expected last run at to be set")
	}
}

func TestTaskService_FindAllEnabledTasks(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create tasks with different enabled status
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
	got, err := svc.FindAllEnabledTasks()
	if err != nil {
		t.Fatalf("find all enabled failed: %v", err)
	}

	// Should return task1 and task2 (enabled with cron)
	// Should NOT return task3 (disabled) or task4 (no cron)
	if len(got) != 2 {
		t.Fatalf("expected 2 enabled tasks, got %d", len(got))
	}
}

func TestTaskService_FindTaskByKey(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create task directly in database
	task := &model.SyncTask{
		Key:  "test-task",
		Name: "Test Task",
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Find by key
	got, err := svc.FindTaskByKey("test-task")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Key != "test-task" {
		t.Errorf("expected key 'test-task', got '%s'", got.Key)
	}

	if got.Name != "Test Task" {
		t.Errorf("expected name 'Test Task', got '%s'", got.Name)
	}
}

func TestTaskService_FindTaskByKey_NotFound(t *testing.T) {
	_, svc := setupTaskServiceTestDB(t)

	// Try to get a non-existent task
	got, err := svc.FindTaskByKey("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent task")
	}
}

func TestTaskService_CleanupOldRuns(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

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
	count, err := svc.CleanupOldRuns(24 * time.Hour)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}

func TestTaskService_CleanupOldRunSteps(t *testing.T) {
	db, svc := setupTaskServiceTestDB(t)

	// Create a run
	run := &model.SyncRun{
		TaskKey:   "task1",
		Status:    "success",
		CreatedAt: time.Now(),
	}

	if err := db.Create(run).Error; err != nil {
		t.Fatalf("create run failed: %v", err)
	}

	// Create run steps with different ages
	oldStep := &model.SyncRunStep{
		RunID:     run.ID,
		StepName:  "fetch",
		Status:    "success",
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}

	newStep := &model.SyncRunStep{
		RunID:     run.ID,
		StepName:  "push",
		Status:    "success",
		CreatedAt: time.Now(),
	}

	if err := db.Create(oldStep).Error; err != nil {
		t.Fatalf("create old step failed: %v", err)
	}

	if err := db.Create(newStep).Error; err != nil {
		t.Fatalf("create new step failed: %v", err)
	}

	// Cleanup older than 24 hours
	count, err := svc.CleanupOldRunSteps(24 * time.Hour)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}
