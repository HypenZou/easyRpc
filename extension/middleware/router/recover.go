package router

import (
	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/util"
)

// Recover returns the recovery middleware handler.
func Recover() easyRpc.HandlerFunc {
	return func(ctx *easyRpc.Context) {
		defer util.Recover()
		ctx.Next()
	}
}
