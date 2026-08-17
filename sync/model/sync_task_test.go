package model

import (
	"testing"
	"time"
)

func TestSyncTask_TableName(t *testing.T) {
	tk := SyncTask{}
	expected := "sync_tasks"

	if tk.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, tk.TableName())
	}
}

func TestSyncTask_Fields(t *testing.T) {
	now := time.Now()
	tk := SyncTask{
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

	if tk.ID != 1 {
		t.Errorf("expected ID 1, got %d", tk.ID)
	}

	if tk.Key != "test-task" {
		t.Errorf("expected key 'test-task', got '%s'", tk.Key)
	}

	if tk.Name != "Test Task" {
		t.Errorf("expected name 'Test Task', got '%s'", tk.Name)
	}

	if tk.SourceRepoKey != "source-repo" {
		t.Errorf("expected source repo key 'source-repo', got '%s'", tk.SourceRepoKey)
	}

	if tk.SourceBranch != "main" {
		t.Errorf("expected source branch 'main', got '%s'", tk.SourceBranch)
	}

	if tk.TargetRepoKey != "target-repo" {
		t.Errorf("expected target repo key 'target-repo', got '%s'", tk.TargetRepoKey)
	}

	if tk.TargetBranch != "main" {
		t.Errorf("expected target branch 'main', got '%s'", tk.TargetBranch)
	}

	if tk.SyncMode != "mirror" {
		t.Errorf("expected sync mode 'mirror', got '%s'", tk.SyncMode)
	}

	if tk.Cron != "0 * * * *" {
		t.Errorf("expected cron '0 * * * *', got '%s'", tk.Cron)
	}

	if tk.WebhookToken != "token123" {
		t.Errorf("expected webhook token 'token123', got '%s'", tk.WebhookToken)
	}

	if !tk.Enabled {
		t.Error("expected enabled to be true")
	}

	if !tk.GitTags {
		t.Error("expected git tags to be true")
	}

	if tk.GitForce {
		t.Error("expected git force to be false")
	}

	if !tk.GitPrune {
		t.Error("expected git prune to be true")
	}

	if tk.LastStatus != "success" {
		t.Errorf("expected last status 'success', got '%s'", tk.LastStatus)
	}
}
