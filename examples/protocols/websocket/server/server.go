package main

import (
	"log"
	"net/http"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpcext/websocket"
)

func main() {
	ln, _ := websocket.NewListener(":8888", nil)
	http.HandleFunc("/ws", ln.(*websocket.Listener).Handler)
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	svr := easyRpc.NewServer()
	svr.Handler.SetBatchRecv(false)
	// register router
	svr.Handler.Handle("/echo", func(ctx *easyRpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		ctx.Write(str)
		log.Printf("/echo: \"%v\", error: %v", str, err)
	})

	svr.Serve(ln)
}
