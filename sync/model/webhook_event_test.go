package model

import (
	"testing"
	"time"
)

func TestWebhookEvent_TableName(t *testing.T) {
	e := WebhookEvent{}
	expected := "webhook_events"

	if e.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, e.TableName())
	}
}

func TestWebhookEvent_Fields(t *testing.T) {
	now := time.Now()
	e := WebhookEvent{
		ID:           1,
		EventID:      "evt-123",
		RepoKey:      "test-repo",
		EventType:    "push",
		Source:       "github",
		ActorName:    "testuser",
		Branch:       "main",
		CommitSHA:    "abc123",
		Status:       "processed",
		ErrorMessage: "",
		ProcessedAt:  &now,
		CreatedAt:    now,
	}

	if e.ID != 1 {
		t.Errorf("expected ID 1, got %d", e.ID)
	}

	if e.EventID != "evt-123" {
		t.Errorf("expected event ID 'evt-123', got '%s'", e.EventID)
	}

	if e.RepoKey != "test-repo" {
		t.Errorf("expected repo key 'test-repo', got '%s'", e.RepoKey)
	}

	if e.EventType != "push" {
		t.Errorf("expected event type 'push', got '%s'", e.EventType)
	}

	if e.Source != "github" {
		t.Errorf("expected source 'github', got '%s'", e.Source)
	}

	if e.ActorName != "testuser" {
		t.Errorf("expected actor name 'testuser', got '%s'", e.ActorName)
	}

	if e.Branch != "main" {
		t.Errorf("expected branch 'main', got '%s'", e.Branch)
	}

	if e.CommitSHA != "abc123" {
		t.Errorf("expected commit SHA 'abc123', got '%s'", e.CommitSHA)
	}

	if e.Status != "processed" {
		t.Errorf("expected status 'processed', got '%s'", e.Status)
	}

	if e.ErrorMessage != "" {
		t.Errorf("expected empty error message, got '%s'", e.ErrorMessage)
	}

	if e.ProcessedAt == nil {
		t.Error("expected processed at to be set")
	}
}
