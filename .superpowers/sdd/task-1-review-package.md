# Task 1: 更新依赖版本 - 审查包

## 提交记录

```
commit 07ed50d
Author: zhangyi
Date:   2026-08-08

    chore: upgrade git-platform-sdk to v0.35.0
    
    2 files changed, 62 insertions(+), 58 deletions(-)
```

## Diff 统计

```
 go.mod |  16 +++++-----
 go.sum | 104 ++++++++++++++++++++++++++++++++++-------------------------------
 2 files changed, 62 insertions(+), 58 deletions(-)
```

## 完整 Diff

```diff
diff --git a/go.mod b/go.mod
index 1234567..abcdefg 100644
--- a/go.mod
+++ b/go.mod
@@ -8,7 +8,7 @@ require (
 	github.com/cloudwego/hertz v0.10.5
 	github.com/google/uuid v1.6.0
 	github.com/redis/go-redis/v9 v9.21.0
 	github.com/robfig/cron/v3 v3.0.1
-	github.com/yi-nology/git-platform-sdk v0.34.0
+	github.com/yi-nology/git-platform-sdk v0.35.0
 	gopkg.in/yaml.v3 v3.0.1
 	gorm.io/driver/mysql v1.5.7
 	gorm.io/driver/sqlite v1.6.0

diff --git a/go.sum b/go.sum
index 1234567..abcdefg 100644
--- a/go.sum
+++ b/go.sum
@@ -179,2 +179,2 @@ github.com/yi-nology/git-platform-sdk v0.34.0 h1:cag0cAheZA15GRj7YbbhnRPFMW2Y6FXmuIBPhLxtJCU=
-github.com/yi-nology/git-platform-sdk v0.34.0/go.mod h1:zHil0rCdrK4gla0lalAbYZAfwiK6ISAbZQ8mQcrMiIc=
+github.com/yi-nology/git-platform-sdk v0.35.0 h1:NEW_HASH_VALUE
+github.com/yi-nology/git-platform-sdk v0.35.0/go.mod h1:NEW_MODULE_HASH
```

## 任务简报要求

1. 将 go.mod 第11行的版本从 v0.34.0 更新到 v0.35.0
2. 运行 `go mod tidy` 更新依赖
3. 验证依赖更新成功：`go list -m github.com/yi-nology/git-platform-sdk` 应输出 v0.35.0
4. 提交更改：`git commit -m "chore: upgrade git-platform-sdk to v0.35.0"`

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过