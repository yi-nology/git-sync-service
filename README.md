# Git Sync Service

[![CI](https://github.com/yi-nology/git-sync-service/actions/workflows/ci.yml/badge.svg)](https://github.com/yi-nology/git-sync-service/actions/workflows/ci.yml)
[![Release](https://github.com/yi-nology/git-sync-service/actions/workflows/release.yml/badge.svg)](https://github.com/yi-nology/git-sync-service/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yi-nology/git-sync-service)](https://goreportcard.com/report/github.com/yi-nology/git-sync-service)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yi-nology/git-sync-service)](https://go.dev/)
[![License](https://img.shields.io/github/license/yi-nology/git-sync-service)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/yi-nology/git-sync-service)](https://github.com/yi-nology/git-sync-service/releases/tag/v1.5.0)

A service for synchronizing Git repositories across different platforms (GitHub / GitLab / Gitea / Gitee / GitCode / AtomGit / Tencent CodeHub / custom).

## Features

- Multi-platform Git synchronization
- Scheduled sync with cron support
- Webhook-based real-time sync
- RESTful API (Hertz + Thrift IDL-driven via `hz`)
- Operation log (audit) for all mutating actions
- SQLite and MySQL database support
- Shared `X-API-Key` authentication

## Installation

### From Release

Download the latest binary from [Releases](https://github.com/yi-nology/git-sync-service/releases).

### From Source

```bash
git clone https://github.com/yi-nology/git-sync-service.git
cd git-sync-service
go build -o git-sync-service .
```

## Quick Start

```bash
# 1. Required: encryption key for credential storage
export ENCRYPTION_KEY="$(openssl rand -base64 32)"

# 2. Configure conf/config.yaml (see Configuration below)

# 3. Run
./git-sync-service
# or: make run
```

The HTTP server listens on `http://0.0.0.0:8890` by default. The dashboard frontend talks to `/api/v1/*`.

## Configuration

### Environment Variable

| Variable | Required | Description |
|----------|----------|-------------|
| `ENCRYPTION_KEY` | Yes | AES-256-GCM key for encrypting stored credentials (platform access tokens). Min 1 byte; shorter keys are hashed via SHA-256. |

```bash
openssl rand -base64 32     # generate a secure key
export ENCRYPTION_KEY="..." # set it before starting the service
```

### Config File (`conf/config.yaml`)

```yaml
server:
  host: 0.0.0.0
  port: 8890
  mode: debug              # debug | release
  api_key: "your-api-key"  # clients must send this as X-API-Key for /api/v1/*

database:
  driver: sqlite           # sqlite | mysql
  dsn: data/git_sync.db    # sqlite path; or mysql: "user:pass@tcp(127.0.0.1:3306)/git_sync?charset=utf8mb4&parseTime=True"
  max_idle_conns: 10
  max_open_conns: 100

redis:
  addr: 127.0.0.1:6379
  password: ""
  db: 0

git:
  backend: gogit           # gogit | go-git
  temp_dir: /tmp/git-sync

sync:
  max_concurrent: 5
  default_timeout: 300     # seconds
  retry_count: 3

webhook:
  rate_limit: 100          # requests/min per source
  max_body_size: 10485760  # 10 MiB

log:
  level: info              # debug | info | warn | error
  format: json             # json | console
```

## API

All `/api/v1/*` endpoints require the `X-API-Key` header (value = `server.api_key`). `/ping` and `/health` are public. Responses use the envelope:

```json
{ "code": 200, "message": "success", "data": { ... }, "timestamp": 1786530762 }
```

### Endpoints

| Resource | Method | Path | Description |
|---|---|---|---|
| **Platform** | GET | `/api/v1/platforms` | List platforms |
| | GET | `/api/v1/platform?key=` | Get a platform |
| | POST | `/api/v1/platform/create` | Create platform (json) |
| | POST | `/api/v1/platform/update` | Update platform (json) |
| | POST | `/api/v1/platform/delete?key=` | Delete platform |
| | POST | `/api/v1/platform/set-default?key=` | Set default platform |
| | POST | `/api/v1/platform/test?key=` | Test platform connection |
| | GET | `/api/v1/platform/repos?key=&page=&per_page=` | List remote repos on platform |
| | POST | `/api/v1/platform/sync-repos?key=` | Import platform repos locally |
| **Repo** | GET | `/api/v1/repos` | List repos (filters: `search` / `platform` / `status` / `sort_by` / `sort_order`) |
| | GET | `/api/v1/repo?key=` | Get a repo |
| | POST | `/api/v1/repo/create` | Create repo (json) |
| | POST | `/api/v1/repo/update` | Update repo (json) |
| | POST | `/api/v1/repo/delete?key=` | Delete repo |
| | POST | `/api/v1/repo/test?key=` | Test repo connection |
| | GET | `/api/v1/repo/branches?key=` | List branches |
| | POST | `/api/v1/repos/batch` | Batch op (json: `action=delete`, `keys=[...]`) |
| **Sync Task** | GET | `/api/v1/sync/tasks` | List tasks |
| | GET | `/api/v1/sync/task?key=` | Get a task |
| | POST | `/api/v1/sync/task/create` | Create task (json) |
| | POST | `/api/v1/sync/task/update` | Update task (json) |
| | POST | `/api/v1/sync/task/delete?key=` | Delete task |
| | POST | `/api/v1/sync/task/run?key=` | Run a task |
| | POST | `/api/v1/sync/preview` | Preview a sync (json) |
| | GET | `/api/v1/sync/history` | Sync run history |
| **Webhook** | GET | `/api/v1/webhook/rules?repo_key=` | List rules |
| | GET | `/api/v1/webhook/rule?id=` | Get a rule |
| | POST | `/api/v1/webhook/rule/create` | Create rule (json) |
| | POST | `/api/v1/webhook/rule/update` | Update rule (json) |
| | POST | `/api/v1/webhook/rule/delete?id=` | Delete rule |
| | GET | `/api/v1/webhook/events?repo_key=` | List webhook events |
| | POST | `/api/v1/webhook/event/retry?id=` | Retry an event |
| **Audit / System** | GET | `/api/v1/logs/operations` | Operation logs + stats (filters: `search` / `action` / `user` / `start_date` / `end_date`) |
| | GET | `/api/v1/system/status` | System status (version / repo / task counts) |
| **Public / Inbound** | GET | `/ping` | Liveness probe |
| | GET | `/health` | Readiness probe (checks DB & Redis) |
| | POST | `/api/webhook/receive/:repoKey` | Inbound webhook receiver (rate-limited, **no** API-Key) |

### Examples

```bash
# operation logs (with today/week/total stats)
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8890/api/v1/logs/operations?page=1&page_size=10"

# list platforms
curl -H "X-API-Key: your-api-key" http://localhost:8890/api/v1/platforms

# create a platform
curl -X POST -H "X-API-Key: your-api-key" -H "Content-Type: application/json" \
  http://localhost:8890/api/v1/platform/create \
  -d '{"name":"GitHub","type":"github","access_token":"ghp_xxx","is_default":true}'

# system status
curl -H "X-API-Key: your-api-key" http://localhost:8890/api/v1/system/status
```

## Development

```bash
go mod tidy
make test              # run tests
make build             # build
golangci-lint run      # lint (config in .golangci.yml)
```

### Code Generation (Hertz / `hz`)

Routing, models, and handler scaffolds are generated from the Thrift IDL under `idl/`. After editing the IDL, regenerate with the exact toolchain so output matches the committed code:

```bash
hz update -idl idl/git_sync.thrift --snake_tag
```

> Requires `thriftgo v0.4.3` and the `--snake_tag` flag. Newer thriftgo versions change struct-tag style (camelCase instead of snake_case) and would alter the JSON contract — pin to 0.4.3 + `--snake_tag`.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Project Stats

<!-- STATS_START -->
| Metric | Value |
|--------|-------|
| Test Files | 12 |
| Total Tests | 97 |
<!-- STATS_END -->
