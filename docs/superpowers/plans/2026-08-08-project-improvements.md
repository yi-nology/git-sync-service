# 项目改进实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复所有严重、高、中、低优先级问题，提升代码质量、安全性和测试覆盖率

**Architecture:** 按优先级分批处理，从严重问题开始，逐步处理到低优先级问题

**Tech Stack:** Go, Hertz, GORM, Redis, Lua, Docker

## Global Constraints

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查

---

## 文件结构

### 修改的文件

1. `biz/router/git_sync/middleware.go` - 实现 API 认证中间件
2. `internal/lock/lock.go` - 修复 Lock context 和 Semaphore 原子操作
3. `internal/service/webhook.go` - 修复 Goroutine context 问题
4. `internal/service/repo.go` - 定义哨兵错误
5. `internal/dao/repo_dao.go` - 修改 NewRepoDAO 返回错误
6. `biz/handler/git_sync/webhook_receive.go` - 实现 Webhook 速率限制
7. `main.go` - 应用日志配置
8. `internal/executor/executor.go` - 实现临时目录清理
9. `.github/workflows/ci.yml` - 修复 CI 覆盖率条件
10. `go.mod` - 处理 Thrift 强制降级
11. `README.md` - 修复端口错误

### 新增的文件

1. `internal/executor/executor_test.go` - Executor 测试
2. `internal/service/task_test.go` - Task Service 测试
3. `internal/service/webhook_test.go` - Webhook Service 测试
4. `internal/service/cron_test.go` - Cron Service 测试
5. `internal/lock/redis_lock_test.go` - Redis Lock 测试
6. `LICENSE` - MIT 许可证文件
7. `Dockerfile` - Docker 构建文件

---

## Task 1: 实现 API 认证中间件

**Files:**
- Modify: `biz/router/git_sync/middleware.go`

**Interfaces:**
- Consumes: `Service.GetAPIKey()` 返回 API Key
- Produces: 认证中间件函数

- [ ] **Step 1: 修改 middleware.go 添加认证逻辑**

```go
package middleware

import (
    "context"
    "net/http"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/yi-nology/git-sync-service/internal/service"
)

func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        apiKey := c.GetHeader("X-API-Key")
        if string(apiKey) != service.GetSyncService().GetAPIKey() {
            c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
            c.Abort()
            return
        }
        c.Next(ctx)
    }
}
```

- [ ] **Step 2: 更新所有中间件函数使用认证**

```go
func _createrepoMw() []app.HandlerFunc {
    return []app.HandlerFunc{AuthMiddleware()}
}
```

- [ ] **Step 3: 测试认证功能**

```bash
go test ./biz/router/git_sync/... -v
```

- [ ] **Step 4: 提交更改**

```bash
git add biz/router/git_sync/middleware.go
git commit -m "feat: implement API authentication middleware"
```

---

## Task 2: 修复 Lock context 值丢失

**Files:**
- Modify: `internal/lock/lock.go`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的 Lock/Unlock 功能

- [ ] **Step 1: 分析 context 值传递逻辑**

```go
// 当前代码（有问题）
func (l *LocalLock) LockWithTTL(ctx context.Context, key string, value string, ttl time.Duration) (bool, context.Context, error) {
    ctx2 := context.WithValue(ctx, lockValueKey{}, value)
    *lockValueFromContext(ctx2) = value  // 这行代码有问题
    // ...
}
```

- [ ] **Step 2: 修复 context 值传递**

```go
func (l *LocalLock) LockWithTTL(ctx context.Context, key string, value string, ttl time.Duration) (bool, context.Context, error) {
    // 创建新的 context 包含锁值
    ctx2 := context.WithValue(ctx, lockValueKey{}, value)
    
    // 尝试获取锁
    l.mu.Lock()
    if l.locks[key] != "" {
        l.mu.Unlock()
        return false, ctx, nil
    }
    l.locks[key] = value
    l.mu.Unlock()
    
    return true, ctx2, nil
}
```

