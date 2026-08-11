# Git Platform SDK 升级设计文档

## 概述

本文档描述了将 `git-platform-sdk` 从 v0.34.0 升级到 v0.35.0 的设计方案，包括 breaking changes 分析、修改计划和实现步骤。

## 1. 升级目标

- 从 `git-platform-sdk v0.34.0` 升级到 `v0.35.0`
- 处理所有 breaking changes
- 保持现有功能正常工作
- 确保代码质量和安全性

## 2. Breaking Changes 分析

### 2.1 SSH 默认行为变更

**变更内容**: `BuildSSHCommand` 现在默认启用 `StrictHostKeyChecking`

**影响评估**: 
- 当前代码中没有直接使用 `BuildSSHCommand`，所以这个变更不影响当前代码
- 无需修改

**解决方案**: 无需修改

### 2.2 平台检测变更

**变更内容**: 对于未知 Git 主机，现在返回 `ErrPlatformNotSupported` 错误，而不再是默认当作 GitLab 处理

**影响评估**: 
- 当前代码使用 `sdkprov.DetectPlatform(req.RemoteURL)`，如果遇到未知平台，现在会返回错误
- 需要更新错误处理逻辑

**解决方案**: 
- 在 `internal/service/repo.go` 中捕获 `ErrPlatformNotSupported` 错误
- 返回合适的错误信息给用户

### 2.3 废弃 API 移除

**变更内容**: 删除了 `EncryptGCM/DecryptGCM`、`NewProviderError/WrapProviderError` 等函数

**影响评估**: 
- 当前代码使用了 `credential.EncryptGCM` 和 `credential.DecryptGCM`
- 这是 breaking change，需要替换为新的 API

**解决方案**: 
- 使用新的 `CryptoManager` API 替换旧的加密函数
- 在 `internal/dao/repo_dao.go` 中使用 `CryptoManager`

### 2.4 依赖适配

**变更内容**: 适配 `gongfeng-sdk-go v0.6.0` 的 API 变更

**影响评估**: 
- 可能影响 TencentCode 后端的使用
- 当前代码可能没有使用 TencentCode 后端

**解决方案**: 
- 检查是否有使用 TencentCode 后端
- 如果有则更新代码

## 3. 修改文件列表

### 3.1 go.mod

**文件路径**: `/Users/zhangyi/my_project/git-sync-service/go.mod`

**修改内容**:
- 将 `github.com/yi-nology/git-platform-sdk v0.34.0` 更新为 `github.com/yi-nology/git-platform-sdk v0.35.0`

### 3.2 internal/dao/repo_dao.go

**文件路径**: `/Users/zhangyi/my_project/git-sync-service/internal/dao/repo_dao.go`

**修改内容**:
- 将 `credential.EncryptGCM` 替换为 `CryptoManager.Encrypt`
- 将 `credential.DecryptGCM` 替换为 `CryptoManager.Decrypt`
- 更新错误处理逻辑

**当前代码**:
```go
import "github.com/yi-nology/git-platform-sdk/pkg/credential"

func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
    // ...
    decrypted, err := credential.DecryptGCM(repo.AccessToken)
    // ...
}

func (d *RepoDAO) Create(repo *model.Repo) error {
    encrypted, err := credential.EncryptGCM(repo.AccessToken)
    // ...
}

func (d *RepoDAO) Update(repo *model.Repo) error {
    if repo.AccessToken != "" {
        encrypted, err := credential.EncryptGCM(repo.AccessToken)
        // ...
    }
    // ...
}
```

**修改后代码**:
```go
import "github.com/yi-nology/git-platform-sdk/pkg/credential"

type RepoDAO struct {
    db  *gorm.DB
    cm  *credential.CryptoManager
}

func NewRepoDAO(db *gorm.DB) *RepoDAO {
    cm, err := credential.NewCryptoManager()
    if err != nil {
        // 处理错误
        panic(err)
    }
    return &RepoDAO{db: db, cm: cm}
}

func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
    // ...
    decrypted, err := d.cm.Decrypt(repo.AccessToken)
    // ...
}

func (d *RepoDAO) Create(repo *model.Repo) error {
    encrypted, err := d.cm.Encrypt(repo.AccessToken)
    // ...
}

func (d *RepoDAO) Update(repo *model.Repo) error {
    if repo.AccessToken != "" {
        encrypted, err := d.cm.Encrypt(repo.AccessToken)
        // ...
    }
    // ...
}
```

### 3.3 internal/service/repo.go

**文件路径**: `/Users/zhangyi/my_project/git-sync-service/internal/service/repo.go`

**修改内容**:
- 更新错误处理逻辑以处理 `ErrPlatformNotSupported`
- 改进错误信息

