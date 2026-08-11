# Task 7: 最终验证和清理

**Files:**
- 所有修改的文件

**Interfaces:**
- 依赖: 所有之前的任务
- 提供: 最终验证

## 步骤

### Step 1: 运行代码质量检查

```bash
go vet ./...
go mod tidy
```

### Step 2: 运行 linter

```bash
golangci-lint run
```

### Step 3: 验证所有测试通过

```bash
go test ./... -v
```

### Step 4: 创建升级总结文档

创建 `UPGRADE_SUMMARY.md` 文档

### Step 5: 提交最终更改

```bash
git add .
git commit -m "chore: complete git-platform-sdk upgrade to v0.35.0"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过