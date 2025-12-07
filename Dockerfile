# 阶段1：编译 Go 代码
FROM golang:1.21-alpine AS builder
WORKDIR /app
# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download
# 复制项目代码
COPY . .
# 编译（禁用 CGO 确保静态链接，适配 alpine）
RUN CGO_ENABLED=0 GOOS=linux go build -o swiftchat main.go

# 阶段2：构建轻量镜像
FROM alpine:latest
# 安装必要工具（如 CA 证书，用于 HTTPS 连接）
RUN apk --no-cache add ca-certificates tzdata
# 设置工作目录
WORKDIR /app
# 从编译阶段复制二进制文件
COPY --from=builder /app/swiftchat .
# 复制配置文件、静态资源、视图模板
COPY config/ ./config/
COPY views/ ./views/
COPY asset/ ./asset/
# 创建上传目录并赋予权限
RUN mkdir -p ./asset/upload && chmod 777 ./asset/upload
# 暴露项目端口（与 config/app.yml 中一致）
EXPOSE 8081
# 启动命令
CMD ["./swiftchat"]