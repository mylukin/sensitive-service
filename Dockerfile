# 使用最新的官方Go镜像作为基础镜像
FROM golang:1.19-alpine

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖项到本地
RUN go mod download

# 复制源代码
COPY . .

# 将所有依赖项打包到 vendor 文件夹
RUN go mod vendor

# 使用本地的 vendor 目录编译应用程序
RUN go build -mod=vendor -o main .

# 使用一个更小的基础镜像来运行应用程序
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制编译后的可执行文件到新镜像中
COPY --from=0 /app/main .

# 复制 vendor 文件夹到新镜像中（如果你的程序在运行时需要这些文件）
COPY --from=0 /app/vendor ./vendor

# 安装 ca-certificates 以确保可以进行 HTTPS 请求
RUN apk --no-cache add ca-certificates

# 暴露应用程序的端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]
