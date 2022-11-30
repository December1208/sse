FROM golang:1.18-alpine AS build

WORKDIR /mnt/bolean
COPY . /mnt/bolean
RUN export GOPROXY=https://proxy.golang.com.cn,direct && \
    go mod vendor && \
    go build -mod vendor

FROM alpine:3.11
WORKDIR /mnt/bolean
COPY --from=build /mnt/bolean /mnt/bolean
ENV LANG='C.UTF-8' LC_ALL='C.UTF-8' TZ='Asia/Shanghai'
CMD ["./sse_demo"]
