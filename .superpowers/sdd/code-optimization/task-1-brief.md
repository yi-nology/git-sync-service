# Task 1: 修复 fetchRepo 认证 Bug

**Files:**
- Modify: `internal/executor/executor.go`

**问题:** fetchRepo 硬编码 AuthNone，私有仓库会失败

## 步骤

### Step 1: 修改 fetchRepo 方法签名

```go
func (e *Executor) fetchRepo(ctx context.Context, dir string, task *model.SyncTask, repo *model.Repo, details *strings.Builder) error {
```

### Step 2: 使用认证配置

```go
_, err := e.backend.Fetch(ctx, gitbackend.FetchOptions{
    RepoPath: dir,
    Remote:   "origin",
    Branches: []string{task.SourceBranch},
    Tags:     task.GitTags,
    Prune:    task.GitPrune,
    Auth:     e.authConfig(repo),
})
```

### Step 3: 更新所有调用点

在 Execute 方法中，调用 fetchRepo 时传入 repo 参数：
- 第一次调用（fetchRepo）需要传入 sourceRepo
- 重试调用也需要传入 sourceRepo

### Step 4: 运行测试

```bash
go test ./internal/executor/... -v
```

### Step 5: 提交

```bash
git add internal/executor/executor.go
git commit -m "fix: use authentication in fetchRepo for private repositories"
```

## 全局约束

- 所有测试必须通过
- 代码必须通过 lint 检查
- 不引入新的 breaking changes