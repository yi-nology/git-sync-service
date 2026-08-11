# Task 5 修复审查包: 添加缺失的单元测试

## 提交记录

```
commit 70c5a0d
Author: zhangyi
Date:   2026-08-08

    test: add unit tests for CryptoManager and CreateRepo
    
    2 files changed, 368 insertions(+)
```

## Diff 统计

```
 internal/dao/repo_dao_test.go | 135 ++++++++++++++++++++++++
 internal/service/repo_test.go | 233 ++++++++++++++++++++++++++++++++++++++++++
 2 files changed, 368 insertions(+)
```

## 测试覆盖

### CryptoManager 单元测试 (`internal/dao/repo_dao_test.go`)

| 测试 | 描述 |
|------|------|
| `TestCryptoManager_EncryptDecrypt_Roundtrip` | 正常加密解密往返，覆盖普通文本、特殊字符、Unicode、长字符串 |
| `TestCryptoManager_EmptyString` | 空字符串输入处理 |
| `TestCryptoManager_WithoutKey` | ENCRYPTION_KEY 未设置时返回错误 |
| `TestCryptoManager_NewCryptoManagerFromKey` | 使用 NewCryptoManagerFromKey 构造函数 |
| `TestDefaultPagination` | 分页参数边界条件（零值、负值、超限） |

### CreateRepo 单元测试 (`internal/service/repo_test.go`)

| 测试 | 描述 |
|------|------|
| `TestCreateRepo_ValidGitHubURL` | GitHub HTTPS URL 正常流程，验证平台检测和加密存储 |
| `TestCreateRepo_ValidGitLabURL` | GitLab HTTPS URL 正常流程 |
| `TestCreateRepo_UnsupportedPlatform` | 不支持的平台返回 `ErrPlatformNotSupported` |
| `TestCreateRepo_InvalidURL` | 无效 URL 的错误处理 |
| `TestCreateRepo_SSHURL` | SSH 格式 URL 的平台检测 |
| `TestCryptoManager_Integration` | 通过 DAO 层的端到端加密解密集成测试 |
| `TestCryptoManager_Direct_EncryptDecrypt_Roundtrip` | CryptoManager 直接测试 |

## 测试结果

全部 12 个测试通过，现有测试无回归。

## 审查发现修复

1. ✅ 添加了 CryptoManager 的直接测试覆盖
2. ✅ 添加了 CreateRepo 的直接测试覆盖
3. ✅ 覆盖了所有关键测试场景

## 任务简报要求

1. 为 CryptoManager 添加单元测试
2. 为 CreateRepo 添加单元测试
3. 运行所有测试
4. 提交更改

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过