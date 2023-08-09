#!/bin/bash

set -ex

make clean
make bootstrap
make env/mysql
go test -v -race -covermode=atomic -coverprofile=coverage.txt ./...
go tool cover -func=coverage.txt
