#!/bin/bash

set -e

[[ -z "$APP" ]] && echo "missing \$APP arguments" && exit 1

binary_name="server"

# output message
echo "start build ${APP}"
    mkdir -p ./dist/bin/${APP}/certs/
    go build -v -o dist/bin/${APP}/${binary_name} ./cmd/${APP}/
    cp ./build/config/${APP}.yaml ./dist/bin/${APP}/config.yaml
    cp ./dist/certs/pudding*.pem ./dist/bin/${APP}/certs/

echo "build ${APP} successfully"