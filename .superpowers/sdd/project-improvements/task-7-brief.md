# Task 7: 添加核心功能测试

**Files:**
- Create: `internal/executor/executor_test.go`
- Create: `internal/service/task_test.go`
- Create: `internal/service/webhook_test.go`
- Create: `internal/service/cron_test.go`
- Create: `internal/lock/redis_lock_test.go`

**Interfaces:**
- Consumes: 无
- Produces: 测试覆盖

## 步骤

### Step 1: 创建 Executor 测试

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

### Step 2: 创建 Task Service 测试

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

### Step 3: 创建 Webhook Service 测试

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

### Step 4: 创建 Cron Service 测试

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

### Step 5: 创建 Redis Lock 测试

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

### Step 6: 运行所有测试

```bash
go test ./... -v
```

### Step 7: 提交更改

```bash
git add internal/executor/executor_test.go internal/service/task_test.go internal/service/webhook_test.go internal/service/cron_test.go internal/lock/redis_lock_test.go
git commit -m "test: add unit tests for core functionality"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查