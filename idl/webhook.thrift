

struct WebhookRuleInfo {
    1: i64 id (api.json="id")
    2: string name (api.json="name")
    3: string repoKey (api.json="repo_key")
    4: string eventType (api.json="event_type")
    5: string branchPattern (api.json="branch_pattern")
    6: string action (api.json="action")
    7: string syncTaskKeys (api.json="sync_task_keys")
    8: i32 minInterval (api.json="min_interval")
    9: bool enabled (api.json="enabled")
    10: string description (api.json="description")
    11: string createdAt (api.json="created_at")
}

struct ListRulesReq {
    1: string repoKey (api.query="repo_key")
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
    2: string repoKey (api.json="repo_key")
    3: string eventType (api.json="event_type")
    4: string branchPattern (api.json="branch_pattern")
    5: string action (api.json="action")
    6: string syncTaskKeys (api.json="sync_task_keys")
    7: i32 minInterval (api.json="min_interval")
    8: bool enabled (api.json="enabled")
    9: string description (api.json="description")
}

struct CreateRuleResp {
    1: WebhookRuleInfo rule (api.json="rule")
}

struct UpdateRuleReq {
    1: i64 id (api.json="id")
    2: string name (api.json="name")
    3: string eventType (api.json="event_type")
    4: string branchPattern (api.json="branch_pattern")
    5: string action (api.json="action")
    6: string syncTaskKeys (api.json="sync_task_keys")
    7: i32 minInterval (api.json="min_interval")
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
    2: string eventId (api.json="event_id")
    3: string repoKey (api.json="repo_key")
    4: string eventType (api.json="event_type")
    5: string source (api.json="source")
    6: string actorName (api.json="actor_name")
    7: string branch (api.json="branch")
    8: string commitSha (api.json="commit_sha")
    9: string status (api.json="status")
    10: string errorMessage (api.json="error_message")
    11: string processedAt (api.json="processed_at")
    12: string createdAt (api.json="created_at")
}

struct ListEventsReq {
    1: string repoKey (api.query="repo_key")
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
