#!/bin/sh

set -e

[[ -z "$APP" ]] && echo "missing \$APP arguments" && exit 1

binary_name="server"

# output message
echo "start build ${APP}"
    mkdir -p ./build/bin/${APP}/certs/
    go build -v -o build/bin/${APP}/${binary_name} ./cmd/${APP}/
    cp ./build/config/${APP}.yaml ./build/bin/${APP}/config.yaml
    cp ./build/bin/certs/pudding*.pem ./build/bin/${APP}/certs/

echo "build ${APP} successfully"