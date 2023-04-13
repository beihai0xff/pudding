WORKSPACE_DIR		?= $(shell pwd)

IMAGE_VERSION 		?= alpha-1
APP 				?= ""
IMAGE_NAME 			= pudding/${APP}:latest

SWAGGER_UI_VERSION	:=v4.15.5


# lint
lint:
	cd api/protobuf-spec && buf mod update && buf lint
	golangci-lint run

test: env/test test/run

.PHONY: clean lint test


# build binary app
build/binary: gen/proto gen/struct_tag gen/swagger-ui
	@echo "build ${APP}"
	APP=$(APP) bash scripts/build_binary.sh

# build docker image
build/docker: clean
	@DOCKER_BUILDKIT=1 docker build \
	--build-arg APP=${APP} \
	-t ${IMAGE_NAME} -f ./build/Dockerfile .

# gen
gen/proto:
	bash -x scripts/gen_proto.sh

gen/struct_tag:
	bash -x scripts/gen_configs_struct_tag.sh

gen/mock:
	bash -x scripts/gen_mock.sh

gen/swagger-ui:
	SWAGGER_UI_VERSION=$(SWAGGER_UI_VERSION) APP=$(APP) bash -x ./scripts/gen_swagger-ui.sh

gen/certs:
	bash -x scripts/gen_certs.sh

.PHONY: build/binary build/docker gen/proto gen/struct_tag gen/mock gen/swagger-ui gen/certs


# clean
clean:
	@echo "clean build dir"
	@rm -rf ./build/bin

.PHONY: clean


# bootstrap
env/mysql:
	go run scripts/init_mysql_env.go

env/dev: bootstrap

# bootstrap the build by downloading additional tools that may be used by devs
bootstrap:
	@go generate -tags tools tools/tools.go

env/test:
	@echo "init unittest docker compose container"
	WORKSPACE_DIR=${WORKSPACE_DIR} \
	docker compose -f ./test/docker-compose.yml up \
	-d --force-recreate --renew-anon-volumes --wait

test/run:
	bash -x ./scripts/run_test.sh

.PHONY: env/dev env/mysql env/test test/run bootstrap

deploy/docker-compose:
	bash scripts/deploy/deploy_docker_compose.sh

.PHONY: deploy/docker-compose