- [ ] **Step 3: 测试锁功能**

```bash
go test ./internal/lock/... -v
```

- [ ] **Step 4: 提交更改**

```bash
git add internal/lock/lock.go
git commit -m "fix: repair Lock context value loss"
```

---

## Task 3: 修复 Goroutine context 取消问题

**Files:**
- Modify: `internal/service/webhook.go`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的 webhook 处理

- [ ] **Step 1: 修改 ReceiveWebhook 函数**

```go
func (s *Service) ReceiveWebhook(ctx context.Context, repoKey string, req *http.Request) (*model.WebhookEvent, error) {
    // ... 前面的代码保持不变 ...
    
    // 使用 context.Background() 替换请求 context
    go s.safeApplyRules(context.Background(), repoKey, whEvent)
    
    return whEvent, nil
}
```

- [ ] **Step 2: 测试 webhook 处理**

```bash
go test ./internal/service/... -v -run TestReceiveWebhook
```

- [ ] **Step 3: 提交更改**

```bash
git add internal/service/webhook.go
git commit -m "fix: use context.Background() for goroutine in ReceiveWebhook"
```

---

## Task 4: 实现 Semaphore 原子操作

**Files:**
- Modify: `internal/lock/lock.go`

**Interfaces:**
- Consumes: 无
- Produces: 原子化的 Semaphore 操作

- [ ] **Step 1: 编写 Lua 脚本**

```go
const semaphoreAcquireScript = `
local key = KEYS[1]
local member = ARGV[1]
local max = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 清理过期的成员
redis.call('ZREMRANGEBYSCORE', key, '-inf', now - 1)

-- 检查当前数量
local current = redis.call('ZCARD', key)
if current >= max then
    return 0
end

-- 添加新成员
redis.call('ZADD', key, now, member)
return 1
`
```

- [ ] **Step 2: 修改 Semaphore.Acquire 方法**

```go
func (s *Semaphore) Acquire(ctx context.Context, key string, value string) (bool, error) {
    script := redis.NewScript(semaphoreAcquireScript)
    result, err := script.Run(ctx, s.client, []string{key}, value, s.max, time.Now().Unix()).Int()
    if err != nil {
        return false, err
    }
    return result == 1, nil
}
```

- [ ] **Step 3: 测试并发场景**

```bash
go test ./internal/lock/... -v -run TestSemaphore
```

- [ ] **Step 4: 提交更改**

```bash
git add internal/lock/lock.go
git commit -m "feat: implement atomic Semaphore operations with Lua script"
```

---

## Task 5: 修复字符串比较错误

**Files:**
- Modify: `internal/service/repo.go`
- Modify: `biz/handler/git_sync/repo_service.go`

**Interfaces:**
- Consumes: 无
- Produces: 类型化错误

- [ ] **Step 1: 定义哨兵错误**

```go
// internal/service/errors.go
package service

import "errors"

var (
    ErrRepoNotFound = errors.New("repo not found")
    ErrTaskNotFound = errors.New("task not found")
)
```

- [ ] **Step 2: 修改错误处理逻辑**

```go
// internal/service/repo.go
func (s *Service) GetRepo(ctx context.Context, key string) (*model.Repo, error) {
    repo, err := s.repoDAO.FindByKey(key)
    if err != nil {
        return nil, err
    }
    if repo == nil {
        return nil, ErrRepoNotFound
    }
    return repo, nil
}
```

- [ ] **Step 3: 更新 handler 使用 errors.Is**

```go
// biz/handler/git_sync/repo_service.go
if errors.Is(err, service.ErrRepoNotFound) {
    c.JSON(http.StatusNotFound, map[string]string{"error": "repo not found"})
    return
}
```

- [ ] **Step 4: 测试错误处理**

```bash
go test ./internal/service/... -v -run TestGetRepo
```

- [ ] **Step 5: 提交更改**

