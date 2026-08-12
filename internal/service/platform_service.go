package service

import (
	"context"
	"fmt"

	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

// PlatformService 平台服务
type PlatformService struct {
	platformDAO *dao.PlatformDAO
	repoDAO     *dao.RepoDAO
	providerMgr *sdkprov.Manager
}

// NewPlatformService 创建 PlatformService
func NewPlatformService(platformDAO *dao.PlatformDAO, repoDAO *dao.RepoDAO, providerMgr *sdkprov.Manager) *PlatformService {
	return &PlatformService{
		platformDAO: platformDAO,
		repoDAO:     repoDAO,
		providerMgr: providerMgr,
	}
}

// CreatePlatform 创建平台
func (s *PlatformService) CreatePlatform(ctx context.Context, platform *model.Platform) error {
	return s.platformDAO.Create(platform)
}

// GetPlatform 获取平台
func (s *PlatformService) GetPlatform(ctx context.Context, key string) (*model.Platform, error) {
	return s.platformDAO.FindByKey(key)
}

// GetPlatformByID 根据 ID 获取平台
func (s *PlatformService) GetPlatformByID(ctx context.Context, id uint) (*model.Platform, error) {
	return s.platformDAO.FindByID(id)
}

// ListPlatforms 列出所有平台
func (s *PlatformService) ListPlatforms(ctx context.Context) ([]*model.Platform, error) {
	return s.platformDAO.FindAll()
}

// UpdatePlatform 更新平台
func (s *PlatformService) UpdatePlatform(ctx context.Context, platform *model.Platform) error {
	return s.platformDAO.Update(platform)
}

// DeletePlatform 删除平台
func (s *PlatformService) DeletePlatform(ctx context.Context, key string) error {
	return s.platformDAO.Delete(key)
}

// SetDefaultPlatform 设置默认平台
func (s *PlatformService) SetDefaultPlatform(ctx context.Context, key string) error {
	return s.platformDAO.SetDefault(key)
}

// UpdatePlatformStatus 更新平台状态
func (s *PlatformService) UpdatePlatformStatus(ctx context.Context, key, status, testResult string) error {
	return s.platformDAO.UpdateStatus(key, status, testResult)
}

// TestPlatformConnection 测试平台连接
func (s *PlatformService) TestPlatformConnection(ctx context.Context, key string) (*sdkprov.TestConnectionResult, error) {
	platform, err := s.platformDAO.FindByKey(key)
	if err != nil {
		return nil, fmt.Errorf("platform not found: %w", err)
	}

	// 创建 provider
	cfg := sdkprov.Config{
		Platform: sdkprov.Platform(platform.Type),
		BaseURL:  platform.APIURL,
		Token:    platform.AccessToken,
		SkipTLS:  platform.SkipTLSVerify,
	}

	provider, err := sdkprov.NewProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("create provider failed: %w", err)
	}

	// 测试连接
	result, err := provider.TestConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("test connection failed: %w", err)
	}

	return result, nil
}

// ListPlatformRepos 列出平台上的仓库
func (s *PlatformService) ListPlatformRepos(ctx context.Context, key, page, perPage string) ([]*sdkprov.PlatformRepo, error) {
	platform, err := s.platformDAO.FindByKey(key)
	if err != nil {
		return nil, fmt.Errorf("platform not found: %w", err)
	}

	// 创建 provider
	cfg := sdkprov.Config{
		Platform: sdkprov.Platform(platform.Type),
		BaseURL:  platform.APIURL,
		Token:    platform.AccessToken,
		SkipTLS:  platform.SkipTLSVerify,
	}

	provider, err := sdkprov.NewProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("create provider failed: %w", err)
	}

	// 获取仓库列表
	opts := sdkprov.ListRepoOptions{}
	repos, err := provider.ListRepos(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("list repos failed: %w", err)
	}

	return repos, nil
}

// SyncPlatformRepos 同步平台仓库到本地
func (s *PlatformService) SyncPlatformRepos(ctx context.Context, key string) (int, error) {
	platform, err := s.platformDAO.FindByKey(key)
	if err != nil {
		return 0, fmt.Errorf("platform not found: %w", err)
	}

	// 创建 provider
	cfg := sdkprov.Config{
		Platform: sdkprov.Platform(platform.Type),
		BaseURL:  platform.APIURL,
		Token:    platform.AccessToken,
		SkipTLS:  platform.SkipTLSVerify,
	}

	provider, err := sdkprov.NewProvider(cfg)
	if err != nil {
		return 0, fmt.Errorf("create provider failed: %w", err)
	}

	// 获取仓库列表
	opts := sdkprov.ListRepoOptions{}
	repos, err := provider.ListRepos(ctx, opts)
	if err != nil {
		return 0, fmt.Errorf("list repos failed: %w", err)
	}

	// 同步到本地
	count := 0
	for _, repo := range repos {
		// 检查是否已存在
		existing, _ := s.repoDAO.FindByCloneURL(repo.CloneURL)
		if existing != nil {
			continue
		}

		// 创建新仓库
		newRepo := &model.Repo{
			Key:           repo.FullName, // 使用全名作为 key
			Name:          repo.Name,
			PlatformID:    platform.ID,
			Platform:      platform.Type,
			PlatformOwner: repo.Owner,
			PlatformRepo:  repo.Name,
			CloneURL:      repo.CloneURL,
			SSHURL:        repo.SSHURL,
			DefaultBranch: repo.DefaultBranch,
			Status:        "active",
		}

		if err := s.repoDAO.Create(newRepo); err != nil {
			continue // 跳过失败的
		}
		count++
	}

	// 更新平台仓库数量
	_ = s.platformDAO.UpdateRepoCount(platform.ID)

	return count, nil
}

// ListReposByPlatform 列出平台下的仓库
func (s *PlatformService) ListReposByPlatform(ctx context.Context, platformKey string) ([]*model.Repo, error) {
	platform, err := s.platformDAO.FindByKey(platformKey)
	if err != nil {
		return nil, err
	}
	return s.repoDAO.FindByPlatformID(platform.ID)
}

// GetProviderByPlatform 根据平台获取 Provider
func (s *PlatformService) GetProviderByPlatform(platform *model.Platform) (sdkprov.Provider, error) {
	cfg := sdkprov.Config{
		Platform: sdkprov.Platform(platform.Type),
		BaseURL:  platform.APIURL,
		Token:    platform.AccessToken,
		SkipTLS:  platform.SkipTLSVerify,
	}

	return sdkprov.NewProvider(cfg)
}
