#!/bin/sh

scheduler_binary_name="scheduler"

#输出信息
echo "start build scheduler linux_amd64"

	go env -w GOARCH=amd64
	go env -w GOOS=linux
	go build -o ../build/bin/${scheduler_binary_name}_linux_amd64 ../cmd/scheduler/

echo "build scheduler linux_amd64 finished"