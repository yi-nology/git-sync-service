package model

import (
	"testing"
	"time"
)

func TestSyncRun_TableName(t *testing.T) {
	r := SyncRun{}
	expected := "sync_runs"

	if r.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, r.TableName())
	}
}

func TestSyncRun_Fields(t *testing.T) {
	now := time.Now()
	r := SyncRun{
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

	if r.ID != 1 {
		t.Errorf("expected ID 1, got %d", r.ID)
	}

	if r.TaskKey != "test-task" {
		t.Errorf("expected task key 'test-task', got '%s'", r.TaskKey)
	}

	if r.TriggerSource != "manual" {
		t.Errorf("expected trigger source 'manual', got '%s'", r.TriggerSource)
	}

	if r.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", r.Status)
	}

	if r.CommitRange != "abc123..def456" {
		t.Errorf("expected commit range 'abc123..def456', got '%s'", r.CommitRange)
	}

	if r.Details != "Sync completed" {
		t.Errorf("expected details 'Sync completed', got '%s'", r.Details)
	}

	if r.ErrorMessage != "" {
		t.Errorf("expected empty error message, got '%s'", r.ErrorMessage)
	}

	if r.ErrorType != "" {
		t.Errorf("expected empty error type, got '%s'", r.ErrorType)
	}

	if r.DurationMs != 1000 {
		t.Errorf("expected duration ms 1000, got %d", r.DurationMs)
	}

	if r.RetryTotal != 0 {
		t.Errorf("expected retry total 0, got %d", r.RetryTotal)
	}

	if r.WebhookEventID != nil {
		t.Error("expected webhook event ID to be nil")
	}
}
