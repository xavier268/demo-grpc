// Package main implements a client for Greeter service.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
	"time"

	"github.com/xavier268/demo-grpc/auto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "Xavier"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	bye  = flag.Bool("bye", false, "Request server to stop")
)

// Communication will be encrypted with a server certificate valid for the provided CA authority certificate.
func GetTlsConfigVerifiedProvided() *tls.Config {
	b, _ := os.ReadFile("certif/ca.cert") // ca authority certificate, needed to validate cert
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		panic("credentials: failed to append CA certificates to local cert pool")
	}

	return &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            cp,
	}
}

// Communication will be encrypted with a server certificate valid for the known System CA certificates.
func GetTlsConfigVerifiedSystem() *tls.Config {

	cp, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("cannot access system certificate pool %v", err)
	}

	return &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            cp,
	}
}

// Communication will be encrypted with a server certificate valid for the known System CA certificates
// and client will be authenticated by server.
func GetTlsConfigAuthenticatedClient() *tls.Config {

	// Load the client certificates from disk
	clientCert, err := tls.LoadX509KeyPair("certif/client.pem", "certif/client.key")
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	cp := x509.NewCertPool()
	ca, err := os.ReadFile("certif/ca.cert")
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the certificates from the CA
	if !cp.AppendCertsFromPEM(ca) {
		log.Fatalf("failed to append ca certs")
	}

	return &tls.Config{
		ServerName:         *addr, // NOTE: this is required!
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            cp,
		InsecureSkipVerify: false,
	}
}

// Communication will be encrypted with whatever server certificate was provided.
func GetTlsConfigUnverified() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(credentials.NewTLS(GetTlsConfigAuthenticatedClient())))
	if err != nil {
		log.Printf("did not connect with credentials - try connecting without ...: %v", err)
		conn, err = grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect even withoutcredentials - try connecting without ...: %v", err)
		}
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := auto.NewGreeterClient(conn)
	p := auto.NewEchoClient(conn)

	// ping the server
	pong, err := p.Echo(ctx, &auto.Ping{Message: "pinging ..."})
	if err != nil {
		log.Fatalf("could not ping: %v", err)
	}
	log.Printf("pong received : %v", pong)
	// Contact the server and print out its response.
	r, err := c.SayHello(ctx, &auto.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	if *bye {
		_, err = c.Bye(ctx, &auto.Empty{})
		if err != nil {
			log.Fatalf("could not stop server: %v", err)
		}
		log.Printf("Server stop requested.")
	}
}
