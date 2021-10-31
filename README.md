# easyRpc - Sync && Async Call supported



## Examples

- server

```golang
package main

import (
	"log"
	"net"

	"github.com/wubbalubbaaa/easyRpc"
)

const (
	addr = ":8888"
)

type HelloReq struct {
	Msg string
}

type HelloRsp struct {
	Msg string
}

func OnHello(ctx *easyRpc.Context) {
	req := &HelloReq{}
	rsp := &HelloRsp{}

	ctx.Bind(req)
	log.Printf("OnHello: \"%v\"", req.Msg)

	rsp.Msg = req.Msg
	ctx.Write(rsp)
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
```

- client

```golang
package main

import (
	"log"
	"net"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

const (
	addr = "localhost:8888"
)

type HelloReq struct {
	Msg string
}

type HelloRsp struct {
	Msg string
}

func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", addr, time.Second*3)
}

func main() {
	client, err := easyRpc.NewClient(dialer)
	if err != nil {
		log.Println("NewClient failed:", err)
		return
	}

	client.Run()
	defer client.Stop()

	req := &HelloReq{Msg: "Hello"}
	rsp := &HelloRsp{}
	err = client.Call("Hello", req, rsp, time.Second*5)
	if err != nil {
		log.Println("Call Hello failed: %v", err)
	} else {
		log.Printf("HelloRsp: \"%v\"", rsp.Msg)
	}
}
```