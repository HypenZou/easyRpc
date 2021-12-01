package main

import (
	"log"
	"net"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpcext/websocket"
)

func main() {
	client, err := easyRpc.NewClient(func() (net.Conn, error) {
		return websocket.Dial("ws://localhost:8888/ws")
	})
	if err != nil {
		panic(err)
	}
	client.Handler.SetBatchRecv(false)

	client.Run()
	defer client.Stop()

	req := "hello"
	rsp := ""
	err = client.Call("/echo", &req, &rsp, time.Second*5)
	if err != nil {
		log.Fatalf("Call failed: %v", err)
	} else {
		log.Printf("Call Response: \"%v\"", rsp)
	}
}
