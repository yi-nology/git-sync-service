# Task 1: 更新依赖版本

**Files:**
- Modify: `go.mod:11`

**Interfaces:**
- 无依赖

## 步骤

### Step 1: 更新 go.mod 文件

将第11行的版本从 v0.34.0 更新到 v0.35.0

```go
github.com/yi-nology/git-platform-sdk v0.35.0
```

### Step 2: 运行 go mod tidy 更新依赖

```bash
go mod tidy
```

### Step 3: 验证依赖更新成功

```bash
go list -m github.com/yi-nology/git-platform-sdk
```

Expected: `github.com/yi-nology/git-platform-sdk v0.35.0`

### Step 4: 提交更改

```bash
git add go.mod go.sum
git commit -m "chore: upgrade git-platform-sdk to v0.35.0"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过