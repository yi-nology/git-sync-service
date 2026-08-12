

include "base.thrift"
include "repo.thrift"
include "sync_task.thrift"
include "webhook.thrift"

service RepoService {
    repo.ListReposResp List(1: repo.ListReposReq req) (api.get="/api/v1/repos")
    repo.GetRepoResp Get(1: repo.GetRepoReq req) (api.get="/api/v1/repo")
    repo.CreateRepoResp Create(1: repo.CreateRepoReq req) (api.post="/api/v1/repo/create")
    repo.UpdateRepoResp Update(1: repo.UpdateRepoReq req) (api.post="/api/v1/repo/update")
    repo.DeleteRepoResp Delete(1: repo.DeleteRepoReq req) (api.post="/api/v1/repo/delete")
    repo.TestConnectionResp Test(1: repo.TestConnectionReq req) (api.post="/api/v1/repo/test")
    repo.ListBranchesResp Branches(1: repo.ListBranchesReq req) (api.get="/api/v1/repo/branches")
}

service SyncTaskService {
    sync_task.ListTasksResp TaskList(1: sync_task.ListTasksReq req) (api.get="/api/v1/sync/tasks")
    sync_task.GetTaskResp TaskGet(1: sync_task.GetTaskReq req) (api.get="/api/v1/sync/task")
    sync_task.CreateTaskResp TaskCreate(1: sync_task.CreateTaskReq req) (api.post="/api/v1/sync/task/create")
    sync_task.UpdateTaskResp TaskUpdate(1: sync_task.UpdateTaskReq req) (api.post="/api/v1/sync/task/update")
    sync_task.DeleteTaskResp TaskDelete(1: sync_task.DeleteTaskReq req) (api.post="/api/v1/sync/task/delete")
    sync_task.RunTaskResp TaskRun(1: sync_task.RunTaskReq req) (api.post="/api/v1/sync/task/run")
    sync_task.PreviewSyncResp TaskPreview(1: sync_task.PreviewSyncReq req) (api.post="/api/v1/sync/preview")
    sync_task.ListHistoryResp TaskHistory(1: sync_task.ListHistoryReq req) (api.get="/api/v1/sync/history")
}

service WebhookService {
    webhook.ListRulesResp RuleList(1: webhook.ListRulesReq req) (api.get="/api/v1/webhook/rules")
    webhook.GetRuleResp RuleGet(1: webhook.GetRuleReq req) (api.get="/api/v1/webhook/rule")
    webhook.CreateRuleResp RuleCreate(1: webhook.CreateRuleReq req) (api.post="/api/v1/webhook/rule/create")
    webhook.UpdateRuleResp RuleUpdate(1: webhook.UpdateRuleReq req) (api.post="/api/v1/webhook/rule/update")
    webhook.DeleteRuleResp RuleDelete(1: webhook.DeleteRuleReq req) (api.post="/api/v1/webhook/rule/delete")
    webhook.ListEventsResp ListEvents(1: webhook.ListEventsReq req) (api.get="/api/v1/webhook/events")
    webhook.RetryEventResp RetryEvent(1: webhook.RetryEventReq req) (api.post="/api/v1/webhook/event/retry")
}
