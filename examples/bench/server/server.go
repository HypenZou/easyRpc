package main

import (
	"log"
	"net"

	"github.com/wubbalubbaaa/easyRpc"
)

const (
	addr = "localhost:8888"
)

// HelloReq .
type HelloReq struct {
	Msg string
}

// HelloRsp .
type HelloRsp struct {
	Msg string
}

// OnHello .
func OnHello(ctx *easyRpc.Context) {
	req := &HelloReq{}
	ctx.Bind(req)
	ctx.Write(&HelloRsp{Msg: req.Msg})
}

func main() {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svr := easyRpc.NewServer()
	svr.Handler.Handle("Hello", OnHello)
	svr.Serve(ln)
}
