namespace go repo

struct RepoInfo {
    1: i64 id (api.json="id")
    2: string key (api.json="key")
    3: string name (api.json="name")
    4: string platform (api.json="platform")
    5: string platform_owner (api.json="platform_owner")
    6: string platform_repo (api.json="platform_repo")
    7: string clone_url (api.json="clone_url")
    8: string ssh_url (api.json="ssh_url")
    9: string default_branch (api.json="default_branch")
    10: string status (api.json="status")
    11: string created_at (api.json="created_at")
}

struct ListReposReq {
    1: i32 page (api.query="page")
    2: i32 page_size (api.query="page_size")
}

struct ListReposResp {
    1: list<RepoInfo> repos (api.json="repos")
    2: i64 total (api.json="total")
}

struct GetRepoReq {
    1: string key (api.query="key")
}

struct GetRepoResp {
    1: RepoInfo repo (api.json="repo")
}

struct CreateRepoReq {
    1: string name (api.json="name")
    2: string remote_url (api.json="remote_url")
    3: string access_token (api.json="access_token")
}

struct CreateRepoResp {
    1: RepoInfo repo (api.json="repo")
}

struct UpdateRepoReq {
    1: string key (api.json="key")
    2: string name (api.json="name")
    3: string access_token (api.json="access_token")
}

struct UpdateRepoResp {
    1: RepoInfo repo (api.json="repo")
}

struct DeleteRepoReq {
    1: string key (api.query="key")
}

struct DeleteRepoResp {
    1: bool success (api.json="success")
}

struct TestConnectionReq {
    1: string key (api.query="key")
}

struct TestConnectionResp {
    1: bool success (api.json="success")
    2: string message (api.json="message")
}

struct ListBranchesReq {
    1: string key (api.query="key")
}

struct ListBranchesResp {
    1: list<string> branches (api.json="branches")
}
