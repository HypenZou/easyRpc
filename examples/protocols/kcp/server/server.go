package main

import (
	"crypto/sha1"
	"log"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	ln, err := kcp.ListenWithOptions(":8888", block, 10, 3)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svr := easyRpc.NewServer()

	// register router
	svr.Handler.Handle("/echo", func(ctx *easyRpc.Context) {
		str := ""
		ctx.Bind(&str)
		ctx.Write(str)
		log.Printf("/echo: \"%v\", error: %v", str, err)
	})

	svr.Serve(ln)
}
