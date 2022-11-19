IMAGE_NAME = pudding-test:latest
IMAGE_VERSION = alpha-1
GO_VERSION ?= 1.18
GOLANG_IMAGE = golang:$(GO_VERSION)

lint:
	buf lint
	golangci-lint run
	go test ./...


.PHONY: build
build: gen
	cd scripts && sh -x ./build.sh


.PHONY: gen
gen:
	cd api && rm -rf gen && rm -rf openapi
	cd api/proto && buf mod update
	buf generate


.PHONY: docker-build
docker-build: clean
	DOCKER_BUILDKIT=0 docker build -t scheduler:${IMAGE_VERSION} -f ./build/scheduler/Dockerfile .


.PHONY: clean
clean:
	rm -rf ./build/bin
	docker image prune
