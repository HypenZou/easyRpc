package main

import (
	"log"

	"github.com/anacrolix/utp"
	"github.com/wubbalubbaaa/easyRpc"
)

func main() {
	ln, err := utp.NewSocket("udp", ":8888")
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
