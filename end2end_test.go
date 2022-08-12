package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestGRPCUnsafe(t *testing.T) {
	go RunServer("-unsafe")
	time.Sleep(2 * time.Second) // some time for server to start ...
	err := RunClient("-bye")    // secured connection on unsecure server should fail ..
	if err == nil {
		t.Fatal("secured connection to unsafe server should fail")
	}
	err = RunClient("-bye", "-unsafe")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second) // some time for server to stop ...

}

func TestGRPCSecure(t *testing.T) {
	go RunServer()
	time.Sleep(2 * time.Second) // some time for server to start ...
	err := RunClient("-unsafe") // should fail in secured mode
	if err == nil {
		t.Fatal("unsecure client should fail on secured server")
	}
	err = RunClient("-bye")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second) // some time for server to stop ...
}

// ---------------- utilities --------------------

func TestRESTGateway(t *testing.T) {
	go RunServer("-unsafe")
	go RunGateway()
	time.Sleep(2 * time.Second) // some time for server to start ...

	// send REST commands - valid command
	st, err := RunRESTClient("http://localhost:8080/greeter/sayhello/xav/toto")
	if err != nil || st != http.StatusOK {
		t.Fatal(st, err)
	}

	// send REST commands - invalid command
	st, err = RunRESTClient("http://localhost:8080/invalidrequest")
	if err != nil || st != http.StatusNotFound {
		t.Fatal(st, err)
	}

	// send REST commands - bye command (does not stop gateway, only grpc server)
	st, err = RunRESTClient("http://localhost:8080/bye")
	if err != nil || st != http.StatusOK {
		t.Fatal(st, err)
	}

	time.Sleep(1 * time.Second) // some time for GRPC server to stop ...

	// send REST commands - valid command, will fail since grpc not running anymore, but gateway stil running.
	st, err = RunRESTClient("http://localhost:8080/greeter/sayhello/xav/toto")
	if err != nil || st != http.StatusServiceUnavailable {
		t.Fatal(st, err)
	}

}

func RunServer(flags ...string) {
	args := []string{"run", "./greeter_server"}
	args = append(args, flags...)
	serv := exec.Command("go", args...)
	res, err := serv.CombinedOutput()
	fmt.Printf("(server)\n%s\n", res)
	if err != nil {
		fmt.Println("Server error : ", err)
		panic(err)
	}
}

func RunClient(flags ...string) error {
	args := []string{"run", "./greeter_client"}
	args = append(args, flags...)
	serv := exec.Command("go", args...)
	res, err := serv.CombinedOutput()
	fmt.Printf("(client)\n%s\n", res)
	return err
}

func RunGateway(flags ...string) {
	args := []string{"run", "./gateway"}
	args = append(args, flags...)
	serv := exec.Command("go", args...)
	res, err := serv.CombinedOutput()
	fmt.Printf("(gateway)\n%s\n", res)
	if err != nil {
		fmt.Println("Gateway error : ", err)
		panic(err)
	}
}

func RunRESTClient(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	fmt.Printf("(RESTClient)\nStatus : %s\n", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	fmt.Printf("Body : %s\n", body)
	return resp.StatusCode, nil
}
