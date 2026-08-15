package git_sync

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/git-sync-service/biz/model/sync_task"
	"github.com/yi-nology/git-sync-service/internal/converter"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
	"github.com/yi-nology/git-sync-service/internal/service"
	syncmodel "github.com/yi-nology/git-sync-service/sync/model"
)

func TaskList(ctx context.Context, c *app.RequestContext) {
	var req sync_task.ListTasksReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	offset, limit := converter.PageToOffset(req.Page, req.PageSize)
	list, total, err := GetSyncService().ListTasks(ctx, req.RepoKey, offset, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &sync_task.ListTasksResp{
		Tasks: converter.ToTaskInfoList(list),
		Total: total,
	})
}

func TaskGet(ctx context.Context, c *app.RequestContext) {
	var req sync_task.GetTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	t, err := GetSyncService().GetTask(ctx, req.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if t == nil {
		response.NotFound(c, "task not found")
		return
	}
	response.Success(c, &sync_task.GetTaskResp{Task: converter.ToTaskInfo(t)})
}

func TaskCreate(ctx context.Context, c *app.RequestContext) {
	var req sync_task.CreateTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Name == "" || req.SourceRepoKey == "" || req.TargetRepoKey == "" {
		response.BadRequest(c, "name, sourceRepoKey and targetRepoKey are required")
		return
	}

	t, err := GetSyncService().CreateTask(ctx, &syncmodel.CreateTaskRequest{
		Name: req.Name, SourceRepoKey: req.SourceRepoKey, SourceBranch: req.SourceBranch,
		TargetRepoKey: req.TargetRepoKey, TargetBranch: req.TargetBranch,
		SyncMode: req.SyncMode, Cron: req.Cron, GitTags: req.GitTags,
		GitForce: req.GitForce, GitPrune: req.GitPrune,
	})
	if err != nil {
		response.InternalError(c, fmt.Sprintf("create task failed: %v", err))
		return
	}
	recordAudit(ctx, c, "create", "task", t.Key, "创建同步任务 "+t.Name)
	response.Created(c, &sync_task.CreateTaskResp{Task: converter.ToTaskInfo(t)})
}

func TaskUpdate(ctx context.Context, c *app.RequestContext) {
	var req sync_task.UpdateTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	t, err := GetSyncService().UpdateTask(ctx, &syncmodel.UpdateTaskRequest{
		Key: req.Key, Name: req.Name, SourceBranch: req.SourceBranch,
		TargetBranch: req.TargetBranch, SyncMode: req.SyncMode, Cron: req.Cron,
		Enabled: req.Enabled, GitTags: req.GitTags, GitForce: req.GitForce,
		GitPrune: req.GitPrune,
	})
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	recordAudit(ctx, c, "update", "task", t.Key, "更新同步任务 "+t.Key)
	response.Success(c, &sync_task.UpdateTaskResp{Task: converter.ToTaskInfo(t)})
}

func TaskDelete(ctx context.Context, c *app.RequestContext) {
	var req sync_task.DeleteTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}
	if err := GetSyncService().DeleteTask(ctx, req.Key); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	recordAudit(ctx, c, "delete", "task", req.Key, "删除同步任务 "+req.Key)
	response.NoContent(c)
}

func TaskRun(ctx context.Context, c *app.RequestContext) {
	var req sync_task.RunTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Key == "" {
		response.BadRequest(c, "key is required")
		return
	}
	if err := GetSyncService().RunTask(ctx, req.Key); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	recordAudit(ctx, c, "run", "task", req.Key, "手动触发同步 "+req.Key)
	response.Success(c, &sync_task.RunTaskResp{Success: true, Message: "task started"})
}

func TaskPreview(ctx context.Context, c *app.RequestContext) {
	var req sync_task.PreviewSyncReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.SourceRepoKey == "" || req.TargetRepoKey == "" {
		response.BadRequest(c, "sourceRepoKey and targetRepoKey are required")
		return
	}

	result, err := GetSyncService().PreviewSync(ctx, &syncmodel.PreviewSyncRequest{
		SourceRepoKey: req.SourceRepoKey, SourceBranch: req.SourceBranch,
		TargetRepoKey: req.TargetRepoKey, TargetBranch: req.TargetBranch,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, &sync_task.PreviewSyncResp{
		CanSync: result.CanSync, SourceExists: result.SourceExists,
		TargetExists: result.TargetExists, Message: result.Message,
	})
}

func TaskHistory(ctx context.Context, c *app.RequestContext) {
	var req sync_task.ListHistoryReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.TaskKey == "" {
		response.BadRequest(c, "taskKey is required")
		return
	}

	offset, _ := strconv.Atoi(c.Query("offset"))
	page := dao.DefaultPagination(offset, int(req.Limit))
	runs, _, err := GetSyncService().ListHistory(ctx, req.TaskKey, page.Offset, page.Limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &sync_task.ListHistoryResp{Runs: converter.ToSyncRunInfoList(runs)})
}
