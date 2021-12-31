package main

import (
	"log"

	"github.com/wubbalubbaaa/easyRpc"
)

func main() {
	svr := easyRpc.NewServer()

	// register router
	svr.Handler.Handle("/echo", func(ctx *easyRpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		ctx.Write(str)
		log.Printf("/echo: \"%v\", error: %v", str, err)
	})

	svr.Run("localhost:8888")
}
