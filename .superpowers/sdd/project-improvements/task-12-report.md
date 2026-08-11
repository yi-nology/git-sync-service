# Task 12: 处理 Thrift 强制降级 - 实现报告

## 状态

DONE

## 实现内容

### 问题分析

`go.mod` 中存在一个 `replace` 指令，将 `github.com/apache/thrift` 强制从 v0.23.0 降级到 v0.13.0：

```
replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
```

根本原因是：项目中由 thriftgo 0.4.3 生成的 Thrift 代码使用了 v0.13.0 的 API（不包含 `context.Context` 参数），而 `go.mod` 中 require 的 v0.23.0 版本的 `TProtocol` 接口已添加了 `context.Context` 参数，两者不兼容。

### 解决方案

将 `go.mod` 中的 require 版本从 `v0.23.0` 改为 `v0.13.0`，然后移除 `replace` 指令。这样既保持了代码兼容性，又消除了 replace 指令的 hack。

### 修改的文件

- `/Users/zhangyi/my_project/git-sync-service/go.mod`
  - `github.com/apache/thrift v0.23.0` -> `github.com/apache/thrift v0.13.0`
  - 移除了 `replace github.com/apache/thrift => github.com/apache/thrift v0.13.0`

## 测试结果

- `go build ./...` - 编译通过
- `go test ./...` - 全部通过（Redis 相关测试因无 Redis 服务被跳过，属于预期行为）

## 提交记录

```
85d448d chore: remove Thrift forced downgrade
```

## 关注点

### 版本锁定说明

当前方案将 thrift 锁定在 v0.13.0。这是一个较旧的版本。如果未来需要升级到更新的 thrift 版本，需要同时使用支持 context.Context 的 thriftgo 版本重新生成所有 Thrift 代码（涉及 `biz/model/` 下的 5 个包）。Makefile 中的 `generate` 目标也需要更新（当前使用了错误的 `--go` 标志，应为 `-g go`）。

### thriftgo 兼容性

thriftgo 0.4.3 生成的代码默认不包含 `context.Context` 参数，与 thrift >= v0.14.0 不兼容。Apache Thrift 从 v0.14.0 开始在 `TProtocol` 接口中添加了 `context.Context` 参数。
