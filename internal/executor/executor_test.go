package executor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/yi-nology/git-platform-sdk/gitbackend"
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

	config := exec.authConfig(repo, nil)
	if config.Type != gitbackend.AuthHTTPBasic {
		t.Errorf("expected HTTP Basic auth, got %q", config.Type)
	}
	if config.Password != "test-token" {
		t.Errorf("expected password 'test-token', got %q", config.Password)
	}
	if config.Username == "" {
		t.Error("expected non-empty placeholder username for HTTP Basic")
	}
}

func TestAuthConfig_WithoutToken(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	repo := &model.Repo{}

	config := exec.authConfig(repo, nil)
	if config.Type != gitbackend.AuthNone {
		t.Errorf("expected auth type none, got %q", config.Type)
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

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"nil error", nil, ""},
		{"auth 401", fmt.Errorf("401 Unauthorized"), model.ErrorAuth},
		{"auth 403", fmt.Errorf("403 Forbidden"), model.ErrorAuth},
		{"auth authentication", fmt.Errorf("authentication failed"), model.ErrorAuth},
		{"auth unauthorized", fmt.Errorf("unauthorized access"), model.ErrorAuth},
		{"auth permission denied", fmt.Errorf("permission denied"), model.ErrorAuth},
		{"auth token", fmt.Errorf("invalid token"), model.ErrorAuth},
		{"auth credential", fmt.Errorf("credential error"), model.ErrorAuth},
		{"network timeout", fmt.Errorf("connection timeout"), model.ErrorNetwork},
		{"network connection refused", fmt.Errorf("connection refused"), model.ErrorNetwork},
		{"network dial tcp", fmt.Errorf("dial tcp: lookup example.com"), model.ErrorNetwork},
		{"network network", fmt.Errorf("network error"), model.ErrorNetwork},
		{"network i/o error", fmt.Errorf("i/o error"), model.ErrorNetwork},
		{"network eof", fmt.Errorf("unexpected eof"), model.ErrorNetwork},
		{"network no route", fmt.Errorf("no route to host"), model.ErrorNetwork},
		{"config not found", fmt.Errorf("repository not found"), model.ErrorConfig},
		{"config does not exist", fmt.Errorf("branch does not exist"), model.ErrorConfig},
		{"config repository not found", fmt.Errorf("repository not found"), model.ErrorConfig},
		{"config branch not found", fmt.Errorf("branch not found"), model.ErrorConfig},
		{"git non-fast-forward", fmt.Errorf("non-fast-forward"), model.ErrorGit},
		{"git conflict", fmt.Errorf("merge conflict"), model.ErrorGit},
		{"git rejected", fmt.Errorf("push rejected"), model.ErrorGit},
		{"git failed to push", fmt.Errorf("failed to push"), model.ErrorGit},
		{"git fetch first", fmt.Errorf("fetch first"), model.ErrorGit},
		{"unknown error", fmt.Errorf("some random error"), model.ErrorUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if result != tt.expected {
				t.Errorf("ClassifyError(%v) = %q, want %q", tt.err, result, tt.expected)
			}
		})
	}
}

func TestExecute_WithContext(t *testing.T) {
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

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	run, err := exec.Execute(ctx, task, "manual", nil)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if run == nil {
		t.Fatal("expected non-nil run even on error")
	}
}

func TestExecute_WithWebhookEventID(t *testing.T) {
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

	webhookEventID := uint(123)
	ctx := context.Background()
	_, _ = exec.Execute(ctx, task, "webhook", &webhookEventID)

	// Verify run was created with webhook event ID
	if len(svc.runs) == 0 {
		t.Fatal("expected run to be created")
	}

	run := svc.runs[0]
	if run.WebhookEventID == nil {
		t.Fatal("expected webhook event ID to be set")
	}

	if *run.WebhookEventID != 123 {
		t.Errorf("expected webhook event ID 123, got %d", *run.WebhookEventID)
	}
}

func TestBeginStep(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	// Create a run
	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		TargetRepoKey: "target-repo",
	}

	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("CreateRun failed: %v", err)
	}

	// Begin a step
	step := exec.beginStep(run.ID, "fetch")
	if step == nil {
		t.Fatal("expected non-nil step")
	}

	if step.RunID != run.ID {
		t.Errorf("expected run ID %d, got %d", run.ID, step.RunID)
	}

	if step.StepName != "fetch" {
		t.Errorf("expected step name 'fetch', got '%s'", step.StepName)
	}

	if step.Status != "running" {
		t.Errorf("expected status 'running', got '%s'", step.Status)
	}
}

