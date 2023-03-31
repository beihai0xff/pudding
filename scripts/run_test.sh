#!/bin/bash

set -e -x

make clean
make bootstrap
make gen/mock
make gen/struct_tag
make env/mysql
go test -v -race -covermode=atomic -coverprofile=coverprofile.cov ./...
go tool cover -func=coverprofile.cov