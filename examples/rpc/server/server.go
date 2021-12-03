package main

import (
	"log"

	"github.com/wubbalubbaaa/arpc"
)

func main() {
	svr := arpc.NewServer()

	// register router
	svr.Handler.Handle("/echo/sync", func(ctx *arpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		ctx.Write(str)
		log.Printf("/echo/sync: \"%v\", error: %v", str, err)
	})

	// register router
	svr.Handler.Handle("/echo/async", func(ctx *arpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		go ctx.Write(str)
		log.Printf("/echo/async: \"%v\", error: %v", str, err)
	})

	svr.Run(":8888")
}