```bash
git add internal/service/errors.go internal/service/repo.go biz/handler/git_sync/repo_service.go
git commit -m "feat: use typed sentinel errors instead of string comparison"
```

---

## Task 6: 修改 NewRepoDAO 返回错误

**Files:**
- Modify: `internal/dao/repo_dao.go`
- Modify: `internal/service/service.go`

**Interfaces:**
- Consumes: 无
- Produces: 返回错误的 NewRepoDAO

- [ ] **Step 1: 修改 NewRepoDAO 函数签名**

```go
func NewRepoDAO(db *gorm.DB) (*RepoDAO, error) {
    cm, err := credential.NewCryptoManager()
    if err != nil {
        return nil, fmt.Errorf("failed to create CryptoManager: %w", err)
    }
    return &RepoDAO{db: db, cm: cm}, nil
}
```

- [ ] **Step 2: 更新所有调用点**

```go
// internal/service/service.go
repoDAO, err := dao.NewRepoDAO(db)
if err != nil {
    return nil, fmt.Errorf("init repo DAO failed: %w", err)
}
```

- [ ] **Step 3: 测试 DAO 功能**

```bash
go test ./internal/dao/... -v
```

- [ ] **Step 4: 提交更改**

```bash
git add internal/dao/repo_dao.go internal/service/service.go
git commit -m "refactor: NewRepoDAO returns error instead of panicking"
```

---

## Task 7: 添加核心功能测试

**Files:**
- Create: `internal/executor/executor_test.go`
- Create: `internal/service/task_test.go`
- Create: `internal/service/webhook_test.go`
- Create: `internal/service/cron_test.go`
- Create: `internal/lock/redis_lock_test.go`

**Interfaces:**
- Consumes: 无
- Produces: 测试覆盖

- [ ] **Step 1: 创建 Executor 测试**

```go
// internal/executor/executor_test.go
package executor

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewExecutor(t *testing.T) {
    // 测试 Executor 初始化
}

func TestExecute(t *testing.T) {
    // 测试同步执行
}
```

- [ ] **Step 2: 创建 Task Service 测试**

```go
// internal/service/task_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
    // 测试创建任务
}

func TestUpdateTask(t *testing.T) {
    // 测试更新任务
}
```

- [ ] **Step 3: 创建 Webhook Service 测试**

```go
// internal/service/webhook_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestReceiveWebhook(t *testing.T) {
    // 测试接收 webhook
}

func TestApplyRules(t *testing.T) {
    // 测试应用规则
}
```

- [ ] **Step 4: 创建 Cron Service 测试**

```go
// internal/service/cron_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAddCronJob(t *testing.T) {
    // 测试添加 cron 任务
}

func TestRemoveCronJob(t *testing.T) {
    // 测试移除 cron 任务
}
```

- [ ] **Step 5: 创建 Redis Lock 测试**

```go
// internal/lock/redis_lock_test.go
package lock

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestRedisLock_TryLock(t *testing.T) {
    // 测试 Redis 锁
}
```

- [ ] **Step 6: 运行所有测试**

```bash
go test ./... -v
```

- [ ] **Step 7: 提交更改**

```bash
git add internal/executor/executor_test.go internal/service/task_test.go internal/service/webhook_test.go internal/service/cron_test.go internal/lock/redis_lock_test.go
git commit -m "test: add unit tests for core functionality"
```

---

## Task 8: 实现 Webhook 速率限制

**Files:**
- Modify: `biz/handler/git_sync/webhook_receive.go`

**Interfaces:**
- Consumes: `cfg.Webhook.RateLimit` 配置
- Produces: 速率限制功能

- [ ] **Step 1: 实现速率限制逻辑**

```go
func RateLimitMiddleware(rateLimit int) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 使用 Redis 或内存实现速率限制
        // 如果超过限制，返回 429
    }
}
```

- [ ] **Step 2: 应用速率限制**

```go
func ReceiveWebhook(ctx context.Context, c *app.RequestContext) {
    // 应用速率限制中间件
}
```

