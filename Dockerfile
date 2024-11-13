FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn,direct \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build 

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 bluebell
RUN go build -o bluebell . 

# 上面是打包镜像， 最终需要的是下面的下的镜像， 其中有上面的打包镜像中已经编译好的文件   
# 创建小的容器 ，最终的这个From就是最终的那个镜像
FROM ubuntu:latest 

COPY ./conf /conf

# 从builder镜像中把可执行文件拷贝到当前目录
COPY --from=builder /build/bluebell /

RUN set -eux \
    && apt-get update \
    && apt-get install -y --no-install-recommends  netcat-openbsd

ENTRYPOINT ["/bluebell"]