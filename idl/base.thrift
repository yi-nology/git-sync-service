namespace go base

struct Empty {
}

struct EmptyResp {
    1: string message (api.json="message")
}

struct DeleteReq {
    1: string key (api.query="key")
}

struct DeleteResp {
    1: bool success (api.json="success")
    2: string message (api.json="message")
}

struct PaginationReq {
    1: i32 page (api.query="page")
    2: i32 page_size (api.query="page_size")
}

struct PaginationResp {
    1: i64 total (api.json="total")
    2: i32 page (api.json="page")
    3: i32 page_size (api.json="page_size")
}
