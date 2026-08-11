# 统一同步任务行为日志设计方案

## 目标

为每次同步执行建立结构化的行为日志系统，记录每个步骤的执行状态、耗时、错误分类，并关联触发事件。

## 当前问题

1. `SyncRun.Details` 是自由文本，无法按步骤查询/过滤
2. `SyncRun.CommitRange` 字段存在但从未填充
3. `WebhookEvent` 与 `SyncRun` 无直接关联，无法追溯触发源
4. 错误信息扁平化，无法区分瞬态/永久性错误
5. 无步骤级耗时统计

## 设计方案

### 1. 新增 `SyncRunStep` 模型

**文件**: `sync/model/sync_run_step.go`

```go
type SyncRunStep struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    RunID       uint      `json:"runId" gorm:"not null;index"`
    StepName    string    `json:"stepName" gorm:"size:50;not null"`   // clone/fetch/ensure_remote/push
    Status      string    `json:"status" gorm:"size:20;not null"`     // running/success/failed
    StartTime   time.Time `json:"startTime"`
    EndTime     *time.Time `json:"endTime"`
    DurationMs  int64     `json:"durationMs"`                         // 预计算耗时(ms)
    ErrorMsg    string    `json:"errorMsg" gorm:"type:text"`
    ErrorType   string    `json:"errorType" gorm:"size:30"`           // auth/network/config/git/unknown
    Output      string    `json:"output" gorm:"type:text"`            // git命令输出摘要
    RetryCount  int       `json:"retryCount"`                         // 该步骤重试次数
    CreatedAt   time.Time `json:"createdAt"`
}
```

步骤名称常量:
- `step_clone` - 首次克隆
- `step_fetch` - 拉取更新
- `step_checkout` - 切换分支
- `step_ensure_remote` - 确保目标远程
- `step_push` - 推送到目标

错误类型常量:
- `error_auth` - 认证失败（token无效/过期）
- `error_network` - 网络错误（超时/连接拒绝）
- `error_config` - 配置错误（仓库不存在/分支不存在）
- `error_git` - Git操作错误（冲突/非快进）
- `error_unknown` - 未知错误

### 2. 增强 `SyncRun` 模型

**文件**: `sync/model/sync_run.go`（修改）

新增字段:
```go
type SyncRun struct {
    // ... 现有字段 ...
    WebhookEventID *uint  `json:"webhookEventId" gorm:"index"`           // 关联的webhook事件
    CommitRange    string `json:"commitRange" gorm:"size:255"`           // 实际同步的提交范围
    DurationMs     int64  `json:"durationMs"`                            // 总耗时(ms)
    ErrorType      string `json:"errorType" gorm:"size:30"`              // 汇总错误类型
    RetryTotal     int    `json:"retryTotal"`                            // 总重试次数
    Steps          []SyncRunStep `json:"steps" gorm:"foreignKey:RunID"`  // 关联步骤
}
```

### 3. 新增错误分类工具函数

**文件**: `internal/executor/error_classify.go`（新建）

```go
func ClassifyError(err error) string {
    // 根据错误内容自动分类:
    // - 包含 "401"/"403"/"authentication"/"token" -> error_auth
    // - 包含 "timeout"/"connection"/"dial"/"network" -> error_network
    // - 包含 "not found"/"does not exist" -> error_config
    // - 包含 "non-fast-forward"/"conflict"/"rejected" -> error_git
    // - 其他 -> error_unknown
}
```

### 4. 重构 Executor 执行流程

**文件**: `internal/executor/executor.go`（修改）

将现有的 `strings.Builder` 方式替换为步骤级记录:

```
Execute(ctx, task, trigger, webhookEventID)
  ├─ CreateRun(task, trigger, webhookEventID)  // 传入关联事件ID
  ├─ Step: clone/fetch
  │   ├─ createStep(run.ID, "clone"/"fetch", StatusRunning)
  │   ├─ 执行操作
  │   ├─ updateStep(step, StatusSuccess/StatusFailed, output/error)
  │   └─ 记录 CommitRange（从fetch结果提取）
  ├─ Step: ensure_remote
  │   └─ (同上)
  ├─ Step: push (含重试循环)
  │   ├─ createStep(run.ID, "push", StatusRunning)
  │   ├─ 每次重试更新 RetryCount
  │   └─ 最终更新状态
  └─ CompleteRun(run)  // 计算 DurationMs, 汇总 ErrorType
```

### 5. 更新 Service 层接口

**文件**: `internal/executor/executor.go`（修改 Service 接口）

```go
type RunManager interface {
    CreateRun(task *model.SyncTask, trigger string, webhookEventID *uint) (*model.SyncRun, error)
    CreateRunStep(step *model.SyncRunStep) error
    UpdateRunStep(step *model.SyncRunStep) error
    CompleteRun(run *model.SyncRun) error
    UpdateTaskLastRun(task *model.SyncTask, run *model.SyncRun) error
}
```

**文件**: `internal/service/task_service.go`（修改）

新增方法:
- `CreateRunStep(step *model.SyncRunStep) error` -> 调用 `runStepDAO.Create`
- `UpdateRunStep(step *model.SyncRunStep) error` -> 调用 `runStepDAO.Update`

