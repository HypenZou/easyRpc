package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

const (
	addr = "localhost:8888"

	method = "Hello"
)

// OnClientCallAsyncResponse .
func OnClientCallAsyncResponse(ctx *easyRpc.Context) {
	ret := ""
	ctx.Bind(&ret)
	log.Printf("OnClientCallAsyncResponse: \"%v\"", ret)
	os.Exit(0)
}

func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", addr, time.Second*3)
}

func main() {
	client, err := easyRpc.NewClient(dialer)
	if err != nil {
		log.Println("NewClient failed:", err)
		return
	}

	client.Run()
	payload := "hello from client.CallAsync"
	client.CallAsync(method, payload, OnClientCallAsyncResponse, time.Second)
	defer client.Stop()

	<-make(chan int)
}
