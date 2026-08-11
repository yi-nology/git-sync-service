# Task 12: 处理 Thrift 强制降级

**Files:**
- Modify: `go.mod`

**Interfaces:**
- Consumes: 无
- Produces: 移除 replace 指令

## 步骤

### Step 1: 检查 Thrift 代码生成

```bash
# 检查是否需要重新生成 Thrift 代码
```

### Step 2: 重新生成代码（如果需要）

```bash
# 使用兼容的 thriftgo 版本重新生成代码
```

### Step 3: 移除 replace 指令

```go
// 删除这一行
replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
```

### Step 4: 测试编译

```bash
go build ./...
```

### Step 5: 提交更改

```bash
git add go.mod
git commit -m "chore: remove Thrift forced downgrade"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查