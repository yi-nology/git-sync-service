# Task 8: 实现 Webhook 速率限制

**Files:**
- Modify: `biz/handler/git_sync/webhook_receive.go`

**Interfaces:**
- Consumes: `cfg.Webhook.RateLimit` 配置
- Produces: 速率限制功能

## 步骤

### Step 1: 实现速率限制逻辑

```go
func RateLimitMiddleware(rateLimit int) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 使用 Redis 或内存实现速率限制
        // 如果超过限制，返回 429
    }
}
```

### Step 2: 应用速率限制

```go
func ReceiveWebhook(ctx context.Context, c *app.RequestContext) {
    // 应用速率限制中间件
}
```

### Step 3: 测试速率限制

```bash
go test ./biz/handler/git_sync/... -v -run TestRateLimit
```

### Step 4: 提交更改

```bash
git add biz/handler/git_sync/webhook_receive.go
git commit -m "feat: implement webhook rate limiting"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查