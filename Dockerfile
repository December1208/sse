FROM golang:1.18-alpine AS build

WORKDIR /code
COPY . /code
RUN export GOPROXY=https://proxy.golang.com.cn,direct && \
    go mod vendor && \
    go build -mod vendor

FROM alpine:3.11
WORKDIR /code
COPY --from=build /code/sse /code/sse
ENV LANG='C.UTF-8' LC_ALL='C.UTF-8' TZ='Asia/Shanghai'
CMD ["./sse"]