**当前代码**:
```go
func (s *Service) CreateRepo(ctx context.Context, req *model.CreateRepoRequest) (*model.Repo, error) {
    result, err := sdkprov.DetectPlatform(req.RemoteURL)
    if err != nil {
        return nil, fmt.Errorf("invalid remote URL: %w", err)
    }
    // ...
}
```

**修改后代码**:
```go
func (s *Service) CreateRepo(ctx context.Context, req *model.CreateRepoRequest) (*model.Repo, error) {
    result, err := sdkprov.DetectPlatform(req.RemoteURL)
    if err != nil {
        if errors.Is(err, sdkprov.ErrPlatformNotSupported) {
            return nil, fmt.Errorf("unsupported platform for URL %s: %w", req.RemoteURL, err)
        }
        return nil, fmt.Errorf("invalid remote URL: %w", err)
    }
    // ...
}
```

## 4. 实现步骤

### 4.1 更新依赖

1. 更新 `go.mod` 文件，将 `github.com/yi-nology/git-platform-sdk` 版本从 v0.34.0 更新到 v0.35.0
2. 运行 `go mod tidy` 更新依赖

### 4.2 修改加密逻辑

1. 修改 `internal/dao/repo_dao.go`
2. 添加 `CryptoManager` 字段到 `RepoDAO` 结构体
3. 更新 `NewRepoDAO` 函数以初始化 `CryptoManager`
4. 替换所有 `credential.EncryptGCM` 和 `credential.DecryptGCM` 调用

### 4.3 更新错误处理

1. 修改 `internal/service/repo.go`
2. 更新错误处理逻辑以处理 `ErrPlatformNotSupported`
3. 改进错误信息

### 4.4 测试验证

1. 运行单元测试确保所有功能正常工作
2. 测试加密和解密功能
3. 测试平台检测功能
4. 测试错误处理逻辑

## 5. 配置要求

### 5.1 环境变量

需要设置 `ENCRYPTION_KEY` 环境变量，用于 `CryptoManager` 的加密密钥。

**示例**:
```bash
export ENCRYPTION_KEY="your-secret-key-here"
```

### 5.2 密钥要求

- 密钥可以是任意长度
- SDK 会使用 SHA-256 派生 32 字节的密钥
- 建议使用强随机密钥

## 6. 风险评估

### 6.1 数据迁移风险

**风险**: 现有的加密数据可能无法使用新密钥解密

**缓解措施**:
- 在升级前备份数据库
- 测试新密钥是否能解密现有数据
- 如果不能，需要数据迁移脚本

### 6.2 功能兼容性风险

**风险**: 新版本可能有其他未发现的 breaking changes

**缓解措施**:
- 全面运行测试套件
- 逐步升级，先在测试环境验证

### 6.3 性能风险

**风险**: 新的加密实现可能有性能差异

**缓解措施**:
- 进行性能测试
- 监控生产环境性能

## 7. 回滚计划

如果升级失败，可以按以下步骤回滚：

1. 恢复 `go.mod` 文件中的版本号
2. 运行 `go mod tidy` 恢复旧依赖
3. 恢复 `internal/dao/repo_dao.go` 中的旧加密逻辑
4. 恢复 `internal/service/repo.go` 中的旧错误处理逻辑
5. 恢复数据库备份（如果进行了数据迁移）

## 8. 监控和告警

### 8.1 关键指标

- 加密/解密操作的成功率
- 平台检测的错误率
- 应用程序的整体错误率

### 8.2 告警规则

- 加密/解密操作失败率超过 1%
- 平台检测错误率超过 5%
- 应用程序错误率显著上升

## 9. 文档更新

### 9.1 更新 README

- 更新 SDK 版本信息
- 添加 `ENCRYPTION_KEY` 环境变量说明

### 9.2 更新部署文档

- 添加环境变量配置说明
- 添加升级步骤说明

## 10. 时间计划

### 10.1 开发阶段

- 第 1 天：更新依赖和修改加密逻辑
- 第 2 天：更新错误处理逻辑
- 第 3 天：单元测试和集成测试

### 10.2 测试阶段

- 第 4-5 天：功能测试和性能测试
- 第 6 天：安全测试

### 10.3 部署阶段

- 第 7 天：部署到测试环境
- 第 8 天：部署到生产环境

## 11. 总结

本设计方案详细描述了将 `git-platform-sdk` 从 v0.34.0 升级到 v0.35.0 的完整过程。主要修改包括：

1. 更新依赖版本
2. 替换废弃的加密 API
3. 更新错误处理逻辑
4. 配置环境变量

通过遵循本设计方案，可以安全、高效地完成升级，同时保持现有功能的正常工作。