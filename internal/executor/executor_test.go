package executor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
)

// mockRunWriter implements RunWriter interface
type mockRunWriter struct {
	runs    []*model.SyncRun
	createErr error
	updateErr error
}

func (m *mockRunWriter) Create(run *model.SyncRun) error {
	if m.createErr != nil {
		return m.createErr
	}
	run.ID = uint(len(m.runs) + 1)
	m.runs = append(m.runs, run)
	return nil
}

func (m *mockRunWriter) Update(run *model.SyncRun) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i, r := range m.runs {
		if r.ID == run.ID {
			m.runs[i] = run
			break
		}
	}
	return nil
}

// mockTaskUpdater implements TaskUpdater interface
type mockTaskUpdater struct {
	tasks     []*model.SyncTask
	updateErr error
}

func (m *mockTaskUpdater) Update(task *model.SyncTask) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i, t := range m.tasks {
		if t.Key == task.Key {
			m.tasks[i] = task
			break
		}
	}
	return nil
}

// mockRepoReader implements RepoReader interface
type mockRepoReader struct {
	repos map[string]*model.Repo
	err   error
}

func (m *mockRepoReader) FindByKey(key string) (*model.Repo, error) {
	if m.err != nil {
		return nil, m.err
	}
	repo, ok := m.repos[key]
	if !ok {
		return nil, fmt.Errorf("repo not found: %s", key)
	}
	return repo, nil
}

// mockService implements Service interface
type mockService struct {
	runWriter   *mockRunWriter
	taskUpdater *mockTaskUpdater
	repoReader  *mockRepoReader
	tempDir     string
	config      *model.Config
}

func (m *mockService) RunDAO() RunWriter {
	return m.runWriter
}

func (m *mockService) TaskDAO() TaskUpdater {
	return m.taskUpdater
}

func (m *mockService) RepoDAO() RepoReader {
	return m.repoReader
}

func (m *mockService) GetTempDir(taskKey string) string {
	return m.tempDir
}

func (m *mockService) GetConfig() *model.Config {
	return m.config
}

func newMockService() *mockService {
	return &mockService{
		runWriter: &mockRunWriter{
			runs: make([]*model.SyncRun, 0),
		},
		taskUpdater: &mockTaskUpdater{
			tasks: make([]*model.SyncTask, 0),
		},
		repoReader: &mockRepoReader{
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
	svc.repoReader.repos = map[string]*model.Repo{} // Empty repos

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
	run, err := exec.Execute(ctx, task, "manual")
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
	svc.repoReader.repos = map[string]*model.Repo{
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
	run, err := exec.Execute(ctx, task, "manual")
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
	_, _ = exec.Execute(ctx, task, "manual")

	// Verify run was created
	if len(svc.runWriter.runs) == 0 {
		t.Fatal("expected run to be created")
	}

	run := svc.runWriter.runs[0]
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
	svc.runWriter.createErr = fmt.Errorf("database error")

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
	_, err = exec.Execute(ctx, task, "manual")
	if err == nil {
		t.Fatal("expected error when run creation fails")
	}
}

func TestPreview_SourceRepoNotFound(t *testing.T) {
	svc := newMockService()
	svc.repoReader.repos = map[string]*model.Repo{}

	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		SourceRepoKey: "missing-repo",
		TargetRepoKey: "target-repo",
	}

	ctx := context.Background()
	preview, err := exec.Preview(ctx, task)
	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}
	if preview.SourceExists {
		t.Error("expected source to not exist")
	}
	if preview.CanSync {
		t.Error("expected canSync to be false")
	}
}

func TestPreview_TargetRepoNotFound(t *testing.T) {
	svc := newMockService()
	svc.repoReader.repos = map[string]*model.Repo{
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
		SourceRepoKey: "source-repo",
		TargetRepoKey: "missing-repo",
	}

	ctx := context.Background()
	preview, err := exec.Preview(ctx, task)
	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}
	if !preview.SourceExists {
		t.Error("expected source to exist")
	}
	if preview.TargetExists {
		t.Error("expected target to not exist")
	}
	if preview.CanSync {
		t.Error("expected canSync to be false")
	}
}

func TestPreview_CanSync(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		SourceRepoKey: "source-repo",
		TargetRepoKey: "target-repo",
	}

	ctx := context.Background()
	preview, err := exec.Preview(ctx, task)
	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}
	if !preview.SourceExists {
		t.Error("expected source to exist")
	}
	if !preview.TargetExists {
		t.Error("expected target to exist")
	}
	if !preview.CanSync {
		t.Error("expected canSync to be true")
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
