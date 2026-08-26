package dao

import (
	"testing"

	"github.com/yi-nology/git-sync-service/sync/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupPlatformTestDB(t *testing.T) (*gorm.DB, *PlatformDAO) {
	t.Helper()
	// Set up encryption key for tests
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Platform{}, &model.Repo{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	d, err := NewPlatformDAO(db)
	if err != nil {
		t.Fatalf("failed to create PlatformDAO: %v", err)
	}

	return db, d
}

func TestPlatformDAO_CreateAndFindByKey(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:         "github",
		Name:        "GitHub",
		Type:        "github",
		InstanceURL: "https://github.com",
		APIURL:      "https://api.github.com",
		Status:      "active",
	}

	if err := d.Create(platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if platform.ID == 0 {
		t.Error("expected platform ID to be set after create")
	}

	// Find by key
	got, err := d.FindByKey("github")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Key != "github" {
		t.Errorf("expected key 'github', got '%s'", got.Key)
	}

	if got.Name != "GitHub" {
		t.Errorf("expected name 'GitHub', got '%s'", got.Name)
	}

	if got.Type != "github" {
		t.Errorf("expected type 'github', got '%s'", got.Type)
	}

	if got.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", got.Status)
	}
}

func TestPlatformDAO_FindByID(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:         "gitlab",
		Name:        "GitLab",
		Type:        "gitlab",
		InstanceURL: "https://gitlab.com",
		Status:      "active",
	}

	if err := d.Create(platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Find by ID
	got, err := d.FindByID(platform.ID)
	if err != nil {
		t.Fatalf("find by ID failed: %v", err)
	}

	if got.Key != "gitlab" {
		t.Errorf("expected key 'gitlab', got '%s'", got.Key)
	}

	if got.Name != "GitLab" {
		t.Errorf("expected name 'GitLab', got '%s'", got.Name)
	}
}