func TestFailStep(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	// Create a run
	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		TargetRepoKey: "target-repo",
	}

	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("CreateRun failed: %v", err)
	}

	// Begin a step
	step := exec.beginStep(run.ID, "fetch")
	if step == nil {
		t.Fatal("expected non-nil step")
	}

	// Fail the step
	exec.failStep(step, fmt.Errorf("fetch failed"))

	if step.Status != "failed" {
		t.Errorf("expected status 'failed', got '%s'", step.Status)
	}

	if step.ErrorMsg != "fetch failed" {
		t.Errorf("expected error message 'fetch failed', got '%s'", step.ErrorMsg)
	}
}

func TestCompleteStep(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	// Create a run
	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "source-repo",
		TargetRepoKey: "target-repo",
	}

	run, err := svc.CreateRun(task, "manual", nil)
	if err != nil {
		t.Fatalf("CreateRun failed: %v", err)
	}

	// Begin a step
	step := exec.beginStep(run.ID, "fetch")
	if step == nil {
		t.Fatal("expected non-nil step")
	}

	// Complete the step
	exec.completeStep(step, "details")

	if step.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", step.Status)
	}
}

func TestAuthConfig_WithPlatform(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	repo := &model.Repo{
		AccessToken: "",
		PlatformID:  1,
	}

	config := exec.authConfig(repo, nil)
	// mockService returns a platform with AccessToken="test-token"
	if config.Password != "test-token" {
		t.Errorf("expected platform token as password, got %q", config.Password)
	}
}

func TestAuthConfig_RepoTokenOverridesPlatform(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	repo := &model.Repo{
		AccessToken: "repo-specific-token",
		PlatformID:  1,
	}

	config := exec.authConfig(repo, nil)
	// Repo token should take precedence
	if config.Password != "repo-specific-token" {
		t.Errorf("expected repo token as password, got %q", config.Password)
	}
}

func TestAuthConfig_NoPlatformID(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	repo := &model.Repo{
		AccessToken: "",
		PlatformID:  0,
	}

	config := exec.authConfig(repo, nil)
	if config.Type != "none" {
		t.Errorf("expected auth type 'none', got '%q'", config.Type)
	}
}

func TestCloneRepo_WithForceFlag(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		SourceBranch: "main",
		GitForce:     true,
	}

	repo := &model.Repo{
		CloneURL: "https://github.com/test/repo.git",
	}

	// This will fail because we don't have a real git backend,
	// but it exercises the depth=0 path
	_ = exec.cloneRepo(context.Background(), "/tmp/test-clone", repo, task, nil)
}

func TestCloneRepo_ShallowClone(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		SourceBranch: "main",
		GitForce:     false,
	}

	repo := &model.Repo{
		CloneURL: "https://github.com/test/repo.git",
	}

	// This will fail because we don't have a real git backend,
	// but it exercises the depth=1 path
	_ = exec.cloneRepo(context.Background(), "/tmp/test-clone", repo, task, nil)
}

func TestNewExecutor_WithBackendConfig(t *testing.T) {
	svc := &mockServiceWithConfig{
		mockService: newMockService(),
		backendType: "gogit",
	}
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}
	if exec == nil {
		t.Fatal("expected non-nil executor")
	}
}

func TestNewExecutor_InvalidBackend(t *testing.T) {
	svc := &mockServiceWithConfig{
		mockService: newMockService(),
		backendType: "invalid-backend-type",
	}
	_, err := NewExecutor(svc)
	if err == nil {
		t.Fatal("expected error for invalid backend type")
	}
}

// mockServiceWithConfig wraps mockService to return a specific config
type mockServiceWithConfig struct {
	*mockService
	backendType string
}

func (m *mockServiceWithConfig) GetConfig() *model.Config {
	return &model.Config{
		Git: model.GitConfig{
			Backend: m.backendType,
		},
		Sync: model.SyncConfig{
			DefaultTimeout: 300,
			RetryCount:     3,
		},
	}
}

func TestExecute_CompleteRunError(t *testing.T) {
	svc := newMockService()
	svc.completeErr = fmt.Errorf("db write failed")
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "missing-repo", // Will fail, triggering complete
		TargetRepoKey: "target-repo",
	}

	ctx := context.Background()
	run, err := exec.Execute(ctx, task, "manual", nil)
	// Should still return the run even with complete error
	if run == nil {
		t.Fatal("expected non-nil run")
	}
	if err == nil {
		t.Fatal("expected error for missing source repo")
	}
}

