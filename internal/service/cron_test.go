package service

import (
	"context"
	"sync"
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupCronTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.SyncTask{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	// sqlite :memory: 每个连接是独立库;强制单连接,确保 cron goroutine 能看到已迁移的表
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	return db
}

func setupCronTestService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()

	db := setupCronTestDB(t)
	taskDAO := dao.NewSyncTaskDAO(db)

	taskService := NewTaskService(taskDAO, nil, nil, nil)

	return &Service{
		tasks:        taskService,
		cron:         cron.New(cron.WithSeconds()),
		cronEntryIDs: make(map[string]cron.EntryID),
		cronMu:       sync.RWMutex{},
	}, db
}

func TestAddCronJob(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "* * * * * *", // Every second
	}

	err := svc.addCronJob(task)
	if err != nil {
		t.Fatalf("addCronJob failed: %v", err)
	}

	// Verify the cron job was added
	if _, exists := svc.cronEntryIDs[task.Key]; !exists {
		t.Error("expected cron job to be added")
	}
}

func TestAddCronJob_UpdateExisting(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "* * * * * *",
	}

	// Add initial cron job
	err := svc.addCronJob(task)
	if err != nil {
		t.Fatalf("addCronJob failed: %v", err)
	}

	initialEntryID := svc.cronEntryIDs[task.Key]

	// Update the cron job
	task.Cron = "*/2 * * * * *"
	err = svc.addCronJob(task)
	if err != nil {
		t.Fatalf("addCronJob update failed: %v", err)
	}

	// Verify the entry ID changed
	newEntryID := svc.cronEntryIDs[task.Key]
	if newEntryID == initialEntryID {
		t.Error("expected entry ID to change after update")
	}
}

func TestAddCronJob_InvalidCron(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "invalid-cron-expression",
	}

	err := svc.addCronJob(task)
	if err == nil {
		t.Fatal("expected error for invalid cron expression")
	}
}

func TestRemoveCronJob(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "* * * * * *",
	}

	// Add a cron job first
	err := svc.addCronJob(task)
	if err != nil {
		t.Fatalf("addCronJob failed: %v", err)
	}

	// Verify it exists
	if _, exists := svc.cronEntryIDs[task.Key]; !exists {
		t.Fatal("expected cron job to exist before removal")
	}

	// Remove the cron job
	svc.removeCronJob(task.Key)

	// Verify it's removed
	if _, exists := svc.cronEntryIDs[task.Key]; exists {
		t.Error("expected cron job to be removed")
	}
}

func TestRemoveCronJob_NonExistent(t *testing.T) {
	svc, _ := setupCronTestService(t)

	// Remove a non-existent cron job (should not panic)
	svc.removeCronJob("non-existent-key")
}

func TestStartCronJobs(t *testing.T) {
	svc, _ := setupCronTestService(t)

	// Create tasks in the database
	tasks := []struct {
		key  string
		cron string
	}{
		{"task-1", "* * * * * *"},
		{"task-2", "*/2 * * * * *"},
	}

	for _, tt := range tasks {
		if _, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
			Name:          tt.key,
			SourceRepoKey: "source",
			SourceBranch:  "main",
			TargetRepoKey: "target",
			TargetBranch:  "main",
			Cron:          tt.cron,
		}); err != nil {
			t.Fatalf("Create task failed: %v", err)
		}
	}

	err := svc.startCronJobs()
	if err != nil {
		t.Fatalf("startCronJobs failed: %v", err)
	}

	// Verify cron jobs were added
	if len(svc.cronEntryIDs) != 2 {
		t.Errorf("expected 2 cron jobs, got %d", len(svc.cronEntryIDs))
	}

	// Cleanup
	svc.stopCronJobs()
}

func TestStartCronJobs_WithDisabledTasks(t *testing.T) {
	svc, db := setupCronTestService(t)

	// Create tasks with some disabled
	task1, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
		Name: "task-1", Cron: "* * * * * *",
		SourceRepoKey: "source", SourceBranch: "main",
		TargetRepoKey: "target", TargetBranch: "main",
	})
	if err != nil {
		t.Fatalf("Create task1 failed: %v", err)
	}
	task2, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
		Name: "task-2", Cron: "*/2 * * * * *",
		SourceRepoKey: "source", SourceBranch: "main",
		TargetRepoKey: "target", TargetBranch: "main",
	})
	if err != nil {
		t.Fatalf("Create task2 failed: %v", err)
	}

	// Update task2 to set enabled=false using the actual key
	if _, err := svc.tasks.UpdateTask(context.TODO(), &model.UpdateTaskRequest{
		Key:     task2.Key,
		Enabled: false,
	}); err != nil {
		t.Fatalf("Update task2 failed: %v", err)
	}

	// Debug: check what's in the database
	var allTasks []*model.SyncTask
	db.Find(&allTasks)
	t.Logf("All tasks in DB: %d", len(allTasks))
	for _, task := range allTasks {
		t.Logf("  Task %s: enabled=%v, cron=%q", task.Key, task.Enabled, task.Cron)
	}

	enabled, err := svc.tasks.FindAllEnabledTasks()
	if err != nil {
		t.Fatalf("FindAllEnabled failed: %v", err)
	}
	t.Logf("Enabled tasks: %d", len(enabled))
	t.Logf("Task1 key: %s, Task2 key: %s", task1.Key, task2.Key)

	err = svc.startCronJobs()
	if err != nil {
		t.Fatalf("startCronJobs failed: %v", err)
	}

	// Only enabled tasks should have cron jobs
	if len(svc.cronEntryIDs) != 1 {
		t.Errorf("expected 1 cron job (only enabled), got %d", len(svc.cronEntryIDs))
	}

	// Cleanup
	svc.stopCronJobs()
}

func TestStartCronJobs_WithTasksWithoutCron(t *testing.T) {
	svc, _ := setupCronTestService(t)

	// Create tasks without cron
	if _, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
		Name: "task-1", Cron: "",
		SourceRepoKey: "source", SourceBranch: "main",
		TargetRepoKey: "target", TargetBranch: "main",
	}); err != nil {
		t.Fatalf("Create task failed: %v", err)
	}

	err := svc.startCronJobs()
	if err != nil {
		t.Fatalf("startCronJobs failed: %v", err)
	}

	// No cron jobs should be added
	if len(svc.cronEntryIDs) != 0 {
		t.Errorf("expected 0 cron jobs, got %d", len(svc.cronEntryIDs))
	}

	// Cleanup
	svc.stopCronJobs()
}

func TestStopCronJobs(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "* * * * * *",
	}

	// Add a cron job
	err := svc.addCronJob(task)
	if err != nil {
		t.Fatalf("addCronJob failed: %v", err)
	}

	// Stop cron jobs (should not panic)
	svc.stopCronJobs()
}
