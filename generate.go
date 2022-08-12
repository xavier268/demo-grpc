//go:generate bash -c "echo Generating from $( pwd )"
//go:generate go version
//go:generate go env GOPATH

//go:generate bash -c "protoc --version"
//go:generate bash -c "$( go env GOPATH )/bin/protoc-gen-go --version"
//go:generate bash -c "$( go env GOPATH )/bin/protoc-gen-go-grpc --version"
//go:generate bash -c "export PATH=$PATH:$(go env GOPATH)/bin;protoc  --go_out=.  --go-grpc_out=. src/*.proto"

//go:generate bash -c "export PATH=$PATH:$(go env GOPATH)/bin;protoc  --grpc-gateway_out=. --grpc-gateway_opt grpc_api_configuration=src/greeter.yaml src/*.proto"
//go:generate go mod tidy
package main
