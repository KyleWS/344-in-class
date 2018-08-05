#!/usr/bin/env bash
set -e
echo "building go server for Linux..."
#Linux users, execut: CGO_ENABLED=0 go build -a
GOOS=linux go build 
docker build -t drstearns/zipsvr .
docker push drstearns/zipsvr
go clean
