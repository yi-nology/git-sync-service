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

# 构建应用（启用 CGO 以支持 SQLite）
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates 和 SQLite 运行时库
RUN apk --no-cache add ca-certificates sqlite-libs

# 复制构建的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/conf ./conf

# 暴露端口
EXPOSE 8890

# 运行应用
CMD ["./main"]
