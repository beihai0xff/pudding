
SCHEDULER_BINARY_NAME="scheduler"

build_linux_amd64:
	@echo "scheduler linux_amd64 Building..."
	go env -w GOARCH=amd64
	go env -w GOOS=linux
	go build -o ../build/$(SCHEDULER_BINARY_NAME)_linux_amd64 ../cmd/scheduler/


build_darwin_arm64:
	@echo "scheduler darwin_arm64 Building..."
	go env -w GOARCH=arm64
	go env -w GOOS=darwin
	go build -o ../build/$(SCHEDULER_BINARY_NAME)_linux_amd64 ../cmd/scheduler/
