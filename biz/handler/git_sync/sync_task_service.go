package git_sync

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	sync_task "github.com/yi-nology/git-sync-service/biz/model/sync_task"
	syncmodel "github.com/yi-nology/git-sync-service/sync/model"
)

func toTaskModel(t *syncmodel.SyncTask) *sync_task.SyncTaskInfo {
	if t == nil {
		return nil
	}
	lastRunAt := ""
	if t.LastRunAt != nil {
		lastRunAt = t.LastRunAt.Format("2006-01-02 15:04:05")
	}
	return &sync_task.SyncTaskInfo{
		ID: int64(t.ID), Key: t.Key, Name: t.Name,
		SourceRepoKey: t.SourceRepoKey, SourceBranch: t.SourceBranch,
		TargetRepoKey: t.TargetRepoKey, TargetBranch: t.TargetBranch,
		SyncMode: t.SyncMode, Cron: t.Cron, WebhookToken: t.WebhookToken,
		Enabled: t.Enabled, GitTags: t.GitTags, GitForce: t.GitForce,
		GitPrune: t.GitPrune, GitNoVerify: t.GitNoVerify,
		LastRunAt: lastRunAt, LastStatus: t.LastStatus,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toSyncRunModel(r *syncmodel.SyncRun) *sync_task.SyncRunInfo {
	if r == nil {
		return nil
	}
	endTime := ""
	if r.EndTime != nil {
		endTime = r.EndTime.Format("2006-01-02 15:04:05")
	}
	return &sync_task.SyncRunInfo{
		ID: int64(r.ID), TaskKey: r.TaskKey, TriggerSource: r.TriggerSource,
		Status: r.Status, StartTime: r.StartTime.Format("2006-01-02 15:04:05"),
		EndTime: endTime, CommitRange: r.CommitRange, Details: r.Details,
		ErrorMessage: r.ErrorMessage, CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ListTasks(ctx context.Context, c *app.RequestContext) {
	var req sync_task.ListTasksReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	offset, limit := pageToOffset(req.Page, req.PageSize)
	list, total, err := GetSyncService().ListTasks(ctx, req.RepoKey, offset, limit)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	tasks := make([]*sync_task.SyncTaskInfo, 0, len(list))
	for _, t := range list {
		tasks = append(tasks, toTaskModel(t))
	}
	c.JSON(consts.StatusOK, &sync_task.ListTasksResp{Tasks: tasks, Total: total})
}

func GetTask(ctx context.Context, c *app.RequestContext) {
	var req sync_task.GetTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}

	t, err := GetSyncService().GetTask(ctx, req.Key)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &sync_task.GetTaskResp{Task: toTaskModel(t)})
}

func CreateTask(ctx context.Context, c *app.RequestContext) {
	var req sync_task.CreateTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" || req.SourceRepoKey == "" || req.TargetRepoKey == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "name, sourceRepoKey and targetRepoKey are required"})
		return
	}

	t, err := GetSyncService().CreateTask(ctx, &syncmodel.CreateTaskRequest{
		Name: req.Name, SourceRepoKey: req.SourceRepoKey, SourceBranch: req.SourceBranch,
		TargetRepoKey: req.TargetRepoKey, TargetBranch: req.TargetBranch,
		SyncMode: req.SyncMode, Cron: req.Cron, GitTags: req.GitTags,
		GitForce: req.GitForce, GitPrune: req.GitPrune,
		GitNoVerify: req.GitNoVerify, PushOptions: req.PushOptions,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": fmt.Sprintf("create task failed: %v", err)})
		return
	}
	c.JSON(consts.StatusOK, &sync_task.CreateTaskResp{Task: toTaskModel(t)})
}

func UpdateTask(ctx context.Context, c *app.RequestContext) {
	var req sync_task.UpdateTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}

	t, err := GetSyncService().UpdateTask(ctx, &syncmodel.UpdateTaskRequest{
		Key: req.Key, Name: req.Name, SourceBranch: req.SourceBranch,
		TargetBranch: req.TargetBranch, SyncMode: req.SyncMode, Cron: req.Cron,
		Enabled: req.Enabled, GitTags: req.GitTags, GitForce: req.GitForce,
		GitPrune: req.GitPrune, GitNoVerify: req.GitNoVerify, PushOptions: req.PushOptions,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &sync_task.UpdateTaskResp{Task: toTaskModel(t)})
}

func DeleteTask(ctx context.Context, c *app.RequestContext) {
	var req sync_task.DeleteTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}
	if err := GetSyncService().DeleteTask(ctx, req.Key); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &sync_task.DeleteTaskResp{Success: true})
}

func RunTask(ctx context.Context, c *app.RequestContext) {
	var req sync_task.RunTaskReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Key == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "key is required"})
		return
	}
	if err := GetSyncService().RunTask(ctx, req.Key); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &sync_task.RunTaskResp{Success: true, Message: "task started"})
}

func PreviewSync(ctx context.Context, c *app.RequestContext) {
	var req sync_task.PreviewSyncReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.SourceRepoKey == "" || req.TargetRepoKey == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "sourceRepoKey and targetRepoKey are required"})
		return
	}

	result, err := GetSyncService().PreviewSync(ctx, &syncmodel.PreviewSyncRequest{
		SourceRepoKey: req.SourceRepoKey, SourceBranch: req.SourceBranch,
		TargetRepoKey: req.TargetRepoKey, TargetBranch: req.TargetBranch,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, &sync_task.PreviewSyncResp{
		CanSync: result.CanSync, SourceExists: result.SourceExists,
		TargetExists: result.TargetExists, CommitCount: int32(result.CommitCount),
		LatestCommit: result.LatestCommit, Message: result.Message,
	})
}

func ListHistory(ctx context.Context, c *app.RequestContext) {
	var req sync_task.ListHistoryReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.TaskKey == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{"error": "taskKey is required"})
		return
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 50
	}
	runs, _, err := GetSyncService().ListHistory(ctx, req.TaskKey, 0, limit)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	result := make([]*sync_task.SyncRunInfo, 0, len(runs))
	for _, r := range runs {
		result = append(result, toSyncRunModel(r))
	}
	c.JSON(consts.StatusOK, &sync_task.ListHistoryResp{Runs: result})
}
