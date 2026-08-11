# Task 3: 修复 Goroutine context 取消问题

**Files:**
- Modify: `internal/service/webhook.go`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的 webhook 处理

## 步骤

### Step 1: 修改 ReceiveWebhook 函数

```go
func (s *Service) ReceiveWebhook(ctx context.Context, repoKey string, req *http.Request) (*model.WebhookEvent, error) {
    // ... 前面的代码保持不变 ...
    
    // 使用 context.Background() 替换请求 context
    go s.safeApplyRules(context.Background(), repoKey, whEvent)
    
    return whEvent, nil
}
```

### Step 2: 测试 webhook 处理

```bash
go test ./internal/service/... -v -run TestReceiveWebhook
```

### Step 3: 提交更改

```bash
git add internal/service/webhook.go
git commit -m "fix: use context.Background() for goroutine in ReceiveWebhook"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查