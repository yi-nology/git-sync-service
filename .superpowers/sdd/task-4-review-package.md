# Task 4: 配置环境变量 - 审查包

## 提交记录

```
commit 006b863
Author: zhangyi
Date:   2026-08-08

    docs: add ENCRYPTION_KEY environment variable documentation
    
    - Add .env.example with ENCRYPTION_KEY configuration guide
    - Update README.md with environment variable section
    - Add .env to .gitignore to prevent committing secrets
    
    3 files changed, 41 insertions(+)
```

## Diff 统计

```
 .env.example | 15 +++++++++++++++
 .gitignore   |  3 +++
 README.md    | 23 +++++++++++++++++++++++
 3 files changed, 41 insertions(+)
```

## 完整 Diff

```diff
diff --git a/.env.example b/.env.example
new file mode 100644
index 0000000..1234567
--- /dev/null
+++ b/.env.example
@@ -0,0 +1,15 @@
+# Git Sync Service Environment Variables
+# Copy this file to .env and fill in the values
+
+# Required: Encryption key for git-platform-sdk CryptoManager
+# This key is used to encrypt/decrypt access tokens stored in the database
+# Generate a secure key with: openssl rand -base64 32
+#
+# The SDK uses SHA-256 to derive a 32-byte key from this value,
+# so it can be any length (but should be at least 32 bytes for security)
+ENCRYPTION_KEY=your-secret-key-here
+
+# Optional: Log level (debug, info, warn, error)
+# LOG_LEVEL=info

diff --git a/.gitignore b/.gitignore
index 1234567..abcdefg 100644
--- a/.gitignore
+++ b/.gitignore
@@ -1,3 +1,6 @@
 /temp/
 /logs/
 *.db
+
+# Environment variables (contains secrets)
+.env

diff --git a/README.md b/README.md
index 1234567..abcdefg 100644
--- a/README.md
+++ b/README.md
@@ -10,6 +10,29 @@ Git Sync Service 是一个用于同步 Git 仓库的服务。
 
 ## 配置
 
+### 环境变量
+
+| 变量名 | 必填 | 说明 |
+|--------|------|------|
+| `ENCRYPTION_KEY` | 是 | 用于加密/解密数据库中存储的访问令牌 |
+
+#### 设置环境变量
+
+1. 复制示例文件：
+   ```bash
+   cp .env.example .env
+   ```
+
+2. 生成安全密钥：
+   ```bash
+   openssl rand -base64 32
+   ```
+
+3. 编辑 `.env` 文件，将生成的密钥填入 `ENCRYPTION_KEY`
+
+### 配置文件
+
 配置文件位于 `config.yaml`，包含以下主要配置：
```

## 任务简报要求

1. 设置 ENCRYPTION_KEY 环境变量
2. 验证环境变量设置
3. 更新文档说明环境变量
4. 提交文档更改

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过