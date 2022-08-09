#!/bin/bash
go version
go env GOPATH
protoc --version
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$(go env GOPATH)/bin/protoc-gen-go --version
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$(go env GOPATH)/bin/protoc-gen-go-grpc --version
echo $PATH

