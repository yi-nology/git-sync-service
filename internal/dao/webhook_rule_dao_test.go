package dao

import (
	"testing"

	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupWebhookRuleTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.WebhookRule{}, &model.WebhookRuleTask{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return db
}

func TestWebhookRuleDAO_CreateAndFindByID(t *testing.T) {
	db := setupWebhookRuleTestDB(t)
	d := NewWebhookRuleDAO(db)

	// Create a rule
	rule := &model.WebhookRule{
		Name:          "Test Rule",
		RepoKey:       "test-repo",
		EventType:     "push",
		BranchPattern: "main",
		Action:        "sync",
		MinInterval:   60,
		Enabled:       true,
		Description:   "Test webhook rule",
	}

	if err := d.Create(rule); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if rule.ID == 0 {
		t.Error("expected rule ID to be set after create")
	}

	// Get by ID
	got, err := d.FindByID(rule.ID)
	if err != nil {
		t.Fatalf("find by ID failed: %v", err)
	}

	if got.Name != "Test Rule" {
		t.Errorf("expected name 'Test Rule', got '%s'", got.Name)
	}

	if got.RepoKey != "test-repo" {
		t.Errorf("expected repo key 'test-repo', got '%s'", got.RepoKey)
	}

	if got.EventType != "push" {
		t.Errorf("expected event type 'push', got '%s'", got.EventType)
	}

	if got.BranchPattern != "main" {
		t.Errorf("expected branch pattern 'main', got '%s'", got.BranchPattern)
	}

	if got.Action != "sync" {
		t.Errorf("expected action 'sync', got '%s'", got.Action)
	}

	if got.MinInterval != 60 {
		t.Errorf("expected min interval 60, got %d", got.MinInterval)
	}

	if !got.Enabled {
		t.Error("expected enabled to be true")
	}

	if got.Description != "Test webhook rule" {
		t.Errorf("expected description 'Test webhook rule', got '%s'", got.Description)
	}
}

func TestWebhookRuleDAO_Update(t *testing.T) {
	db := setupWebhookRuleTestDB(t)
	d := NewWebhookRuleDAO(db)

	// Create a rule
	rule := &model.WebhookRule{
		Name:        "Test Rule",
		RepoKey:     "test-repo",
		EventType:   "push",
		Action:      "sync",
		MinInterval: 60,
	}

	if err := d.Create(rule); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update the rule
	rule.Name = "Updated Rule"
	rule.MinInterval = 120

	if err := d.Update(rule); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Get the updated rule
	got, err := d.FindByID(rule.ID)
	if err != nil {
		t.Fatalf("find by ID failed: %v", err)
	}

	if got.Name != "Updated Rule" {
		t.Errorf("expected name 'Updated Rule', got '%s'", got.Name)
	}

	if got.MinInterval != 120 {
		t.Errorf("expected min interval 120, got %d", got.MinInterval)
	}
}

func TestWebhookRuleDAO_Delete(t *testing.T) {
	db := setupWebhookRuleTestDB(t)
	d := NewWebhookRuleDAO(db)

	// Create a rule
	rule := &model.WebhookRule{
		Name:    "Test Rule",
		RepoKey: "test-repo",
	}

	if err := d.Create(rule); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete the rule
	if err := d.Delete(rule.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Try to get the deleted rule
	got, err := d.FindByID(rule.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting deleted rule")
	}
}

func TestWebhookRuleDAO_FindByID_NotFound(t *testing.T) {
	db := setupWebhookRuleTestDB(t)
	d := NewWebhookRuleDAO(db)

	// Try to get a non-existent rule
	got, err := d.FindByID(999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent rule")
	}
}

func TestWebhookRuleDAO_FindByRepoKey(t *testing.T) {
	db := setupWebhookRuleTestDB(t)
	d := NewWebhookRuleDAO(db)

	// Create rules with different repo keys
	rules := []*model.WebhookRule{
		{Name: "Rule 1", RepoKey: "repo1", EventType: "push"},
		{Name: "Rule 2", RepoKey: "repo1", EventType: "push"},
		{Name: "Rule 3", RepoKey: "repo2", EventType: "push"},
	}

	for _, rule := range rules {
		if err := d.Create(rule); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find by repo key
	got, err := d.FindByRepoKey("repo1")
	if err != nil {
		t.Fatalf("find by repo key failed: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(got))
	}
}

func TestWebhookRuleDAO_FindTasksByRuleID(t *testing.T) {
	db := setupWebhookRuleTestDB(t)
	d := NewWebhookRuleDAO(db)

	// Create a rule
	rule := &model.WebhookRule{
		Name:    "Test Rule",
		RepoKey: "test-repo",
	}

	if err := d.Create(rule); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Create tasks for the rule
	tasks := []model.WebhookRuleTask{
		{RuleID: rule.ID, TaskKey: "task1"},
		{RuleID: rule.ID, TaskKey: "task2"},
	}

	for _, task := range tasks {
		if err := db.Create(&task).Error; err != nil {
			t.Fatalf("create task failed: %v", err)
		}
	}

	// Find tasks by rule ID
	got, err := d.FindTasksByRuleID(rule.ID)
	if err != nil {
		t.Fatalf("find tasks by rule ID failed: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(got))
	}
}

func TestWebhookRuleDAO_Fields(t *testing.T) {
	rule := model.WebhookRule{
		ID:            1,
		Name:          "Test Rule",
		RepoKey:       "test-repo",
		EventType:     "push",
		BranchPattern: "main",
		Action:        "sync",
		MinInterval:   60,
		Enabled:       true,
		Description:   "Test description",
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
