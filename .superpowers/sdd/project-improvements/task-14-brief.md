# Task 14: 创建 Dockerfile

**Files:**
- Create: `Dockerfile`

**Interfaces:**
- Consumes: 无
- Produces: Docker 构建文件

## 步骤

### Step 1: 创建 Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 复制构建的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/conf ./conf

# 暴露端口
EXPOSE 8890

# 运行应用
CMD ["./main"]
```

### Step 2: 测试 Docker 构建

```bash
docker build -t git-sync-service .
docker run --rm git-sync-service --help
```

### Step 3: 提交更改

```bash
git add Dockerfile
git commit -m "feat: add Dockerfile for containerized deployment"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查