package model

import (
	"testing"
	"time"
)

func TestWebhookRule_GetTaskKeys(t *testing.T) {
	rule := &WebhookRule{
		ID:   1,
		Name: "Test Rule",
		Tasks: []WebhookRuleTask{
			{TaskKey: "task1"},
			{TaskKey: "task2"},
			{TaskKey: "task3"},
		},
	}

	keys := rule.GetTaskKeys()

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}

	if keys[0] != "task1" {
		t.Errorf("expected key 'task1', got '%s'", keys[0])
	}

	if keys[1] != "task2" {
		t.Errorf("expected key 'task2', got '%s'", keys[1])
	}

	if keys[2] != "task3" {
		t.Errorf("expected key 'task3', got '%s'", keys[2])
	}
}

func TestWebhookRule_GetTaskKeys_Empty(t *testing.T) {
	rule := &WebhookRule{
		ID:    1,
		Name:  "Test Rule",
		Tasks: []WebhookRuleTask{},
	}

	keys := rule.GetTaskKeys()

	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestWebhookRule_GetTaskKeys_Nil(t *testing.T) {
	rule := &WebhookRule{
		ID:    1,
		Name:  "Test Rule",
		Tasks: nil,
	}

	keys := rule.GetTaskKeys()

	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestWebhookRule_SetTaskKeys(t *testing.T) {
	rule := &WebhookRule{
		ID:   1,
		Name: "Test Rule",
	}

	rule.SetTaskKeys([]string{"task1", "task2", "task3"})

	if len(rule.Tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(rule.Tasks))
	}

	if rule.Tasks[0].TaskKey != "task1" {
		t.Errorf("expected task key 'task1', got '%s'", rule.Tasks[0].TaskKey)
	}

	if rule.Tasks[1].TaskKey != "task2" {
		t.Errorf("expected task key 'task2', got '%s'", rule.Tasks[1].TaskKey)
	}

	if rule.Tasks[2].TaskKey != "task3" {
		t.Errorf("expected task key 'task3', got '%s'", rule.Tasks[2].TaskKey)
	}
}

func TestWebhookRule_SetTaskKeys_Empty(t *testing.T) {
	rule := &WebhookRule{
		ID:   1,
		Name: "Test Rule",
		Tasks: []WebhookRuleTask{
			{TaskKey: "old-task"},
		},
	}

	rule.SetTaskKeys([]string{})

	if len(rule.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(rule.Tasks))
	}
}

func TestWebhookRule_SetTaskKeys_Nil(t *testing.T) {
	rule := &WebhookRule{
		ID:   1,
		Name: "Test Rule",
		Tasks: []WebhookRuleTask{
			{TaskKey: "old-task"},
		},
	}

	rule.SetTaskKeys(nil)

	if len(rule.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(rule.Tasks))
	}
}

func TestWebhookRule_TableName(t *testing.T) {
	rule := WebhookRule{}
	expected := "webhook_rules"

	if rule.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, rule.TableName())
	}
}

func TestWebhookRuleTask_TableName(t *testing.T) {
	task := WebhookRuleTask{}
	expected := "webhook_rule_tasks"

	if task.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, task.TableName())
	}
}

func TestWebhookRule_Fields(t *testing.T) {
	now := time.Now()
	rule := WebhookRule{
		ID:            1,
		Name:          "Test Rule",
		RepoKey:       "test-repo",
		EventType:     "push",
		BranchPattern: "main",
		Action:        "sync",
		MinInterval:   60,
		Enabled:       true,
		Description:   "Test description",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if rule.ID != 1 {
		t.Errorf("expected ID 1, got %d", rule.ID)
	}

	if rule.Name != "Test Rule" {
		t.Errorf("expected name 'Test Rule', got '%s'", rule.Name)
	}

	if rule.RepoKey != "test-repo" {
		t.Errorf("expected repo key 'test-repo', got '%s'", rule.RepoKey)
	}

	if rule.EventType != "push" {
		t.Errorf("expected event type 'push', got '%s'", rule.EventType)
	}

	if rule.BranchPattern != "main" {
		t.Errorf("expected branch pattern 'main', got '%s'", rule.BranchPattern)
	}

	if rule.Action != "sync" {
		t.Errorf("expected action 'sync', got '%s'", rule.Action)
	}

	if rule.MinInterval != 60 {
		t.Errorf("expected min interval 60, got %d", rule.MinInterval)
	}

	if !rule.Enabled {
		t.Error("expected enabled to be true")
	}

	if rule.Description != "Test description" {
		t.Errorf("expected description 'Test description', got '%s'", rule.Description)
	}
}

func TestWebhookRuleTask_Fields(t *testing.T) {
	now := time.Now()
	task := WebhookRuleTask{
		ID:        1,
		RuleID:    1,
		TaskKey:   "task1",
		CreatedAt: now,
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.RuleID != 1 {
		t.Errorf("expected rule ID 1, got %d", task.RuleID)
	}

	if task.TaskKey != "task1" {
		t.Errorf("expected task key 'task1', got '%s'", task.TaskKey)
	}

	if task.CreatedAt != now {
		t.Errorf("expected created at %v, got %v", now, task.CreatedAt)
	}
}
