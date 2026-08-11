# Task 3: 更新错误处理逻辑 - 审查包

## 提交记录

```
commit f95fcc1
Author: zhangyi
Date:   2026-08-08

    feat: handle ErrPlatformNotSupported in platform detection
    
    1 file changed, 4 insertions(+)
```

## Diff 统计

```
 internal/service/repo.go | 4 ++++
 1 file changed, 4 insertions(+)
```

## 完整 Diff

```diff
diff --git a/internal/service/repo.go b/internal/service/repo.go
index 1234567..abcdefg 100644
--- a/internal/service/repo.go
+++ b/internal/service/repo.go
@@ -3,6 +3,7 @@ package service
 import (
 	"context"
+	"errors"
 	"fmt"
 
 	"github.com/google/uuid"
@@ -23,7 +24,11 @@ func (s *Service) CreateRepo(ctx context.Context, req *model.CreateRepoRequest)
 	result, err := sdkprov.DetectPlatform(req.RemoteURL)
 	if err != nil {
-		return nil, fmt.Errorf("invalid remote URL: %w", err)
+		if errors.Is(err, sdkprov.ErrPlatformNotSupported) {
+			return nil, fmt.Errorf("unsupported platform for URL %s: %w", req.RemoteURL, err)
+		}
+		return nil, fmt.Errorf("invalid remote URL: %w", err)
 	}
 
 	repo := &model.Repo{
```

## 任务简报要求

1. 添加 errors 包导入
2. 更新 CreateRepo 中的错误处理
3. 运行测试验证错误处理
4. 提交更改

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过