package dao

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yi-nology/git-platform-sdk/pkg/credential"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func TestCryptoManager_EncryptDecrypt_Roundtrip(t *testing.T) {
	// Set up a test encryption key
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	cm, err := credential.NewCryptoManager()
	require.NoError(t, err, "failed to create CryptoManager")

	testCases := []struct {
		name      string
		plaintext string
	}{
		{"normal text", "my-secret-token-123"},
		{"with special chars", "token!@#$%^&*()_+-=[]{}|;':\",./<>?`~"},
		{"unicode", "token-with-unicode-\u4e2d\u6587"},
		{"long string", "a]very-long-token-that-is-quite-substantial-in-length-and-contains-many-characters-for-testing-purposes"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := cm.Encrypt(tc.plaintext)
			require.NoError(t, err, "Encrypt failed")

			assert.NotEqual(t, tc.plaintext, encrypted, "encrypted text should differ from plaintext")

			decrypted, err := cm.Decrypt(encrypted)
			require.NoError(t, err, "Decrypt failed")

			assert.Equal(t, tc.plaintext, decrypted, "decrypt roundtrip failed")
		})
	}
}

func TestCryptoManager_EmptyString(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	t.Setenv("ENCRYPTION_KEY", key)

	cm, err := credential.NewCryptoManager()
	require.NoError(t, err, "failed to create CryptoManager")

	// Test encrypting empty string
	encrypted, err := cm.Encrypt("")
	require.NoError(t, err, "Encrypt empty string failed")

	// Empty string should still produce some output (or handle gracefully)
	decrypted, err := cm.Decrypt(encrypted)
	require.NoError(t, err, "Decrypt empty string failed")

	assert.Equal(t, "", decrypted, "expected empty string")
}

func TestCryptoManager_WithoutKey(t *testing.T) {
	// Ensure ENCRYPTION_KEY is not set
	_ = os.Unsetenv("ENCRYPTION_KEY")

	_, err := credential.NewCryptoManager()
	assert.Error(t, err, "expected error when ENCRYPTION_KEY is not set")
}

func TestCryptoManager_NewCryptoManagerFromKey(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	cm := credential.NewCryptoManagerFromKey(key)

	plaintext := "test-token-value"
	encrypted, err := cm.Encrypt(plaintext)
	require.NoError(t, err, "Encrypt failed")

	decrypted, err := cm.Decrypt(encrypted)
	require.NoError(t, err, "Decrypt failed")

	assert.Equal(t, plaintext, decrypted, "decrypt roundtrip failed")
}

func TestDefaultPagination(t *testing.T) {
	tests := []struct {
		name                        string
		offset, limit               int
		wantOff, wantLim            int
	}{
		{"normal values", 10, 50, 10, 50},
		{"zero limit defaults to 50", 0, 0, 0, 50},
		{"negative limit defaults to 50", 5, -1, 5, 50},
		{"limit over 200 capped", 0, 300, 0, 200},
		{"negative offset becomes 0", -5, 50, 0, 50},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := DefaultPagination(tc.offset, tc.limit)
			assert.Equal(t, tc.wantOff, p.Offset, "Offset mismatch")
			assert.Equal(t, tc.wantLim, p.Limit, "Limit mismatch")
		})
	}
}

func setupRepoTestDB(t *testing.T) (*gorm.DB, *RepoDAO) {
	t.Helper()
	// Set up encryption key for tests
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open test db")

	err = db.AutoMigrate(&model.Repo{}, &model.Platform{})
	require.NoError(t, err, "failed to migrate test db")

	d, err := NewRepoDAO(db)
	require.NoError(t, err, "failed to create RepoDAO")

	return db, d
}

func TestRepoDAO_CreateAndFindByKey(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create a repo
	repo := &model.Repo{
		Key:           "test-repo",
		Name:          "Test Repo",
		Platform:      "github",
		PlatformOwner: "testuser",
		PlatformRepo:  "test-repo",
		CloneURL:      "https://github.com/testuser/test-repo.git",
		DefaultBranch: "main",
		Status:        "active",
	}

	err := d.Create(repo)
	require.NoError(t, err, "create failed")

	assert.NotZero(t, repo.ID, "expected repo ID to be set after create")

	// Find by key
	got, err := d.FindByKey("test-repo")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "test-repo", got.Key, "expected key 'test-repo'")
	assert.Equal(t, "Test Repo", got.Name, "expected name 'Test Repo'")
	assert.Equal(t, "github", got.Platform, "expected platform 'github'")
	assert.Equal(t, "active", got.Status, "expected status 'active'")
}

func TestRepoDAO_FindByCloneURL(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create a repo
	repo := &model.Repo{
		Key:      "test-repo",
		Name:     "Test Repo",
		CloneURL: "https://github.com/testuser/test-repo.git",
	}

	err := d.Create(repo)
	require.NoError(t, err, "create failed")

	// Find by clone URL
	got, err := d.FindByCloneURL("https://github.com/testuser/test-repo.git")
	require.NoError(t, err, "find by clone URL failed")

	assert.Equal(t, "test-repo", got.Key, "expected key 'test-repo'")
	assert.Equal(t, "Test Repo", got.Name, "expected name 'Test Repo'")
}

func TestRepoDAO_FindAll(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create multiple repos
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", Platform: "github"},
		{Key: "repo2", Name: "Repo 2", Platform: "gitlab"},
		{Key: "repo3", Name: "Repo 3", Platform: "bitbucket"},
	}

	for _, repo := range repos {
		err := d.Create(repo)
		require.NoError(t, err, "create failed")
	}

	// Find all
	got, total, err := d.FindAll(DefaultPagination(0, 50))
	require.NoError(t, err, "find all failed")

	require.Equal(t, int64(3), total, "expected 3 repos")
	require.Len(t, got, 3, "expected 3 repos")
}

