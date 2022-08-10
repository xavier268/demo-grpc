// Package main implements a server for Greeter service.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/xavier268/demo-grpc/auto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// greetserver is used to implement helloworld.GreeterServer.
type greetserver struct {
	serv *grpc.Server // needs to access gloabl server to stop it !
	auto.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *greetserver) SayHello(ctx context.Context, in *auto.HelloRequest) (*auto.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &auto.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *greetserver) Bye(ctx context.Context, _ *auto.Empty) (*auto.Empty, error) {
	log.Println("Server stop request received ...")
	go s.serv.GracefulStop()
	return &auto.Empty{}, nil
}

type echoserver struct {
	auto.UnimplementedEchoServer
}

func (s *echoserver) Echo(ctx context.Context, ping *auto.Ping) (*auto.Pong, error) {
	return &auto.Pong{Message: ping.Message}, nil

}

// Setup server credentials - no client authentication.
func GetTransportCredentials() credentials.TransportCredentials {
	creds, err := credentials.NewServerTLSFromFile("certif/service.pem", "certif/service.key")
	if err != nil {
		log.Fatalf("failed to setup TLS %v", err)
	}
	return creds
}

// Setup server credentials - with client authentication.
func GetTransportCredentialsClientAuth() credentials.TransportCredentials {
	// access th server certificate
	certif, err := tls.LoadX509KeyPair("certif/service.pem", "certif/service.key")
	if err != nil {
		log.Fatalf("failed to load server key/cert %v", err)
	}
	// Load certificate authority
	cp := x509.NewCertPool()
	ca, err := os.ReadFile("certif/ca.cert")
	if err != nil {
		log.Fatalf("failed to load CA cert : %v", err)
	}
	if !cp.AppendCertsFromPEM(ca) {
		log.Printf("failed to append CA certs to cert pool")
	}
	// make credentials
	tlsc := tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certif},
		ClientCAs:    cp,
	}
	creds := credentials.NewTLS(&tlsc)
	return creds
}

func main() {

	flag.Parse()

	// open connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// register services
	globalServer := grpc.NewServer(grpc.Creds(GetTransportCredentialsClientAuth()))
	auto.RegisterGreeterServer(globalServer, &greetserver{serv: globalServer})
	auto.RegisterEchoServer(globalServer, &echoserver{})
	log.Printf("server listening at %v", lis.Addr())

	// run server
	if err := globalServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
