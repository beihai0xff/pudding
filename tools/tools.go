//go:build tools

// Package tools ensure tool dependencies are kept in sync.
// This is the recommended way of doing this according to
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// To install the following tools at the version used by this repo run:
// $ make bootstrap
// or
// $ go generate -x -tags tools tools/tools.go
package tools

//go:generate go install github.com/bufbuild/buf/cmd/buf@v1.25.1
//go:generate go install github.com/fatih/gomodifytags@latest
//go:generate go install go.uber.org/mock/mockgen
//go:generate cd .. && make gen/proto
