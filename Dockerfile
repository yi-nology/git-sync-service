# 构建阶段
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 安装依赖（包括 gcc 和 musl-dev 用于 CGO）
RUN apk add --no-cache git gcc musl-dev

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（启用 CGO 以支持 SQLite）；版本号取构建时间戳,可被 --build-arg VERSION= 覆盖
ARG VERSION=docker-dev
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X github.com/yi-nology/git-sync-service/internal/version.Version=${VERSION}" \
    -o main .

# 运行阶段
FROM alpine:3.20

WORKDIR /app

# 安装 ca-certificates 和 SQLite 运行时库
RUN apk --no-cache add ca-certificates sqlite-libs wget

# 创建非 root 用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 复制构建的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/conf ./conf

# 创建数据目录并设置权限
RUN mkdir -p /app/data && chown appuser:appgroup /app/data

# 声明数据卷
VOLUME /app/data

# 暴露端口
EXPOSE 8890

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8890/ping || exit 1

# 切换到非 root 用户
USER appuser

# 运行应用
CMD ["./main"]
