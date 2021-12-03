package router

import (
	"github.com/wubbalubbaaa/arpc"
	"github.com/wubbalubbaaa/arpc/util"
)

func Recover(ctx *arpc.Context) {
	defer util.Recover()
	ctx.Next()
}
