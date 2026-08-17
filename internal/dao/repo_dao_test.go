package dao

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/yi-nology/git-platform-sdk/pkg/credential"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func TestCryptoManager_EncryptDecrypt_Roundtrip(t *testing.T) {
	// Set up a test encryption key
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	cm, err := credential.NewCryptoManager()
	if err != nil {
		t.Fatalf("failed to create CryptoManager: %v", err)
	}

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
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			if encrypted == tc.plaintext {
				t.Errorf("encrypted text should differ from plaintext")
			}

			decrypted, err := cm.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Errorf("decrypt roundtrip failed: got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestCryptoManager_EmptyString(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	t.Setenv("ENCRYPTION_KEY", key)

	cm, err := credential.NewCryptoManager()
	if err != nil {
		t.Fatalf("failed to create CryptoManager: %v", err)
	}

	// Test encrypting empty string
	encrypted, err := cm.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt empty string failed: %v", err)
	}

	// Empty string should still produce some output (or handle gracefully)
	decrypted, err := cm.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt empty string failed: %v", err)
	}

	if decrypted != "" {
		t.Errorf("expected empty string, got %q", decrypted)
	}
}

func TestCryptoManager_WithoutKey(t *testing.T) {
	// Ensure ENCRYPTION_KEY is not set
	_ = os.Unsetenv("ENCRYPTION_KEY")

	_, err := credential.NewCryptoManager()
	if err == nil {
		t.Error("expected error when ENCRYPTION_KEY is not set, got nil")
	}
}

func TestCryptoManager_NewCryptoManagerFromKey(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	cm := credential.NewCryptoManagerFromKey(key)

	plaintext := "test-token-value"
	encrypted, err := cm.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := cm.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypt roundtrip failed: got %q, want %q", decrypted, plaintext)
	}
}

func TestDefaultPagination(t *testing.T) {
	tests := []struct {
		name           string
		offset, limit  int
		wantOff, wantLim int
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
			if p.Offset != tc.wantOff {
				t.Errorf("Offset = %d, want %d", p.Offset, tc.wantOff)
			}
			if p.Limit != tc.wantLim {
				t.Errorf("Limit = %d, want %d", p.Limit, tc.wantLim)
			}
		})
	}
}

func setupRepoTestDB(t *testing.T) (*gorm.DB, *RepoDAO) {
	t.Helper()
	// Set up encryption key for tests
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Repo{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	d, err := NewRepoDAO(db)
	if err != nil {
		t.Fatalf("failed to create RepoDAO: %v", err)
	}

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

	if err := d.Create(repo); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if repo.ID == 0 {
		t.Error("expected repo ID to be set after create")
	}

	// Find by key
	got, err := d.FindByKey("test-repo")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Key != "test-repo" {
		t.Errorf("expected key 'test-repo', got '%s'", got.Key)
	}

	if got.Name != "Test Repo" {
		t.Errorf("expected name 'Test Repo', got '%s'", got.Name)
	}

	if got.Platform != "github" {
		t.Errorf("expected platform 'github', got '%s'", got.Platform)
	}

	if got.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", got.Status)
	}
}

func TestRepoDAO_FindByCloneURL(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create a repo
	repo := &model.Repo{
		Key:      "test-repo",
		Name:     "Test Repo",
		CloneURL: "https://github.com/testuser/test-repo.git",
	}

	if err := d.Create(repo); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Find by clone URL
	got, err := d.FindByCloneURL("https://github.com/testuser/test-repo.git")
	if err != nil {
		t.Fatalf("find by clone URL failed: %v", err)
	}

	if got.Key != "test-repo" {
		t.Errorf("expected key 'test-repo', got '%s'", got.Key)
	}

	if got.Name != "Test Repo" {
		t.Errorf("expected name 'Test Repo', got '%s'", got.Name)
	}
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
		if err := d.Create(repo); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find all
	got, total, err := d.FindAll(DefaultPagination(0, 50))
	if err != nil {
		t.Fatalf("find all failed: %v", err)
	}

	if total != 3 {
		t.Fatalf("expected 3 repos, got %d", total)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 repos, got %d", len(got))
	}
}

