package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

// OnClientCallAsyncResponse .
func OnClientCallAsyncResponse(ctx *easyRpc.Context) {
	ret := ""
	err := ctx.Bind(&ret)
	log.Printf("OnClientCallAsyncResponse: \"%v\", error: %v", ret, err)
	os.Exit(0)
}

func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", "localhost:8888", time.Second*3)
}

func main() {
	client, err := easyRpc.NewClient(dialer)
	if err != nil {
		log.Println("NewClient failed:", err)
		return
	}

	client.Run()
	payload := "hello from client.CallAsync"
	client.CallAsync("/echo", payload, OnClientCallAsyncResponse, time.Second)
	defer client.Stop()

	<-make(chan int)
}
