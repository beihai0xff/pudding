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

# output message
echo "start build ${app}"
    cp -r api/http-spec/pudding/${app}/v1/* third_party/OpenAPI/
    make gen_swagger-ui
    go build -v -o build/bin/${binary_name} ./cmd/${app}/

echo "build ${app}"