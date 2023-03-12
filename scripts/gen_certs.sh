#!/bin/bash

set -e

go install github.com/cloudflare/cfssl/cmd/cfssl@v1.16.0
go install github.com/cloudflare/cfssl/cmd/cfssljson@v1.16.0
mkdir -p ./build/bin/certs/
cfssl gencert -initca ./scripts/certs/ca-csr.json | cfssljson -bare ./build/bin/certs/ca
# shellcheck disable=SC1101
cfssl gencert \
  -ca=./build/bin/certs/ca.pem \
  -ca-key=./build/bin/certs/ca-key.pem \
  -config=./scripts/certs/ca-config.json \
  -profile=pudding \
  ./scripts/certs/pudding-csr.json | cfssljson -bare ./build/bin/certs/pudding