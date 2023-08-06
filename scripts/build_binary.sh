#!/bin/bash

set -ex

[[ -z "$APP" ]] && echo "missing \$APP arguments" && exit 1

binary_name="server"

# output message
echo "start build ${APP} binary"
mkdir -p ./dist/bin/${APP}/
go build -v -o dist/bin/${APP}/${binary_name} ./cmd/${APP}/
cp ./build/config/${APP}.yaml ./dist/bin/${APP}/config.yaml

echo "build ${APP} successfully"
