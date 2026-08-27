package dao

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func setupPlatformTestDB(t *testing.T) (*gorm.DB, *PlatformDAO) {
	t.Helper()
	// Set up encryption key for tests
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	t.Setenv("ENCRYPTION_KEY", key)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open test db")

	err = db.AutoMigrate(&model.Platform{}, &model.Repo{})
	require.NoError(t, err, "failed to migrate test db")

	d, err := NewPlatformDAO(db)
	require.NoError(t, err, "failed to create PlatformDAO")

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

	err := d.Create(platform)
	require.NoError(t, err, "create failed")

	assert.NotZero(t, platform.ID, "expected platform ID to be set after create")

	// Find by key
	got, err := d.FindByKey("github")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "github", got.Key, "expected key 'github'")
	assert.Equal(t, "GitHub", got.Name, "expected name 'GitHub'")
	assert.Equal(t, "github", got.Type, "expected type 'github'")
	assert.Equal(t, "active", got.Status, "expected status 'active'")
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

	err := d.Create(platform)
	require.NoError(t, err, "create failed")

	// Find by ID
	got, err := d.FindByID(platform.ID)
	require.NoError(t, err, "find by ID failed")

	assert.Equal(t, "gitlab", got.Key, "expected key 'gitlab'")
	assert.Equal(t, "GitLab", got.Name, "expected name 'GitLab'")
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
		err := d.Create(p)
		require.NoError(t, err, "create failed")
	}

	// Find all
	got, err := d.FindAll()
	require.NoError(t, err, "find all failed")

	require.Len(t, got, 3, "expected 3 platforms")
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

	err := d.Create(platform)
	require.NoError(t, err, "create failed")

	// Update the platform
	platform.Name = "GitHub Enterprise"
	platform.InstanceURL = "https://github.example.com"

	err = d.Update(platform)
	require.NoError(t, err, "update failed")

	// Get the updated platform
	got, err := d.FindByKey("github")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "GitHub Enterprise", got.Name, "expected name 'GitHub Enterprise'")
	assert.Equal(t, "https://github.example.com", got.InstanceURL, "expected instance URL")
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

	err := d.Create(platform)
	require.NoError(t, err, "create failed")

	// Delete the platform
	err = d.Delete("github")
	require.NoError(t, err, "delete failed")

	// Try to get the deleted platform - should return (nil, nil)
	p, err := d.FindByKey("github")
	require.NoError(t, err, "unexpected error")
	assert.Nil(t, p, "expected nil platform for deleted key")
}

func TestPlatformDAO_FindByKey_NotFound(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Try to get a non-existent platform - should return (nil, nil)
	p, err := d.FindByKey("nonexistent")
	require.NoError(t, err, "unexpected error")
	assert.Nil(t, p, "expected nil platform for non-existent key")
}

func TestPlatformDAO_FindByID_NotFound(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Try to get a non-existent platform - should return (nil, nil)
	p, err := d.FindByID(999)
	require.NoError(t, err, "unexpected error")
	assert.Nil(t, p, "expected nil platform for non-existent ID")
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

	err := d.Create(platform1)
	require.NoError(t, err, "create platform1 failed")

	err = d.Create(platform2)
	require.NoError(t, err, "create platform2 failed")

	// Set platform2 as default
	err = d.SetDefault("gitlab")
	require.NoError(t, err, "set default failed")

	// Verify platform2 is now default
	got2, err := d.FindByKey("gitlab")
	require.NoError(t, err, "find by key failed")

	assert.True(t, got2.IsDefault, "expected platform2 to be default")

	// Verify platform1 is no longer default
	got1, err := d.FindByKey("github")
	require.NoError(t, err, "find by key failed")

	assert.False(t, got1.IsDefault, "expected platform1 to not be default")
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

	err := d.Create(platform)
	require.NoError(t, err, "create failed")

	// Update status
	err = d.UpdateStatus("github", "inactive", "connection failed")
	require.NoError(t, err, "update status failed")

	// Verify the status was updated
	got, err := d.FindByKey("github")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "inactive", got.Status, "expected status 'inactive'")
	assert.Equal(t, "connection failed", got.LastTestResult, "expected last test result 'connection failed'")
}

func TestPlatformDAO_Count(t *testing.T) {
	_, d := setupPlatformTestDB(t)

	// Create platforms
	platforms := []*model.Platform{
		{Key: "github", Name: "GitHub", Type: "github", Status: "active"},
		{Key: "gitlab", Name: "GitLab", Type: "gitlab", Status: "active"},
	}

	for _, p := range platforms {
		err := d.Create(p)
		require.NoError(t, err, "create failed")
	}

	// Count
	count, err := d.Count()
	require.NoError(t, err, "count failed")

	assert.Equal(t, int64(2), count, "expected count 2")
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

	err := d.Create(platform1)
	require.NoError(t, err, "create platform1 failed")

	err = d.Create(platform2)
	require.NoError(t, err, "create platform2 failed")

	// Find default
	got, err := d.FindDefault()
	require.NoError(t, err, "find default failed")

	assert.Equal(t, "gitlab", got.Key, "expected key 'gitlab'")
	assert.True(t, got.IsDefault, "expected platform to be default")
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
		err := d.Create(p)
		require.NoError(t, err, "create failed")
	}

	// Find by type
	got, err := d.FindByType("github")
	require.NoError(t, err, "find by type failed")

	require.Len(t, got, 2, "expected 2 platforms")
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

	err := d.Create(platform)
	require.NoError(t, err, "create failed")

	// Update specific fields
	fields := map[string]interface{}{
		"name":        "GitHub Enterprise",
		"instance_url": "https://github.example.com",
	}

	err = d.UpdateFields("github", fields)
	require.NoError(t, err, "update fields failed")

	// Verify the fields were updated
	got, err := d.FindByKey("github")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, "GitHub Enterprise", got.Name, "expected name 'GitHub Enterprise'")
	assert.Equal(t, "https://github.example.com", got.InstanceURL, "expected instance URL")

	// Verify other fields were not changed
	assert.Equal(t, "github", got.Type, "expected type 'github'")
	assert.Equal(t, "active", got.Status, "expected status 'active'")
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

	err := d.Create(platform)
	require.NoError(t, err, "create failed")

	// Create repos for the platform
	repos := []*model.Repo{
		{Key: "repo1", Name: "Repo 1", PlatformID: platform.ID},
		{Key: "repo2", Name: "Repo 2", PlatformID: platform.ID},
		{Key: "repo3", Name: "Repo 3", PlatformID: platform.ID},
	}

	for _, repo := range repos {
		err := db.Create(repo).Error
		require.NoError(t, err, "create repo failed")
	}

	// Update repo count
	err = d.UpdateRepoCount(platform.ID)
	require.NoError(t, err, "update repo count failed")

	// Verify the repo count was updated
	got, err := d.FindByKey("github")
	require.NoError(t, err, "find by key failed")

	assert.Equal(t, 3, got.RepoCount, "expected repo count 3")
}
