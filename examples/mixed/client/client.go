package main

import (
	"log"
	"net"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

const (
	addr = "localhost:8888"
)

// HelloReq .
type HelloReq struct {
	Msg string
}

// HelloRsp .
type HelloRsp struct {
	Msg string
}

// OnServerHello .
func OnServerHello(ctx *easyRpc.Context) {
	req := &HelloReq{}
	rsp := &HelloRsp{}

	ctx.Bind(req)
	log.Printf("OnServerHello: \"%v\"", req.Msg)

	rsp.Msg = req.Msg
	ctx.Write(rsp)
}

// OnServerNotify .
func OnServerNotify(ctx *easyRpc.Context) {
	str := ""
	ctx.Bind(&str)
	log.Printf("OnServerNotify: \"%v\"", str)
}

// OnServerNotifyRefMessage .
func OnServerNotifyRefMessage(ctx *easyRpc.Context) {
	str := ""
	ctx.Bind(&str)
	log.Printf("OnServerNotifyRefMessage: \"%v\"", str)
}

// OnServerCallAsync .
func OnServerCallAsync(ctx *easyRpc.Context) {
	str := ""
	ctx.Bind(&str)
	log.Printf("OnServerCallAsync: \"%v\"", str)
	ctx.Write(str)
}

// OnClientCallAsyncRsp .
func OnClientCallAsyncRsp(ctx *easyRpc.Context) {
	str := ""
	ctx.Bind(&str)
	log.Printf("OnClientCallAsyncRsp: \"%v\"", str)
}

// InitClient .
func InitClient(client *easyRpc.Client) {
	req := &HelloReq{Msg: "ClientHello"}
	rsp := &HelloRsp{}
	err := client.Call("ClientHello", req, rsp, time.Second*5)
	if err != nil {
		log.Printf("ClientHello Call failed: %v", err)
	} else {
		log.Printf("ClientHello Call Rsp: \"%v\"", rsp.Msg)
	}

	err = client.Call("ClientWantError", "ClientWantError", nil, time.Second*5)
	if err != nil {
		log.Printf("ClientWantError Call success: %v", err)
	} else {
		log.Printf("ClientWantError Call failed: \"%v\"", rsp.Msg)
	}

	client.Notify("ClientNotify", "ClientNotify", time.Second)

	client.CallAsync("ClientCallAsync", "ClientCallAsync", OnClientCallAsyncRsp, time.Second)
}

func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", addr, time.Second*3)
}

func main() {
	clientNum := 1
	coroutineNum := 1
	loopNum := 1
	clients := make([]*easyRpc.Client, clientNum)
	for i := 0; i < clientNum; i++ {
		client, err := easyRpc.NewClient(dialer)
		if err != nil {
			log.Println("NewClient failed:", err)
			return
		}

		client.OnConnected(InitClient)
		client.Handler.Handle("ServerHello", OnServerHello)
		client.Handler.Handle("ServerNotify", OnServerNotify)
		client.Handler.Handle("ServerNotifyRefMessage", OnServerNotifyRefMessage)
		client.Handler.Handle("ServerCallAsync", OnServerCallAsync)

		client.Run()
		// defer client.Stop()
		clients[i] = client
	}

	for i := 0; i < clientNum; i++ {
		client := clients[i]
		for j := 0; j < coroutineNum; j++ {
			go func() {
				for k := 0; k < loopNum; k++ {
					InitClient(client)
				}
			}()
		}
	}
	<-make(chan int)
}
