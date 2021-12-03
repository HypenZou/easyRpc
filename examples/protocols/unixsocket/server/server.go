package main

import (
	"log"
	"net"

	"github.com/wubbalubbaaa/easyRpc"
)

func main() {
	addr, err := net.ResolveUnixAddr("unix", "bench.unixsock")
	if err != nil {
		log.Fatalf("failed to ResolveUnixAddr: %v", err)
	}
	ln, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatalf("failed to ListenUnix: %v", err)
	}

	svr := easyRpc.NewServer()

	// register router
	svr.Handler.Handle("/echo", func(ctx *easyRpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		ctx.Write(str)
		log.Printf("/echo: \"%v\", error: %v", str, err)
	})

	svr.Serve(ln)
}
