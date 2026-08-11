# Task 4: 实现 Semaphore 原子操作

**Files:**
- Modify: `internal/lock/lock.go`

**Interfaces:**
- Consumes: 无
- Produces: 原子化的 Semaphore 操作

## 步骤

### Step 1: 编写 Lua 脚本

```go
const semaphoreAcquireScript = `
local key = KEYS[1]
local member = ARGV[1]
local max = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 清理过期的成员
redis.call('ZREMRANGEBYSCORE', key, '-inf', now - 1)

-- 检查当前数量
local current = redis.call('ZCARD', key)
if current >= max then
    return 0
end

-- 添加新成员
redis.call('ZADD', key, now, member)
return 1
`
```

### Step 2: 修改 Semaphore.Acquire 方法

```go
func (s *Semaphore) Acquire(ctx context.Context, key string, value string) (bool, error) {
    script := redis.NewScript(semaphoreAcquireScript)
    result, err := script.Run(ctx, s.client, []string{key}, value, s.max, time.Now().Unix()).Int()
    if err != nil {
        return false, err
    }
    return result == 1, nil
}
```

### Step 3: 测试并发场景

```bash
go test ./internal/lock/... -v -run TestSemaphore
```

### Step 4: 提交更改

```bash
git add internal/lock/lock.go
git commit -m "feat: implement atomic Semaphore operations with Lua script"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查