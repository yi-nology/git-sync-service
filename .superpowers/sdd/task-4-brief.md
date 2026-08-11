# Task 4: 配置环境变量

**Files:**
- Modify: 无（环境变量配置）

**Interfaces:**
- 依赖: `ENCRYPTION_KEY` 环境变量
- 提供: 加密密钥

## 步骤

### Step 1: 设置 ENCRYPTION_KEY 环境变量

```bash
export ENCRYPTION_KEY="your-secret-key-here"
```

### Step 2: 验证环境变量设置

```bash
echo $ENCRYPTION_KEY
```

### Step 3: 更新文档说明环境变量

在 README.md 中添加环境变量说明

### Step 4: 提交文档更改

```bash
git add README.md
git commit -m "docs: add ENCRYPTION_KEY environment variable documentation"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过