package router

import (
	"time"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/log"
)

// Logger returns the logger middleware.
func Logger() easyRpc.HandlerFunc {
	return func(ctx *easyRpc.Context) {
		t := time.Now()

		ctx.Next()

		cmd := ctx.Message.Cmd()
		method := ctx.Message.Method()
		addr := ctx.Client.Conn.RemoteAddr()
		cost := time.Since(t).Milliseconds()

		switch cmd {
		case easyRpc.CmdRequest, easyRpc.CmdNotify:
			log.Info("'%v',\t%v,\t%v ms cost", method, addr, cost)
			break
		default:
			log.Error("invalid cmd: %d,\tdropped", cmd)
			ctx.Done()
			break
		}
	}
}
