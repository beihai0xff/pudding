IMAGE_VERSION = alpha-1
GO_VERSION ?= 1.18
GOLANG_IMAGE = golang:$(GO_VERSION)
app = ""
IMAGE_NAME = pudding.${app}:latest

# lint
lint:
	cd api/protobuf-spec && buf mod update && buf lint
	golangci-lint run
	go test ./...


# build
# build binary app
build: gen_proto gen_struct_tag gen_mock
	echo ${app}
	sh -x scripts/docker-build.sh -a ${app}

# build docker image
docker-build: clean
	DOCKER_BUILDKIT=0 docker build -t pudding.${app}:${IMAGE_VERSION} -f ./build/Dockerfile . --build-arg app=${app}


# gen
gen_proto:
	sh -x scripts/gen_proto.sh

gen_struct_tag:
	sh -x scripts/gen_configs_struct_tag.sh

gen_mock:
	sh -x scripts/gen_mock.sh


# clean
docker-clean:
	docker image prune

clean:
	rm -rf ./build/bin


.PHONY: build docker-build  gen_proto gen_struct_tag gen_mock  docker-clean clean