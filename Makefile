include ./deployments/Makefile

WORKSPACE_DIR		?= $(shell pwd)

TAG 				?= latest
APP 				?= ""
IMAGE_NAME 			= pudding/${APP}:${TAG}

SWAGGER_UI_VERSION	:=v4.15.5


submodule:
	git submodule update --init --recursive


# clean
clean:
	@echo "clean build dir"
	@rm -rf ./build/bin dist
	@rm -rf api/gen
	@rm -rf coverage.txt

# lint
lint: gen/proto gen/mock gen/struct_tag
	cd api/protobuf-spec && buf mod update && buf lint
	golangci-lint run

.PHONY: clean lint submodule


# build binary app
build/binary: gen/proto gen/struct_tag gen/swagger-ui
	@echo "build ${APP}"
	APP=$(APP) bash scripts/build_binary.sh

# build docker image
build/docker: clean
	docker build \
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

install/tools:
	@go generate -x -tags tools tools/tools.go

# bootstrap the build by downloading additional tools that may be used by devs
bootstrap: install/tools gen/proto gen/struct_tag gen/mock


.PHONY: build/binary build/docker gen/proto gen/struct_tag gen/mock gen/swagger-ui install/tools bootstrap


# bootstrap
env/mysql:
	go run scripts/init_mysql_env.go

env/dev: bootstrap

env/test:
	@echo "init unittest docker compose container"
	WORKSPACE_DIR=${WORKSPACE_DIR} \
	docker compose -f ./test/docker-compose.yml up \
	-d --force-recreate --renew-anon-volumes --wait

test/run:
	bash -x ./scripts/run_test.sh


test: env/test test/run

.PHONY: env/dev env/mysql env/test test/run test

