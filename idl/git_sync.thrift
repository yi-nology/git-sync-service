namespace go git_sync

include "base.thrift"
include "repo.thrift"
include "sync_task.thrift"
include "webhook.thrift"

service RepoService {
    repo.ListReposResp ListRepos(1: repo.ListReposReq req) (api.get="/api/v1/repos")
    repo.GetRepoResp GetRepo(1: repo.GetRepoReq req) (api.get="/api/v1/repo")
    repo.CreateRepoResp CreateRepo(1: repo.CreateRepoReq req) (api.post="/api/v1/repo/create")
    repo.UpdateRepoResp UpdateRepo(1: repo.UpdateRepoReq req) (api.post="/api/v1/repo/update")
    repo.DeleteRepoResp DeleteRepo(1: repo.DeleteRepoReq req) (api.post="/api/v1/repo/delete")
    repo.TestConnectionResp TestConnection(1: repo.TestConnectionReq req) (api.post="/api/v1/repo/test")
    repo.ListBranchesResp ListBranches(1: repo.ListBranchesReq req) (api.get="/api/v1/repo/branches")
}

service SyncTaskService {
    sync_task.ListTasksResp ListTasks(1: sync_task.ListTasksReq req) (api.get="/api/v1/sync/tasks")
    sync_task.GetTaskResp GetTask(1: sync_task.GetTaskReq req) (api.get="/api/v1/sync/task")
    sync_task.CreateTaskResp CreateTask(1: sync_task.CreateTaskReq req) (api.post="/api/v1/sync/task/create")
    sync_task.UpdateTaskResp UpdateTask(1: sync_task.UpdateTaskReq req) (api.post="/api/v1/sync/task/update")
    sync_task.DeleteTaskResp DeleteTask(1: sync_task.DeleteTaskReq req) (api.post="/api/v1/sync/task/delete")
    sync_task.RunTaskResp RunTask(1: sync_task.RunTaskReq req) (api.post="/api/v1/sync/task/run")
    sync_task.PreviewSyncResp PreviewSync(1: sync_task.PreviewSyncReq req) (api.post="/api/v1/sync/preview")
    sync_task.ListHistoryResp ListHistory(1: sync_task.ListHistoryReq req) (api.get="/api/v1/sync/history")
}

service WebhookService {
    webhook.ListRulesResp ListRules(1: webhook.ListRulesReq req) (api.get="/api/v1/webhook/rules")
    webhook.GetRuleResp GetRule(1: webhook.GetRuleReq req) (api.get="/api/v1/webhook/rule")
    webhook.CreateRuleResp CreateRule(1: webhook.CreateRuleReq req) (api.post="/api/v1/webhook/rule/create")
    webhook.UpdateRuleResp UpdateRule(1: webhook.UpdateRuleReq req) (api.post="/api/v1/webhook/rule/update")
    webhook.DeleteRuleResp DeleteRule(1: webhook.DeleteRuleReq req) (api.post="/api/v1/webhook/rule/delete")
    webhook.ListEventsResp ListEvents(1: webhook.ListEventsReq req) (api.get="/api/v1/webhook/events")
    webhook.RetryEventResp RetryEvent(1: webhook.RetryEventReq req) (api.post="/api/v1/webhook/event/retry")
}
