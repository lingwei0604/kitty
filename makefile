SHA := $(shell git rev-parse --short=10 HEAD)

MAKEFILE_PATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
VERSION_DATE := $(shell $(MAKEFILE_PATH)/commit_date.sh)

install:
	go mod tidy
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc \
        github.com/gogo/protobuf/protoc-gen-gogo \
        github.com/gogo/protobuf/protoc-gen-gogofaster \
        github.com/gogo/protobuf/proto \
        github.com/envoyproxy/protoc-gen-validate \
        github.com/google/wire
	go get -d github.com/Reasno/trs/...@v0.7.4
	go install -ldflags '-X "main.version=$(SHA)" -X "main.date=$(VERSION_DATE)"' github.com/Reasno/trs/cmd/trs

default:
	go generate ./... && go build -o ./bin/kitty .

