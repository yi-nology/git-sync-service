package service

import (
	"context"
	"fmt"
	"log/slog"

	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

// RegisterPlatformWebhook 在平台侧为仓库注册 Webhook,回调指向本服务的接收端点
// (/api/webhook/receive/:repoKey)。secret 非空时持久化到仓库记录(加密存储),
// 用于入站事件验签。events 为空时默认订阅 push。
func (s *Service) RegisterPlatformWebhook(ctx context.Context, repoKey, callbackURL, secret string, events []string) (*sdkprov.PlatformWebhook, error) {
	repo, prov, err := s.repos.GetRepoWithProvider(repoKey)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		events = []string{"push"}
	}

	wh, err := prov.CreateWebhook(ctx, sdkprov.CreateWebhookOptions{
		Owner:  repo.PlatformOwner,
		Repo:   repo.PlatformRepo,
		URL:    callbackURL,
		Secret: secret,
		Events: events,
	})
	if err != nil {
		return nil, fmt.Errorf("create platform webhook failed: %w", err)
	}

	if secret != "" {
		if err := s.repos.SetWebhookSecret(repoKey, secret); err != nil {
			// webhook 已在平台侧生效,密钥持久化失败仅告警,不回滚注册
			slog.Warn("persist webhook secret failed", "repoKey", repoKey, "error", err)
		}
	}

	return wh, nil
}

// ListPlatformWebhooks 列出平台侧已注册的 Webhook。
func (s *Service) ListPlatformWebhooks(ctx context.Context, repoKey string) ([]*sdkprov.PlatformWebhook, error) {
	repo, prov, err := s.repos.GetRepoWithProvider(repoKey)
	if err != nil {
		return nil, err
	}

	return prov.ListWebhooks(ctx, repo.PlatformOwner, repo.PlatformRepo)
}

// DeletePlatformWebhook 删除平台侧的 Webhook。
func (s *Service) DeletePlatformWebhook(ctx context.Context, repoKey string, webhookID int64) error {
	repo, prov, err := s.repos.GetRepoWithProvider(repoKey)
	if err != nil {
		return err
	}

	return prov.DeleteWebhook(ctx, repo.PlatformOwner, repo.PlatformRepo, webhookID)
}
