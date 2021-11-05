package main

import (
	"log"
	"net"

	"github.com/wubbalubbaaa/easyRpc"
)

const (
	addr = ":8888"

	method = "Hello"
)

// OnClientHello .
func OnClientHello(ctx *easyRpc.Context) {
	str := ""
	ctx.Bind(&str)

	log.Printf("OnClientHello: \"%v\"", str)

	// async response should Clone a Context to Write
	go ctx.Clone().Write(str)
}

func main() {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svr := easyRpc.NewServer()
	svr.Handler.Handle(method, OnClientHello)

	svr.Serve(ln)
}
