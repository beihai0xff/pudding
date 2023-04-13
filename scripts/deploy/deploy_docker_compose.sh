#!/usr/bin/env bash

set -ex

echo "starting infra and prom containers..."
ADMIN_USER='admin' ADMIN_PASSWORD='admin' ADMIN_PASSWORD_HASH='$2a$14$1l.IozJx7xQRVmlkEQ32OeEEfP5mRxTpbDTCTcXRqn19gXD8YK1pO' \
	docker compose -f ./deployments/docker-compose/dockprom/docker-compose.yml \
	-f ./deployments/docker-compose/infra.yml \
	-p pudding-infra up --force-recreate --wait --wait-timeout 120

echo "infra and prom containers start successfully"
echo "start pudding containers..."
docker compose -p pudding-service -f ./deployments/docker-compose/pudding.yml up --force-recreate -V --abort-on-container-exit
