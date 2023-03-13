#!/bin/bash

set -e -x

make clean
make bootstrap
go test -v ./...