#!/usr/bin/env bash

set -e
GOOS=linux go build
docker build -t cnguyen0/zipsvr .
docker push cnguyen0/zipsvr
go clean