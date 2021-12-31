package main

import (
	"log"
	"net"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

const (
	addr = "localhost:8888"

	methodHello  = "Hello"
	methodNotify = "Notify"
)

// OnClientHello .
func OnClientHello(ctx *easyRpc.Context) {
	str := ""
	ctx.Bind(&str)
	ctx.Write(str)

	log.Printf("OnClientHello: \"%v\"", str)

	client := ctx.Client
	// send 3 notify messages
	go func() {
		notifyPayload := "notify from server, nonblock"
		client.Notify(methodNotify, notifyPayload, easyRpc.TimeZero)

		notifyPayload = "notify from server, block"
		client.Notify(methodNotify, notifyPayload, easyRpc.TimeForever)

		notifyPayload = "notify from server, with 1 second timeout"
		client.Notify(methodNotify, notifyPayload, time.Second)
	}()
}

func main() {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svr := easyRpc.NewServer()
	svr.Handler.Handle(methodHello, OnClientHello)

	svr.Serve(ln)
}
