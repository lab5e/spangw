//go:build tools
// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/mgechev/revive"
	_ "golang.org/x/lint/golint"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
