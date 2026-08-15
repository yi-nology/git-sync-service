package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	// 注册全部平台后端,使 NewProvider 在测试二进制中可用
	_ "github.com/yi-nology/git-platform-sdk/backends/all"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

func setupRegisterTestService(t *testing.T, platformAPIURL string) (*Service, *gorm.DB) {
	t.Helper()
	t.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Repo{}, &model.Platform{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repoDAO, err := dao.NewRepoDAO(db)
	if err != nil {
		t.Fatalf("NewRepoDAO: %v", err)
	}
	platformDAO, err := dao.NewPlatformDAO(db)
	if err != nil {
		t.Fatalf("NewPlatformDAO: %v", err)
	}

	// 平台记录指向 fake gitea 服务器
	platform := &model.Platform{
		Key:  "gitea-test",
		Name: "Gitea Test",
		Type: string(sdkprov.PlatformGitea),
		APIURL: platformAPIURL,
	}
	if err := db.Create(platform).Error; err != nil {
		t.Fatalf("create platform: %v", err)
	}

	if err := repoDAO.Create(&model.Repo{
		Key:           "demo-repo",
		Name:          "demo",
		PlatformID:    platform.ID,
		Platform:      platform.Type,
		PlatformOwner: "octo",
		PlatformRepo:  "demo",
		CloneURL:      platformAPIURL + "/octo/demo.git",
		Status:        model.RepoStatusActive,
	}); err != nil {
		t.Fatalf("create repo: %v", err)
	}

	svc := &Service{
		repos: NewRepoService(repoDAO, platformDAO, sdkprov.NewManager(0)),
	}
	return svc, db
}

// newFakeGitea 模拟 gitea hooks API:创建/列表/删除。
func newFakeGitea(t *testing.T) (*httptest.Server, *[]string) {
	t.Helper()
	var calls []string
	mux := http.NewServeMux()
	// gitea 后端初始化时会探测版本
	mux.HandleFunc("/api/v1/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"version": "1.22.0"})
	})
	mux.HandleFunc("/api/v1/repos/octo/demo/hooks", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 7, "url": "http://cb", "events": []string{"push"}, "active": true})
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": 7, "url": "http://cb", "events": []string{"push"}, "active": true}})
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("/api/v1/repos/octo/demo/hooks/7", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, &calls
}

func TestRegisterPlatformWebhook(t *testing.T) {
	srv, calls := newFakeGitea(t)
	svc, _ := setupRegisterTestService(t, srv.URL)
	ctx := context.Background()

	wh, err := svc.RegisterPlatformWebhook(ctx, "demo-repo", "http://cb", "s3cret", nil)
	if err != nil {
		t.Fatalf("RegisterPlatformWebhook: %v", err)
	}
	if wh.ID != 7 {
		t.Errorf("expected webhook id 7, got %d", wh.ID)
	}
	if len(*calls) == 0 || !strings.HasPrefix((*calls)[0], "POST") {
		t.Errorf("expected POST to platform, calls=%v", *calls)
	}

	// secret 应已加密持久化:读取可解密回原文
	repo, err := svc.repos.GetRepoByKey("demo-repo")
	if err != nil || repo == nil {
		t.Fatalf("repo reload: %v %v", repo, err)
	}
	if repo.WebhookSecret != "s3cret" {
		t.Errorf("expected webhook secret persisted, got %q", repo.WebhookSecret)
	}
}

func TestRegisterPlatformWebhook_RepoNotFound(t *testing.T) {
	srv, _ := newFakeGitea(t)
	svc, _ := setupRegisterTestService(t, srv.URL)

	if _, err := svc.RegisterPlatformWebhook(context.Background(), "nope", "http://cb", "", nil); err != ErrRepoNotFound {
		t.Errorf("expected ErrRepoNotFound, got %v", err)
	}
}

func TestListPlatformWebhooks(t *testing.T) {
	srv, _ := newFakeGitea(t)
	svc, _ := setupRegisterTestService(t, srv.URL)

	list, err := svc.ListPlatformWebhooks(context.Background(), "demo-repo")
	if err != nil {
		t.Fatalf("ListPlatformWebhooks: %v", err)
	}
	if len(list) != 1 || list[0].ID != 7 {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestDeletePlatformWebhook(t *testing.T) {
	srv, calls := newFakeGitea(t)
	svc, _ := setupRegisterTestService(t, srv.URL)

	if err := svc.DeletePlatformWebhook(context.Background(), "demo-repo", 7); err != nil {
		t.Fatalf("DeletePlatformWebhook: %v", err)
	}
	joined := strings.Join(*calls, ";")
	if !strings.Contains(joined, "DELETE /api/v1/repos/octo/demo/hooks/7") {
		t.Errorf("expected DELETE call, calls=%v", *calls)
	}
}
