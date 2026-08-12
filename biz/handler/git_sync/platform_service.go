package git_sync

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// PlatformCreateRequest 创建平台请求
type PlatformCreateRequest struct {
	Name          string `json:"name" vd:"required"`
	Type          string `json:"type" vd:"required"`
	InstanceURL   string `json:"instance_url"`
	APIURL        string `json:"api_url"`
	AccessToken   string `json:"access_token" vd:"required"`
	SkipTLSVerify bool   `json:"skip_tls_verify"`
	CACertPath    string `json:"ca_cert_path"`
	ProxyURL      string `json:"proxy_url"`
	IsDefault     bool   `json:"is_default"`
}

// PlatformUpdateRequest 更新平台请求
type PlatformUpdateRequest struct {
	Key           string `json:"key" vd:"required"`
	Name          string `json:"name"`
	InstanceURL   string `json:"instance_url"`
	APIURL        string `json:"api_url"`
	AccessToken   string `json:"access_token"`
	SkipTLSVerify *bool  `json:"skip_tls_verify"`
	CACertPath    string `json:"ca_cert_path"`
	ProxyURL      string `json:"proxy_url"`
	IsDefault     *bool  `json:"is_default"`
}

// PlatformListResponse 平台列表响应
type PlatformListResponse struct {
	Platforms []*model.Platform `json:"platforms"`
	Total     int64             `json:"total"`
}

// CreatePlatform 创建平台
func CreatePlatform(ctx context.Context, c *app.RequestContext) {
	var req PlatformCreateRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 生成 API URL（如果未提供）
	apiURL := req.APIURL
	if apiURL == "" && req.Type != model.PlatformTypeCustom {
		apiURL = model.GetAPIURL(req.Type, req.InstanceURL)
	}

	platform := &model.Platform{
		Key:           uuid.New().String(),
		Name:          req.Name,
		Type:          req.Type,
		InstanceURL:   req.InstanceURL,
		APIURL:        apiURL,
		AccessToken:   req.AccessToken,
		SkipTLSVerify: req.SkipTLSVerify,
		CACertPath:    req.CACertPath,
		ProxyURL:      req.ProxyURL,
		IsDefault:     req.IsDefault,
		Status:        model.PlatformStatusActive,
	}

	if err := GetSyncService().CreatePlatform(ctx, platform); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, map[string]interface{}{
		"platform": platform,
	})
}

// GetPlatform 获取单个平台
func GetPlatform(ctx context.Context, c *app.RequestContext) {
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	platform, err := GetSyncService().GetPlatform(ctx, key)
	if err != nil {
		response.NotFound(c, "platform not found")
		return
	}

	response.Success(c, map[string]interface{}{
		"platform": platform,
	})
}

// ListPlatforms 列出所有平台
func ListPlatforms(ctx context.Context, c *app.RequestContext) {
	platforms, err := GetSyncService().ListPlatforms(ctx)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &PlatformListResponse{
		Platforms: platforms,
		Total:     int64(len(platforms)),
	})
}

// UpdatePlatform 更新平台
func UpdatePlatform(ctx context.Context, c *app.RequestContext) {
	var req PlatformUpdateRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	platform, err := GetSyncService().GetPlatform(ctx, req.Key)
	if err != nil {
		response.NotFound(c, "platform not found")
		return
	}

	// 更新字段
	if req.Name != "" {
		platform.Name = req.Name
	}
	if req.InstanceURL != "" {
		platform.InstanceURL = req.InstanceURL
	}
	if req.APIURL != "" {
		platform.APIURL = req.APIURL
	}
	if req.AccessToken != "" {
		platform.AccessToken = req.AccessToken
	}
	if req.SkipTLSVerify != nil {
		platform.SkipTLSVerify = *req.SkipTLSVerify
	}
	if req.CACertPath != "" {
		platform.CACertPath = req.CACertPath
	}
	if req.ProxyURL != "" {
		platform.ProxyURL = req.ProxyURL
	}
	if req.IsDefault != nil {
		platform.IsDefault = *req.IsDefault
	}

	if err := GetSyncService().UpdatePlatform(ctx, platform); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"platform": platform,
	})
}

// DeletePlatform 删除平台
func DeletePlatform(ctx context.Context, c *app.RequestContext) {
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	// 检查是否有关联的仓库
	repos, err := GetSyncService().ListReposByPlatform(ctx, key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if len(repos) > 0 {
		response.BadRequest(c, "该平台下还有仓库，请先删除关联的仓库")
		return
	}

	if err := GetSyncService().DeletePlatform(ctx, key); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// SetDefaultPlatform 设置默认平台
func SetDefaultPlatform(ctx context.Context, c *app.RequestContext) {
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	if err := GetSyncService().SetDefaultPlatform(ctx, key); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"message": "设置成功",
	})
}

// TestPlatformConnection 测试平台连接
func TestPlatformConnection(ctx context.Context, c *app.RequestContext) {
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	result, err := GetSyncService().TestPlatformConnection(ctx, key)
	if err != nil {
		// 更新平台状态为错误
		_ = GetSyncService().UpdatePlatformStatus(ctx, key, model.PlatformStatusError, err.Error())
		response.InternalError(c, err.Error())
		return
	}

	// 更新平台状态为正常
	_ = GetSyncService().UpdatePlatformStatus(ctx, key, model.PlatformStatusActive, "connection successful")

	response.Success(c, map[string]interface{}{
		"result": result,
	})
}

// ListPlatformRepos 列出平台上的仓库（从远程 API 获取）
func ListPlatformRepos(ctx context.Context, c *app.RequestContext) {
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "30")

	repos, err := GetSyncService().ListPlatformRepos(ctx, key, page, perPage)
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
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	count, err := GetSyncService().SyncPlatformRepos(ctx, key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"message":     "同步成功",
		"syncedCount": count,
	})
}
