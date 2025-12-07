# 阶段 1：编译 Go 代码（仅定义一次 builder 阶段）
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制依赖文件并下载
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 复制项目所有代码
COPY . .

# 编译 Go 程序（禁用 CGO 确保静态链接）
RUN CGO_ENABLED=0 GOOS=linux go build -o swiftchat main.go


# 阶段 2：构建运行镜像（仅引用已定义的 builder 阶段，无循环）
FROM alpine:3.18

# 设置工作目录
WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件（关键：builder 已在阶段 1 定义，无循环）
COPY --from=builder /app/swiftchat .

# 复制项目必需的配置、静态资源和视图
COPY config/ ./config/
COPY views/ ./views/
COPY asset/ ./asset/
COPY --from=builder /app/index.html ./

# 创建上传目录并授权（对应 attach.go 中的上传路径）
RUN mkdir -p ./asset/upload && chmod 777 ./asset/upload

# 暴露端口（与 config/app.yml 中的 port.server 一致）
EXPOSE 8081

# 启动命令
CMD ["./swiftchat"]