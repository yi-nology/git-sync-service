package service

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
	"gorm.io/gorm"
)

func setupRepoServiceTestDB(t *testing.T) (*gorm.DB, *RepoService) {
	t.Helper()
	// Set up encryption key for tests
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Repo{}, &model.Platform{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	repoDAO, err := dao.NewRepoDAO(db)
	if err != nil {
		t.Fatalf("failed to create RepoDAO: %v", err)
	}

	platformDAO, err := dao.NewPlatformDAO(db)
	if err != nil {
		t.Fatalf("failed to create PlatformDAO: %v", err)
	}

	providerMgr := sdkprov.NewManager(0)
	svc := NewRepoService(repoDAO, platformDAO, providerMgr)

	return db, svc
}

func TestRepoService_List(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)
	ctx := context.Background()

	// Create repos directly in database
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", Platform: "github"},
		{Key: "repo2", Name: "Repo 2", Platform: "gitlab"},
		{Key: "repo3", Name: "Repo 3", Platform: "bitbucket"},
	}

	for _, repo := range repos {
		if err := db.Create(repo).Error; err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List all
	got, total, err := svc.ListRepos(ctx, 0, 50)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if total != 3 {
		t.Fatalf("expected 3 repos, got %d", total)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 repos, got %d", len(got))
	}
}

func TestRepoService_GetNotFound(t *testing.T) {
	_, svc := setupRepoServiceTestDB(t)

	// Try to get a non-existent repo
	got, err := svc.GetRepoByKey("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent repo")
	}
}

func TestRepoService_Count(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)

	// Create repos directly in database
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1"},
		{Key: "repo2", Name: "Repo 2"},
	}

	for _, repo := range repos {
		if err := db.Create(repo).Error; err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Count
	count, err := svc.CountRepos()
	if err != nil {
		t.Fatalf("count failed: %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestRepoService_GetRepoByKey(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)

	// Create repo directly in database
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	if err := db.Create(repo).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Get by key
	got, err := svc.GetRepoByKey("test-repo")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.Key != "test-repo" {
		t.Errorf("expected key 'test-repo', got '%s'", got.Key)
	}

	if got.Name != "Test Repo" {
		t.Errorf("expected name 'Test Repo', got '%s'", got.Name)
	}
}

func TestRepoService_Delete(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)
	ctx := context.Background()

	// Create repo directly in database
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	if err := db.Create(repo).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete the repo
	if err := svc.DeleteRepo(ctx, "test-repo"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Try to get the deleted repo
	got, err := svc.GetRepoByKey("test-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting deleted repo")
	}
}

func TestRepoService_GetRepo(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)
	ctx := context.Background()

	// Create repo directly in database
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	if err := db.Create(repo).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Get repo
	got, err := svc.GetRepo(ctx, "test-repo")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.Key != "test-repo" {
		t.Errorf("expected key 'test-repo', got '%s'", got.Key)
	}

	if got.Name != "Test Repo" {
		t.Errorf("expected name 'Test Repo', got '%s'", got.Name)
	}
}

func TestRepoService_ListReposWithFilter(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)
	ctx := context.Background()

	// Create repos with different platforms
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", Platform: "github"},
		{Key: "repo2", Name: "Repo 2", Platform: "github"},
		{Key: "repo3", Name: "Repo 3", Platform: "gitlab"},
	}

	for _, repo := range repos {
		if err := db.Create(repo).Error; err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List with filter
	filter := &dao.RepoFilter{
		Platform: "github",
	}

	got, total, err := svc.ListReposWithFilter(ctx, 0, 50, filter)
	if err != nil {
		t.Fatalf("list with filter failed: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected 2 repos, got %d", total)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(got))
	}
}
