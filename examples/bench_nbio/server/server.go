package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/codec"
	"github.com/wubbalubbaaa/easyRpc/log"
	"github.com/wubbalubbaaa/nbio"
	nlog "github.com/wubbalubbaaa/nbio/logging"
)

var (
	addr    = "localhost:8888"
	handler = easyRpc.NewHandler()
)

// HelloReq .
type HelloReq struct {
	Msg string
}

// HelloRsp .
type HelloRsp struct {
	Msg string
}

// Session .
type Session struct {
	Client *easyRpc.Client
	Buffer []byte
}

func onOpen(c *nbio.Conn) {
	client := &easyRpc.Client{Conn: c, Codec: codec.DefaultCodec, Handler: handler}
	session := &Session{
		Client: client,
		Buffer: nil,
	}
	c.SetSession(session)
}

func onData(c *nbio.Conn, data []byte) {
	iSession := c.Session()
	if iSession == nil {
		c.Close()
		return
	}
	session := iSession.(*Session)
	session.Buffer = append(session.Buffer, data...)
	if len(session.Buffer) < easyRpc.HeadLen {
		return
	}

	headBuf := session.Buffer[:4]
	header := easyRpc.Header(headBuf)
	if len(session.Buffer) < easyRpc.HeadLen+header.BodyLen() {
		return
	}

	msg := &easyRpc.Message{Buffer: session.Buffer[:easyRpc.HeadLen+header.BodyLen()]}
	session.Buffer = session.Buffer[easyRpc.HeadLen+header.BodyLen():]
	handler.OnMessage(session.Client, msg)
}

func main() {
	nlog.SetLogger(log.DefaultLogger)

	handler.SetAsyncWrite(false)

	// register router
	handler.Handle("Hello", func(ctx *easyRpc.Context) {
		req := &HelloReq{}
		ctx.Bind(req)
		ctx.Write(&HelloRsp{Msg: req.Msg})
	})

	g := nbio.NewGopher(nbio.Config{
		Network: "tcp",
		Addrs:   []string{addr},
	})

	g.OnOpen(onOpen)
	g.OnData(onData)

	err := g.Start()
	if err != nil {
		log.Error("Start failed: %v", err)
	}
	defer g.Stop()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
