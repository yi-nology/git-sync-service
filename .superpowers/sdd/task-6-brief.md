# Task 6: 更新文档和配置

**Files:**
- Modify: `README.md`
- Create: `.env.example`

**Interfaces:**
- 依赖: 所有之前的任务
- 提供: 文档和配置示例

## 步骤

### Step 1: 更新 README.md 中的 SDK 版本信息

```markdown
## Dependencies

- git-platform-sdk v0.35.0
```

### Step 2: 创建 .env.example 文件

```bash
# Encryption key for git-platform-sdk CryptoManager
ENCRYPTION_KEY=your-secret-key-here
```

### Step 3: 更新部署文档

在 README.md 中添加部署说明

### Step 4: 提交文档更改

```bash
git add README.md .env.example
git commit -m "docs: update documentation for SDK upgrade"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过