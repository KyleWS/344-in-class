#!/usr/bin/env bash
set -e
echo "building linux executable"
GOOS=linux go build
docker build -t drstearns/testserver .
docker push drstearns/testserver
go clean
