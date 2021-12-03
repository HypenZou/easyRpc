package main

import (
	"time"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/log"
	"github.com/wubbalubbaaa/easyRpc/middleware/router"
)

func main() {
	svr := easyRpc.NewServer()

	svr.Handler.Use(router.Recover)
	svr.Handler.Use(router.Logger)

	// register router
	svr.Handler.Handle("/panic", func(ctx *easyRpc.Context) {
		ctx.Write(ctx.Body())
		log.Info("/panic handler")
		panic(string(ctx.Body()))
	})

	// register router
	svr.Handler.Handle("/logger", func(ctx *easyRpc.Context) {
		ctx.Write(ctx.Body())
		log.Info("/logger handler")
		time.Sleep(time.Millisecond)
	})

	svr.Run(":8888")
}
