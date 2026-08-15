package service

import (
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
