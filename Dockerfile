# 第一阶段：构建阶段
FROM golang:latest AS builder

WORKDIR /app

# 下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源代码并构建静态二进制文件
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o monitor .

# 第二阶段：运行阶段
FROM alpine:latest

WORKDIR /app

# 将第一阶段的二进制文件复制到运行阶段
COPY --from=builder /app/monitor .

# 复制配置文件
COPY config.toml .

# 设置可执行权限并运行
CMD ["./monitor", "-config", "config.toml"]