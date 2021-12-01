package main

import (
	"log"
	"net/http"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpcext/websocket"
)

func main() {
	ln, _ := websocket.Listen(":8888", nil)
	http.HandleFunc("/ws", ln.(*websocket.Listener).Handler)
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	svr := easyRpc.NewServer()
	// register router
	svr.Handler.Handle("/call/echo", func(ctx *easyRpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		ctx.Write(str)
		log.Printf("/call/echo: \"%v\", error: %v", str, err)
	})

	svr.Handler.Handle("/notify", func(ctx *easyRpc.Context) {
		str := ""
		err := ctx.Bind(&str)
		log.Printf("/notify: \"%v\", error: %v", str, err)
	})

	svr.Handler.HandleConnected(func(c *easyRpc.Client) {
		// go c.Call("/server/call", "server call", 0)
		go c.Notify("/server/notify", time.Now().Format("Welcome! Now Is: 2006-01-02 15:04:05.000"), 0)
	})

	svr.Serve(ln)
}
