package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

// RepoService handles repository-related operations.
type RepoService struct {
	repoDAO     *dao.RepoDAO
	providerMgr *sdkprov.Manager
}

// NewRepoService creates a new RepoService instance.
func NewRepoService(repoDAO *dao.RepoDAO, providerMgr *sdkprov.Manager) *RepoService {
	return &RepoService{
		repoDAO:     repoDAO,
		providerMgr: providerMgr,
	}
}

// ListRepos returns a paginated list of repositories.
func (rs *RepoService) ListRepos(ctx context.Context, offset, limit int) ([]*model.Repo, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return rs.repoDAO.FindAll(page)
}

// GetRepo returns a repository by key.
func (rs *RepoService) GetRepo(ctx context.Context, key string) (*model.Repo, error) {
	return rs.repoDAO.FindByKey(key)
}

// CreateRepo creates a new repository.
func (rs *RepoService) CreateRepo(ctx context.Context, req *model.CreateRepoRequest) (*model.Repo, error) {
	result, err := sdkprov.DetectPlatform(req.RemoteURL)
	if err != nil {
		if errors.Is(err, sdkprov.ErrPlatformNotSupported) {
			return nil, fmt.Errorf("unsupported platform for URL %s: %w", req.RemoteURL, err)
		}
		return nil, fmt.Errorf("invalid remote URL: %w", err)
	}

	repo := &model.Repo{
		Key:           uuid.New().String(),
		Name:          req.Name,
		Platform:      string(result.Platform),
		PlatformOwner: result.Owner,
		PlatformRepo:  result.Repo,
		CloneURL:      req.RemoteURL,
		AccessToken:   req.AccessToken,
		Status:        model.RepoStatusActive,
	}

	if err := rs.repoDAO.Create(repo); err != nil {
		return nil, err
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
	return rs.repoDAO.Delete(key)
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

	prov, err := rs.providerMgr.GetByURL(repo.CloneURL, repo.AccessToken)
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

	prov, err := rs.providerMgr.GetByURL(repo.CloneURL, repo.AccessToken)
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
