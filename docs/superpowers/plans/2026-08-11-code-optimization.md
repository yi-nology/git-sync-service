# 代码优化实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 修复所有 Bug、消除重复代码、清理死代码、修复配置问题

**Architecture:** 按优先级分批处理，从 Bug 开始，逐步处理到配置问题

**Tech Stack:** Go, GORM, Hertz

## Global Constraints

- 所有测试必须通过
- 代码必须通过 lint 检查
- 不引入新的 breaking changes

---

## Task 1: 修复 fetchRepo 认证 Bug

**Files:**
- Modify: `internal/executor/executor.go`

**问题:** fetchRepo 硬编码 AuthNone，私有仓库会失败

- [ ] **Step 1: 修改 fetchRepo 方法签名**

```go
func (e *Executor) fetchRepo(ctx context.Context, dir string, task *model.SyncTask, repo *model.Repo, details *strings.Builder) error {
```

- [ ] **Step 2: 使用认证配置**

```go
_, err := e.backend.Fetch(ctx, gitbackend.FetchOptions{
    RepoPath: dir,
    Remote:   "origin",
    Branches: []string{task.SourceBranch},
    Tags:     task.GitTags,
    Prune:    task.GitPrune,
    Auth:     e.authConfig(repo),
})
```

- [ ] **Step 3: 更新所有调用点**

- [ ] **Step 4: 运行测试**

- [ ] **Step 5: 提交**

---

## Task 2: 修复忽略的错误

**Files:**
- Modify: `internal/service/webhook.go`
- Modify: `internal/service/webhook_service.go`

**问题:** FindEventByEventID 错误被忽略

- [ ] **Step 1: 处理错误**

```go
existing, err := ws.eventDAO.FindByEventID(event.ID)
if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
    return err
}
```

- [ ] **Step 2: 运行测试**

- [ ] **Step 3: 提交**

---

## Task 3: 修复 cleanup goroutine 泄漏

**Files:**
- Modify: `internal/service/cleanup.go`
- Modify: `internal/service/service.go`

**问题:** cleanupTriggerTimes 没有停止机制

- [ ] **Step 1: 添加 done channel**

```go
type Service struct {
    // ... 其他字段
    cleanupDone chan struct{}
}
```

- [ ] **Step 2: 修改 cleanupTriggerTimes**

```go
func (s *Service) cleanupTriggerTimes() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            // cleanup logic
        case <-s.cleanupDone:
            return
        }
    }
}
```

- [ ] **Step 3: 在 Stop 中关闭 channel**

```go
func (s *Service) Stop() {
    close(s.cleanupDone)
    // ... 其他清理逻辑
}
```

- [ ] **Step 4: 运行测试**

- [ ] **Step 5: 提交**

---

## Task 4: 提取 DAO 分页通用函数

**Files:**
- Create: `internal/dao/dao_helper.go`
- Modify: 所有 DAO 文件

**问题:** 分页模式重复 6 次

- [ ] **Step 1: 创建通用分页函数**

```go
func Paginate[T any](db *gorm.DB, page Pagination, dest *[]*T) (int64, error) {
    var total int64
    if err := db.Model(new(T)).Count(&total).Error; err != nil {
        return 0, err
    }
    err := db.Offset(page.Offset).Limit(page.Limit).Order("id DESC").Find(dest).Error
    return total, err
}
```

- [ ] **Step 2: 创建通用 FindByID 函数**

```go
func FindByID[T any](db *gorm.DB, id uint) (*T, error) {
    var result T
    err := db.Where("id = ?", id).First(&result).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &result, nil
}
```

- [ ] **Step 3: 更新所有 DAO 文件使用通用函数**

- [ ] **Step 4: 运行测试**

- [ ] **Step 5: 提交**

---

## Task 5: 简化 Service 委托层

**Files:**
- Modify: `internal/service/repo.go`
- Modify: `internal/service/webhook.go`
- Modify: `internal/service/service.go`

**问题:** 纯委托层添加无价值

- [ ] **Step 1: 移除 repo.go 中的委托方法**

- [ ] **Step 2: 移除 webhook.go 中的委托方法**

- [ ] **Step 3: 在 Service 中嵌入子服务**

```go
type Service struct {
    *RepoService
    *TaskService
    *WebhookService
    // ... 其他字段
}
```

- [ ] **Step 4: 更新所有调用点**

- [ ] **Step 5: 运行测试**

- [ ] **Step 6: 提交**

