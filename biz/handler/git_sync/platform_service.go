package git_sync

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	platform "github.com/yi-nology/git-sync-service/biz/model/platform"
	"github.com/yi-nology/git-sync-service/internal/converter"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// CreatePlatform 创建平台
func CreatePlatform(ctx context.Context, c *app.RequestContext) {
	var req platform.CreatePlatformReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Name == "" || req.Type == "" || req.AccessToken == "" {
		response.BadRequest(c, "name, type, access_token are required")
		return
	}

	// 生成 API URL（如果未提供）
	apiURL := req.ApiUrl
	if apiURL == "" && req.Type != model.PlatformTypeCustom {
		apiURL = model.GetAPIURL(req.Type, req.InstanceUrl)
	}

	p := &model.Platform{
		Key:           uuid.New().String(),
		Name:          req.Name,
		Type:          req.Type,
		InstanceURL:   req.InstanceUrl,
		APIURL:        apiURL,
		AccessToken:   req.AccessToken,
		SkipTLSVerify: req.SkipTlsVerify,
		CACertPath:    req.CaCertPath,
		ProxyURL:      req.ProxyUrl,
		IsDefault:     req.IsDefault,
		Status:        model.PlatformStatusActive,
	}

	if err := GetSyncService().CreatePlatform(ctx, p); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	recordAudit(ctx, c, "create", "platform", p.Key, "创建平台 "+p.Name)
	response.Created(c, &platform.CreatePlatformResp{
		Platform: converter.ToPlatformInfo(p),
	})
}

// GetPlatform 获取单个平台
func GetPlatform(ctx context.Context, c *app.RequestContext) {
	var req platform.GetPlatformReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	p, err := GetSyncService().GetPlatform(ctx, req.Key)
	if err != nil {
		response.NotFound(c, "platform not found")
		return
	}

	response.Success(c, &platform.GetPlatformResp{
		Platform: converter.ToPlatformInfo(p),
	})
}

// ListPlatforms 列出所有平台
func ListPlatforms(ctx context.Context, c *app.RequestContext) {
	platforms, err := GetSyncService().ListPlatforms(ctx)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &platform.ListPlatformsResp{
		Platforms: converter.ToPlatformList(platforms),
		Total:     int64(len(platforms)),
	})
}

// UpdatePlatform 更新平台
func UpdatePlatform(ctx context.Context, c *app.RequestContext) {
	var req platform.UpdatePlatformReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	p, err := GetSyncService().GetPlatform(ctx, req.Key)
	if err != nil {
		response.NotFound(c, "platform not found")
		return
	}

	// 更新字段（空值表示不更新）
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.InstanceUrl != "" {
		p.InstanceURL = req.InstanceUrl
	}
	if req.ApiUrl != "" {
		p.APIURL = req.ApiUrl
	}
	if req.AccessToken != "" {
		p.AccessToken = req.AccessToken
	}
	if req.SkipTlsVerify != nil {
		p.SkipTLSVerify = *req.SkipTlsVerify
	}
	if req.CaCertPath != "" {
		p.CACertPath = req.CaCertPath
	}
	if req.ProxyUrl != "" {
		p.ProxyURL = req.ProxyUrl
	}
	if req.IsDefault != nil {
		p.IsDefault = *req.IsDefault
	}

	if err := GetSyncService().UpdatePlatform(ctx, p); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	recordAudit(ctx, c, "update", "platform", req.Key, "更新平台 "+p.Name)
	response.Success(c, &platform.UpdatePlatformResp{
		Platform: converter.ToPlatformInfo(p),
	})
}

// DeletePlatform 删除平台
func DeletePlatform(ctx context.Context, c *app.RequestContext) {
	var req platform.DeletePlatformReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	// 检查是否有关联的仓库
	repos, err := GetSyncService().ListReposByPlatform(ctx, req.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if len(repos) > 0 {
		response.BadRequest(c, "该平台下还有仓库，请先删除关联的仓库")
		return
	}

	if err := GetSyncService().DeletePlatform(ctx, req.Key); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	recordAudit(ctx, c, "delete", "platform", req.Key, "删除平台 "+req.Key)
	response.NoContent(c)
}

// SetDefaultPlatform 设置默认平台
func SetDefaultPlatform(ctx context.Context, c *app.RequestContext) {
	var req platform.SetDefaultPlatformReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	if err := GetSyncService().SetDefaultPlatform(ctx, req.Key); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	recordAudit(ctx, c, "update", "platform", req.Key, "设置默认平台 "+req.Key)
	response.Success(c, &platform.SetDefaultPlatformResp{
		Message: "设置成功",
	})
}

// TestPlatformConnection 测试平台连接
func TestPlatformConnection(ctx context.Context, c *app.RequestContext) {
	var req platform.TestPlatformConnectionReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	result, err := GetSyncService().TestPlatformConnection(ctx, req.Key)
	if err != nil {
		// 更新平台状态为错误
		_ = GetSyncService().UpdatePlatformStatus(ctx, req.Key, model.PlatformStatusError, err.Error())
		response.InternalError(c, err.Error())
		return
	}

	// 更新平台状态为正常
	_ = GetSyncService().UpdatePlatformStatus(ctx, req.Key, model.PlatformStatusActive, "connection successful")

	response.Success(c, map[string]interface{}{
		"result": result,
	})
}

// ListPlatformRepos 列出平台上的仓库（从远程 API 获取）
func ListPlatformRepos(ctx context.Context, c *app.RequestContext) {
	var req platform.ListPlatformReposReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	page := req.Page
	if page == "" {
		page = "1"
	}
	perPage := req.PerPage
	if perPage == "" {
		perPage = "30"
	}

	repos, err := GetSyncService().ListPlatformRepos(ctx, req.Key, page, perPage)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"repos": repos,
		"total": len(repos),
	})
}

// SyncPlatformRepos 同步平台仓库到本地
func SyncPlatformRepos(ctx context.Context, c *app.RequestContext) {
	var req platform.SyncPlatformReposReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	count, err := GetSyncService().SyncPlatformRepos(ctx, req.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	recordAudit(ctx, c, "sync", "platform", req.Key, fmt.Sprintf("同步平台仓库 %s（导入 %d 个）", req.Key, count))
	response.Success(c, &platform.SyncPlatformReposResp{
		Message:     "同步成功",
		SyncedCount: int32(count),
	})
}
