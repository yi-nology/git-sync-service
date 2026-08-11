# Task 2 Report: 修改加密逻辑

**Status:** DONE

## 实现内容

修改了 `internal/dao/repo_dao.go` 文件，将旧的 `credential.EncryptGCM/DecryptGCM` API 替换为新的 `CryptoManager` API：

1. **添加 CryptoManager 字段** - 在 `RepoDAO` 结构体中添加了 `cm *credential.CryptoManager` 字段
2. **更新 NewRepoDAO 函数** - 初始化时调用 `credential.NewCryptoManager()` 创建 CryptoManager 实例，如果 `ENCRYPTION_KEY` 环境变量未设置则 panic
3. **替换 FindByKey 中的解密逻辑** - 将 `credential.DecryptGCM(repo.AccessToken)` 替换为 `d.cm.Decrypt(repo.AccessToken)`
4. **替换 Create 中的加密逻辑** - 将 `credential.EncryptGCM(repo.AccessToken)` 替换为 `d.cm.Encrypt(repo.AccessToken)`
5. **替换 Update 中的加密逻辑** - 将 `credential.EncryptGCM(repo.AccessToken)` 替换为 `d.cm.Encrypt(repo.AccessToken)`

## 测试结果

- 编译成功：`go build ./...` 通过
- 所有测试通过：`go test ./...` 通过
- 无残留引用：项目中不再有 `EncryptGCM/DecryptGCM` 的引用

## 提交记录

```
9b22948 feat: replace EncryptGCM/DecryptGCM with CryptoManager
```

## 关注点

- `NewRepoDAO` 中使用 `panic` 处理 `ENCRYPTION_KEY` 未设置的情况，这与任务简报中的设计一致
- 项目中没有 DAO 层的单元测试文件，建议后续补充
