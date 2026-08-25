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

func setupPlatformServiceTestDB(t *testing.T) (*gorm.DB, *PlatformService) {
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

	platformDAO, err := dao.NewPlatformDAO(db)
	if err != nil {
		t.Fatalf("failed to create PlatformDAO: %v", err)
	}

	repoDAO, err := dao.NewRepoDAO(db)
	if err != nil {
		t.Fatalf("failed to create RepoDAO: %v", err)
	}

	providerMgr := sdkprov.NewManager(0)
	svc := NewPlatformService(platformDAO, repoDAO, providerMgr)

	return db, svc
}

func TestPlatformService_CreateAndGet(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Create a platform
	platform := &model.Platform{
		Key:         "github",
		Name:        "GitHub",
		Type:        "github",
		InstanceURL: "https://github.com",
		APIURL:      "https://api.github.com",
		Status:      "active",
	}

	if err := svc.CreatePlatform(ctx, platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if platform.ID == 0 {
		t.Error("expected platform ID to be set after create")
	}

	// Get by key
	got, err := svc.GetPlatform(ctx, "github")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.Key != "github" {
		t.Errorf("expected key 'github', got '%s'", got.Key)
	}

	if got.Name != "GitHub" {
		t.Errorf("expected name 'GitHub', got '%s'", got.Name)
	}
}

func TestPlatformService_GetByID(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Create a platform
	platform := &model.Platform{
		Key:  "gitlab",
		Name: "GitLab",
		Type: "gitlab",
	}

	if err := svc.CreatePlatform(ctx, platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Get by ID
	got, err := svc.GetPlatformByID(ctx, platform.ID)
	if err != nil {
		t.Fatalf("get by ID failed: %v", err)
	}

	if got.Key != "gitlab" {
		t.Errorf("expected key 'gitlab', got '%s'", got.Key)
	}

	if got.Name != "GitLab" {
		t.Errorf("expected name 'GitLab', got '%s'", got.Name)
	}
}

func TestPlatformService_List(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Create multiple platforms
	platforms := []*model.Platform{
		{Key: "github", Name: "GitHub", Type: "github"},
		{Key: "gitlab", Name: "GitLab", Type: "gitlab"},
		{Key: "bitbucket", Name: "Bitbucket", Type: "bitbucket"},
	}

	for _, p := range platforms {
		if err := svc.CreatePlatform(ctx, p); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List all
	got, err := svc.ListPlatforms(ctx)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 platforms, got %d", len(got))
	}
}

func TestPlatformService_Update(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Create a platform
	platform := &model.Platform{
		Key:         "github",
		Name:        "GitHub",
		Type:        "github",
		InstanceURL: "https://github.com",
	}

	if err := svc.CreatePlatform(ctx, platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update the platform
	platform.Name = "GitHub Enterprise"
	platform.InstanceURL = "https://github.example.com"

	if err := svc.UpdatePlatform(ctx, platform); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Get the updated platform
	got, err := svc.GetPlatform(ctx, "github")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.Name != "GitHub Enterprise" {
		t.Errorf("expected name 'GitHub Enterprise', got '%s'", got.Name)
	}

	if got.InstanceURL != "https://github.example.com" {
		t.Errorf("expected instance URL 'https://github.example.com', got '%s'", got.InstanceURL)
	}
}

func TestPlatformService_Delete(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Create a platform
	platform := &model.Platform{
		Key:  "github",
		Name: "GitHub",
	}

	if err := svc.CreatePlatform(ctx, platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete the platform
	if err := svc.DeletePlatform(ctx, "github"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Try to get the deleted platform - should return error
	_, err := svc.GetPlatform(ctx, "github")
	if err == nil {
		t.Error("expected error when getting deleted platform")
	}
}

func TestPlatformService_SetDefault(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Create two platforms
	platform1 := &model.Platform{
		Key:       "github",
		Name:      "GitHub",
		IsDefault: true,
	}

	platform2 := &model.Platform{
		Key:       "gitlab",
		Name:      "GitLab",
		IsDefault: false,
	}

	if err := svc.CreatePlatform(ctx, platform1); err != nil {
		t.Fatalf("create platform1 failed: %v", err)
	}

	if err := svc.CreatePlatform(ctx, platform2); err != nil {
		t.Fatalf("create platform2 failed: %v", err)
	}

	// Set platform2 as default
	if err := svc.SetDefaultPlatform(ctx, "gitlab"); err != nil {
		t.Fatalf("set default failed: %v", err)
	}

	// Verify platform2 is now default
	got2, err := svc.GetPlatform(ctx, "gitlab")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if !got2.IsDefault {
		t.Error("expected platform2 to be default")
	}

	// Verify platform1 is no longer default
	got1, err := svc.GetPlatform(ctx, "github")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got1.IsDefault {
		t.Error("expected platform1 to not be default")
	}
}

func TestPlatformService_UpdateStatus(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Create a platform
	platform := &model.Platform{
		Key:    "github",
		Name:   "GitHub",
		Status: "active",
	}

	if err := svc.CreatePlatform(ctx, platform); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update status
	if err := svc.UpdatePlatformStatus(ctx, "github", "inactive", "connection failed"); err != nil {
		t.Fatalf("update status failed: %v", err)
	}

	// Verify the status was updated
	got, err := svc.GetPlatform(ctx, "github")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.Status != "inactive" {
		t.Errorf("expected status 'inactive', got '%s'", got.Status)
	}

	if got.LastTestResult != "connection failed" {
		t.Errorf("expected last test result 'connection failed', got '%s'", got.LastTestResult)
	}
}

func TestPlatformService_GetNotFound(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Try to get a non-existent platform - should return error
	_, err := svc.GetPlatform(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error when getting non-existent platform")
	}
}

func TestPlatformService_GetByIDNotFound(t *testing.T) {
	_, svc := setupPlatformServiceTestDB(t)
	ctx := context.Background()

	// Try to get a non-existent platform - should return error
	_, err := svc.GetPlatformByID(ctx, 999)
	if err == nil {
		t.Error("expected error when getting non-existent platform")
	}
}

func TestRewriteCloneHost(t *testing.T) {
	cases := []struct {
		name       string
		rawURL     string
		instance   string
		want       string
	}{
		{"私有实例重写公网地址", "https://gitcode.com/yi-nology/iam-web.git", "gitcode.kylinos.cn", "https://gitcode.kylinos.cn/yi-nology/iam-web.git"},
		{"实例地址带 scheme", "https://gitcode.com/a/b.git", "https://gitcode.kylinos.cn", "https://gitcode.kylinos.cn/a/b.git"},
		{"host 相同不重写", "https://om-gitlab.kylinos.cn/obs/x.git", "om-gitlab.kylinos.cn", "https://om-gitlab.kylinos.cn/obs/x.git"},
		{"公网平台无实例不重写", "https://github.com/o/r.git", "", "https://github.com/o/r.git"},
		{"空地址不处理", "", "gitcode.kylinos.cn", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := &model.Platform{InstanceURL: c.instance}
			if got := rewriteCloneHost(c.rawURL, p); got != c.want {
				t.Errorf("rewriteCloneHost(%q, %q) = %q, want %q", c.rawURL, c.instance, got, c.want)
			}
		})
	}
}