- [ ] **Step 3: 测试速率限制**

```bash
go test ./biz/handler/git_sync/... -v -run TestRateLimit
```

- [ ] **Step 4: 提交更改**

```bash
git add biz/handler/git_sync/webhook_receive.go
git commit -m "feat: implement webhook rate limiting"
```

---

## Task 9: 应用日志配置

**Files:**
- Modify: `main.go`

**Interfaces:**
- Consumes: `cfg.Log.Level` 和 `cfg.Log.Format` 配置
- Produces: 配置化的日志

- [ ] **Step 1: 修改 main.go 初始化日志**

```go
func main() {
    // 加载配置
    cfg, err := model.LoadConfig("conf/config.yaml")
    if err != nil {
        panic(err)
    }
    
    // 配置日志
    var logLevel slog.Level
    switch cfg.Log.Level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    }
    
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    }))
    slog.SetDefault(logger)
    
    // ... 其余代码保持不变 ...
}
```

- [ ] **Step 2: 测试日志配置**

```bash
go run main.go
```

- [ ] **Step 3: 提交更改**

```bash
git add main.go
git commit -m "feat: apply log configuration from config file"
```

---

## Task 10: 实现临时目录清理

**Files:**
- Modify: `internal/executor/executor.go`

**Interfaces:**
- Consumes: 无
- Produces: 临时目录清理

- [ ] **Step 1: 修改 Execute 函数添加清理逻辑**

```go
func (e *Executor) Execute(ctx context.Context, task *model.SyncTask, trigger string) (*model.SyncRun, error) {
    // ... 前面的代码保持不变 ...
    
    workDir := e.service.GetTempDir(task.Key)
    if err := os.MkdirAll(workDir, 0o755); err != nil {
        // ... 错误处理 ...
    }
    
    // 添加清理逻辑
    defer func() {
        if err := os.RemoveAll(workDir); err != nil {
            slog.Error("failed to cleanup temp dir", "error", err, "dir", workDir)
        }
    }()
    
    // ... 其余代码保持不变 ...
}
```

- [ ] **Step 2: 测试目录清理**

```bash
go test ./internal/executor/... -v -run TestExecute
```

- [ ] **Step 3: 提交更改**

```bash
git add internal/executor/executor.go
git commit -m "feat: cleanup temporary directories after sync execution"
```

---

## Task 11: 修复 CI 覆盖率条件

**Files:**
- Modify: `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的 CI 配置

- [ ] **Step 1: 修改 CI 配置**

```yaml
- name: Upload coverage
  if: matrix.go-version == '1.26'
  uses: codecov/codecov-action@v3
```

- [ ] **Step 2: 测试 CI 流程**

```bash
# 推送到 GitHub 并检查 CI 运行
```

- [ ] **Step 3: 提交更改**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: fix coverage upload condition to match Go version"
```

---

## Task 12: 处理 Thrift 强制降级

**Files:**
- Modify: `go.mod`

**Interfaces:**
- Consumes: 无
- Produces: 移除 replace 指令

- [ ] **Step 1: 检查 Thrift 代码生成**

```bash
# 检查是否需要重新生成 Thrift 代码
```

- [ ] **Step 2: 重新生成代码（如果需要）**

```bash
# 使用兼容的 thriftgo 版本重新生成代码
```

- [ ] **Step 3: 移除 replace 指令**

```go
// 删除这一行
replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
```

- [ ] **Step 4: 测试编译**

```bash
go build ./...
```

- [ ] **Step 5: 提交更改**

```bash
git add go.mod
git commit -m "chore: remove Thrift forced downgrade"
```

---

## Task 13: 创建 LICENSE 文件

**Files:**
- Create: `LICENSE`

**Interfaces:**
- Consumes: 无
- Produces: MIT 许可证文件

- [ ] **Step 1: 创建 LICENSE 文件**

