package service

import (
	"context"

	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// ListRepos returns a paginated list of repositories.
func (s *Service) ListRepos(ctx context.Context, offset, limit int) ([]*model.Repo, int64, error) {
	return s.repos.ListRepos(ctx, offset, limit)
}

// ListReposWithFilter returns a filtered, sorted, paginated list of repositories.
func (s *Service) ListReposWithFilter(ctx context.Context, offset, limit int, filter dao.RepoFilter) ([]*model.Repo, int64, error) {
	return s.repos.ListReposWithFilter(ctx, offset, limit, filter)
}

// GetRepo returns a repository by key.
func (s *Service) GetRepo(ctx context.Context, key string) (*model.Repo, error) {
	return s.repos.GetRepo(ctx, key)
}

// CreateRepo creates a new repository.
func (s *Service) CreateRepo(ctx context.Context, req *model.CreateRepoRequest) (*model.Repo, error) {
	return s.repos.CreateRepo(ctx, req)
}

// UpdateRepo updates an existing repository.
func (s *Service) UpdateRepo(ctx context.Context, req *model.UpdateRepoRequest) (*model.Repo, error) {
	return s.repos.UpdateRepo(ctx, req)
}

// DeleteRepo deletes a repository by key.
func (s *Service) DeleteRepo(ctx context.Context, key string) error {
	return s.repos.DeleteRepo(ctx, key)
}

// TestConnection tests the connection to a repository.
func (s *Service) TestConnection(ctx context.Context, repoKey string) (*model.TestConnectionResult, error) {
	return s.repos.TestConnection(ctx, repoKey)
}

// ListBranches returns a list of branches for a repository.
func (s *Service) ListBranches(ctx context.Context, repoKey string) ([]string, error) {
	return s.repos.ListBranches(ctx, repoKey)
}
