IMAGE_VERSION 		?= alpha-1
GO_VERSION 			?= 1.18
GOLANG_IMAGE_NAME 	?= golang:$(GO_VERSION)
APP 				?= ""
IMAGE_NAME 			= pudding.${APP}:latest

SWAGGER_UI_VERSION	:=v4.15.5

# lint
lint:
	cd api/protobuf-spec && buf mod update && buf lint
	golangci-lint run

test: gen/mock
	go test ./...


# build binary app
build/binary: gen/proto gen/struct_tag gen/swagger-ui
	@echo "build ${APP}"
	APP=$(APP) bash scripts/build_binary.sh

# build docker image
build/docker: clean/build
	DOCKER_BUILDKIT=1 docker build \
	--build-arg APP=${APP} \
	-t ${IMAGE_NAME} -f ./build/Dockerfile .

# gen
gen/proto:
	bash scripts/gen_proto.sh

gen/struct_tag:
	bash scripts/gen_configs_struct_tag.sh

gen/mock:
	bash scripts/gen_mock.sh

gen/swagger-ui:
	SWAGGER_UI_VERSION=$(SWAGGER_UI_VERSION) APP=$(APP) bash -x ./scripts/gen_swagger-ui.sh

gen/certs:
	bash scripts/gen_certs.sh

.PHONY: build/binary build/docker gen/proto gen/struct_tag gen/mock gen/swagger-ui gen/certs

# clean
clean/docker:
	docker image prune

clean/build:
	rm -rf ./build/bin


.PHONY: clean/build clean/docker lint test

# bootstrap

env/dev: bootstrap

# bootstrap the build by downloading additional tools that may be used by devs
bootstrap:
	go generate -tags tools tools/tools.go

.PHONY: env/dev bootstrap