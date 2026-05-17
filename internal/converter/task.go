package converter

import (
	taskmodel "github.com/yi-nology/git-sync-service/biz/model/sync_task"
	"github.com/yi-nology/git-sync-service/sync/model"
)

func ToTaskInfo(t *model.SyncTask) *taskmodel.SyncTaskInfo {
	if t == nil {
		return nil
	}
	lastRunAt := ""
	if t.LastRunAt != nil {
		lastRunAt = t.LastRunAt.Format("2006-01-02 15:04:05")
	}
	return &taskmodel.SyncTaskInfo{
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

func ToTaskInfoList(tasks []*model.SyncTask) []*taskmodel.SyncTaskInfo {
	result := make([]*taskmodel.SyncTaskInfo, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, ToTaskInfo(t))
	}
	return result
}

func ToSyncRunInfo(r *model.SyncRun) *taskmodel.SyncRunInfo {
	if r == nil {
		return nil
	}
	endTime := ""
	if r.EndTime != nil {
		endTime = r.EndTime.Format("2006-01-02 15:04:05")
	}
	return &taskmodel.SyncRunInfo{
		ID: int64(r.ID), TaskKey: r.TaskKey, TriggerSource: r.TriggerSource,
		Status: r.Status, StartTime: r.StartTime.Format("2006-01-02 15:04:05"),
		EndTime: endTime, CommitRange: r.CommitRange, Details: r.Details,
		ErrorMessage: r.ErrorMessage, CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToSyncRunInfoList(runs []*model.SyncRun) []*taskmodel.SyncRunInfo {
	result := make([]*taskmodel.SyncRunInfo, 0, len(runs))
	for _, r := range runs {
		result = append(result, ToSyncRunInfo(r))
	}
	return result
}

func PageToOffset(page, pageSize int32) (int, int) {
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
