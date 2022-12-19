FROM golang:1.18.2-alpine AS gobuild

MAINTAINER Sato Kenta <g.g.satken@gmail.com>

COPY main.go /

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /redis-sentinel-proxy -a -installsuffix cgo /main.go

FROM busybox

ENV LISTEN_ADDRESS=:6379
ENV SENTINEL_ADDRESS=sentinel:26379
ENV REDIS_MASTER_NAME=master

COPY --from=gobuild /redis-sentinel-proxy /usr/local/bin/redis-sentinel-proxy

WORKDIR /usr/local/bin

CMD redis-sentinel-proxy \
      -listen $LISTEN_ADDRESS \
      -sentinel $SENTINEL_ADDRESS \
      -master $REDIS_MASTER_NAME