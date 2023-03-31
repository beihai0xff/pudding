#!/bin/bash

set -e

go install github.com/cloudflare/cfssl/cmd/cfssl@latest
go install github.com/cloudflare/cfssl/cmd/cfssljson@latest
mkdir -p ./dist/certs/
cfssl gencert -initca ./scripts/certs/ca-csr.json | cfssljson -bare ./dist/certs/ca
# shellcheck disable=SC1101
cfssl gencert \
  -ca=./dist/certs/ca.pem \
  -ca-key=./dist/certs/ca-key.pem \
  -config=./scripts/certs/ca-config.json \
  -profile=pudding \
  ./scripts/certs/pudding-csr.json | cfssljson -bare ./dist/certs/pudding