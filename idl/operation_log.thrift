// Operation log (audit) IDL. Mirrors the converter.OperationLogInfo DTO shape
// (id/time/user/action/resource/details/ip) used by the audit feature.

struct OperationLogInfo {
    1: i64 id (api.json="id")
    2: string time (api.json="time")
    3: string user (api.json="user")
    4: string action (api.json="action")
    5: string resource (api.json="resource")
    6: string details (api.json="details")
    7: string ip (api.json="ip")
}

struct OperationLogStats {
    1: i64 today (api.json="today")
    2: i64 week (api.json="week")
    3: i64 total (api.json="total")
}

struct ListOperationLogsReq {
    1: i32 page (api.query="page")
    2: i32 pageSize (api.query="page_size")
    3: string search (api.query="search")
    4: string action (api.query="action")
    5: string user (api.query="user")
    6: string startDate (api.query="start_date")
    7: string endDate (api.query="end_date")
}

struct ListOperationLogsResp {
    1: list<OperationLogInfo> list (api.json="list")
    2: i32 page (api.json="page")
    3: i32 pageSize (api.json="page_size")
    4: i64 total (api.json="total")
    5: i32 totalPages (api.json="total_pages")
    6: OperationLogStats stats (api.json="stats")
}
