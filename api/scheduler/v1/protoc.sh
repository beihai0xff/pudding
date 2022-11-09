#!/bin/sh

protoc *.proto --go_out=. --go-grpc_out=.