### 6. 新增 DAO

**文件**: `internal/dao/sync_run_step_dao.go`（新建）

```go
type SyncRunStepDAO struct {
    db *gorm.DB
}

func (d *SyncRunStepDAO) Create(step *model.SyncRunStep) error
func (d *SyncRunStepDAO) Update(step *model.SyncRunStep) error
func (d *SyncRunStepDAO) FindByRunID(runID uint) ([]*model.SyncRunStep, error)
func (d *SyncRunStepDAO) CleanupOlderThan(olderThan time.Duration) (int64, error)
```

### 7. 更新 Webhook 触发链路

**文件**: `internal/service/webhook.go`（修改）

在 `safeApplyRules` 中传递 `WebhookEvent.ID`:

```go
// 现在的调用
runTaskFn(ctx, taskKey, model.TriggerWebhook)

// 改为
runTaskFn(ctx, taskKey, model.TriggerWebhook, &event.ID)
```

**文件**: `internal/service/task.go`（修改）

```go
func (s *Service) RunTaskWithTrigger(ctx context.Context, taskKey, trigger string, webhookEventID *uint) error {
    // ...
    _, err = s.executor.Execute(ctx, task, trigger, webhookEventID)
    return err
}
```

### 8. 更新 DB Migration

**文件**: `sync/model/init.go`（修改）

在 `AutoMigrate` 中添加 `&SyncRunStep{}`:

```go
db.AutoMigrate(
    &Repo{},
    &SyncTask{},
    &SyncRun{},
    &SyncRunStep{},  // 新增
    &WebhookRule{},
    &WebhookRuleTask{},
    &WebhookEvent{},
)
```

### 9. 更新 Cleanup 逻辑

**文件**: `internal/service/cleanup.go`（修改）

在 `CleanupOldData` 中同步清理 `SyncRunStep`:

```go
func (s *Service) CleanupOldData(ctx context.Context, maxAge time.Duration) (events, runs, steps int64, err error) {
    // 现有逻辑...
    steps, err = s.tasks.CleanupOldRunSteps(maxAge)
    return
}
```

### 10. 更新 API 响应（可选增强）

**文件**: `biz/handler/git_sync/sync_task_service.go`（修改）

`ListHistory` 返回中包含步骤摘要:

```go
type RunInfo struct {
    // ... 现有字段 ...
    Steps      []RunStepInfo `json:"steps"`
    DurationMs int64         `json:"durationMs"`
    ErrorType  string        `json:"errorType"`
}
```

---

## 文件变更清单

| 操作 | 文件 | 说明 |
|------|------|------|
| 新建 | `sync/model/sync_run_step.go` | 步骤日志模型 |
| 新建 | `internal/dao/sync_run_step_dao.go` | 步骤日志 DAO |
| 新建 | `internal/executor/error_classify.go` | 错误分类工具 |
| 修改 | `sync/model/sync_run.go` | 增加关联字段 |
| 修改 | `sync/model/constants.go` | 增加步骤名/错误类型常量 |
| 修改 | `sync/model/init.go` | AutoMigrate 新表 |
| 修改 | `internal/executor/executor.go` | 重构为步骤级记录 |
| 修改 | `internal/service/task_service.go` | 增加步骤 CRUD 方法 |
| 修改 | `internal/service/task.go` | 更新 RunTaskWithTrigger 签名 |
| 修改 | `internal/service/webhook.go` | 传递 WebhookEventID |
| 修改 | `internal/service/webhook_service.go` | 更新 ApplyRules 签名 |
| 修改 | `internal/service/cleanup.go` | 清理步骤日志 |
| 修改 | `internal/service/service.go` | 更新接口适配 |

## 数据库表结构

### sync_run_steps（新增）

| 列名 | 类型 | 说明 |
|------|------|------|
| id | uint (PK) | 自增主键 |
| run_id | uint (FK, index) | 关联 sync_runs.id |
| step_name | varchar(50) | 步骤名称 |
| status | varchar(20) | running/success/failed |
| start_time | datetime | 开始时间 |
| end_time | datetime | 结束时间 |
| duration_ms | bigint | 耗时毫秒 |
| error_msg | text | 错误信息 |
| error_type | varchar(30) | 错误分类 |
| output | text | 操作输出摘要 |
| retry_count | int | 重试次数 |
| created_at | datetime | 创建时间 |

### sync_runs（修改）

| 新增列 | 类型 | 说明 |
|--------|------|------|
| webhook_event_id | uint (nullable, index) | 关联 webhook_events.id |
| duration_ms | bigint | 总耗时毫秒 |
| error_type | varchar(30) | 汇总错误类型 |
| retry_total | int | 总重试次数 |

## 执行顺序

1. 新建模型和常量 (`sync_run_step.go`, `constants.go`)
2. 新建 DAO (`sync_run_step_dao.go`)
3. 修改 `sync_run.go` 增加字段
4. 修改 `init.go` 注册新表
5. 新建错误分类 (`error_classify.go`)
6. 修改 Service 层 (`task_service.go`, `task.go`, `webhook.go`, `webhook_service.go`, `service.go`)
7. 重构 Executor (`executor.go`)
8. 更新 Cleanup (`cleanup.go`)
9. 运行测试验证