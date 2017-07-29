#!/bin/bash

go get .
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o redis-sentinel-proxy .