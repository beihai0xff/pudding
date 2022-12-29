//go:build tools

// Package tools ensures tool dependencies are kept in sync.  This is the
// recommended way of doing this according to
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// To install the following tools at the version used by this repo run:
// $ make bootstrap
// or
// $ go generate -tags tools tools/tools.go
package tools

//go:generate go install github.com/bufbuild/buf/cmd/buf
//go:generate go install github.com/fatih/gomodifytags
//go:generate go install github.com/golang/mock/mockgen
//go:generate go install github.com/cloudflare/cfssl/cmd/cfssl
//go:generate go install github.com/cloudflare/cfssl/cmd/cfssljson
import (
	_ "github.com/bufbuild/buf/cmd/buf"

	_ "github.com/fatih/gomodifytags"

	_ "github.com/golang/mock/mockgen"

	_ "github.com/cloudflare/cfssl/cmd/cfssl"
	_ "github.com/cloudflare/cfssl/cmd/cfssljson"
)
