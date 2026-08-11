# Task 1: 实现 API 认证中间件

**Files:**
- Modify: `biz/router/git_sync/middleware.go`

**Interfaces:**
- Consumes: `Service.GetAPIKey()` 返回 API Key
- Produces: 认证中间件函数

## 步骤

### Step 1: 修改 middleware.go 添加认证逻辑

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

### Step 2: 更新所有中间件函数使用认证

```go
func _createrepoMw() []app.HandlerFunc {
    return []app.HandlerFunc{AuthMiddleware()}
}
```

### Step 3: 测试认证功能

```bash
go test ./biz/router/git_sync/... -v
```

### Step 4: 提交更改

```bash
git add biz/router/git_sync/middleware.go
git commit -m "feat: implement API authentication middleware"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查