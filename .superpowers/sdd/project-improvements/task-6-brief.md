# Task 6: 修改 NewRepoDAO 返回错误

**Files:**
- Modify: `internal/dao/repo_dao.go`
- Modify: `internal/service/service.go`

**Interfaces:**
- Consumes: 无
- Produces: 返回错误的 NewRepoDAO

## 步骤

### Step 1: 修改 NewRepoDAO 函数签名

```go
func NewRepoDAO(db *gorm.DB) (*RepoDAO, error) {
    cm, err := credential.NewCryptoManager()
    if err != nil {
        return nil, fmt.Errorf("failed to create CryptoManager: %w", err)
    }
    return &RepoDAO{db: db, cm: cm}, nil
}
```

### Step 2: 更新所有调用点

```go
// internal/service/service.go
repoDAO, err := dao.NewRepoDAO(db)
if err != nil {
    return nil, fmt.Errorf("init repo DAO failed: %w", err)
}
```

### Step 3: 测试 DAO 功能

```bash
go test ./internal/dao/... -v
```

### Step 4: 提交更改

```bash
git add internal/dao/repo_dao.go internal/service/service.go
git commit -m "refactor: NewRepoDAO returns error instead of panicking"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查