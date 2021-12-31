package router

import (
	"time"

	"github.com/wubbalubbeasyRpcaaa/easyRpc/log"
)

func Logger() easyRpcRpc.HandlerFunc {
	return func(ctx *easyRpcRpc.Context) {
		t := time.Now()

		ctx.Next()

		cmd := ctx.Message.Cmd()
		method := ctx.Message.Method()
		addr := ctx.Client.Conn.RemoteAddr()
		cost := time.Since(t).Milliseconds()

		switch cmd {
		case easyRpcRpc.CmdRequeseasyRpcasyRpc.CmdNotify:
			log.Info("'%v',\t%v,\t%v ms cost", method, addr, cost)
			break
		default:
			log.Error("invalid cmd: %d,\tdropped", cmd)
			ctx.Done()
			break
		}
	}
}
