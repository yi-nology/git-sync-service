# Task 5: 修复字符串比较错误

**Files:**
- Modify: `internal/service/repo.go`
- Modify: `biz/handler/git_sync/repo_service.go`

**Interfaces:**
- Consumes: 无
- Produces: 类型化错误

## 步骤

### Step 1: 定义哨兵错误

```go
// internal/service/errors.go
package service

import "errors"

var (
    ErrRepoNotFound = errors.New("repo not found")
    ErrTaskNotFound = errors.New("task not found")
)
```

### Step 2: 修改错误处理逻辑

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

### Step 3: 更新 handler 使用 errors.Is

```go
// biz/handler/git_sync/repo_service.go
if errors.Is(err, service.ErrRepoNotFound) {
    c.JSON(http.StatusNotFound, map[string]string{"error": "repo not found"})
    return
}
```

### Step 4: 测试错误处理

```bash
go test ./internal/service/... -v -run TestGetRepo
```

### Step 5: 提交更改

```bash
git add internal/service/errors.go internal/service/repo.go biz/handler/git_sync/repo_service.go
git commit -m "feat: use typed sentinel errors instead of string comparison"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查