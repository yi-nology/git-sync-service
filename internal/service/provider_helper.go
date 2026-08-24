package service

import (
	"context"
	"fmt"
	"strconv"

	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// providerConfig 由平台记录构建 SDK provider 配置(token 可覆盖平台默认)。
func providerConfig(p *model.Platform, token string) sdkprov.Config {
	return sdkprov.Config{
		Platform: sdkprov.Platform(p.Type),
		BaseURL:  p.APIURL,
		Token:    token,
		SkipTLS:  p.SkipTLSVerify,
	}
}

// platformProvider 返回平台对应的 provider,统一经 Manager 缓存,
// 替代散落各处的 Config+NewProvider 样板。
func platformProvider(mgr *sdkprov.Manager, p *model.Platform) (sdkprov.Provider, error) {
	prov, err := mgr.Get(providerConfig(p, p.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("create provider failed: %w", err)
	}
	return prov, nil
}

// parsePageOpts 解析分页参数字符串并归一化(空/非法回落 SDK 默认值)。
func parsePageOpts(page, perPage string) (int, int) {
	p, _ := strconv.Atoi(page)
	pp, _ := strconv.Atoi(perPage)
	return sdkprov.NormalizePageOpts(p, pp)
}

// fetchAllPlatformRepos 按最大页大小循环翻页,拉取平台全部仓库。
// 返回条数不足一页时认为已到末页;maxPages 作为安全上限防止异常平台无限翻页。
func fetchAllPlatformRepos(ctx context.Context, provider sdkprov.Provider) ([]*sdkprov.PlatformRepo, error) {
	const maxPages = 100
	perPage := sdkprov.MaxPerPage

	all := make([]*sdkprov.PlatformRepo, 0, perPage)
	seen := make(map[string]struct{}, perPage)
	for page := 1; page <= maxPages; page++ {
		repos, err := provider.ListRepos(ctx, sdkprov.ListRepoOptions{Page: page, PerPage: perPage})
		if err != nil {
			return nil, err
		}
		for _, r := range repos {
			id := r.CloneURL
			if id == "" {
				id = r.FullName
			}
			if _, dup := seen[id]; dup {
				continue
			}
			seen[id] = struct{}{}
			all = append(all, r)
		}
		if len(repos) < perPage {
			break
		}
	}
	return all, nil
}
