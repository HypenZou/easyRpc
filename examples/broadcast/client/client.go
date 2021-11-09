package main

import (
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

var notifyCount int32

// OnBroadcast .
func OnBroadcast(ctx *easyRpc.Context) {
	ret := ""
	ctx.Bind(&ret)
	log.Printf("OnServerNotify: \"%v\"", ret)
	if atomic.AddInt32(&notifyCount, 1) >= 20 {
		os.Exit(0)
	}
}

func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", "localhost:8888", time.Second*3)
}

func main() {
	var clients []*easyRpc.Client

	easyRpc.DefaultHandler.Handle("/broadcast", OnBroadcast)

	for i := 0; i < 10; i++ {
		client, err := easyRpc.NewClient(dialer)
		if err != nil {
			log.Println("NewClient failed:", err)
			return
		}

		client.Run()
		defer client.Stop()

		clients = append(clients, client)
	}

	for i := 0; i < 10; i++ {
		client := clients[i]
		go func() {

			passwd := "123qwe"
			response := ""
			err := client.Call("/enter", passwd, &response, time.Second*5)
			if err != nil {
				log.Printf("Call failed: %v", err)
			} else {
				log.Printf("Call Response: \"%v\"", response)
			}
		}()
	}

	<-make(chan int)
}
