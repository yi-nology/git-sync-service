# Task 3: 修复 cleanup goroutine 泄漏

**Files:**
- Modify: `internal/service/cleanup.go`
- Modify: `internal/service/service.go`

**问题:** cleanupTriggerTimes 没有停止机制

## 步骤

### Step 1: 添加 done channel

在 Service 结构体中添加 cleanupDone 字段：
```go
type Service struct {
    // ... 其他字段
    cleanupDone chan struct{}
}
```

### Step 2: 初始化 cleanupDone

在 NewService 中初始化：
```go
svc := &Service{
    // ... 其他字段
    cleanupDone: make(chan struct{}),
}
```

### Step 3: 修改 cleanupTriggerTimes

```go
func (s *Service) cleanupTriggerTimes() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            s.CleanupOldEvents(7 * 24 * time.Hour)
        case <-s.cleanupDone:
            return
        }
    }
}
```

### Step 4: 在 Stop 中关闭 channel

```go
func (s *Service) Stop() {
    close(s.cleanupDone)
    s.stopCronJobs()
    if closer, ok := s.lock.(interface{ Close() error }); ok {
        if err := closer.Close(); err != nil {
            slog.Error("failed to close lock", "error", err)
        }
    }
}
```

### Step 5: 运行测试

```bash
go test ./internal/service/... -v
```

### Step 6: 提交

```bash
git add internal/service/cleanup.go internal/service/service.go
git commit -m "fix: prevent cleanup goroutine leak with done channel"
```

## 全局约束

- 所有测试必须通过
- 代码必须通过 lint 检查
- 不引入新的 breaking changes