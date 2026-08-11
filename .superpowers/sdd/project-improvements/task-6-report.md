# Task 6: 修改 NewRepoDAO 返回错误

**状态:** DONE

## 实现内容

已修改 `NewRepoDAO` 函数，使其返回 `(*RepoDAO, error)` 而不是直接 panic。

### 修改的文件

1. **`internal/dao/repo_dao.go`**
   - 添加了 `"fmt"` 导入
   - 修改 `NewRepoDAO` 函数签名：`func NewRepoDAO(db *gorm.DB) *RepoDAO` -> `func NewRepoDAO(db *gorm.DB) (*RepoDAO, error)`
   - 移除了 `panic(err)`，改为返回 `nil, fmt.Errorf("failed to create CryptoManager: %w", err)`
   - 成功时返回 `&RepoDAO{db: db, cm: cm}, nil`

2. **`internal/service/service.go`**
   - 更新 `NewService` 函数中的调用点
   - 将 `repoDAO: dao.NewRepoDAO(db)` 改为先赋值给变量并检查错误
   - 错误时返回 `nil, fmt.Errorf("init repo DAO failed: %w", err)`

3. **`internal/service/repo_test.go`**
   - 更新 `setupTestService` 函数中的调用点（第38行）
   - 更新 `TestCryptoManager_Integration` 测试中的调用点（第177行）
   - 两处都添加了错误检查：`if err != nil { t.Fatalf("failed to create RepoDAO: %v", err) }`

## 测试结果

所有测试通过：
- `go test ./internal/dao/...` - PASS
- `go test ./internal/service/...` - PASS
- `go test ./...` - PASS
- `go vet ./...` - PASS
- `go build ./...` - PASS

## 提交记录

```
commit 234fcdd
Author: zhangyi
Date: 2026-08-11

    refactor: NewRepoDAO returns error instead of panicking

    3 files changed, 18 insertions(+), 6 deletions(-)
```

## 关注点

无。所有更改符合任务要求，代码通过了所有测试和静态分析。