func TestPlatformDAO_FindAll(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create multiple platforms
	platforms := []*model.Platform{
		{Key: "github", Name: "GitHub", Type: "github", Status: "active"},
		{Key: "gitlab", Name: "GitLab", Type: "gitlab", Status: "active"},
		{Key: "bitbucket", Name: "Bitbucket", Type: "bitbucket", Status: "inactive"},
	}

	for _, p := range platforms {
		if err := d.Create(p); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find all
	got, err := d.FindAll()
	if err != nil {
		t.Fatalf("find all failed: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 platforms, got %d", len(got))
	}
}

func TestPlatformDAO_Update(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:         "github",
		Name:        "GitHub",
		Type:        "github",
		InstanceURL: "https://github.com",
		Status:      "active",
	}

	if err := d.Create(platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update the platform
	platform.Name = "GitHub Enterprise"
	platform.InstanceURL = "https://github.example.com"

	if err := d.Update(platform); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Get the updated platform
	got, err := d.FindByKey("github")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Name != "GitHub Enterprise" {
		t.Errorf("expected name 'GitHub Enterprise', got '%s'", got.Name)
	}

	if got.InstanceURL != "https://github.example.com" {
		t.Errorf("expected instance URL 'https://github.example.com', got '%s'", got.InstanceURL)
	}
}

func TestPlatformDAO_Delete(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:         "github",
		Name:        "GitHub",
		Type:        "github",
		Status:      "active",
	}

	if err := d.Create(platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete the platform
	if err := d.Delete("github"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Try to get the deleted platform - should return (nil, nil)
	p, err := d.FindByKey("github")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != nil {
		t.Error("expected nil platform for deleted key")
	}
}

func TestPlatformDAO_FindByKey_NotFound(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Try to get a non-existent platform - should return (nil, nil)
	p, err := d.FindByKey("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != nil {
		t.Error("expected nil platform for non-existent key")
	}
}

func TestPlatformDAO_FindByID_NotFound(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Try to get a non-existent platform - should return (nil, nil)
	p, err := d.FindByID(999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != nil {
		t.Error("expected nil platform for non-existent ID")
	}
}

func TestPlatformDAO_SetDefault(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create two platforms
	platform1 := &model.Platform{
		Key:       "github",
		Name:      "GitHub",
		Type:      "github",
		IsDefault: true,
		Status:    "active",
	}

	platform2 := &model.Platform{
		Key:       "gitlab",
		Name:      "GitLab",
		Type:      "gitlab",
		IsDefault: false,
		Status:    "active",
	}

	if err := d.Create(platform1); err != nil {
		t.Fatalf("create platform1 failed: %v", err)
	}

	if err := d.Create(platform2); err != nil {
		t.Fatalf("create platform2 failed: %v", err)
	}

	// Set platform2 as default
	if err := d.SetDefault("gitlab"); err != nil {
		t.Fatalf("set default failed: %v", err)
	}

	// Verify platform2 is now default
	got2, err := d.FindByKey("gitlab")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if !got2.IsDefault {
		t.Error("expected platform2 to be default")
	}

	// Verify platform1 is no longer default
	got1, err := d.FindByKey("github")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got1.IsDefault {
		t.Error("expected platform1 to not be default")
	}
}

func TestPlatformDAO_UpdateStatus(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:    "github",
		Name:   "GitHub",
		Type:   "github",
		Status: "active",
	}

	if err := d.Create(platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update status
	if err := d.UpdateStatus("github", "inactive", "connection failed"); err != nil {
		t.Fatalf("update status failed: %v", err)
	}

	// Verify the status was updated
	got, err := d.FindByKey("github")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Status != "inactive" {
		t.Errorf("expected status 'inactive', got '%s'", got.Status)
	}

	if got.LastTestResult != "connection failed" {
		t.Errorf("expected last test result 'connection failed', got '%s'", got.LastTestResult)
	}
}

func TestPlatformDAO_Count(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create platforms
	platforms := []*model.Platform{
		{Key: "github", Name: "GitHub", Type: "github", Status: "active"},
		{Key: "gitlab", Name: "GitLab", Type: "gitlab", Status: "active"},
	}

	for _, p := range platforms {
		if err := d.Create(p); err != nil {
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

func TestPlatformDAO_FindDefault(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create platforms
	platform1 := &model.Platform{
		Key:       "github",
		Name:      "GitHub",
		Type:      "github",
		IsDefault: false,
		Status:    "active",
	}

	platform2 := &model.Platform{
		Key:       "gitlab",
		Name:      "GitLab",
		Type:      "gitlab",
		IsDefault: true,
		Status:    "active",
	}

	if err := d.Create(platform1); err != nil {
		t.Fatalf("create platform1 failed: %v", err)
	}

	if err := d.Create(platform2); err != nil {
		t.Fatalf("create platform2 failed: %v", err)
	}

	// Find default
	got, err := d.FindDefault()
	if err != nil {
		t.Fatalf("find default failed: %v", err)
	}

	if got.Key != "gitlab" {
		t.Errorf("expected key 'gitlab', got '%s'", got.Key)
	}

	if !got.IsDefault {
		t.Error("expected platform to be default")
	}
}

func TestPlatformDAO_FindByType(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create platforms
	platforms := []*model.Platform{
		{Key: "github1", Name: "GitHub1", Type: "github", Status: "active"},
		{Key: "github2", Name: "GitHub2", Type: "github", Status: "active"},
		{Key: "gitlab", Name: "GitLab", Type: "gitlab", Status: "active"},
	}

	for _, p := range platforms {
		if err := d.Create(p); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// Find by type
	got, err := d.FindByType("github")
	if err != nil {
		t.Fatalf("find by type failed: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(got))
	}
}

func TestPlatformDAO_UpdateFields(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:         "github",
		Name:        "GitHub",
		Type:        "github",
		InstanceURL: "https://github.com",
		Status:      "active",
	}

	if err := d.Create(platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update specific fields
	fields := map[string]interface{}{
		"name":        "GitHub Enterprise",
		"instance_url": "https://github.example.com",
	}

	if err := d.UpdateFields("github", fields); err != nil {
		t.Fatalf("update fields failed: %v", err)
	}

	// Verify the fields were updated
	got, err := d.FindByKey("github")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.Name != "GitHub Enterprise" {
		t.Errorf("expected name 'GitHub Enterprise', got '%s'", got.Name)
	}

	if got.InstanceURL != "https://github.example.com" {
		t.Errorf("expected instance URL 'https://github.example.com', got '%s'", got.InstanceURL)
	}

	// Verify other fields were not changed
	if got.Type != "github" {
		t.Errorf("expected type 'github', got '%s'", got.Type)
	}

	if got.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", got.Status)
	}
}

func TestPlatformDAO_UpdateRepoCount(t *testing.T) {
	db, d := setupPlatformTestDB(t)

	// Create a platform
	platform := &model.Platform{
		Key:       "github",
		Name:      "GitHub",
		Type:      "github",
		RepoCount: 0,
	}

	if err := d.Create(platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Create repos for the platform
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", PlatformID: platform.ID},
		{Key: "repo2", Name: "Repo 2", PlatformID: platform.ID},
		{Key: "repo3", Name: "Repo 3", PlatformID: platform.ID},
	}

	for _, repo := range repos {
		if err := db.Create(repo).Error; err != nil {
			t.Fatalf("create repo failed: %v", err)
		}
	}

	// Update repo count
	if err := d.UpdateRepoCount(platform.ID); err != nil {
		t.Fatalf("update repo count failed: %v", err)
	}

	// Verify the repo count was updated
	got, err := d.FindByKey("github")
	if err != nil {
		t.Fatalf("find by key failed: %v", err)
	}

	if got.RepoCount != 3 {
		t.Errorf("expected repo count 3, got %d", got.RepoCount)
	}
}
