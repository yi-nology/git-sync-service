package service

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err, "failed to open test db")

	err = db.AutoMigrate(&model.Repo{}, &model.Platform{})
	require.NoError(t, err, "failed to migrate test db")

	repoDAO, err := dao.NewRepoDAO(db)
	require.NoError(t, err, "failed to create RepoDAO")

	platformDAO, err := dao.NewPlatformDAO(db)
	require.NoError(t, err, "failed to create PlatformDAO")

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
		err := db.Create(repo).Error
		require.NoError(t, err, "create failed")
	}

	// List all
	got, total, err := svc.ListRepos(ctx, 0, 50)
	require.NoError(t, err, "list failed")

	require.Equal(t, int64(3), total, "expected 3 repos")
	require.Len(t, got, 3, "expected 3 repos")
}

func TestRepoService_GetNotFound(t *testing.T) {
	_, svc := setupRepoServiceTestDB(t)

	// Try to get a non-existent repo
	got, err := svc.GetRepoByKey("nonexistent")
	require.NoError(t, err, "unexpected error")

	assert.Nil(t, got, "expected nil when getting non-existent repo")
}

func TestRepoService_Count(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)

	// Create repos directly in database
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1"},
		{Key: "repo2", Name: "Repo 2"},
	}

	for _, repo := range repos {
		err := db.Create(repo).Error
		require.NoError(t, err, "create failed")
	}

	// Count
	count, err := svc.CountRepos()
	require.NoError(t, err, "count failed")

	assert.Equal(t, int64(2), count, "expected count 2")
}

func TestRepoService_GetRepoByKey(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)

	// Create repo directly in database
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	err := db.Create(repo).Error
	require.NoError(t, err, "create failed")

	// Get by key
	got, err := svc.GetRepoByKey("test-repo")
	require.NoError(t, err, "get failed")

	assert.Equal(t, "test-repo", got.Key, "expected key 'test-repo'")
	assert.Equal(t, "Test Repo", got.Name, "expected name 'Test Repo'")
}

func TestRepoService_Delete(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)
	ctx := context.Background()

	// Create repo directly in database
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	err := db.Create(repo).Error
	require.NoError(t, err, "create failed")

	// Delete the repo
	err = svc.DeleteRepo(ctx, "test-repo")
	require.NoError(t, err, "delete failed")

	// Try to get the deleted repo
	got, err := svc.GetRepoByKey("test-repo")
	require.NoError(t, err, "unexpected error")

	assert.Nil(t, got, "expected nil when getting deleted repo")
}

func TestRepoService_GetRepo(t *testing.T) {
	db, svc := setupRepoServiceTestDB(t)
	ctx := context.Background()

	// Create repo directly in database
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	err := db.Create(repo).Error
	require.NoError(t, err, "create failed")

	// Get repo
	got, err := svc.GetRepo(ctx, "test-repo")
	require.NoError(t, err, "get failed")

	assert.Equal(t, "test-repo", got.Key, "expected key 'test-repo'")
	assert.Equal(t, "Test Repo", got.Name, "expected name 'Test Repo'")
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
		err := db.Create(repo).Error
		require.NoError(t, err, "create failed")
	}

	// List with filter
	filter := &dao.RepoFilter{
		Platform: "github",
	}

	got, total, err := svc.ListReposWithFilter(ctx, 0, 50, filter)
	require.NoError(t, err, "list with filter failed")

	require.Equal(t, int64(2), total, "expected 2 repos")
	require.Len(t, got, 2, "expected 2 repos")
}
