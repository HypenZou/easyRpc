# easyRpc - More Effective Network Communication 

[![GoDoc][1]][2] [![MIT licensed][3]][4] [![Build Status][5]][6] [![Go Report Card][7]][8] [![Coverage Statusd][9]][10]

[1]: https://godoc.org/github.com/wubbalubbaaa/easyRpc?status.svg
[2]: https://godoc.org/github.com/wubbalubbaaa/easyRpc
[3]: https://img.shields.io/badge/license-MIT-blue.svg
[4]: LICENSE
[5]: https://travis-ci.org/wubbalubbaaa/easyRpc.svg?branch=master
[6]: https://travis-ci.org/wubbalubbaaa/easyRpc
[7]: https://goreportcard.com/badge/github.com/wubbalubbaaa/easyRpc
[8]: https://goreportcard.com/report/github.com/wubbalubbaaa/easyRpc
[9]: https://codecov.io/gh/wubbalubbaaa/easyRpc/branch/master/graph/badge.svg
[10]: https://codecov.io/gh/wubbalubbaaa/easyRpc




## Contents

- [easyRpc - More Effective Network Communication](#easyRpc---more-effective-network-communication)
	- [Contents](#contents)
	- [Features](#features)
	- [Performance](#performance)
	- [Header Layout](#header-layout)
	- [Installation](#installation)
	- [Quick start](#quick-start)
	- [API Examples](#api-examples)
		- [Register Routers](#register-routers)
		- [Use Middleware](#use-middleware)
		- [Client Call, CallAsync, Notify](#client-call-callasync-notify)
		- [Server Call, CallAsync, Notify](#server-call-callasync-notify)
		- [Broadcast - Notify](#broadcast---notify)
		- [Async Response](#async-response)
		- [Handle New Connection](#handle-new-connection)
		- [Handle Disconnected](#handle-disconnected)
		- [Handle Client's send queue overstock](#handle-clients-send-queue-overstock)
		- [Custom Net Protocol](#custom-net-protocol)
		- [Custom Codec](#custom-codec)
		- [Custom Logger](#custom-logger)
		- [Custom operations before conn's recv and send](#custom-operations-before-conns-recv-and-send)
		- [Custom easyRpc.Client's Reader by wrapping net.Conn](#custom-easyRpcclients-reader-by-wrapping-netconn)
		- [Custom easyRpc.Client's send queue capacity](#custom-easyRpcclients-send-queue-capacity)
	- [Pub/Sub Examples](#pubsub-examples)
	- [More Examples](#more-examples)

## Features
- [x] Two-Way Calling
- [x] Two-Way Notify
- [x] Sync and Async Calling
- [x] Sync and Async Response
- [x] Batch Write | Writev | net.Buffers 
- [x] Broadcast
- [x] Middleware
- [x] Pub/Sub

| Pattern | Interactive Directions       | Description              |
| ------- | ---------------------------- | ------------------------ |
| call    | two-way:<br>c -> s<br>s -> c | request and response     |
| notify  | two-way:<br>c -> s<br>s -> c | request without response |


## Performance

- simple echo load testing

| Framework | Protocol        | codec.Codec   | Configuration                                             | Connection Num | Goroutine Num | Qps     |
| --------- | --------------- | ------------- | --------------------------------------------------------- | -------------- | ------------- | ------- |
| easyRpc      | tcp/localhost   | encoding/json | os: VMWare Ubuntu 18.04<br>cpu: AMD 3500U 4c8t<br>mem: 2G | 8              | 10            | 80-100k |
| grpc      | http2/localhost | protobuf      | os: VMWare Ubuntu 18.04<br>cpu: AMD 3500U 4c8t<br>mem: 2G | 8              | 10            | 20-30k  |


## Header Layout

- LittleEndian

| bodyLen | reserved | cmd    | flag    | methodLen | sequence | method          | body                    |
| ------- | -------- | ------ | ------- | --------- | -------- | --------------- | ----------------------- |
| 4 bytes | 1 byte   | 1 byte | 1 bytes | 1 bytes   | 8 bytes  | methodLen bytes | bodyLen-methodLen bytes |



## Installation

1. Get and install easyRpc

```sh
$ go get -u github.com/wubbalubbaaa/easyRpc
```

2. Import in your code:

```go
import "github.com/wubbalubbaaa/easyRpc"
```


## Quick start
 
- start a [server](https://github.com/wubbalubbaaa/easyRpc/blob/master/examples/rpc_sync/server/server.go)

```go
package main

import "github.com/wubbalubbaaa/easyRpc"

func main() {
	server := easyRpc.NewServer()

	// register router
	server.Handler.Handle("/echo", func(ctx *easyRpc.Context) {
		str := ""
		if err := ctx.Bind(&str); err == nil {
			ctx.Write(str)
		}
	})

	server.Run(":8888")
}
```

- start a [client](https://github.com/wubbalubbaaa/easyRpc/blob/master/examples/rpc/client/client.go)

```go
package main

import (
	"log"
	"net"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

func main() {
	client, err := easyRpc.NewClient(func() (net.Conn, error) {
		return net.DialTimeout("tcp", "localhost:8888", time.Second*3)
	})
	if err != nil {
		panic(err)
	}

	client.Run()
	defer client.Stop()

	req := "hello"
	rsp := ""
	err = client.Call("/echo", &req, &rsp, time.Second*5)
	if err != nil {
		log.Fatalf("Call failed: %v", err)
	} else {
		log.Printf("Call Response: \"%v\"", rsp)
	}
}
```



## API Examples

### Register Routers

```golang
var handler easyRpc.Handler

// package
handler = easyRpc.DefaultHandler
// server
handler = server.Handler
// client
handler = client.Handler

// message would be default handled one by one  in the same conn reader goroutine
handler.Handle("/route", func(ctx *easyRpc.Context) { ... })
handler.Handle("/route2", func(ctx *easyRpc.Context) { ... })

// this make message handled by a new goroutine
async := true
handler.Handle("/asyncResponse", func(ctx *easyRpc.Context) { ... }, async)
```

### Use Middleware

```golang
var handler easyRpc.Handler

// package
handler = easyRpc.DefaultHandler
// server
handler = server.Handler
// client
handler = client.Handler

handler.Use(func(ctx *easyRpc.Context) { ... })
handler.Handle("/echo", func(ctx *easyRpc.Context) { ... })
handler.Use(func(ctx *easyRpc.Context) { ... })
```

### Client Call, CallAsync, Notify

1. Call (Block, with timeout/context)

```golang
request := &Echo{...}
response := &Echo{}
timeout := time.Second*5
err := client.Call("/call/echo", request, response, timeout)
// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// defer cancel()
// err := client.CallWith(ctx, "/call/echo", request, response)
```

2. CallAsync (Nonblock, with callback and timeout/context)

```golang
request := &Echo{...}

timeout := time.Second*5
err := client.CallAsync("/call/echo", request, func(ctx *easyRpc.Context) {
	response := &Echo{}
	ctx.Bind(response)
	...	
}, timeout)
```

3. Notify (same as CallAsync with timeout/context, without callback)

```golang
data := &Notify{...}
client.Notify("/notify", data, time.Second)
// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// defer cancel()
// client.NotifyWith(ctx, "/notify", data)
```

### Server Call, CallAsync, Notify

1. Get client and keep it in your application

```golang
var client *easyRpc.Client
server.Handler.Handle("/route", func(ctx *easyRpc.Context) {
	client = ctx.Client
	// release client
	client.OnDisconnected(func(c *easyRpc.Client){
		client = nil
	})
})

go func() {
	for {
		time.Sleep(time.Second)
		if client != nil {
			client.Call(...)
			client.CallAsync(...)
			client.Notify(...)
		}
	}
}()
```

2. Then Call/CallAsync/Notify

- [See Previous](#client-call-callasync-notify)

### Broadcast - Notify

- for more details:	[**server**](https://github.com/wubbalubbaaa/easyRpc/blob/master/examples/broadcast/server/server.go) [**client**](https://github.com/wubbalubbaaa/easyRpc/blob/master/examples/broadcast/client/client.go)

```golang
var mux = sync.RWMutex{}
var clientMap = make(map[*easyRpc.Client]struct{})

func broadcast() {
	var svr *easyRpc.Server = ... 
	msg := svr.NewMessage(easyRpc.CmdNotify, "/broadcast", fmt.Sprintf("broadcast msg %d", i))
	mux.RLock()
	for client := range clientMap {
		client.PushMsg(msg, easyRpc.TimeZero)
	}
	mux.RUnlock()
}
```

### Async Response

```golang
var handler easyRpc.Handler

// package
handler = easyRpc.DefaultHandler
// server
handler = server.Handler
// client
handler = client.Handler

handler.Handle("/echo", func(ctx *easyRpc.Context) {
	req := ...
	err := ctx.Bind(req)
	if err == nil {
		// async response
		go ctx.Write(data)
	}
})
```


### Handle New Connection

```golang
// package
easyRpc.DefaultHandler.HandleConnected(func(c *easyRpc.Client) {
	...
})

// server
svr := easyRpc.NewServer()
svr.Handler.HandleConnected(func(c *easyRpc.Client) {
	...
})

// client
client, err := easyRpc.NewClient(...)
client.Handler.HandleConnected(func(c *easyRpc.Client) {
	...
})
```

### Handle Disconnected

```golang
// package
easyRpc.DefaultHandler.HandleDisconnected(func(c *easyRpc.Client) {
	...
})

// server
svr := easyRpc.NewServer()
svr.Handler.HandleDisconnected(func(c *easyRpc.Client) {
	...
})

// client
client, err := easyRpc.NewClient(...)
client.Handler.HandleDisconnected(func(c *easyRpc.Client) {
	...
})
```

### Handle Client's send queue overstock

```golang
// package
easyRpc.DefaultHandler.HandleOverstock(func(c *easyRpc.Client) {
	...
})

// server
svr := easyRpc.NewServer()
svr.Handler.HandleOverstock(func(c *easyRpc.Client) {
	...
})

// client
client, err := easyRpc.NewClient(...)
client.Handler.HandleOverstock(func(c *easyRpc.Client) {
	...
})
```

### Custom Net Protocol

```golang
// server
var ln net.Listener = ...
svr := easyRpc.NewServer()
svr.Serve(ln)

// client
dialer := func() (net.Conn, error) { 
	return ... 
}
client, err := easyRpc.NewClient(dialer)
```
 
### Custom Codec

```golang
import "github.com/wubbalubbaaa/easyRpc/codec"

var codec easyRpc.Codec = ...

// package
codec.Defaultcodec = codec

// server
svr := easyRpc.NewServer()
svr.Codec = codec

// client
client, err := easyRpc.NewClient(...)
client.Codec = codec
```

### Custom Logger

```golang
import "github.com/wubbalubbaaa/easyRpc/log"

var logger easyRpc.Logger = ...
log.SetLogger(logger) // log.DefaultLogger = logger
``` 

### Custom operations before conn's recv and send

```golang
easyRpc.DefaultHandler.BeforeRecv(func(conn net.Conn) error) {
	// ...
})

easyRpc.DefaultHandler.BeforeSend(func(conn net.Conn) error) {
	// ...
})
```

### Custom easyRpc.Client's Reader by wrapping net.Conn 

```golang
easyRpc.DefaultHandler.SetReaderWrapper(func(conn net.Conn) io.Reader) {
	// ...
})
```

### Custom easyRpc.Client's send queue capacity 

```golang
easyRpc.DefaultHandler.SetSendQueueSize(4096)
```

## Pub/Sub Examples

- start a server
```golang
import "github.com/wubbalubbaaa/easyRpc/pubsub"

var (
	address = "localhost:8888"

	password = "123qwe"

	topicName = "Broadcast"
)

func main() {
	s := pubsub.NewServer()
	s.Password = password

	// server publish to all clients
	go func() {
		for i := 0; true; i++ {
			time.Sleep(time.Second)
			s.Publish(topicName, fmt.Sprintf("message from server %v", i))
		}
	}()

	s.Run(address)
}
```

- start a subscribe client
```golang
import "github.com/wubbalubbaaa/easyRpc/log"
import "github.com/wubbalubbaaa/easyRpc/pubsub"

var (
	address = "localhost:8888"

	password = "123qwe"

	topicName = "Broadcast"
)

func onTopic(topic *pubsub.Topic) {
	log.Info("[OnTopic] [%v] \"%v\", [%v]",
		topic.Name,
		string(topic.Data),
		time.Unix(topic.Timestamp/1000000000, topic.Timestamp%1000000000).Format("2006-01-02 15:04:05.000"))
}

func main() {
	client, err := pubsub.NewClient(func() (net.Conn, error) {
		return net.DialTimeout("tcp", address, time.Second*3)
	})
	if err != nil {
		panic(err)
	}
	client.Password = password
	client.Run()

	// authentication
	err = client.Authenticate()
	if err != nil {
		panic(err)
	}

	// subscribe topic
	if err := client.Subscribe(topicName, onTopic, time.Second); err != nil {
		panic(err)
	}

	<-make(chan int)
}
```

- start a publish client
```golang
import "github.com/wubbalubbaaa/easyRpc/pubsub"

var (
	address = "localhost:8888"

	password = "123qwe"

	topicName = "Broadcast"
)

func main() {
	client, err := pubsub.NewClient(func() (net.Conn, error) {
		return net.DialTimeout("tcp", address, time.Second*3)
	})
	if err != nil {
		panic(err)
	}
	client.Password = password
	client.Run()

	// authentication
	err = client.Authenticate()
	if err != nil {
		panic(err)
	}

	for i := 0; true; i++ {
		if i%5 == 0 {
			// publish msg to all clients
			client.Publish(topicName, fmt.Sprintf("message from client %d", i), time.Second)
		} else {
			// publish msg to only one client
			client.PublishToOne(topicName, fmt.Sprintf("message from client %d", i), time.Second)
		}
		time.Sleep(time.Second)
	}
}
```


## More Examples

- See [examples](https://github.com/wubbalubbaaa/easyRpc/tree/master/examples)