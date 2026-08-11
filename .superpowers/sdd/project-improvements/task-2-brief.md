# Task 2: 修复 Lock context 值丢失

**Files:**
- Modify: `internal/lock/lock.go`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的 Lock/Unlock 功能

## 步骤

### Step 1: 分析 context 值传递逻辑

```go
// 当前代码（有问题）
func (l *LocalLock) LockWithTTL(ctx context.Context, key string, value string, ttl time.Duration) (bool, context.Context, error) {
    ctx2 := context.WithValue(ctx, lockValueKey{}, value)
    *lockValueFromContext(ctx2) = value  // 这行代码有问题
    // ...
}
```

### Step 2: 修复 context 值传递

```go
func (l *LocalLock) LockWithTTL(ctx context.Context, key string, value string, ttl time.Duration) (bool, context.Context, error) {
    // 创建新的 context 包含锁值
    ctx2 := context.WithValue(ctx, lockValueKey{}, value)
    
    // 尝试获取锁
    l.mu.Lock()
    if l.locks[key] != "" {
        l.mu.Unlock()
        return false, ctx, nil
    }
    l.locks[key] = value
    l.mu.Unlock()
    
    return true, ctx2, nil
}
```

### Step 3: 测试锁功能

```bash
go test ./internal/lock/... -v
```

### Step 4: 提交更改

```bash
git add internal/lock/lock.go
git commit -m "fix: repair Lock context value loss"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查