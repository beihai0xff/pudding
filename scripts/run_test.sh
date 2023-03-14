#!/bin/bash

set -e -x

make clean
make bootstrap
make gen/mock
make env/mysql
go test -v -covermode=count -coverprofile=coverprofile.cov ./...
go tool cover -func=coverprofile.cov