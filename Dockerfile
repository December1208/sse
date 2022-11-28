FROM golang:1.14.9-alpine AS build

COPY . /go/src/video-player
WORKDIR /go/src/video-player
RUN export GOPROXY=http://10.10.20.88:8081/repository/golang-group/ && \
    go mod vendor && \
    go build -mod vendor

FROM alpine:3.11
COPY --from=build /go/src/video-player/video-player /webserver/video-player/
WORKDIR /webserver/video-player
ENV LD_LIBRARY_PATH=/webserver/video-player/lib
ENV LANG='C.UTF-8' LC_ALL='C.UTF-8' TZ='Asia/Shanghai'

RUN echo "nameserver 223.5.5.5" > /etc/resolv.conf \
    && wget http://b.aonn.vip/ffmpeg-release-amd64-static.tar.xz \
    && tar -C /usr/sbin/ -xf ffmpeg-release-amd64-static.tar.xz $(tar -tf ffmpeg-release-amd64-static.tar.xz |grep ffmpeg$) --strip-components=1 \
    && rm ffmpeg-release-amd64-static.tar.xz
CMD ["./video-player"]
