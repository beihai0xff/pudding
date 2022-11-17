IMAGE_NAME = pudding-test:latest
GO_VERSION ?= 1.18
GOLANG_IMAGE = golang:$(GO_VERSION)

lint:
	buf lint
	golangci-lint run
	go test ./...


.PHONY: build
build: lint clean
	cd scripts && chmod 777 build.sh && ./build.sh


.PHONY: gen
gen:
	cd api && rm -rf gen && rm -rf openapi
	cd api/proto && buf mod update
	buf generate


container:
	docker build -t ${IMAGE_NAME} --build-arg GOLANG_IMAGE="${GOLANG_IMAGE}" \
	    --build-arg PULSAR_IMAGE="${PULSAR_IMAGE}" .

.PHONY: clean
clean:
	rm -rf ./build/bin
