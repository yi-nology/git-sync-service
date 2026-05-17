package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/internal/provider"
	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

func (s *Service) ListRepos() ([]*model.Repo, error) {
	return s.repoDAO.FindAll()
}

func (s *Service) GetRepo(key string) (*model.Repo, error) {
	return s.repoDAO.FindByKey(key)
}

func (s *Service) CreateRepo(req *model.CreateRepoRequest) (*model.Repo, error) {
	result, err := sdkprov.DetectPlatform(req.RemoteURL)
	if err != nil {
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
		Status:        "active",
	}

	if err := s.repoDAO.Create(repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *Service) UpdateRepo(req *model.UpdateRepoRequest) (*model.Repo, error) {
	repo, err := s.repoDAO.FindByKey(req.Key)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, fmt.Errorf("repo not found")
	}

	if req.Name != "" {
		repo.Name = req.Name
	}
	if req.AccessToken != "" {
		repo.AccessToken = req.AccessToken
	}

	if err := s.repoDAO.Update(repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *Service) DeleteRepo(key string) error {
	return s.repoDAO.Delete(key)
}

func (s *Service) TestConnection(ctx context.Context, repoKey string) (*model.TestConnectionResult, error) {
	repo, err := s.repoDAO.FindByKey(repoKey)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return &model.TestConnectionResult{Success: false, Message: "repo not found"}, nil
	}

	prov, err := s.providerMgr.GetProvider(repo)
	if err != nil {
		return &model.TestConnectionResult{Success: false, Message: err.Error()}, nil
	}

	result, err := prov.TestConnection(ctx)
	if err != nil {
		return &model.TestConnectionResult{Success: false, Message: err.Error()}, nil
	}

	return &model.TestConnectionResult{Success: result.Connected, Message: result.Message}, nil
}

func (s *Service) ListBranches(ctx context.Context, repoKey string) ([]string, error) {
	repo, err := s.repoDAO.FindByKey(repoKey)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, fmt.Errorf("repo not found")
	}

	prov, err := s.providerMgr.GetProvider(repo)
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

var _ provider.GitProvider = (*provider.SDKProviderAdapter)(nil)
