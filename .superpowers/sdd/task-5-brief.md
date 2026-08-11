# Task 5: 全面测试验证

**Files:**
- Test: 所有测试文件

**Interfaces:**
- 依赖: 所有之前的任务
- 提供: 验证所有功能正常工作

## 步骤

### Step 1: 运行所有单元测试

```bash
go test ./... -v
```

### Step 2: 运行集成测试

```bash
go test ./... -v -tags=integration
```

### Step 3: 验证加密解密功能

创建测试脚本验证加密解密功能

### Step 4: 验证平台检测功能

创建测试脚本验证平台检测功能

### Step 5: 提交测试结果

```bash
git add .
git commit -m "test: verify all functionality after SDK upgrade"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过