# Git Platform SDK 升级实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 `git-platform-sdk` 从 v0.34.0 升级到 v0.35.0，处理所有 breaking changes

**Architecture:** 使用 CryptoManager 替换旧的加密 API，更新错误处理逻辑以处理平台检测变更

**Tech Stack:** Go, git-platform-sdk, GORM, AES-256-GCM

## Global Constraints

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过

---

## 文件结构

### 修改的文件

1. **go.mod** - 更新 SDK 版本
2. **internal/dao/repo_dao.go** - 替换加密 API
3. **internal/service/repo.go** - 更新错误处理逻辑

### 新增的文件

1. **docs/superpowers/plans/2026-08-08-git-platform-sdk-upgrade.md** - 本实现计划

---

## Task 1: 更新依赖版本

**Files:**
- Modify: `go.mod:11`

**Interfaces:**
- 无依赖

- [ ] **Step 1: 更新 go.mod 文件**

```go
// 将第11行的版本从 v0.34.0 更新到 v0.35.0
github.com/yi-nology/git-platform-sdk v0.35.0
```

- [ ] **Step 2: 运行 go mod tidy 更新依赖**

```bash
go mod tidy
```

- [ ] **Step 3: 验证依赖更新成功**

```bash
go list -m github.com/yi-nology/git-platform-sdk
```

Expected: `github.com/yi-nology/git-platform-sdk v0.35.0`

- [ ] **Step 4: 提交更改**

```bash
git add go.mod go.sum
git commit -m "chore: upgrade git-platform-sdk to v0.35.0"
```

---

## Task 2: 修改加密逻辑

**Files:**
- Modify: `internal/dao/repo_dao.go:1-85`

**Interfaces:**
- 依赖: `credential.NewCryptoManager()` 返回 `*CryptoManager`
- 提供: `RepoDAO` 结构体，包含 `cm *credential.CryptoManager` 字段

- [ ] **Step 1: 添加 CryptoManager 字段到 RepoDAO 结构体**

```go
type RepoDAO struct {
    db  *gorm.DB
    cm  *credential.CryptoManager
}
```

- [ ] **Step 2: 更新 NewRepoDAO 函数**

```go
func NewRepoDAO(db *gorm.DB) *RepoDAO {
    cm, err := credential.NewCryptoManager()
    if err != nil {
        panic(err)
    }
    return &RepoDAO{db: db, cm: cm}
}
```

- [ ] **Step 3: 替换 FindByKey 中的解密逻辑**

```go
func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
    var repo model.Repo
    err := d.db.Where("`key` = ?", key).First(&repo).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    decrypted, err := d.cm.Decrypt(repo.AccessToken)
    if err != nil {
        return &repo, nil
    }
    repo.AccessToken = decrypted
    return &repo, nil
}
```

- [ ] **Step 4: 替换 Create 中的加密逻辑**

```go
func (d *RepoDAO) Create(repo *model.Repo) error {
    encrypted, err := d.cm.Encrypt(repo.AccessToken)
    if err != nil {
        return err
    }
    repo.AccessToken = encrypted
    return d.db.Create(repo).Error
}
```

- [ ] **Step 5: 替换 Update 中的加密逻辑**

```go
func (d *RepoDAO) Update(repo *model.Repo) error {
    if repo.AccessToken != "" {
        encrypted, err := d.cm.Encrypt(repo.AccessToken)
        if err != nil {
            return err
        }
        repo.AccessToken = encrypted
    }
    return d.db.Save(repo).Error
}
```

- [ ] **Step 6: 运行测试验证加密逻辑**

```bash
go test ./internal/dao/... -v
```

- [ ] **Step 7: 提交更改**

```bash
git add internal/dao/repo_dao.go
git commit -m "feat: replace EncryptGCM/DecryptGCM with CryptoManager"
```

---

## Task 3: 更新错误处理逻辑

**Files:**
- Modify: `internal/service/repo.go:22-26`

**Interfaces:**
- 依赖: `sdkprov.ErrPlatformNotSupported` 错误类型
- 提供: 改进的错误信息

- [ ] **Step 1: 添加 errors 包导入**

```go
import (
    "context"
    "errors"
    "fmt"
    // ... 其他导入
)
```

