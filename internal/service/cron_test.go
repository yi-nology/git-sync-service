package service

import (
	"context"
	"sync"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func setupCronTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open test db")

	err = db.AutoMigrate(&model.SyncTask{})
	require.NoError(t, err, "failed to migrate test db")

	// sqlite :memory: 每个连接是独立库;强制单连接,确保 cron goroutine 能看到已迁移的表
	sqlDB, err := db.DB()
	require.NoError(t, err, "failed to get underlying db")
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
		cron:         cron.New(), // 与生产一致:标准 5 字段解析
		cronEntryIDs: make(map[string]cron.EntryID),
		cronMu:       sync.RWMutex{},
	}, db
}

func TestAddCronJob(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "* * * * *", // 标准五字段:每分钟
	}

	err := svc.addCronJob(task)
	require.NoError(t, err, "addCronJob failed")

	// Verify the cron job was added
	_, exists := svc.cronEntryIDs[task.Key]
	assert.True(t, exists, "expected cron job to be added")
}

func TestAddCronJob_UpdateExisting(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "*/5 * * * *",
	}

	// Add initial cron job
	err := svc.addCronJob(task)
	require.NoError(t, err, "addCronJob failed")

	initialEntryID := svc.cronEntryIDs[task.Key]

	// Update the cron job
	task.Cron = "*/2 * * * *"
	err = svc.addCronJob(task)
	require.NoError(t, err, "addCronJob update failed")

	// Verify the entry ID changed
	newEntryID := svc.cronEntryIDs[task.Key]
	assert.NotEqual(t, initialEntryID, newEntryID, "expected entry ID to change after update")
}

func TestAddCronJob_InvalidCron(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "invalid-cron-expression",
	}

	err := svc.addCronJob(task)
	require.Error(t, err, "expected error for invalid cron expression")
}

func TestRemoveCronJob(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "*/5 * * * *",
	}

	// Add a cron job first
	err := svc.addCronJob(task)
	require.NoError(t, err, "addCronJob failed")

	// Verify it exists
	_, exists := svc.cronEntryIDs[task.Key]
	require.True(t, exists, "expected cron job to exist before removal")

	// Remove the cron job
	svc.removeCronJob(task.Key)

	// Verify it's removed
	_, exists = svc.cronEntryIDs[task.Key]
	assert.False(t, exists, "expected cron job to be removed")
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
		{"task-1", "* * * * *"},
		{"task-2", "*/2 * * * *"},
	}

	for _, tt := range tasks {
		_, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
			Name:          tt.key,
			SourceRepoKey: "source",
			SourceBranch:  "main",
			TargetRepoKey: "target",
			TargetBranch:  "main",
			Cron:          tt.cron,
		})
		require.NoError(t, err, "Create task failed")
	}

	err := svc.startCronJobs()
	require.NoError(t, err, "startCronJobs failed")

	// Verify cron jobs were added
	assert.Equal(t, 2, len(svc.cronEntryIDs), "expected 2 cron jobs")

	// Cleanup
	svc.stopCronJobs()
}

func TestStartCronJobs_WithDisabledTasks(t *testing.T) {
	svc, db := setupCronTestService(t)

	// Create tasks with some disabled
	task1, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
		Name: "task-1", Cron: "*/5 * * * *",
		SourceRepoKey: "source", SourceBranch: "main",
		TargetRepoKey: "target", TargetBranch: "main",
	})
	require.NoError(t, err, "Create task1 failed")

	task2, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
		Name: "task-2", Cron: "*/2 * * * *",
		SourceRepoKey: "source", SourceBranch: "main",
		TargetRepoKey: "target", TargetBranch: "main",
	})
	require.NoError(t, err, "Create task2 failed")

	// Update task2 to set enabled=false using the actual key
	_, err = svc.tasks.UpdateTask(context.TODO(), &model.UpdateTaskRequest{
		Key:     task2.Key,
		Enabled: false,
	})
	require.NoError(t, err, "Update task2 failed")

	// Debug: check what's in the database
	var allTasks []*model.SyncTask
	db.Find(&allTasks)
	t.Logf("All tasks in DB: %d", len(allTasks))
	for _, task := range allTasks {
		t.Logf("  Task %s: enabled=%v, cron=%q", task.Key, task.Enabled, task.Cron)
	}

	enabled, err := svc.tasks.FindAllEnabledTasks()
	require.NoError(t, err, "FindAllEnabled failed")
	t.Logf("Enabled tasks: %d", len(enabled))
	t.Logf("Task1 key: %s, Task2 key: %s", task1.Key, task2.Key)

	err = svc.startCronJobs()
	require.NoError(t, err, "startCronJobs failed")

	// Only enabled tasks should have cron jobs
	assert.Equal(t, 1, len(svc.cronEntryIDs), "expected 1 cron job (only enabled)")

	// Cleanup
	svc.stopCronJobs()
}

func TestStartCronJobs_WithTasksWithoutCron(t *testing.T) {
	svc, _ := setupCronTestService(t)

	// Create tasks without cron
	_, err := svc.tasks.CreateTask(context.TODO(), &model.CreateTaskRequest{
		Name: "task-1", Cron: "",
		SourceRepoKey: "source", SourceBranch: "main",
		TargetRepoKey: "target", TargetBranch: "main",
	})
	require.NoError(t, err, "Create task failed")

	err = svc.startCronJobs()
	require.NoError(t, err, "startCronJobs failed")

	// No cron jobs should be added
	assert.Equal(t, 0, len(svc.cronEntryIDs), "expected 0 cron jobs")

	// Cleanup
	svc.stopCronJobs()
}

func TestStopCronJobs(t *testing.T) {
	svc, _ := setupCronTestService(t)

	task := &model.SyncTask{
		Key:  "test-task-1",
		Cron: "*/5 * * * *",
	}

	// Add a cron job
	err := svc.addCronJob(task)
	require.NoError(t, err, "addCronJob failed")

	// Stop cron jobs (should not panic)
	svc.stopCronJobs()
}

func TestAddCronJob_StandardFiveFieldExpressions(t *testing.T) {
	// 用户与前端预设输入的都是标准 5 字段 crontab;曾因生产用 6 字段
	// (WithSeconds)解析器导致 "expected exactly 6 fields" 全部失败。
	for _, expr := range []string{"* * * * *", "*/5 * * * *", "0 0 * * *", "0 9 * * 1-5", "30 2 1 * *"} {
		svc, _ := setupCronTestService(t)
		task := &model.SyncTask{Key: "five-field-" + expr, Cron: expr}
		err := svc.addCronJob(task)
		assert.NoError(t, err, "标准五字段表达式 %q 注册失败", expr)
	}
}
