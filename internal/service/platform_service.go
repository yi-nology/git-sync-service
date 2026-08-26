package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

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
		return nil, fmt.Errorf("query platform failed: %w", err)
	}
	if platform == nil {
		return nil, fmt.Errorf("platform not found: %s", key)
	}

	provider, err := platformProvider(s.providerMgr, platform)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("query platform failed: %w", err)
	}
	if platform == nil {
		return nil, fmt.Errorf("platform not found: %s", key)
	}

	provider, err := platformProvider(s.providerMgr, platform)
	if err != nil {
		return nil, err
	}

	// 分页参数归一化后真正下传给 SDK
	p, pp := parsePageOpts(page, perPage)
	repos, err := provider.ListRepos(ctx, sdkprov.ListRepoOptions{Page: p, PerPage: pp})
	if err != nil {
		return nil, fmt.Errorf("list repos failed: %w", err)
	}

	return repos, nil
}

// SyncPlatformRepos 同步平台仓库到本地
func (s *PlatformService) SyncPlatformRepos(ctx context.Context, key string) (int, error) {
	platform, err := s.platformDAO.FindByKey(key)
	if err != nil {
		return 0, fmt.Errorf("query platform failed: %w", err)
	}
	if platform == nil {
		return 0, fmt.Errorf("platform not found: %s", key)
	}

	provider, err := platformProvider(s.providerMgr, platform)
	if err != nil {
		return 0, err
	}

	// 翻页拉取全部仓库,不能只取第一页,否则超过一页的仓库永远不会被同步
	repos, err := fetchAllPlatformRepos(ctx, provider)
	if err != nil {
		return 0, fmt.Errorf("list repos failed: %w", err)
	}

	// 一次加载该平台所有已有仓库到内存,避免 N 次 DB 查询。
	existingRepos, _ := s.repoDAO.FindByPlatformID(platform.ID)
	existingMap := make(map[string]*model.Repo, len(existingRepos))
	for _, r := range existingRepos {
		existingMap[r.PlatformRepo] = r
	}

	// 同步到本地
	count := 0
	for _, repo := range repos {
		// 私有部署实例的 API 可能返回公网 clone 地址(如 gitcode.kylinos.cn
		// 返回 gitcode.com),内网执行器不可达;按平台实例地址重写 host。
		cloneURL := rewriteCloneHost(repo.CloneURL, platform)
		sshURL := rewriteCloneHost(repo.SSHURL, platform)

		// 内存查重:比逐条 DB 查询快一个数量级。
		existing := existingMap[repo.Name]
		if existing != nil {
			// 已存在:实例地址变化时刷新存量 clone/ssh 地址,保证下次执行可用
			if existing.CloneURL != cloneURL {
				existing.CloneURL = cloneURL
				existing.SSHURL = sshURL
				if err := s.repoDAO.Update(existing); err != nil {
					slog.Error("sync repo: update clone URL failed", "repo", repo.FullName, "error", err)
					continue
				}
				count++
			}
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
			CloneURL:      cloneURL,
			SSHURL:        sshURL,
			DefaultBranch: repo.DefaultBranch,
			Status:        "active",
		}

		if err := s.repoDAO.Create(newRepo); err != nil {
			slog.Error("sync repo: create failed", "repo", repo.FullName, "error", err)
			continue
		}
		count++
	}

	// 更新平台仓库数量
	_ = s.platformDAO.UpdateRepoCount(platform.ID)

	return count, nil
}

// rewriteCloneHost 将仓库 clone/ssh 地址的 scheme+host 替换为平台实例地址。
// 仅当平台配置了私有实例地址(instance_url)且与地址 host 不同时重写,
// 避免 GitHub 等公网平台(api.github.com 与 github.com host 天然不同)被误改。
func rewriteCloneHost(rawURL string, platform *model.Platform) string {
	if rawURL == "" || platform == nil || platform.InstanceURL == "" {
		return rawURL
	}
	instance := platform.InstanceURL
	if !strings.Contains(instance, "://") {
		instance = "https://" + instance
	}
	iu, err := url.Parse(instance)
	if err != nil || iu.Host == "" {
		return rawURL
	}
	cu, err := url.Parse(rawURL)
	if err != nil || cu.Host == "" || cu.Host == iu.Host {
		return rawURL
	}
	cu.Scheme = iu.Scheme
	cu.Host = iu.Host
	return cu.String()
}

// ListReposByPlatform 列出平台下的仓库
func (s *PlatformService) ListReposByPlatform(ctx context.Context, platformKey string) ([]*model.Repo, error) {
	platform, err := s.platformDAO.FindByKey(platformKey)
	if err != nil {
		return nil, err
	}
	return s.repoDAO.FindByPlatformID(platform.ID)
}
