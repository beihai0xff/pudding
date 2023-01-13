#!/bin/bash

rm -rf api/gen && rm -rf api/http-spec
cd api/protobuf-spec && buf mod update
cd ../.. && buf generate