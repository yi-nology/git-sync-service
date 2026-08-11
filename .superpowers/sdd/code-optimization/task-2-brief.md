# Task 2: 修复忽略的错误

**Files:**
- Modify: `internal/service/webhook.go`
- Modify: `internal/service/webhook_service.go`

**问题:** FindEventByEventID 错误被忽略

## 步骤

### Step 1: 处理错误

在 webhook.go 和 webhook_service.go 中，将：
```go
existing, _ := ws.eventDAO.FindByEventID(event.ID)
```

改为：
```go
existing, err := ws.eventDAO.FindByEventID(event.ID)
if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
    return err
}
```

### Step 2: 运行测试

```bash
go test ./internal/service/... -v
```

### Step 3: 提交

```bash
git add internal/service/webhook.go internal/service/webhook_service.go
git commit -m "fix: handle FindEventByEventID errors properly"
```

## 全局约束

- 所有测试必须通过
- 代码必须通过 lint 检查
- 不引入新的 breaking changes