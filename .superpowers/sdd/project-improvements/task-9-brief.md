# Task 9: 应用日志配置

**Files:**
- Modify: `main.go`

**Interfaces:**
- Consumes: `cfg.Log.Level` 和 `cfg.Log.Format` 配置
- Produces: 配置化的日志

## 步骤

### Step 1: 修改 main.go 初始化日志

```go
func main() {
    // 加载配置
    cfg, err := model.LoadConfig("conf/config.yaml")
    if err != nil {
        panic(err)
    }
    
    // 配置日志
    var logLevel slog.Level
    switch cfg.Log.Level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    }
    
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    }))
    slog.SetDefault(logger)
    
    // ... 其余代码保持不变 ...
}
```

### Step 2: 测试日志配置

```bash
go run main.go
```

### Step 3: 提交更改

```bash
git add main.go
git commit -m "feat: apply log configuration from config file"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查