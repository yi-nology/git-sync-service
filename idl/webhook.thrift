namespace go webhook

struct WebhookRuleInfo {
    1: i64 id (api.json="id")
    2: string name (api.json="name")
    3: string repo_key (api.json="repo_key")
    4: string event_type (api.json="event_type")
    5: string branch_pattern (api.json="branch_pattern")
    6: string action (api.json="action")
    7: string sync_task_keys (api.json="sync_task_keys")
    8: i32 min_interval (api.json="min_interval")
    9: bool enabled (api.json="enabled")
    10: string description (api.json="description")
    11: string created_at (api.json="created_at")
}

struct ListRulesReq {
    1: string repo_key (api.query="repo_key")
}

struct ListRulesResp {
    1: list<WebhookRuleInfo> rules (api.json="rules")
}

struct GetRuleReq {
    1: i64 id (api.query="id")
}

struct GetRuleResp {
    1: WebhookRuleInfo rule (api.json="rule")
}

struct CreateRuleReq {
    1: string name (api.json="name")
    2: string repo_key (api.json="repo_key")
    3: string event_type (api.json="event_type")
    4: string branch_pattern (api.json="branch_pattern")
    5: string action (api.json="action")
    6: string sync_task_keys (api.json="sync_task_keys")
    7: i32 min_interval (api.json="min_interval")
    8: bool enabled (api.json="enabled")
    9: string description (api.json="description")
}

struct CreateRuleResp {
    1: WebhookRuleInfo rule (api.json="rule")
}

struct UpdateRuleReq {
    1: i64 id (api.json="id")
    2: string name (api.json="name")
    3: string event_type (api.json="event_type")
    4: string branch_pattern (api.json="branch_pattern")
    5: string action (api.json="action")
    6: string sync_task_keys (api.json="sync_task_keys")
    7: i32 min_interval (api.json="min_interval")
    8: bool enabled (api.json="enabled")
    9: string description (api.json="description")
}

struct UpdateRuleResp {
    1: WebhookRuleInfo rule (api.json="rule")
}

struct DeleteRuleReq {
    1: i64 id (api.query="id")
}

struct DeleteRuleResp {
    1: bool success (api.json="success")
}

struct WebhookEventInfo {
    1: i64 id (api.json="id")
    2: string event_id (api.json="event_id")
    3: string repo_key (api.json="repo_key")
    4: string event_type (api.json="event_type")
    5: string source (api.json="source")
    6: string actor_name (api.json="actor_name")
    7: string branch (api.json="branch")
    8: string commit_sha (api.json="commit_sha")
    9: string status (api.json="status")
    10: string error_message (api.json="error_message")
    11: string processed_at (api.json="processed_at")
    12: string created_at (api.json="created_at")
}

struct ListEventsReq {
    1: string repo_key (api.query="repo_key")
    2: i32 limit (api.query="limit")
}

struct ListEventsResp {
    1: list<WebhookEventInfo> events (api.json="events")
}

struct RetryEventReq {
    1: i64 id (api.query="id")
}

struct RetryEventResp {
    1: bool success (api.json="success")
    2: string message (api.json="message")
}
