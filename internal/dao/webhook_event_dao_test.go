package dao

import (
	"testing"
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupWebhookEventTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.WebhookEvent{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return db
}

func TestWebhookEventDAO_CreateAndFindByID(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Create an event
	event := &model.WebhookEvent{
		EventID:   "evt-123",
		RepoKey:   "test-repo",
		EventType: "push",
		Source:    "github",
		ActorName: "testuser",
		Branch:    "main",
		CommitSHA: "abc123",
		Status:    "pending",
	}

	if err := d.Create(event); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if event.ID == 0 {
		t.Error("expected event ID to be set after create")
	}

	// Find by ID
	got, err := d.FindByID(event.ID)
	if err != nil {
		t.Fatalf("find by ID failed: %v", err)
	}

	if got.EventID != "evt-123" {
		t.Errorf("expected event ID 'evt-123', got '%s'", got.EventID)
	}

	if got.RepoKey != "test-repo" {
		t.Errorf("expected repo key 'test-repo', got '%s'", got.RepoKey)
	}

	if got.EventType != "push" {
		t.Errorf("expected event type 'push', got '%s'", got.EventType)
	}

	if got.Source != "github" {
		t.Errorf("expected source 'github', got '%s'", got.Source)
	}

	if got.ActorName != "testuser" {
		t.Errorf("expected actor name 'testuser', got '%s'", got.ActorName)
	}

	if got.Branch != "main" {
		t.Errorf("expected branch 'main', got '%s'", got.Branch)
	}

	if got.CommitSHA != "abc123" {
		t.Errorf("expected commit SHA 'abc123', got '%s'", got.CommitSHA)
	}

	if got.Status != "pending" {
		t.Errorf("expected status 'pending', got '%s'", got.Status)
	}
}

func TestWebhookEventDAO_FindByRepoKey(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Create events with different repo keys
	events := []*model.WebhookEvent{
		{EventID: "evt-1", RepoKey: "repo1", EventType: "push", Status: "pending"},
		{EventID: "evt-2", RepoKey: "repo1", EventType: "push", Status: "processed"},
		{EventID: "evt-3", RepoKey: "repo2", EventType: "tag", Status: "failed"},
	}

	for _, event := range events {
		if err := d.Create(event); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find by repo key
	got, total, err := d.FindByRepoKey("repo1", DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find by repo key failed: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected 2 events, got %d", total)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
}

func TestWebhookEventDAO_FindByEventID(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Create an event
	event := &model.WebhookEvent{
		EventID:   "evt-123",
		RepoKey:   "test-repo",
		EventType: "push",
		Status:    "pending",
	}

	if err := d.Create(event); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Find by event ID
	got, err := d.FindByEventID("evt-123")
	if err != nil {
		t.Fatalf("find by event ID failed: %v", err)
	}

	if got.EventID != "evt-123" {
		t.Errorf("expected event ID 'evt-123', got '%s'", got.EventID)
	}

	if got.RepoKey != "test-repo" {
		t.Errorf("expected repo key 'test-repo', got '%s'", got.RepoKey)
	}
}

func TestWebhookEventDAO_Update(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Create an event
	event := &model.WebhookEvent{
		EventID:   "evt-123",
		RepoKey:   "test-repo",
		EventType: "push",
		Status:    "pending",
	}

	if err := d.Create(event); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update the event
	now := time.Now()
	event.Status = "processed"
	event.ProcessedAt = &now

	if err := d.Update(event); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Find by ID to verify update
	got, err := d.FindByID(event.ID)
	if err != nil {
		t.Fatalf("find by ID failed: %v", err)
	}

	if got.Status != "processed" {
		t.Errorf("expected status 'processed', got '%s'", got.Status)
	}

	if got.ProcessedAt == nil {
		t.Error("expected processed at to be set")
	}
}

func TestWebhookEventDAO_FindByID_NotFound(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Try to get a non-existent event
	got, err := d.FindByID(999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent event")
	}
}

func TestWebhookEventDAO_FindRecent(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Create events
	events := []*model.WebhookEvent{
		{EventID: "evt-1", RepoKey: "repo1", Status: "pending"},
		{EventID: "evt-2", RepoKey: "repo1", Status: "processed"},
		{EventID: "evt-3", RepoKey: "repo2", Status: "failed"},
	}

	for _, event := range events {
		if err := d.Create(event); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find recent for repo1
	got, err := d.FindRecent("repo1", DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find recent failed: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
}

func TestWebhookEventDAO_CountByRepoKey(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Create events with different repo keys
	events := []*model.WebhookEvent{
		{EventID: "evt-1", RepoKey: "repo1", Status: "pending"},
		{EventID: "evt-2", RepoKey: "repo1", Status: "processed"},
		{EventID: "evt-3", RepoKey: "repo2", Status: "failed"},
	}

	for _, event := range events {
		if err := d.Create(event); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Count by repo key
	count, err := d.CountByRepoKey("repo1")
	if err != nil {
		t.Fatalf("count by repo key failed: %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestWebhookEventDAO_CleanupOlderThan(t *testing.T) {
	db := setupWebhookEventTestDB(t)
	d := NewWebhookEventDAO(db)

	// Create events with different ages
	oldEvent := &model.WebhookEvent{
		EventID:   "evt-1",
		RepoKey:   "repo1",
		Status:    "processed",
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}

	newEvent := &model.WebhookEvent{
		EventID:   "evt-2",
		RepoKey:   "repo1",
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := db.Create(oldEvent).Error; err != nil {
		t.Fatalf("create old event failed: %v", err)
	}

	if err := db.Create(newEvent).Error; err != nil {
		t.Fatalf("create new event failed: %v", err)
	}

	// Cleanup older than 24 hours
	count, err := d.CleanupOlderThan(24 * time.Hour)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Verify only new event remains
	got, err := d.FindRecent("repo1", DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find recent failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}

	if got[0].EventID != "evt-2" {
		t.Errorf("expected event ID 'evt-2', got '%s'", got[0].EventID)
	}
}

func TestWebhookEventDAO_Fields(t *testing.T) {
	now := time.Now()
	event := model.WebhookEvent{
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

	if event.ID != 1 {
		t.Errorf("expected ID 1, got %d", event.ID)
	}

	if event.EventID != "evt-123" {
		t.Errorf("expected event ID 'evt-123', got '%s'", event.EventID)
	}

	if event.RepoKey != "test-repo" {
		t.Errorf("expected repo key 'test-repo', got '%s'", event.RepoKey)
	}

	if event.EventType != "push" {
		t.Errorf("expected event type 'push', got '%s'", event.EventType)
	}

	if event.Source != "github" {
		t.Errorf("expected source 'github', got '%s'", event.Source)
	}

	if event.ActorName != "testuser" {
		t.Errorf("expected actor name 'testuser', got '%s'", event.ActorName)
	}

	if event.Branch != "main" {
		t.Errorf("expected branch 'main', got '%s'", event.Branch)
	}

	if event.CommitSHA != "abc123" {
		t.Errorf("expected commit SHA 'abc123', got '%s'", event.CommitSHA)
	}

	if event.Status != "processed" {
		t.Errorf("expected status 'processed', got '%s'", event.Status)
	}

	if event.ErrorMessage != "" {
		t.Errorf("expected empty error message, got '%s'", event.ErrorMessage)
	}

	if event.ProcessedAt == nil {
		t.Error("expected processed at to be set")
	}
}
