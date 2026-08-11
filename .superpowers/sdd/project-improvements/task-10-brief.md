# Task 10: 实现临时目录清理

**Files:**
- Modify: `internal/executor/executor.go`

**Interfaces:**
- Consumes: 无
- Produces: 临时目录清理

## 步骤

### Step 1: 修改 Execute 函数添加清理逻辑

```go
func (e *Executor) Execute(ctx context.Context, task *model.SyncTask, trigger string) (*model.SyncRun, error) {
    // ... 前面的代码保持不变 ...
    
    workDir := e.service.GetTempDir(task.Key)
    if err := os.MkdirAll(workDir, 0o755); err != nil {
        // ... 错误处理 ...
    }
    
    // 添加清理逻辑
    defer func() {
        if err := os.RemoveAll(workDir); err != nil {
            slog.Error("failed to cleanup temp dir", "error", err, "dir", workDir)
        }
    }()
    
    // ... 其余代码保持不变 ...
}
```

### Step 2: 测试目录清理

```bash
go test ./internal/executor/... -v -run TestExecute
```

### Step 3: 提交更改

```bash
git add internal/executor/executor.go
git commit -m "feat: cleanup temporary directories after sync execution"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查