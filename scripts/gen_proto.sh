#!/bin/bash

set -ex

rm -rf api/gen api/http-spec
cd api/protobuf-spec && buf mod update
cd ../.. && buf generate
