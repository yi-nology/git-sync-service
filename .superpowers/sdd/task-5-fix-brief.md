# Task 5 修复: 添加缺失的单元测试

## 审查发现

1. 缺少加密功能的直接测试覆盖：`internal/dao/` 目录下没有任何测试文件
2. 缺少平台检测功能的直接测试覆盖：`internal/service/repo.go` 中 CreateRepo 方法没有测试覆盖
3. 集成测试步骤为空执行：代码库中没有任何带 `//go:build integration` 标签的文件

## 修复要求

### 1. 为 CryptoManager 添加单元测试

创建文件：`internal/dao/repo_dao_test.go`

测试用例：
- 正常加密解密往返
- ENCRYPTION_KEY 未设置时的行为
- 空字符串输入处理

### 2. 为 CreateRepo 添加单元测试

在 `internal/service/` 目录下添加测试

测试用例：
- 有效平台 URL 的正常流程
- 不支持平台的 ErrPlatformNotSupported 返回
- 无效 URL 的错误处理

### 3. 运行所有测试

```bash
go test ./... -v
```

### 4. 提交更改

```bash
git add internal/dao/repo_dao_test.go internal/service/
git commit -m "test: add unit tests for CryptoManager and CreateRepo"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过