```
MIT License

Copyright (c) 2026 yi-nology

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 2: 提交更改**

```bash
git add LICENSE
git commit -m "docs: add MIT license file"
```

---

## Task 14: 创建 Dockerfile

**Files:**
- Create: `Dockerfile`

**Interfaces:**
- Consumes: 无
- Produces: Docker 构建文件

- [ ] **Step 1: 创建 Dockerfile**

```dockerfile
# 构建阶段
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 复制构建的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/conf ./conf

# 暴露端口
EXPOSE 8890

# 运行应用
CMD ["./main"]
```

- [ ] **Step 2: 测试 Docker 构建**

```bash
docker build -t git-sync-service .
docker run --rm git-sync-service --help
```

- [ ] **Step 3: 提交更改**

```bash
git add Dockerfile
git commit -m "feat: add Dockerfile for containerized deployment"
```

---

## Task 15: 修复 README 端口错误

**Files:**
- Modify: `README.md`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的文档

- [ ] **Step 1: 修改 README.md 中的端口引用**

```markdown
API 文档地址: http://localhost:8890
```

- [ ] **Step 2: 验证端口配置**

```bash
grep -n "8080" README.md
```

- [ ] **Step 3: 提交更改**

```bash
git add README.md
git commit -m "docs: fix port reference in README (8080 -> 8890)"
```

---

## Task 16: 拆分 Service God Object

**Files:**
- Create: `internal/service/repo_service.go`
- Create: `internal/service/task_service.go`
- Create: `internal/service/webhook_service.go`
- Modify: `internal/service/service.go`

**Interfaces:**
- Consumes: 无
- Produces: 拆分后的服务结构

- [ ] **Step 1: 创建 RepoService**

```go
// internal/service/repo_service.go
package service

type RepoService struct {
    repoDAO    *dao.RepoDAO
    providerMgr *sdkprov.Manager
}

func NewRepoService(repoDAO *dao.RepoDAO, providerMgr *sdkprov.Manager) *RepoService {
    return &RepoService{
        repoDAO:    repoDAO,
        providerMgr: providerMgr,
    }
}
```

- [ ] **Step 2: 创建 TaskService**

```go
// internal/service/task_service.go
package service

type TaskService struct {
    taskDAO *dao.SyncTaskDAO
    runDAO  *dao.SyncRunDAO
}

func NewTaskService(taskDAO *dao.SyncTaskDAO, runDAO *dao.SyncRunDAO) *TaskService {
    return &TaskService{
        taskDAO: taskDAO,
        runDAO:  runDAO,
    }
}
```

- [ ] **Step 3: 创建 WebhookService**

```go
// internal/service/webhook_service.go
package service

type WebhookService struct {
    ruleDAO  *dao.WebhookRuleDAO
    eventDAO *dao.WebhookEventDAO
}

func NewWebhookService(ruleDAO *dao.WebhookRuleDAO, eventDAO *dao.WebhookEventDAO) *WebhookService {
    return &WebhookService{
        ruleDAO:  ruleDAO,
        eventDAO: eventDAO,
    }
}
```

- [ ] **Step 4: 重构 Service 结构体**

```go
// internal/service/service.go
type Service struct {
    config      *Config
    db          *gorm.DB
    repoService *RepoService
    taskService *TaskService
    webhookService *WebhookService
    // ... 其他字段保持不变 ...
}
```

- [ ] **Step 5: 测试重构后的代码**

```bash
go test ./... -v
```

- [ ] **Step 6: 提交更改**

```bash
git add internal/service/repo_service.go internal/service/task_service.go internal/service/webhook_service.go internal/service/service.go
git commit -m "refactor: split Service god object into focused services"
```

---

## 执行选项

**Plan complete and saved to `docs/superpowers/plans/2026-08-08-project-improvements.md`. Two execution options:**

**1. Subagent-Driven (recommended)** - 我为每个任务分发一个新的子代理，任务之间进行审查，快速迭代

**2. Inline Execution** - 在本会话中执行任务，使用 executing-plans 进行批量执行和检查点

**您选择哪种方式？**