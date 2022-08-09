// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/xavier268/demo-grpc/auto"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	serv *grpc.Server // needs to access gloabl server to stop it !
	auto.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *auto.HelloRequest) (*auto.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &auto.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) Bye(ctx context.Context, _ *auto.Empty) (*auto.Empty, error) {
	log.Println("Server stop request received ...")
	go s.serv.GracefulStop()
	return &auto.Empty{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	globalServer := grpc.NewServer()
	auto.RegisterGreeterServer(globalServer, &server{serv: globalServer})
	log.Printf("server listening at %v", lis.Addr())
	if err := globalServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
