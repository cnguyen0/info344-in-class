#!/usr/bin/env bash

set -e
GOOS=linux go build
docker build -t cnguyen0/testserver .
docker push cnguyen0/testserver