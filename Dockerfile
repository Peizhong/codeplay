ARG APP_NAME=codeplay

FROM golang:1.21-alpine AS builder
ARG APP_NAME
ENV CGO_ENABLED=0 GOPROXY="https://goproxy.cn,direct"
RUN apk add git
RUN go install github.com/google/gops@latest
WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build \
	-ldflags="-s -w -X 'github.com/peizhong/codeplay/cmd.BuildDate=`date +%s`' -X 'github.com/peizhong/codeplay/cmd.GitBranch=`git branch --show-current`' -X 'github.com/peizhong/codeplay/cmd.GitCommit=`git rev-parse --short HEAD`'" \
    -o $APP_NAME main.go

FROM alpine
ARG APP_NAME
WORKDIR /app
COPY --from=builder /build/$APP_NAME /app/$APP_NAME
COPY --from=builder /go/bin/gops /usr/local/bin/gops
CMD /app/$APP_NAME

# sudo docker build -f fast.Dockerfile -t registry.cn-shenzhen.aliyuncs.com/peizhong/codeplay:v0.0.1 .