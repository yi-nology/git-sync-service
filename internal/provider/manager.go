package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

type GitProvider interface {
	GetRepository(ctx context.Context, owner, name string) (*sdkprov.PlatformRepo, error)
	ListBranches(ctx context.Context, owner, name string) ([]*sdkprov.PlatformBranch, error)
	ValidateWebhookSignature(r *http.Request, secret string) error
	ParseWebhookEvent(r *http.Request, secret string) (*sdkprov.NormalizedEvent, error)
	TestConnection(ctx context.Context) (*sdkprov.TestConnectionResult, error)
	Platform() sdkprov.Platform
}

type SDKProviderAdapter struct {
	provider sdkprov.Provider
}

func (a *SDKProviderAdapter) GetRepository(ctx context.Context, owner, name string) (*sdkprov.PlatformRepo, error) {
	return a.provider.GetRepo(ctx, owner, name)
}

func (a *SDKProviderAdapter) ListBranches(ctx context.Context, owner, name string) ([]*sdkprov.PlatformBranch, error) {
	return a.provider.ListBranches(ctx, owner, name)
}

func (a *SDKProviderAdapter) ValidateWebhookSignature(r *http.Request, secret string) error {
	return a.provider.ValidateWebhookSignature(r, secret)
}

func (a *SDKProviderAdapter) ParseWebhookEvent(r *http.Request, secret string) (*sdkprov.NormalizedEvent, error) {
	return a.provider.ParseWebhookEvent(r, secret)
}

func (a *SDKProviderAdapter) TestConnection(ctx context.Context) (*sdkprov.TestConnectionResult, error) {
	return a.provider.TestConnection(ctx)
}

func (a *SDKProviderAdapter) Platform() sdkprov.Platform {
	return a.provider.Platform()
}

type ProviderManager struct {
	providers map[string]GitProvider
	mu        sync.RWMutex
}

func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: make(map[string]GitProvider),
	}
}

func (m *ProviderManager) GetProvider(repo *model.Repo) (GitProvider, error) {
	key := fmt.Sprintf("%s:%s", repo.Platform, repo.Key)

	m.mu.RLock()
	if p, ok := m.providers[key]; ok {
		m.mu.RUnlock()
		return p, nil
	}
	m.mu.RUnlock()

	p, err := m.newProvider(repo)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.providers[key] = p
	m.mu.Unlock()

	return p, nil
}

func (m *ProviderManager) newProvider(repo *model.Repo) (GitProvider, error) {
	platform, err := parsePlatform(repo.Platform)
	if err != nil {
		return nil, err
	}

	baseURL := getPlatformBaseURL(repo.Platform, repo.CloneURL)

	sdkProvider, err := sdkprov.NewProvider(sdkprov.Config{
		Platform: platform,
		BaseURL:  baseURL,
		Token:    repo.AccessToken,
		SkipTLS:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("create provider failed: %w", err)
	}

	return &SDKProviderAdapter{provider: sdkProvider}, nil
}

func parsePlatform(platform string) (sdkprov.Platform, error) {
	switch platform {
	case "github":
		return sdkprov.PlatformGitHub, nil
	case "gitlab":
		return sdkprov.PlatformGitLab, nil
	case "gitea":
		return sdkprov.PlatformGitea, nil
	case "forgejo":
		return sdkprov.PlatformForgejo, nil
	case "tencent_code":
		return sdkprov.PlatformTencentCode, nil
	case "gitee":
		return sdkprov.PlatformGitee, nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
}

func getPlatformBaseURL(platform, cloneURL string) string {
	if cloneURL == "" {
		return getDefaultBaseURL(platform)
	}

	u, err := parseCloneURL(cloneURL)
	if err != nil {
		return getDefaultBaseURL(platform)
	}

	baseURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	switch platform {
	case "github":
		if strings.Contains(u.Host, "github.com") {
			return ""
		}
	case "gitlab":
		if strings.Contains(u.Host, "gitlab.com") {
			return ""
		}
		baseURL = fmt.Sprintf("%s/api/v4", baseURL)
	case "gitea":
		baseURL = fmt.Sprintf("%s/api/v1", baseURL)
	case "forgejo":
		baseURL = fmt.Sprintf("%s/api/v1", baseURL)
	}

	return baseURL
}

func getDefaultBaseURL(platform string) string {
	switch platform {
	case "gitlab":
		return "https://gitlab.com/api/v4"
	case "gitea":
		return "https://gitea.com/api/v1"
	case "forgejo":
		return "https://codeberg.org/api/v1"
	case "tencent_code":
		return "https://e.coding.net/open-api"
	case "gitee":
		return "https://gitee.com/api/v5"
	default:
		return ""
	}
}

type URLParts struct {
	Scheme string
	Host   string
	Owner  string
	Repo   string
}

func parseCloneURL(cloneURL string) (*URLParts, error) {
	if strings.HasPrefix(cloneURL, "git@") {
		return parseSSHURL(cloneURL)
	}
	return parseHTTPSURL(cloneURL)
}

func parseSSHURL(cloneURL string) (*URLParts, error) {
	parts := strings.SplitN(strings.TrimPrefix(cloneURL, "git@"), ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SSH URL")
	}

	pathParts := strings.SplitN(strings.TrimSuffix(parts[1], ".git"), "/", 2)
	if len(pathParts) != 2 {
		return nil, fmt.Errorf("invalid SSH URL path")
	}

	return &URLParts{
		Scheme: "https",
		Host:   parts[0],
		Owner:  pathParts[0],
		Repo:   pathParts[1],
	}, nil
}

func parseHTTPSURL(cloneURL string) (*URLParts, error) {
	withoutScheme := strings.TrimPrefix(strings.TrimPrefix(cloneURL, "https://"), "http://")
	parts := strings.SplitN(withoutScheme, "/", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid HTTPS URL")
	}

	return &URLParts{
		Scheme: "https",
		Host:   parts[0],
		Owner:  parts[1],
		Repo:   strings.TrimSuffix(parts[2], ".git"),
	}, nil
}

func ExtractOwnerRepoFromURL(cloneURL string) (string, string, error) {
	parts, err := parseCloneURL(cloneURL)
	if err != nil {
		return "", "", err
	}
	return parts.Owner, parts.Repo, nil
}
