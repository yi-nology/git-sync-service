package executor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
)

// mockService implements Service interface for testing.
type mockService struct {
	runs        []*model.SyncRun
	tasks       []*model.SyncTask
	repos       map[string]*model.Repo
	createErr   error
	completeErr error
	updateErr   error
	repoErr     error
	tempDir     string
	config      *model.Config
}

func (m *mockService) CreateRun(task *model.SyncTask, trigger string, webhookEventID *uint) (*model.SyncRun, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	run := &model.SyncRun{
		TaskKey:        task.Key,
		TriggerSource:  trigger,
		Status:         "running",
		StartTime:      time.Now(),
		WebhookEventID: webhookEventID,
	}
	run.ID = uint(len(m.runs) + 1)
	m.runs = append(m.runs, run)
	return run, nil
}

func (m *mockService) CreateRunStep(step *model.SyncRunStep) error {
	return nil
}

func (m *mockService) UpdateRunStep(step *model.SyncRunStep) error {
	return nil
}

func (m *mockService) CompleteRun(run *model.SyncRun) error {
	if m.completeErr != nil {
		return m.completeErr
	}
	for i, r := range m.runs {
		if r.ID == run.ID {
			m.runs[i] = run
			break
		}
	}
	return nil
}

func (m *mockService) UpdateTaskLastRun(task *model.SyncTask, run *model.SyncRun) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	task.LastRunAt = run.EndTime
	task.LastStatus = run.Status
	return nil
}

func (m *mockService) GetRepoByKey(key string) (*model.Repo, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	repo, ok := m.repos[key]
	if !ok {
		return nil, nil
	}
	return repo, nil
}

func (m *mockService) GetTempDir(taskKey string) string {
	return m.tempDir
}

func (m *mockService) GetConfig() *model.Config {
	return m.config
}

func (m *mockService) GetPlatformByID(id uint) (*model.Platform, error) {
	return &model.Platform{
		ID:          id,
		Key:         "test-platform",
		Name:        "Test Platform",
		Type:        "github",
		AccessToken: "test-token",
	}, nil
}

func newMockService() *mockService {
	return &mockService{
		runs:  make([]*model.SyncRun, 0),
		tasks: make([]*model.SyncTask, 0),
		repos: map[string]*model.Repo{
			"source-repo": {
				Key:         "source-repo",
				Name:        "Source Repo",
				CloneURL:    "https://github.com/source/repo.git",
				AccessToken: "source-token",
			},
			"target-repo": {
				Key:         "target-repo",
				Name:        "Target Repo",
				CloneURL:    "https://github.com/target/repo.git",
				AccessToken: "target-token",
			},
		},
		tempDir: "/tmp/git-sync-test",
		config: &model.Config{
			Sync: model.SyncConfig{
				DefaultTimeout: 300,
				RetryCount:     3,
			},
		},
	}
}

func TestNewExecutor(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}
	if exec == nil {
		t.Fatal("expected non-nil executor")
	}
	if exec.service != svc {
		t.Error("expected service to be set")
	}
	if exec.backend == nil {
		t.Error("expected backend to be set")
	}
}

func TestExecute_SourceRepoNotFound(t *testing.T) {
	svc := newMockService()
	svc.repos = map[string]*model.Repo{} // Empty repos

	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "missing-repo",
		TargetRepoKey: "target-repo",
		SourceBranch:  "main",
		TargetBranch:  "main",
	}

	ctx := context.Background()
	run, err := exec.Execute(ctx, task, "manual", nil)
	if err == nil {
		t.Fatal("expected error for missing source repo")
	}
	if run == nil {
		t.Fatal("expected non-nil run even on error")
	}
	if run.Status != "failed" {
		t.Errorf("expected status 'failed', got %q", run.Status)
	}
}

func TestExecute_TargetRepoNotFound(t *testing.T) {
	svc := newMockService()
	svc.repos = map[string]*model.Repo{
		"source-repo": {
			Key:      "source-repo",
			CloneURL: "https://github.com/source/repo.git",
		},
	}

	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		TargetRepoKey: "missing-repo",
		SourceBranch:  "main",
		TargetBranch:  "main",
	}

	ctx := context.Background()
	run, err := exec.Execute(ctx, task, "manual", nil)
	if err == nil {
		t.Fatal("expected error for missing target repo")
	}
	if run == nil {
		t.Fatal("expected non-nil run even on error")
	}
	if run.Status != "failed" {
		t.Errorf("expected status 'failed', got %q", run.Status)
	}
}

func TestExecute_RunCreated(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		TargetRepoKey: "target-repo",
		SourceBranch:  "main",
		TargetBranch:  "main",
	}

	ctx := context.Background()
	_, _ = exec.Execute(ctx, task, "manual", nil)

	// Verify run was created
	if len(svc.runs) == 0 {
		t.Fatal("expected run to be created")
	}

	run := svc.runs[0]
	if run.TaskKey != "test-task" {
		t.Errorf("expected task key 'test-task', got %q", run.TaskKey)
	}
	if run.TriggerSource != "manual" {
		t.Errorf("expected trigger 'manual', got %q", run.TriggerSource)
	}
	if run.StartTime.IsZero() {
		t.Error("expected start time to be set")
	}
}

func TestExecute_CreateRunError(t *testing.T) {
	svc := newMockService()
	svc.createErr = fmt.Errorf("database error")

	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		TargetRepoKey: "target-repo",
		SourceBranch:  "main",
		TargetBranch:  "main",
	}

	ctx := context.Background()
	_, err = exec.Execute(ctx, task, "manual", nil)
	if err == nil {
		t.Fatal("expected error when run creation fails")
	}
}

func TestAuthConfig_WithToken(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	repo := &model.Repo{
		AccessToken: "test-token",
	}

	config := exec.authConfig(repo)
	if config.Token != "test-token" {
		t.Errorf("expected token 'test-token', got %q", config.Token)
	}
}

func TestAuthConfig_WithoutToken(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	repo := &model.Repo{}

	config := exec.authConfig(repo)
	if config.Token != "" {
		t.Errorf("expected empty token, got %q", config.Token)
	}
}

func TestTimePtr(t *testing.T) {
	now := time.Now()
	ptr := timePtr(now)
	if ptr == nil {
		t.Fatal("expected non-nil time pointer")
	}
	if *ptr != now {
		t.Errorf("expected %v, got %v", now, *ptr)
	}
}
