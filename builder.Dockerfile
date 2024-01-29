FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN  CGO_ENABLED=0 GOPROXY="https://goproxy.cn,direct" go mod tidy

# sudo docker build -f builder.Dockerfile -t codeplay:v0-builder .