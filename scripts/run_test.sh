#!/bin/bash

set -e -x

make clean
make bootstrap
make gen/mock
go test -v ./...