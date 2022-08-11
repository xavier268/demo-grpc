package main

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestUnsafe(t *testing.T) {
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

func TestSecure(t *testing.T) {
	go RunServer()
	time.Sleep(2 * time.Second) // some time for server to start ...
	err := RunClient("-unsafe") // should fail in secured mode
	if err == nil {
		t.Fatal("unsecure cleint should fail on secured server")
	}
	err = RunClient("-bye")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second) // some time for server to stop ...
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
