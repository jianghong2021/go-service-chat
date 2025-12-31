FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./ 

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ./main install && ./main server