func TestRepoDAO_Update(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create a repo
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	err := d.Create(repo)
	require.NoError(t, err, "create failed")

	// Update the repo
	repo.Name = "Updated Repo"
	repo.Status = "inactive"

	err = d.Update(repo)
	require.NoError(t, err, "update failed")

	// Get the updated repo
	got, err := d.FindByKey("test-repo")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "Updated Repo", got.Name, "expected name 'Updated Repo'")
	assert.Equal(t, "inactive", got.Status, "expected status 'inactive'")
}

func TestRepoDAO_Delete(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create a repo
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	err := d.Create(repo)
	require.NoError(t, err, "create failed")

	// Delete the repo
	err = d.Delete("test-repo")
	require.NoError(t, err, "delete failed")

	// Try to get the deleted repo
	got, err := d.FindByKey("test-repo")
	require.NoError(t, err, "unexpected error")

	assert.Nil(t, got, "expected nil when getting deleted repo")
}

func TestRepoDAO_FindByKey_NotFound(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Try to get a non-existent repo
	got, err := d.FindByKey("nonexistent")
	require.NoError(t, err, "unexpected error")

	assert.Nil(t, got, "expected nil when getting non-existent repo")
}

func TestRepoDAO_FindByCloneURL_NotFound(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Try to get a non-existent repo
	got, err := d.FindByCloneURL("https://github.com/nonexistent/repo.git")
	require.NoError(t, err, "unexpected error")

	assert.Nil(t, got, "expected nil when getting non-existent repo")
}

func TestRepoDAO_Count(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create repos
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1"},
		{Key: "repo2", Name: "Repo 2"},
	}

	for _, repo := range repos {
		err := d.Create(repo)
		require.NoError(t, err, "create failed")
	}

	// Count
	count, err := d.Count()
	require.NoError(t, err, "count failed")

	assert.Equal(t, int64(2), count, "expected count 2")
}

func TestRepoDAO_Fields(t *testing.T) {
	repo := model.Repo{
		ID:            1,
		Key:           "test-repo",
		Name:          "Test Repo",
		Platform:      "github",
		PlatformOwner: "testuser",
		PlatformRepo:  "test-repo",
		CloneURL:      "https://github.com/testuser/test-repo.git",
		SSHURL:        "git@github.com:testuser/test-repo.git",
		DefaultBranch: "main",
		Status:        "active",
	}

	assert.Equal(t, uint(1), repo.ID, "expected ID 1")
	assert.Equal(t, "test-repo", repo.Key, "expected key 'test-repo'")
	assert.Equal(t, "Test Repo", repo.Name, "expected name 'Test Repo'")
	assert.Equal(t, "github", repo.Platform, "expected platform 'github'")
	assert.Equal(t, "testuser", repo.PlatformOwner, "expected platform owner 'testuser'")
	assert.Equal(t, "test-repo", repo.PlatformRepo, "expected platform repo 'test-repo'")
	assert.Equal(t, "https://github.com/testuser/test-repo.git", repo.CloneURL, "expected clone URL")
	assert.Equal(t, "git@github.com:testuser/test-repo.git", repo.SSHURL, "expected SSH URL")
	assert.Equal(t, "main", repo.DefaultBranch, "expected default branch 'main'")
	assert.Equal(t, "active", repo.Status, "expected status 'active'")
}

func TestRepoDAO_FindByPlatformID(t *testing.T) {
	db, d := setupRepoTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:  "github",
		Name: "GitHub",
	}
	err := db.Create(platform).Error
	require.NoError(t, err, "create platform failed")

	// Create repos with different platform IDs
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", PlatformID: platform.ID},
		{Key: "repo2", Name: "Repo 2", PlatformID: platform.ID},
		{Key: "repo3", Name: "Repo 3"}, // No platform
	}

	for _, repo := range repos {
		err := d.Create(repo)
		require.NoError(t, err, "create failed")
	}

	// Find by platform ID
	got, err := d.FindByPlatformID(platform.ID)
	require.NoError(t, err, "find by platform ID failed")

	require.Len(t, got, 2, "expected 2 repos")
}

func TestRepoDAO_ListWithFilter(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create repos with different platforms
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", Platform: "github", Status: "active"},
		{Key: "repo2", Name: "Repo 2", Platform: "github", Status: "active"},
		{Key: "repo3", Name: "Repo 3", Platform: "gitlab", Status: "inactive"},
	}

	for _, repo := range repos {
		err := d.Create(repo)
		require.NoError(t, err, "create failed")
	}

	// List with filter
	filter := &RepoFilter{
		Platform: "github",
	}

	got, total, err := d.ListWithFilter(DefaultPagination(0, 50), filter)
	require.NoError(t, err, "list with filter failed")

	require.Equal(t, int64(2), total, "expected 2 repos")
	require.Len(t, got, 2, "expected 2 repos")
}

func TestRepoDAO_ListWithFilter_Status(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create repos with different statuses
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", Platform: "github", Status: "active"},
		{Key: "repo2", Name: "Repo 2", Platform: "github", Status: "inactive"},
		{Key: "repo3", Name: "Repo 3", Platform: "gitlab", Status: "active"},
	}

	for _, repo := range repos {
		err := d.Create(repo)
		require.NoError(t, err, "create failed")
	}

	// List with status filter
	filter := &RepoFilter{
		Status: "active",
	}

	got, total, err := d.ListWithFilter(DefaultPagination(0, 50), filter)
	require.NoError(t, err, "list with filter failed")

	require.Equal(t, int64(2), total, "expected 2 repos")
	require.Len(t, got, 2, "expected 2 repos")
}
