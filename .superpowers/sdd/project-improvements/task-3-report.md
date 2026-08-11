# Task 3: 修复 Goroutine context 取消问题

## 状态

DONE

## 实现内容

修改了 `internal/service/webhook.go` 中的 `ReceiveWebhook` 函数，将 goroutine 中使用的请求 context 替换为 `context.Background()`。

**修改前:**
```go
go s.safeApplyRules(ctx, repoKey, whEvent)
```

**修改后:**
```go
go s.safeApplyRules(context.Background(), repoKey, whEvent)
```

**问题说明:**
- 原代码使用 HTTP 请求的 context 启动 goroutine
- 当 HTTP 请求完成后，context 会被取消
- 这会导致 goroutine 中的 `applyRules` 被中断，无法完成规则处理
- 使用 `context.Background()` 确保 goroutine 可以独立运行直到完成

## 测试结果

所有现有测试通过：
```
=== RUN   TestCreateRepo_ValidGitHubURL
--- PASS: TestCreateRepo_ValidGitHubURL (0.00s)
=== RUN   TestCreateRepo_ValidGitLabURL
--- PASS: TestCreateRepo_ValidGitLabURL (0.00s)
=== RUN   TestCreateRepo_UnsupportedPlatform
--- PASS: TestCreateRepo_UnsupportedPlatform (0.00s)
=== RUN   TestCreateRepo_InvalidURL
--- PASS: TestCreateRepo_InvalidURL (0.00s)
=== RUN   TestCreateRepo_SSHURL
--- PASS: TestCreateRepo_SSHURL (0.00s)
=== RUN   TestCryptoManager_Integration
--- PASS: TestCryptoManager_Integration (0.00s)
=== RUN   TestCryptoManager_Direct_EncryptDecrypt_Roundtrip
--- PASS: TestCryptoManager_Direct_EncryptDecrypt_Roundtrip (0.00s)
=== RUN   TestMatchBranch
--- PASS: TestMatchBranch (0.00s)
=== RUN   TestMatchEventType
--- PASS: TestMatchEventType (0.00s)
PASS
ok  	github.com/yi-nology/git-sync-service/internal/service	0.296s
```

## 提交记录

```
commit 1e45bec
Author: zhangyi
Date:   Mon Aug 11 2026

    fix: use context.Background() for goroutine in ReceiveWebhook
```

## 关注点

无。修改简单明确，不影响现有功能。
