package provider

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

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

type cachedProvider struct {
	provider  GitProvider
	createdAt time.Time
}

type ProviderManager struct {
	providers map[string]cachedProvider
	mu        sync.RWMutex
	ttl       time.Duration
}

func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: make(map[string]cachedProvider),
		ttl:       30 * time.Minute,
	}
}

func (m *ProviderManager) GetProvider(repo *model.Repo) (GitProvider, error) {
	key := fmt.Sprintf("%s:%s", repo.Platform, repo.Key)

	m.mu.RLock()
	if cp, ok := m.providers[key]; ok {
		if time.Since(cp.createdAt) < m.ttl {
			m.mu.RUnlock()
			return cp.provider, nil
		}
	}
	m.mu.RUnlock()

	p, err := m.newProvider(repo)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.providers[key] = cachedProvider{provider: p, createdAt: time.Now()}
	m.mu.Unlock()

	return p, nil
}

func (m *ProviderManager) newProvider(repo *model.Repo) (GitProvider, error) {
	result, err := sdkprov.DetectPlatform(repo.CloneURL)
	if err != nil {
		return nil, fmt.Errorf("detect platform failed: %w", err)
	}

	sdkProvider, err := sdkprov.NewProvider(sdkprov.Config{
		Platform: result.Platform,
		BaseURL:  result.BaseURL,
		Token:    repo.AccessToken,
		SkipTLS:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("create provider failed: %w", err)
	}

	return &SDKProviderAdapter{provider: sdkProvider}, nil
}

func ExtractOwnerRepoFromURL(cloneURL string) (string, string, error) {
	result, err := sdkprov.DetectPlatform(cloneURL)
	if err != nil {
		return "", "", err
	}
	return result.Owner, result.Repo, nil
}