func TestRepoDAO_Update(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create a repo
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	if err := d.Create(repo); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update the repo
	repo.Name = "Updated Repo"
	repo.Status = "inactive"

	if err := d.Update(repo); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Get the updated repo
	got, err := d.FindByKey("test-repo")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Name != "Updated Repo" {
		t.Errorf("expected name 'Updated Repo', got '%s'", got.Name)
	}

	if got.Status != "inactive" {
		t.Errorf("expected status 'inactive', got '%s'", got.Status)
	}
}

func TestRepoDAO_Delete(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create a repo
	repo := &model.Repo{
		Key:  "test-repo",
		Name: "Test Repo",
	}

	if err := d.Create(repo); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete the repo
	if err := d.Delete("test-repo"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Try to get the deleted repo
	got, err := d.FindByKey("test-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting deleted repo")
	}
}

func TestRepoDAO_FindByKey_NotFound(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Try to get a non-existent repo
	got, err := d.FindByKey("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent repo")
	}
}

func TestRepoDAO_FindByCloneURL_NotFound(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Try to get a non-existent repo
	got, err := d.FindByCloneURL("https://github.com/nonexistent/repo.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Error("expected nil when getting non-existent repo")
	}
}

func TestRepoDAO_Count(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create repos
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1"},
		{Key: "repo2", Name: "Repo 2"},
	}

	for _, repo := range repos {
		if err := d.Create(repo); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Count
	count, err := d.Count()
	if err != nil {
		t.Fatalf("count failed: %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
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

	if repo.ID != 1 {
		t.Errorf("expected ID 1, got %d", repo.ID)
	}

	if repo.Key != "test-repo" {
		t.Errorf("expected key 'test-repo', got '%s'", repo.Key)
	}

	if repo.Name != "Test Repo" {
		t.Errorf("expected name 'Test Repo', got '%s'", repo.Name)
	}

	if repo.Platform != "github" {
		t.Errorf("expected platform 'github', got '%s'", repo.Platform)
	}

	if repo.PlatformOwner != "testuser" {
		t.Errorf("expected platform owner 'testuser', got '%s'", repo.PlatformOwner)
	}

	if repo.PlatformRepo != "test-repo" {
		t.Errorf("expected platform repo 'test-repo', got '%s'", repo.PlatformRepo)
	}

	if repo.CloneURL != "https://github.com/testuser/test-repo.git" {
		t.Errorf("expected clone URL 'https://github.com/testuser/test-repo.git', got '%s'", repo.CloneURL)
	}

	if repo.SSHURL != "git@github.com:testuser/test-repo.git" {
		t.Errorf("expected SSH URL 'git@github.com:testuser/test-repo.git', got '%s'", repo.SSHURL)
	}

	if repo.DefaultBranch != "main" {
		t.Errorf("expected default branch 'main', got '%s'", repo.DefaultBranch)
	}

	if repo.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", repo.Status)
	}
}

func TestRepoDAO_FindByPlatformID(t *testing.T) {
	db, d := setupRepoTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:  "github",
		Name: "GitHub",
	}
	if err := db.Create(platform).Error; err != nil {
		t.Fatalf("create platform failed: %v", err)
	}

	// Create repos with different platform IDs
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", PlatformID: platform.ID},
		{Key: "repo2", Name: "Repo 2", PlatformID: platform.ID},
		{Key: "repo3", Name: "Repo 3"}, // No platform
	}

	for _, repo := range repos {
		if err := d.Create(repo); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find by platform ID
	got, err := d.FindByPlatformID(platform.ID)
	if err != nil {
		t.Fatalf("find by platform ID failed: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(got))
	}
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
		if err := d.Create(repo); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List with filter
	filter := &RepoFilter{
		Platform: "github",
	}

	got, total, err := d.ListWithFilter(DefaultPagination(0, 50), filter)
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

func TestRepoDAO_ListWithFilter_Status(t *testing.T) {
	_, d := setupRepoTestDB(t)

	// Create repos with different statuses
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", Platform: "github", Status: "active"},
		{Key: "repo2", Name: "Repo 2", Platform: "github", Status: "inactive"},
		{Key: "repo3", Name: "Repo 3", Platform: "gitlab", Status: "active"},
	}

	for _, repo := range repos {
		if err := d.Create(repo); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List with status filter
	filter := &RepoFilter{
		Status: "active",
	}

	got, total, err := d.ListWithFilter(DefaultPagination(0, 50), filter)
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
