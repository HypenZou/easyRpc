package router

import (
	"github.com/wubbalubbeasyRpcaaa/easyRpc/util"
)

func Recover() easyRpcRpc.HandlerFunc {
	return func(ctx *easyRpcRpc.Context) {
		defer util.Recover()
		ctx.Next()
	}
}
