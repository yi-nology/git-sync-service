package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

// RepoService handles repository-related operations.
type RepoService struct {
	repoDAO     *dao.RepoDAO
	platformDAO *dao.PlatformDAO
	providerMgr *sdkprov.Manager
}

// NewRepoService creates a new RepoService instance.
func NewRepoService(repoDAO *dao.RepoDAO, platformDAO *dao.PlatformDAO, providerMgr *sdkprov.Manager) *RepoService {
	return &RepoService{
		repoDAO:     repoDAO,
		platformDAO: platformDAO,
		providerMgr: providerMgr,
	}
}

// ListRepos returns a paginated list of repositories.
func (rs *RepoService) ListRepos(ctx context.Context, offset, limit int) ([]*model.Repo, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return rs.repoDAO.FindAll(page)
}

// ListReposWithFilter returns a filtered, sorted, paginated list of repositories.
func (rs *RepoService) ListReposWithFilter(ctx context.Context, offset, limit int, filter *dao.RepoFilter) ([]*model.Repo, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return rs.repoDAO.ListWithFilter(page, *filter)
}

// GetRepo returns a repository by key.
func (rs *RepoService) GetRepo(ctx context.Context, key string) (*model.Repo, error) {
	return rs.repoDAO.FindByKey(key)
}

// CreateRepo creates a new repository.
func (rs *RepoService) CreateRepo(ctx context.Context, req *model.CreateRepoRequest) (*model.Repo, error) {
	// Try to detect platform from URL
	result, err := sdkprov.DetectPlatform(req.RemoteURL)
	if err != nil && req.PlatformID == 0 {
		// If detection fails and no platform_id provided, return error
		if errors.Is(err, sdkprov.ErrPlatformNotSupported) {
			return nil, fmt.Errorf("unsupported platform for URL %s: %w", req.RemoteURL, err)
		}
		return nil, fmt.Errorf("invalid remote URL: %w", err)
	}

	// If platform_id is provided, use the platform's type
	var platformType string
	var platformOwner, platformRepo string
	if result != nil {
		platformType = string(result.Platform)
		platformOwner = result.Owner
		platformRepo = result.Repo
	}

	if req.PlatformID > 0 && rs.platformDAO != nil {
		platform, err := rs.platformDAO.FindByID(req.PlatformID)
		if err == nil && platform != nil {
			platformType = platform.Type
		}
	}

	// Parse owner/repo from URL if not detected
	if platformOwner == "" || platformRepo == "" {
		// Try to parse from URL: https://host/owner/repo.git
		url := req.RemoteURL
		url = strings.TrimSuffix(url, ".git")
		url = strings.TrimSuffix(url, "/")
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			platformRepo = parts[len(parts)-1]
			platformOwner = parts[len(parts)-2]
		}
	}

	repo := &model.Repo{
		Key:           uuid.New().String(),
		Name:          req.Name,
		PlatformID:    req.PlatformID,
		Platform:      platformType,
		PlatformOwner: platformOwner,
		PlatformRepo:  platformRepo,
		CloneURL:      req.RemoteURL,
		AccessToken:   req.AccessToken,
		Status:        model.RepoStatusActive,
	}

	if err := rs.repoDAO.Create(repo); err != nil {
		return nil, err
	}

	// Update platform repo count
	if repo.PlatformID > 0 && rs.platformDAO != nil {
		rs.platformDAO.UpdateRepoCount(repo.PlatformID)
	}

	return repo, nil
}

// UpdateRepo updates an existing repository.
func (rs *RepoService) UpdateRepo(ctx context.Context, req *model.UpdateRepoRequest) (*model.Repo, error) {
	repo, err := rs.repoDAO.FindByKey(req.Key)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, ErrRepoNotFound
	}

	if req.Name != "" {
		repo.Name = req.Name
	}
	if req.AccessToken != "" {
		repo.AccessToken = req.AccessToken
	}

	if err := rs.repoDAO.Update(repo); err != nil {
		return nil, err
	}

	return repo, nil
}

