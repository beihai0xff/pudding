ADMIN_USER				?= 'admin'
ADMIN_PASSWORD			?= 'admin'
ADMIN_PASSWORD_HASH		?= '$$2a$$14$$1l.IozJx7xQRVmlkEQ32OeEEfP5mRxTpbDTCTcXRqn19gXD8YK1pO'

aa := $(SHELL pwd)

HOST_IP					?= $(shell ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $$2}'|tr -d "addr:")

.PHONY: deploy/docker-compose
deploy/docker-compose: deploy/docker-compose-infra deploy/docker-compose-pudding

.PHONY: deploy/docker-compose-infra
deploy/docker-compose-infra:
	@echo "starting pudding prom containers..."
	ADMIN_USER=$(ADMIN_USER) ADMIN_PASSWORD=$(ADMIN_PASSWORD) ADMIN_PASSWORD_HASH=$(ADMIN_PASSWORD_HASH) \
		docker compose -f ./deployments/docker-compose/dockprom/docker-compose.yml \
		-f ./deployments/docker-compose/dockprom-override.yml \
		-p pudding-prom up --wait --wait-timeout 120
	@echo "pudding prom containers start successfully"
	@echo "starting pudding infra containers..."
	HOST_IP=$(HOST_IP) docker compose -f ./deployments/docker-compose/infra.yml \
		-p pudding-infra up --wait --wait-timeout 120
	@echo "pudding infra containers start successfully"

.PHONY: deploy/docker-compose-pudding
deploy/docker-compose-pudding:
	@echo "starting pudding service containers..."
	docker compose -f ./deployments/docker-compose/pudding.yml \
		-p pudding-service up --force-recreate -V --abort-on-container-exit
	@echo "pudding service containers start successfully"