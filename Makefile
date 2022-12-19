IMAGE_VERSION = alpha-1
GO_VERSION ?= 1.18
GOLANG_IMAGE = golang:$(GO_VERSION)
app = ""
IMAGE_NAME = pudding.${app}:latest

lint:
	cd api/proto && buf mod update && buf lint
	golangci-lint run
	go test ./...


.PHONY: build
build: gen
	echo ${app}
	cd scripts && sh -x ./build.sh -a ${app}


.PHONY: gen
gen:
	cd api && rm -rf gen && rm -rf http-spec
	cd api/protobuf-spec && buf mod update
	buf generate
	sh -x scripts/gen_configs_struct_tag.sh
	sh -x scripts/gen_mock.sh

.PHONY: docker-build
docker-build: clean
	DOCKER_BUILDKIT=0 docker build -t pudding.${app}:${IMAGE_VERSION} -f ./build/Dockerfile . --build-arg app=${app}

.PHONY:docker-clean
docker-clean:
	docker image prune

.PHONY: clean
clean:
	rm -rf ./build/bin
