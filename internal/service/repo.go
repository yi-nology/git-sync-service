package service

import (
	"context"

	"github.com/yi-nology/git-sync-service/sync/model"
)

func (s *Service) ListRepos(ctx context.Context, offset, limit int) ([]*model.Repo, int64, error) {
	return s.repoService.ListRepos(ctx, offset, limit)
}

func (s *Service) GetRepo(ctx context.Context, key string) (*model.Repo, error) {
	return s.repoService.GetRepo(ctx, key)
}

func (s *Service) CreateRepo(ctx context.Context, req *model.CreateRepoRequest) (*model.Repo, error) {
	return s.repoService.CreateRepo(ctx, req)
}

func (s *Service) UpdateRepo(ctx context.Context, req *model.UpdateRepoRequest) (*model.Repo, error) {
	return s.repoService.UpdateRepo(ctx, req)
}

func (s *Service) DeleteRepo(ctx context.Context, key string) error {
	return s.repoService.DeleteRepo(ctx, key)
}

func (s *Service) TestConnection(ctx context.Context, repoKey string) (*model.TestConnectionResult, error) {
	return s.repoService.TestConnection(ctx, repoKey)
}

func (s *Service) ListBranches(ctx context.Context, repoKey string) ([]string, error) {
	return s.repoService.ListBranches(ctx, repoKey)
}
