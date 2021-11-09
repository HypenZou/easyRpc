package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

var mux = sync.RWMutex{}
var clientMap = make(map[*easyRpc.Client]struct{})

func main() {
	server := easyRpc.NewServer()
	server.Handler.Handle("/enter", func(ctx *easyRpc.Context) {
		passwd := ""
		ctx.Bind(&passwd)
		if passwd == "123qwe" {
			// keep client
			mux.Lock()
			clientMap[ctx.Client] = struct{}{}
			mux.Unlock()

			// release client
			ctx.Client.OnDisconnected(func(c *easyRpc.Client) {
				mux.Lock()
				delete(clientMap, c)
				mux.Unlock()
			})

			log.Printf("enter success")
		} else {
			log.Printf("enter failed invalid passwd: %v", passwd)
			ctx.Client.Stop()
		}
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		for i := 0; true; i++ {
			<-ticker.C
			broadcast(i)
		}
	}()

	server.Run(":8888")
}

func broadcast(i int) {
	msg := easyRpc.NewRefMessage(easyRpc.CmdNotify, easyRpc.DefaultCodec, "/broadcast", fmt.Sprintf("broadcast msg %d", i))
	defer msg.Release()

	mux.RLock()
	for client := range clientMap {
		client.PushMsg(msg, easyRpc.TimeZero)
	}
	mux.RUnlock()
}
