IMAGE_VERSION = alpha-1
GO_VERSION ?= 1.18
GOLANG_IMAGE = golang:$(GO_VERSION)
app = ""
IMAGE_NAME = pudding.${app}:latest

SWAGGER_UI_VERSION:=v4.15.5

# lint
lint:
	cd api/protobuf-spec && buf mod update && buf lint
	golangci-lint run

test: gen/mock
	go test ./...


# build binary app
build/binary: gen/proto gen/struct_tag gen/swagger-ui
	echo build ${app}
	sh -x scripts/docker-build.sh -a ${app}

# build docker image
build/docker: clean/build
	DOCKER_BUILDKIT=0 docker build -t pudding.${app}:${IMAGE_VERSION} -f ./build/Dockerfile . --build-arg app=${app}


# gen
gen/proto:
	sh -x scripts/gen_proto.sh

gen/struct_tag:
	sh -x scripts/gen_configs_struct_tag.sh

gen/mock:
	sh -x scripts/gen_mock.sh

gen/swagger-ui:
	SWAGGER_UI_VERSION=$(SWAGGER_UI_VERSION) app=$(app) sh -x ./scripts/gen_swagger-ui.sh

gen/certs:
	sh -x scripts/gen_certs.sh

.PHONY: build/binary build/docker gen/proto gen/struct_tag gen/mock gen/swagger-ui gen/certs
# clean
clean/docker:
	docker image prune

clean/build:
	rm -rf ./build/bin


# bootstrap

dev: bootstrap

# bootstrap the build by downloading additional tools that may be used by devs
bootstrap:
	go generate -tags tools tools/tools.go

.PHONY:   clean/build clean/docker lint test
.PHONY: bootstrap