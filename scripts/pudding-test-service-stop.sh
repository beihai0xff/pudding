#!/bin/bash

set -e

SRC_DIR=$(git rev-parse --show-toplevel)
cd $SRC_DIR

IMAGE_NAME=pudding-go-test:latest

if [[ -f /.dockerenv ]]; then
    # When running tests inside docker
    /pulsar/bin/pulsar-daemon stop standalone
else
    docker kill pulsar-client-go-test
fi

echo "Stopped Test Pulsar Service"
