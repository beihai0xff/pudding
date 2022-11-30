#!/bin/sh

binary_name="server"
app=

while getopts "a:" opt
do
    case $opt in
        a)
        echo "app=$OPTARG"
        app=$OPTARG
        ;;
        ?)
        echo "unknown argument"
    esac
done

#输出信息
echo "start build ${app}"

	go build -v -o ../build/bin/${binary_name} ../cmd/${app}/

echo "build ${app}"