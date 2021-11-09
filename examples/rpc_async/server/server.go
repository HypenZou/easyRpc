package main

import "github.com/wubbalubbaaa/easyRpc"

func main() {
	svr := easyRpc.NewServer()

	// register router
	svr.Handler.Handle("/echo", func(ctx *easyRpc.Context) {
		str := ""
		ctx.Bind(&str)

		// async response should Clone a Context to Write and Release after used
		ctxCopy := ctx.Clone()
		go func() {
			defer ctxCopy.Release()
			ctxCopy.Write(str)
		}()
	})

	svr.Run(":8888")
}
