#!/bin/bash

set -e

SRC_DIR=$(git rev-parse --show-toplevel)
cd $SRC_DIR

IMAGE_NAME=pudding-go-test:latest

if [[ -f /.dockerenv ]]; then
    # When running tests inside docker
    PULSAR_ADMIN=/pulsar/bin/pulsar-admin
    cat /pulsar/conf/standalone.conf
    /pulsar/bin/pulsar-daemon start standalone --no-functions-worker --no-stream-storage
else
    docker build -t ${IMAGE_NAME} .

    docker kill pudding-go-test || true
    docker run -d --rm --name pudding-go-test \
                -p 8080:8080 \
                -p 6650:6650 \
                -p 8443:8443 \
                -p 6651:6651 \
                ${IMAGE_NAME} \
                /pulsar/bin/pulsar standalone \
                    --no-functions-worker --no-stream-storage

    PULSAR_ADMIN="docker exec -it pudding-go-test /pulsar/bin/pulsar-admin"
fi

echo "-- Wait for Pulsar service to be ready"
until curl http://localhost:8080/metrics > /dev/null 2>&1 ; do sleep 1; done

echo "-- Pulsar service is ready -- Configure permissions"

# Create "standalone" cluster
$PULSAR_ADMIN clusters create \
        standalone \
        --url http://localhost:8080/ \
        --url-secure https://localhost:8443/ \
        --broker-url pulsar://localhost:6650/ \
        --broker-url-secure pulsar+ssl://localhost:6651/

# Create "public" tenant
$PULSAR_ADMIN tenants create public -r "anonymous" -c "standalone"

# Create "public/default" with no auth required
$PULSAR_ADMIN namespaces create public/default
$PULSAR_ADMIN namespaces grant-permission public/default \
                        --actions produce,consume \
                        --role "anonymous"

# Create "private" tenant
$PULSAR_ADMIN tenants create private

# Create "private/auth" with required authentication
$PULSAR_ADMIN namespaces create private/auth
$PULSAR_ADMIN namespaces grant-permission private/auth \
                        --actions produce,consume \
                        --role "token-principal"

echo "-- Ready to start tests"
