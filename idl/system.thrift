// System status (admin dashboard) IDL.

struct SystemStatusReq {}

struct SystemStatusData {
    1: string status (api.json="status")
    2: string version (api.json="version")
    3: i64 uptime (api.json="uptime")
    4: i32 repoCount (api.json="repo_count")
    5: i32 taskCount (api.json="task_count")
    6: i32 runningTask (api.json="running_task")
    7: i32 failedTask (api.json="failed_task")
    8: string goVersion (api.json="go_version")
    9: string platform (api.json="platform")
}
