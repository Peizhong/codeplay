FROM codeplay:v0.0.1-builder AS builder
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 GOPROXY="https://goproxy.cn,direct"
RUN go build -ldflags "-s -w" -o codeplay ./main.go

FROM alpine
WORKDIR /app
COPY --from=builder /build/codeplay /app/codeplay

ENTRYPOINT ["/app/codeplay"]

# sudo docker build -f app.Dockerfile -t 10.10.10.1:5000/codeplay:v0 .