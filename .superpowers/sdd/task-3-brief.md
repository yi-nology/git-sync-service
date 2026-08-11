# Task 3: 更新错误处理逻辑

**Files:**
- Modify: `internal/service/repo.go:22-26`

**Interfaces:**
- 依赖: `sdkprov.ErrPlatformNotSupported` 错误类型
- 提供: 改进的错误信息

## 步骤

### Step 1: 添加 errors 包导入

```go
import (
    "context"
    "errors"
    "fmt"
    // ... 其他导入
)
```

### Step 2: 更新 CreateRepo 中的错误处理

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

### Step 3: 运行测试验证错误处理

```bash
go test ./internal/service/... -v
```

### Step 4: 提交更改

```bash
git add internal/service/repo.go
git commit -m "feat: handle ErrPlatformNotSupported in platform detection"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过