func TestExecute_UpdateTaskLastRunError(t *testing.T) {
	svc := newMockService()
	svc.updateErr = fmt.Errorf("db update failed")
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "missing-repo", // Will fail, triggering update
		TargetRepoKey: "target-repo",
	}

	ctx := context.Background()
	run, _ := exec.Execute(ctx, task, "manual", nil)
	if run == nil {
		t.Fatal("expected non-nil run")
	}
}

func TestExecute_SourceRepoQueryError(t *testing.T) {
	svc := newMockService()
	svc.repoErr = fmt.Errorf("database connection lost")
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
	run, err := exec.Execute(ctx, task, "manual", nil)
	if err == nil {
		t.Fatal("expected error when repo query fails")
	}
	if run == nil {
		t.Fatal("expected non-nil run")
	}
	if run.Status != model.StatusFailed {
		t.Errorf("expected status 'failed', got '%q'", run.Status)
	}
	if run.ErrorType != model.ErrorUnknown {
		t.Errorf("expected error type '%s', got '%q'", model.ErrorUnknown, run.ErrorType)
	}
}

func TestFailRun(t *testing.T) {
	run := &model.SyncRun{
		ID:     1,
		Status: model.StatusRunning,
	}

	err := fmt.Errorf("401 authentication failed")
	resultRun, resultErr := failRun(run, err)

	if resultRun.Status != model.StatusFailed {
		t.Errorf("expected status 'failed', got '%q'", resultRun.Status)
	}
	if resultRun.ErrorMessage != "401 authentication failed" {
		t.Errorf("expected error message '401 authentication failed', got '%q'", resultRun.ErrorMessage)
	}
	if resultRun.ErrorType != model.ErrorAuth {
		t.Errorf("expected error type '%s', got '%q'", model.ErrorAuth, resultRun.ErrorType)
	}
	if resultErr == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestFetchRepo(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		SourceBranch: "main",
		GitTags:      true,
		GitPrune:     true,
	}

	repo := &model.Repo{
		CloneURL: "https://github.com/test/repo.git",
	}

	// This will fail because we don't have a real git repo,
	// but it exercises the fetchRepo code path
	_ = exec.fetchRepo(context.Background(), "/tmp/nonexistent", task, repo, nil)
}

func TestEnsureRemote(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	repo := &model.Repo{
		CloneURL: "https://github.com/test/target.git",
	}

	// This will fail because we don't have a real git repo,
	// but it exercises the ensureRemote code path
	_ = exec.ensureRemote(context.Background(), "/tmp/nonexistent", repo)
}

func TestPush(t *testing.T) {
	svc := newMockService()
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		SourceBranch: "main",
		TargetBranch: "main",
		GitForce:     false,
	}

	repo := &model.Repo{
		CloneURL: "https://github.com/test/target.git",
	}

	// This will fail because we don't have a real git repo,
	// but it exercises the push code path
	_ = exec.push(context.Background(), "/tmp/nonexistent", task, repo, nil)
}

func TestExecute_WithDefaultTimeout(t *testing.T) {
	svc := &mockServiceWithConfig{
		mockService: newMockService(),
		backendType: "",
	}
	svc.config = &model.Config{
		Sync: model.SyncConfig{
			DefaultTimeout: 0, // Should default to 300
			RetryCount:     1,
		},
	}
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "missing-repo",
		TargetRepoKey: "target-repo",
	}

	ctx := context.Background()
	run, _ := exec.Execute(ctx, task, "manual", nil)
	if run == nil {
		t.Fatal("expected non-nil run")
	}
}

func TestExecute_WithZeroRetryCount(t *testing.T) {
	svc := &mockServiceWithConfig{
		mockService: newMockService(),
		backendType: "",
	}
	svc.config = &model.Config{
		Sync: model.SyncConfig{
			DefaultTimeout: 300,
			RetryCount:     0, // Should default to 3
		},
	}
	exec, err := NewExecutor(svc)
	if err != nil {
		t.Fatalf("NewExecutor failed: %v", err)
	}

	task := &model.SyncTask{
		Key:           "test-task",
		Name:          "Test Task",
		SourceRepoKey: "missing-repo",
		TargetRepoKey: "target-repo",
	}

	ctx := context.Background()
	run, _ := exec.Execute(ctx, task, "manual", nil)
	if run == nil {
		t.Fatal("expected non-nil run")
	}
}
