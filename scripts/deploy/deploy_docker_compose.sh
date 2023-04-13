#!/usr/bin/env bash

set -ex

ADMIN_USER='admin' ADMIN_PASSWORD='admin' ADMIN_PASSWORD_HASH='$2a$14$1l.IozJx7xQRVmlkEQ32OeEEfP5mRxTpbDTCTcXRqn19gXD8YK1pO' \
	docker compose -f ./deployments/docker-compose/dockprom/docker-compose.yml \
	-f ./deployments/docker-compose/infra.yml \
	-f ./deployments/docker-compose/pudding.yml \
	-p pudding up --always-recreate-deps -V
