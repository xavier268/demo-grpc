#!/bin/bash
go version
go env GOPATH

# protobuf-compiler should be installed from linux distribution
protoc --version 

# grpc generators
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$(go env GOPATH)/bin/protoc-gen-go --version
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$(go env GOPATH)/bin/protoc-gen-go-grpc --version

# REST-Gateway generators (+ openapiv2)
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
$(go env GOPATH)/bin/protoc-gen-grpc-gateway --version
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
$(go env GOPATH)/bin/protoc-gen-openapiv2 --version

echo $PATH