// DeleteRepo deletes a repository by key.
func (rs *RepoService) DeleteRepo(ctx context.Context, key string) error {
	// Get repo first to know the platform_id
	repo, err := rs.repoDAO.FindByKey(key)
	if err != nil {
		return err
	}
	if repo == nil {
		return ErrRepoNotFound
	}

	platformID := repo.PlatformID

	if err := rs.repoDAO.Delete(key); err != nil {
		return err
	}

	// Update platform repo count
	if platformID > 0 && rs.platformDAO != nil {
		rs.platformDAO.UpdateRepoCount(platformID)
	}

	return nil
}

// GetProvider returns a git platform provider for the given repository.
func (rs *RepoService) GetProvider(cloneURL, accessToken string) (sdkprov.Provider, error) {
	return rs.providerMgr.GetByURL(cloneURL, accessToken)
}

// GetRepoByKey returns a repository by key.
func (rs *RepoService) GetRepoByKey(key string) (*model.Repo, error) {
	return rs.repoDAO.FindByKey(key)
}

// TestConnection tests the connection to a repository.
func (rs *RepoService) TestConnection(ctx context.Context, repoKey string) (*model.TestConnectionResult, error) {
	repo, err := rs.repoDAO.FindByKey(repoKey)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return &model.TestConnectionResult{Success: false, Message: "repo not found"}, nil
	}

	// Use repo's access token, or fall back to platform's token
	token := repo.AccessToken
	var platform *model.Platform
	if repo.PlatformID > 0 && rs.platformDAO != nil {
		platform, _ = rs.platformDAO.FindByID(repo.PlatformID)
		if platform != nil && token == "" {
			token = platform.AccessToken
		}
	}

	// Create provider using platform config if available
	var prov sdkprov.Provider
	if platform != nil {
		cfg := sdkprov.Config{
			Platform: sdkprov.Platform(platform.Type),
			BaseURL:  platform.APIURL,
			Token:    token,
			SkipTLS:  platform.SkipTLSVerify,
		}
		prov, err = sdkprov.NewProvider(cfg)
	} else {
		prov, err = rs.providerMgr.GetByURL(repo.CloneURL, token)
	}
	if err != nil {
		return &model.TestConnectionResult{Success: false, Message: err.Error()}, nil
	}

	result, err := prov.TestConnection(ctx)
	if err != nil {
		return &model.TestConnectionResult{Success: false, Message: err.Error()}, nil
	}

	return &model.TestConnectionResult{Success: result.Connected, Message: result.Message}, nil
}

// ListBranches returns a list of branches for a repository.
func (rs *RepoService) ListBranches(ctx context.Context, repoKey string) ([]string, error) {
	repo, err := rs.repoDAO.FindByKey(repoKey)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, ErrRepoNotFound
	}

	// Use repo's access token, or fall back to platform's token
	token := repo.AccessToken
	var platform *model.Platform
	if repo.PlatformID > 0 && rs.platformDAO != nil {
		platform, _ = rs.platformDAO.FindByID(repo.PlatformID)
		if platform != nil && token == "" {
			token = platform.AccessToken
		}
	}

	// Create provider using platform config if available
	var prov sdkprov.Provider
	if platform != nil {
		cfg := sdkprov.Config{
			Platform: sdkprov.Platform(platform.Type),
			BaseURL:  platform.APIURL,
			Token:    token,
			SkipTLS:  platform.SkipTLSVerify,
		}
		prov, err = sdkprov.NewProvider(cfg)
	} else {
		prov, err = rs.providerMgr.GetByURL(repo.CloneURL, token)
	}
	if err != nil {
		return nil, err
	}

	branches, err := prov.ListBranches(ctx, repo.PlatformOwner, repo.PlatformRepo)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, b := range branches {
		result = append(result, b.Name)
	}

	return result, nil
}
