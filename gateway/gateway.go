// Gateway is a REST server, listening to REST clients ad forwarding requests to the GRPC server,
// and then returning REST reformated response to REST client.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/xavier268/demo-grpc/auto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint with the REST gateway
	// Note: Make sure the gRPC server is already running properly and accessible
	// Note: As per documentation, do not use the orher register functions, do not attempt non unary functions.
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := auto.RegisterGreeterHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	// Note how the same grpc mux is already a valid http.Handler !

	// A  approach for testing - gateway stops after 10 seconds ...
	server := &http.Server{Addr: "localhost:8080", Handler: mux}
	go func() { // kill gateway after than 10 seconds, for testing purposes ...
		time.Sleep(10 * time.Second)
		server.Shutdown(ctx)
	}()

	return server.ListenAndServe()

}

func main() {
	flag.Parse()
	log.Printf("launching gateway")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