- [ ] **Step 2: 更新 CreateRepo 中的错误处理**

```go
func (s *Service) CreateRepo(ctx context.Context, req *model.CreateRepoRequest) (*model.Repo, error) {
    result, err := sdkprov.DetectPlatform(req.RemoteURL)
    if err != nil {
        if errors.Is(err, sdkprov.ErrPlatformNotSupported) {
            return nil, fmt.Errorf("unsupported platform for URL %s: %w", req.RemoteURL, err)
        }
        return nil, fmt.Errorf("invalid remote URL: %w", err)
    }
    // ... 其余代码保持不变
}
```

- [ ] **Step 3: 运行测试验证错误处理**

```bash
go test ./internal/service/... -v
```

- [ ] **Step 4: 提交更改**

```bash
git add internal/service/repo.go
git commit -m "feat: handle ErrPlatformNotSupported in platform detection"
```

---

## Task 4: 配置环境变量

**Files:**
- Modify: 无（环境变量配置）

**Interfaces:**
- 依赖: `ENCRYPTION_KEY` 环境变量
- 提供: 加密密钥

- [ ] **Step 1: 设置 ENCRYPTION_KEY 环境变量**

```bash
export ENCRYPTION_KEY="your-secret-key-here"
```

- [ ] **Step 2: 验证环境变量设置**

```bash
echo $ENCRYPTION_KEY
```

- [ ] **Step 3: 更新文档说明环境变量**

在 README.md 中添加环境变量说明

- [ ] **Step 4: 提交文档更改**

```bash
git add README.md
git commit -m "docs: add ENCRYPTION_KEY environment variable documentation"
```

---

## Task 5: 全面测试验证

**Files:**
- Test: 所有测试文件

**Interfaces:**
- 依赖: 所有之前的任务
- 提供: 验证所有功能正常工作

- [ ] **Step 1: 运行所有单元测试**

```bash
go test ./... -v
```

- [ ] **Step 2: 运行集成测试**

```bash
go test ./... -v -tags=integration
```

- [ ] **Step 3: 验证加密解密功能**

创建测试脚本验证加密解密功能

- [ ] **Step 4: 验证平台检测功能**

创建测试脚本验证平台检测功能

- [ ] **Step 5: 提交测试结果**

```bash
git add .
git commit -m "test: verify all functionality after SDK upgrade"
```

---

## Task 6: 更新文档和配置

**Files:**
- Modify: `README.md`
- Create: `.env.example`

**Interfaces:**
- 依赖: 所有之前的任务
- 提供: 文档和配置示例

- [ ] **Step 1: 更新 README.md 中的 SDK 版本信息**

```markdown
## Dependencies

- git-platform-sdk v0.35.0
```

- [ ] **Step 2: 创建 .env.example 文件**

```bash
# Encryption key for git-platform-sdk CryptoManager
ENCRYPTION_KEY=your-secret-key-here
```

- [ ] **Step 3: 更新部署文档**

在 README.md 中添加部署说明

- [ ] **Step 4: 提交文档更改**

```bash
git add README.md .env.example
git commit -m "docs: update documentation for SDK upgrade"
```

---

## Task 7: 最终验证和清理

**Files:**
- 所有修改的文件

**Interfaces:**
- 依赖: 所有之前的任务
- 提供: 最终验证

- [ ] **Step 1: 运行代码质量检查**

```bash
go vet ./...
go mod tidy
```

- [ ] **Step 2: 运行 linter**

```bash
golangci-lint run
```

- [ ] **Step 3: 验证所有测试通过**

```bash
go test ./... -v
```

- [ ] **Step 4: 创建升级总结文档**

创建 `UPGRADE_SUMMARY.md` 文档

- [ ] **Step 5: 提交最终更改**

```bash
git add .
git commit -m "chore: complete git-platform-sdk upgrade to v0.35.0"
```

---

## 执行选项

**Plan complete and saved to `docs/superpowers/plans/2026-08-08-git-platform-sdk-upgrade.md`. Two execution options:**

**1. Subagent-Driven (recommended)** - 我为每个任务分发一个新的子代理，任务之间进行审查，快速迭代

**2. Inline Execution** - 在本会话中执行任务，使用 executing-plans 进行批量执行和检查点

**您选择哪种方式？**