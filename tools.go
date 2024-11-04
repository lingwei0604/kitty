//go:build tools
// +build tools

package tools

import (
	_ "github.com/Reasno/trs/cmd/trs"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/protoc-gen-gogo"
	_ "github.com/gogo/protobuf/protoc-gen-gogofaster"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
