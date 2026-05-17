# Git-Sync-Service 改进计划

> 基于 v0.4.4 SDK 升级的全量改进方案

## 一、SDK 集成优化（消除 ~170 行重复代码）

| # | 当前实现 | SDK 替代 | 涉及文件 |
|---|---------|---------|---------|
| 1 | `rule.go` 全部（`matchBranch`/`matchSinglePattern`/`splitAndTrim`） | `branchfilter.New(pattern).Match(branch)` | `internal/service/rule.go` → 删除 |
| 2 | `executor.go:buildAuthURL` | `credential.NewManager().BuildAuthURL(...)` | `internal/executor/executor.go` |
| 3 | `manager.go` 的 `parsePlatform`+`getPlatformBaseURL`+`getDefaultBaseURL`+`parseCloneURL` 等 7 个函数 | `sdkprov.DetectPlatform(remoteURL)` | `internal/provider/manager.go` |
| 4 | `repo.go:parseRemoteURL` | `sdkprov.DetectPlatform(remoteURL)` | `internal/service/repo.go` |
| 5 | Token 明文存储 | `credential.EncryptGCM`/`DecryptGCM` | `internal/dao/repo_dao.go` |

## 二、P0 安全问题（必须修复）

| # | 问题 | 修复方案 | 涉及文件 |
|---|------|---------|---------|
| 1 | 24 个 API 端点无认证 | 添加 API Key 中间件，webhook receive 端点除外 | `biz/router/git_sync/middleware.go` |
| 2 | AccessToken 明文存储在数据库 | 使用 `credential.EncryptGCM` 加密，读取时 `DecryptGCM` 解密 | `internal/dao/repo_dao.go` |
| 3 | goroutine 无 panic recovery | `go s.applyRules(...)` 包裹 `recover()` | `internal/service/webhook.go` |
| 4 | Redis 锁 Unlock 无所有权校验 | 使用 Lua 脚本：`if redis.call('get',KEYS[1])==ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end` | `internal/lock/lock.go` |
| 5 | 无请求体大小限制 | webhook receive 加 10MB 限制 | `biz/handler/git_sync/webhook_receive.go` |
| 6 | 前端代理端口错误（8888 vs 8890） | 修正 `vite.config.ts` target | `frontend/vite.config.ts` |
| 7 | 无 CORS 中间件 | 添加 CORS 中间件 | `biz/router/git_sync/middleware.go` |

## 三、P1 正确性 Bug

| # | 问题 | 修复方案 | 涉及文件 |
|---|------|---------|---------|
| 8 | `ListTasks("")` 只返回 cron 启用的任务 | 改用 `FindAll()` 返回全部任务 | `internal/service/task.go` |
| 9 | `MinInterval` 防抖未实现 | 使用 `sync.Map` + 时间戳检查 | `internal/service/webhook.go` |
| 10 | `LocalLock` TTL 存了但不检查 | `TryLock` 中检查过期时间 | `internal/lock/lock.go` |
| 11 | `cronEntryIDs` map 并发竞争 | 加 `sync.RWMutex` | `internal/service/service.go` |
| 12 | `parseRemoteURL` 未知平台默认 GitHub | 应返回错误 | `internal/service/repo.go` |
| 13 | Repo Key 截断 `[:8]`，WebhookToken 截断 `[:16]` | 使用完整 UUID | `internal/service/repo.go`, `internal/service/task.go` |
| 14 | `SyncConfigEntry` 未 AutoMigrate | 从 init.go 移除未使用引用 | `sync/model/init.go` |

## 四、P2 生产就绪

| # | 问题 | 修复方案 | 涉及文件 |
|---|------|---------|---------|
| 15 | `fmt.Printf` 日志 | 替换为 `log/slog` | 全局 |
| 16 | 列表无分页 | DAO 加 `Limit`/`Offset` | `internal/dao/*.go` |
| 17 | 无优雅关闭 | `signal.NotifyContext` + `h.Shutdown(ctx)` | `main.go` |
| 18 | Provider 缓存无淘汰 | 加 TTL 过期淘汰 | `internal/provider/manager.go` |
| 19 | `go.mod` 的 `replace` 指令 | 改为正式引用 `v0.4.4` | `go.mod` |
| 20 | 事件/运行记录无限增长 | 加保留策略和清理方法 | `internal/dao/webhook_event_dao.go`, `sync_run_dao.go` |
| 21 | 浅克隆每次从头开始 | 不再 `RemoveAll` 工作目录，增量 fetch | `internal/executor/executor.go` |

## 执行顺序

1. ✅ 创建改进计划文档
2. 更新 go.mod → 引用 v0.4.4
3. 用 SDK 消除重复代码
4. P0 安全修复
5. P1 Bug 修复
6. P2 生产就绪
7. 编译验证