---

## Task 6: 提取 Converter 通用函数

**Files:**
- Create: `internal/converter/converter_helper.go`
- Modify: 所有 converter 文件

**问题:** 列表转换函数重复

- [ ] **Step 1: 创建通用 MapSlice 函数**

```go
func MapSlice[T any, R any](src []T, fn func(T) R) []R {
    result := make([]R, 0, len(src))
    for _, v := range src {
        result = append(result, fn(v))
    }
    return result
}
```

- [ ] **Step 2: 更新所有 converter 文件**

- [ ] **Step 3: 运行测试**

- [ ] **Step 4: 提交**

---

## Task 7: 简化 Handler 请求验证

**Files:**
- Create: `biz/handler/git_sync/handler_helper.go`
- Modify: 所有 handler 文件

**问题:** 请求验证模式重复

- [ ] **Step 1: 创建通用请求处理函数**

```go
func handleRequest[T any](c *app.RequestContext, validate func(*T) string, fn func(*T) (interface{}, error)) {
    var req T
    if err := c.BindAndValidate(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    if msg := validate(&req); msg != "" {
        response.BadRequest(c, msg)
        return
    }
    result, err := fn(&req)
    if err != nil {
        // handle error
        return
    }
    response.Success(c, result)
}
```

- [ ] **Step 2: 更新所有 handler 文件**

- [ ] **Step 3: 运行测试**

- [ ] **Step 4: 提交**

---

## Task 8: 提取 TaskKeys 解析函数

**Files:**
- Modify: `biz/handler/git_sync/webhook_service.go`

**问题:** TaskKeys 解析重复

- [ ] **Step 1: 创建解析函数**

```go
func parseTaskKeys(s string) []string {
    if s == "" {
        return nil
    }
    parts := strings.Split(s, ",")
    result := make([]string, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p != "" {
            result = append(result, p)
        }
    }
    return result
}
```

- [ ] **Step 2: 更新 CreateRule 和 UpdateRule**

- [ ] **Step 3: 运行测试**

- [ ] **Step 4: 提交**

---

## Task 9: 消除 ReceiveWebhook 重复

**Files:**
- Modify: `internal/service/webhook.go`
- Modify: `internal/service/webhook_service.go`

**问题:** ReceiveWebhook 逻辑重复

- [ ] **Step 1: 让 Service.ReceiveWebhook 调用 WebhookService.ReceiveWebhook**

- [ ] **Step 2: 移除 Service.ReceiveWebhook 中的重复逻辑**

- [ ] **Step 3: 运行测试**

- [ ] **Step 4: 提交**

---

## Task 10: 清理死代码

**Files:**
- Modify: `internal/service/webhook_service.go`
- Modify: `internal/service/task_service.go`
- Modify: `internal/executor/executor.go`

**问题:** 多处死代码

- [ ] **Step 1: 移除 WebhookParser 接口和相关结构体**

- [ ] **Step 2: 移除 RunTaskWithTrigger 回调版本**

- [ ] **Step 3: 移除 SyncPreview 和 Preview 方法**

- [ ] **Step 4: 移除重复的 FindRulesByRepoKey**

- [ ] **Step 5: 移除直接传递的 FindEventByID/FindEventByEventID**

- [ ] **Step 6: 运行测试**

- [ ] **Step 7: 提交**

---

## Task 11: 修复配置问题

**Files:**
- Modify: `biz/handler/git_sync/webhook_receive.go`
- Modify: `internal/service/webhook.go`

**问题:** MaxBodySize 配置未使用、matchEventType 不支持逗号分隔

- [ ] **Step 1: 使用配置的 MaxBodySize**

- [ ] **Step 2: 修改 matchEventType 支持逗号分隔**

```go
func matchEventType(pattern, actual string) bool {
    if pattern == "" || pattern == "*" {
        return true
    }
    for _, p := range strings.Split(pattern, ",") {
        if strings.TrimSpace(p) == actual {
            return true
        }
    }
    return false
}
```

- [ ] **Step 3: 运行测试**

- [ ] **Step 4: 提交**

---

## Task 12: 最终验证

- [ ] **Step 1: 运行所有测试**

```bash
go test ./... -v
```

- [ ] **Step 2: 运行 lint 检查**

```bash
go vet ./...
golangci-lint run
```

- [ ] **Step 3: 提交最终更改**

```bash
git commit -m "chore: code optimization and bug fixes"
```
