# Task 16: 拆分 Service God Object

**Files:**
- Create: `internal/service/repo_service.go`
- Create: `internal/service/task_service.go`
- Create: `internal/service/webhook_service.go`
- Modify: `internal/service/service.go`

**Interfaces:**
- Consumes: 无
- Produces: 拆分后的服务结构

## 步骤

### Step 1: 创建 RepoService

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

### Step 2: 创建 TaskService

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

### Step 3: 创建 WebhookService

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

### Step 4: 重构 Service 结构体

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

### Step 5: 测试重构后的代码

```bash
go test ./... -v
```

### Step 6: 提交更改

```bash
git add internal/service/repo_service.go internal/service/task_service.go internal/service/webhook_service.go internal/service/service.go
git commit -m "refactor: split Service god object into focused services"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查