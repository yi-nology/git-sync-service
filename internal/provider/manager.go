package provider

import (
	"fmt"
	"sync"
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

type cachedProvider struct {
	provider  sdkprov.Provider
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

func (m *ProviderManager) GetProvider(repo *model.Repo) (sdkprov.Provider, error) {
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

func (m *ProviderManager) newProvider(repo *model.Repo) (sdkprov.Provider, error) {
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

	return sdkProvider, nil
}

func ExtractOwnerRepoFromURL(cloneURL string) (string, string, error) {
	result, err := sdkprov.DetectPlatform(cloneURL)
	if err != nil {
		return "", "", err
	}
	return result.Owner, result.Repo, nil
}
