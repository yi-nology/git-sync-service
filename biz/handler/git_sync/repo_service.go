package git_sync

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	repo "github.com/yi-nology/git-sync-service/biz/model/repo"
	syncmodel "github.com/yi-nology/git-sync-service/sync/model"
)

func toRepoModel(r *syncmodel.Repo) *repo.RepoInfo {
	if r == nil {
		return nil
	}
	return &repo.RepoInfo{
		ID:            int64(r.ID),
		Key:           r.Key,
		Name:          r.Name,
		Platform:      r.Platform,
		PlatformOwner: r.PlatformOwner,
		PlatformRepo:  r.PlatformRepo,
		CloneURL:      r.CloneURL,
		SSHURL:        r.SSHURL,
		DefaultBranch: r.DefaultBranch,
		Status:        r.Status,
		CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func pageToOffset(page, pageSize int32) (int, int) {
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	if page <= 0 {
		page = 1
	}
	return int((page - 1) * pageSize), int(pageSize)
}

func ListRepos(ctx context.Context, c *app.RequestContext) {
	var req repo.ListReposReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	offset, limit := pageToOffset(req.Page, req.PageSize)
	list, total, err := GetSyncService().ListRepos(ctx, offset, limit)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	repos := make([]*repo.RepoInfo, 0, len(list))
	for _, r := range list {
		repos = append(repos, toRepoModel(r))
	}

	c.JSON(consts.StatusOK, &repo.ListReposResp{Repos: repos, Total: total})
}

func GetRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.GetRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}

	r, err := GetSyncService().GetRepo(ctx, req.Key)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &repo.GetRepoResp{Repo: toRepoModel(r)})
}

func CreateRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.CreateRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" || req.RemoteURL == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "name and remoteUrl are required"})
		return
	}

	r, err := GetSyncService().CreateRepo(ctx, &syncmodel.CreateRepoRequest{
		Name: req.Name, RemoteURL: req.RemoteURL, AccessToken: req.AccessToken,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": fmt.Sprintf("create repo failed: %v", err)})
		return
	}
	c.JSON(consts.StatusOK, &repo.CreateRepoResp{Repo: toRepoModel(r)})
}

func UpdateRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.UpdateRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}

	r, err := GetSyncService().UpdateRepo(ctx, &syncmodel.UpdateRepoRequest{
		Key: req.Key, Name: req.Name, AccessToken: req.AccessToken,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &repo.UpdateRepoResp{Repo: toRepoModel(r)})
}

func DeleteRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.DeleteRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}

	if err := GetSyncService().DeleteRepo(ctx, req.Key); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &repo.DeleteRepoResp{Success: true})
}

func TestConnection(ctx context.Context, c *app.RequestContext) {
	var req repo.TestConnectionReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}

	result, err := GetSyncService().TestConnection(ctx, req.Key)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &repo.TestConnectionResp{Success: result.Success, Message: result.Message})
}

func ListBranches(ctx context.Context, c *app.RequestContext) {
	var req repo.ListBranchesReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}

	branches, err := GetSyncService().ListBranches(ctx, req.Key)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &repo.ListBranchesResp{Branches: branches})
}

func parseInt32OrDefault(s string, def int32) int32 {
	if s == "" {
		return def
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return def
	}
	return int32(v)
}
