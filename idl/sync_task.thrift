namespace go sync_task

struct SyncTaskInfo {
    1: i64 id (api.json="id")
    2: string key (api.json="key")
    3: string name (api.json="name")
    4: string source_repo_key (api.json="source_repo_key")
    5: string source_branch (api.json="source_branch")
    6: string target_repo_key (api.json="target_repo_key")
    7: string target_branch (api.json="target_branch")
    8: string sync_mode (api.json="sync_mode")
    9: string cron (api.json="cron")
    10: string webhook_token (api.json="webhook_token")
    11: bool enabled (api.json="enabled")
    12: bool git_tags (api.json="git_tags")
    13: bool git_force (api.json="git_force")
    14: bool git_prune (api.json="git_prune")
    15: bool git_no_verify (api.json="git_no_verify")
    16: string last_run_at (api.json="last_run_at")
    17: string last_status (api.json="last_status")
    18: string created_at (api.json="created_at")
}

struct ListTasksReq {
    1: string repo_key (api.query="repo_key")
    2: i32 page (api.query="page")
    3: i32 page_size (api.query="page_size")
}

struct ListTasksResp {
    1: list<SyncTaskInfo> tasks (api.json="tasks")
    2: i64 total (api.json="total")
}

struct GetTaskReq {
    1: string key (api.query="key")
}

struct GetTaskResp {
    1: SyncTaskInfo task (api.json="task")
}

struct CreateTaskReq {
    1: string name (api.json="name")
    2: string source_repo_key (api.json="source_repo_key")
    3: string source_branch (api.json="source_branch")
    4: string target_repo_key (api.json="target_repo_key")
    5: string target_branch (api.json="target_branch")
    6: string sync_mode (api.json="sync_mode")
    7: string cron (api.json="cron")
    8: bool git_tags (api.json="git_tags")
    9: bool git_force (api.json="git_force")
    10: bool git_prune (api.json="git_prune")
    11: bool git_no_verify (api.json="git_no_verify")
    12: string push_options (api.json="push_options")
}

struct CreateTaskResp {
    1: SyncTaskInfo task (api.json="task")
}

struct UpdateTaskReq {
    1: string key (api.json="key")
    2: string name (api.json="name")
    3: string source_branch (api.json="source_branch")
    4: string target_branch (api.json="target_branch")
    5: string sync_mode (api.json="sync_mode")
    6: string cron (api.json="cron")
    7: bool enabled (api.json="enabled")
    8: bool git_tags (api.json="git_tags")
    9: bool git_force (api.json="git_force")
    10: bool git_prune (api.json="git_prune")
    11: bool git_no_verify (api.json="git_no_verify")
    12: string push_options (api.json="push_options")
}

struct UpdateTaskResp {
    1: SyncTaskInfo task (api.json="task")
}

struct DeleteTaskReq {
    1: string key (api.query="key")
}

struct DeleteTaskResp {
    1: bool success (api.json="success")
}

struct RunTaskReq {
    1: string key (api.query="key")
}

struct RunTaskResp {
    1: bool success (api.json="success")
    2: string message (api.json="message")
}

struct PreviewSyncReq {
    1: string source_repo_key (api.json="source_repo_key")
    2: string source_branch (api.json="source_branch")
    3: string target_repo_key (api.json="target_repo_key")
    4: string target_branch (api.json="target_branch")
}

struct PreviewSyncResp {
    1: bool can_sync (api.json="can_sync")
    2: bool source_exists (api.json="source_exists")
    3: bool target_exists (api.json="target_exists")
    4: i32 commit_count (api.json="commit_count")
    5: string latest_commit (api.json="latest_commit")
    6: string message (api.json="message")
}

struct SyncRunInfo {
    1: i64 id (api.json="id")
    2: string task_key (api.json="task_key")
    3: string trigger_source (api.json="trigger_source")
    4: string status (api.json="status")
    5: string start_time (api.json="start_time")
    6: string end_time (api.json="end_time")
    7: string commit_range (api.json="commit_range")
    8: string details (api.json="details")
    9: string error_message (api.json="error_message")
    10: string created_at (api.json="created_at")
    11: optional i64 webhook_event_id (api.json="webhook_event_id")
    12: i64 duration_ms (api.json="duration_ms")
    13: string error_type (api.json="error_type")
    14: i32 retry_total (api.json="retry_total")
    15: list<SyncRunStepInfo> steps (api.json="steps")
}

struct SyncRunStepInfo {
    1: i64 id (api.json="id")
    2: string step_name (api.json="step_name")
    3: string status (api.json="status")
    4: string start_time (api.json="start_time")
    5: string end_time (api.json="end_time")
    6: i64 duration_ms (api.json="duration_ms")
    7: string error_msg (api.json="error_msg")
    8: string error_type (api.json="error_type")
    9: i32 retry_count (api.json="retry_count")
}

struct ListHistoryReq {
    1: string task_key (api.query="task_key")
    2: i32 limit (api.query="limit")
}

struct ListHistoryResp {
    1: list<SyncRunInfo> runs (api.json="runs")
}
