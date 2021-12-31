package main

import (
	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/extension/middleware/coder/gzip"
	"github.com/wubbalubbaaa/easyRpc/log"
)

func main() {
	svr := easyRpc.NewServer()

	svr.Handler.UseCoder(gzip.New(1024))

	// register router
	svr.Handler.Handle("/echo", func(ctx *easyRpc.Context) {
		ctx.Write(ctx.Body())
		log.Info("/echo")
	})

	svr.Run("localhost:8888")
}
