package git_sync

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/git-sync-service/internal/converter"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
	repo "github.com/yi-nology/git-sync-service/biz/model/repo"
	syncmodel "github.com/yi-nology/git-sync-service/sync/model"
	sdkprov "github.com/yi-nology/git-platform-sdk/provider"
)

func ListRepos(ctx context.Context, c *app.RequestContext) {
	var req repo.ListReposReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	offset, limit := converter.PageToOffset(req.Page, req.PageSize)
	list, total, err := GetSyncService().ListRepos(ctx, offset, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &repo.ListReposResp{
		Repos: converter.ToRepoInfoList(list),
		Total: total,
	})
}

func GetRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.GetRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	r, err := GetSyncService().GetRepo(ctx, req.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if r == nil {
		response.NotFound(c, "repo not found")
		return
	}
	response.Success(c, &repo.GetRepoResp{Repo: converter.ToRepoInfo(r)})
}

func CreateRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.CreateRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Name == "" || req.RemoteURL == "" {
		response.BadRequest(c, "name and remoteUrl are required")
		return
	}

	r, err := GetSyncService().CreateRepo(ctx, &syncmodel.CreateRepoRequest{
		Name: req.Name, RemoteURL: req.RemoteURL, AccessToken: req.AccessToken,
	})
	if err != nil {
		if sdkprov.IsInvalidInput(err) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, fmt.Sprintf("create repo failed: %v", err))
		return
	}
	response.Created(c, &repo.CreateRepoResp{Repo: converter.ToRepoInfo(r)})
}

func UpdateRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.UpdateRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	r, err := GetSyncService().UpdateRepo(ctx, &syncmodel.UpdateRepoRequest{
		Key: req.Key, Name: req.Name, AccessToken: req.AccessToken,
	})
	if err != nil {
		if err.Error() == "repo not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, &repo.UpdateRepoResp{Repo: converter.ToRepoInfo(r)})
}

func DeleteRepo(ctx context.Context, c *app.RequestContext) {
	var req repo.DeleteRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	if err := GetSyncService().DeleteRepo(ctx, req.Key); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.NoContent(c)
}

func TestConnection(ctx context.Context, c *app.RequestContext) {
	var req repo.TestConnectionReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	result, err := GetSyncService().TestConnection(ctx, req.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, &repo.TestConnectionResp{Success: result.Success, Message: result.Message})
}

func ListBranches(ctx context.Context, c *app.RequestContext) {
	var req repo.ListBranchesReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	branches, err := GetSyncService().ListBranches(ctx, req.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, &repo.ListBranchesResp{Branches: branches})
}
