package model

import (
	"testing"
	"time"
)

func TestSyncRunStep_TableName(t *testing.T) {
	s := SyncRunStep{}
	expected := "sync_run_steps"

	if s.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, s.TableName())
	}
}

func TestSyncRunStep_Fields(t *testing.T) {
	now := time.Now()
	s := SyncRunStep{
		ID:          1,
		RunID:       1,
		StepName:    "fetch",
		Status:      "success",
		StartTime:   now,
		EndTime:     &now,
		DurationMs:  500,
		ErrorMsg:    "",
		ErrorType:   "",
		RetryCount:  0,
		CreatedAt:   now,
	}

	if s.ID != 1 {
		t.Errorf("expected ID 1, got %d", s.ID)
	}

	if s.RunID != 1 {
		t.Errorf("expected run ID 1, got %d", s.RunID)
	}

	if s.StepName != "fetch" {
		t.Errorf("expected step name 'fetch', got '%s'", s.StepName)
	}

	if s.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", s.Status)
	}

	if s.DurationMs != 500 {
		t.Errorf("expected duration ms 500, got %d", s.DurationMs)
	}

	if s.ErrorMsg != "" {
		t.Errorf("expected empty error message, got '%s'", s.ErrorMsg)
	}

	if s.ErrorType != "" {
		t.Errorf("expected empty error type, got '%s'", s.ErrorType)
	}

	if s.RetryCount != 0 {
		t.Errorf("expected retry count 0, got %d", s.RetryCount)
	}
}
