package git_sync

import (
	"context"
	"log/slog"

	"github.com/cloudwego/hertz/pkg/app"
	syncmodel "github.com/yi-nology/git-sync-service/sync/model"
)

// recordAudit 记录一条审计日志（best-effort：写入失败仅告警，不影响主流程）。
//
// 系统采用共享 API Key 鉴权，无用户登录概念，因此 actor 取自请求头 X-User，
// 缺省为 "admin"；待接入多用户鉴权后自然细化为真实用户。
func recordAudit(ctx context.Context, c *app.RequestContext, action, resourceType, resourceKey, resource string) {
	actor := string(c.GetHeader("X-User"))
	if actor == "" {
		actor = "admin"
	}
	entry := &syncmodel.OperationLog{
		Action:       action,
		ResourceType: resourceType,
		ResourceKey:  resourceKey,
		Resource:     resource,
		Actor:        actor,
		IP:           c.ClientIP(),
		Status:       syncmodel.StatusSuccess,
	}
	if err := GetSyncService().RecordOperation(ctx, entry); err != nil {
		slog.Warn("record audit log failed", "error", err, "action", action, "resource", resource)
	}
}
