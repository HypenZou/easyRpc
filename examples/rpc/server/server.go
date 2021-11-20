package main

import (
	"log"

	"github.com/wubbalubbaaa/easyRpc"
)

func main() {
	svr := easyRpc.NewServer()

	// register router
	svr.Handler.Handle("/echo/sync", func(ctx *easyRpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		ctx.Write(str)
		log.Printf("/echo/sync: \"%v\", error: %v", str, err)
	})

	// register router
	svr.Handler.Handle("/echo/async", func(ctx *easyRpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		go ctx.Write(str)
		log.Printf("/echo/async: \"%v\", error: %v", str, err)
	})

	svr.Run(":8888")
}
