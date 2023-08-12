# 使用基于Golang官方镜像
FROM golang:1.20 as builder

# 设置工作目录
WORKDIR /app

# 复制go mod和go sum文件
COPY go.mod go.sum ./

# 下载所有依赖项
RUN go mod download

# 将源代码复制到容器内部
COPY . .

# 编译Go项目
RUN go build -o mydocuments MyDocuments

# 使用Ubuntu 22.04作为基础镜像
FROM ubuntu:22.04

WORKDIR /root/

# 将构建的二进制文件复制到新容器中
COPY --from=builder /app/MyDocuments .

# 对外暴露的端口
EXPOSE 8080

# 更新CA证书
RUN apt-get -qq update && apt-get -qq install -y --no-install-recommends ca-certificates curl

# 容器启动时执行的命令
CMD ["./mydocuments"]
