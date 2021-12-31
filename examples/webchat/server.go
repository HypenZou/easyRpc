package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/extension/protocol/websocket"
	"github.com/wubbalubbaaa/easyRpc/log"
)

// Message .
type Message struct {
	User      uint64 `json:"user"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// NewMessage .
func NewMessage(user uint64, msg string) *Message {
	return &Message{
		User:      user,
		Message:   msg,
		Timestamp: time.Now().UnixNano(),
	}
}

// Room .
type Room struct {
	users       map[*easyRpc.Client]uint64
	chEnterRoom chan *easyRpc.Client
	chLeaveRoom chan *easyRpc.Client
	chBroadcast chan *Message
	chStop      chan struct{}
}

// Enter .
func (room *Room) Enter(cli *easyRpc.Client) {
	room.chEnterRoom <- cli
}

// Leave .
func (room *Room) Leave(cli *easyRpc.Client) {
	room.chLeaveRoom <- cli
}

// Broadcast .
func (room *Room) Broadcast(msg *Message) {
	room.chBroadcast <- msg
}

// Run .
func (room *Room) Run() *Room {
	go func() {
		for userCnt := uint64(10000); true; userCnt++ {
			select {
			case cli := <-room.chEnterRoom:
				room.users[cli] = userCnt
				cli.Set("user", userCnt)
				userid := fmt.Sprintf("%v", userCnt)
				cli.Notify("/chat/server/userid", userid, 0)
				room.broadcastMsg("/chat/server/userenter", NewMessage(userCnt, ""))
				userCnt++
				log.Info("[user_%v] enter room", userid)
			case cli := <-room.chLeaveRoom:
				delete(room.users, cli)
				user, ok := cli.Get("user")
				if ok {
					userid, ok := user.(uint64)
					if ok {
						room.broadcastMsg("/chat/server/userleave", NewMessage(userid, ""))
						log.Info("[user_%v] leave room", userid)
					}
				}
			case msg := <-room.chBroadcast:
				room.broadcastMsg("/chat/server/broadcast", msg)
			case <-room.chStop:
				room.broadcastMsg("/chat/server/shutdown", nil)
				return
			}
		}
	}()
	return room
}

// Stop .
func (room *Room) Stop() *Room {
	close(room.chStop)
	return room
}

func (room *Room) broadcastMsg(method string, msg *Message) {
	for cli := range room.users {
		cli.Notify(method, msg, 0)
	}
}

// NewRoom .
func NewRoom() *Room {
	return &Room{
		users:       map[*easyRpc.Client]uint64{},
		chEnterRoom: make(chan *easyRpc.Client, 1024),
		chLeaveRoom: make(chan *easyRpc.Client, 1024),
		chBroadcast: make(chan *Message, 1024),
		chStop:      make(chan struct{}),
	}
}

// NewServer .
func NewServer(room *Room) *easyRpc.Server {
	ln, _ := websocket.Listen("localhost:8888", nil)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("url: %v", r.URL.String())
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "chat.html")
		} else if r.URL.Path == "/easyRpc.js" {
			http.ServeFile(w, r, "easyRpc.js")
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/ws", ln.(*websocket.Listener).Handler)
	go func() {
		err := http.ListenAndServe("localhost:8888", nil)
		if err != nil {
			log.Error("ListenAndServe: ", err)
			panic(err)
		}
	}()

	svr := easyRpc.NewServer()

	svr.Handler.Handle("/chat/user/say", func(ctx *easyRpc.Context) {
		if user, ok := ctx.Client.Get("user"); ok {
			if userid, ok := user.(uint64); ok {
				msg := &Message{User: userid}
				err := ctx.Bind(&msg.Message)
				if err == nil {
					room.Broadcast(msg)
				}
			}
		}
	})

	svr.Handler.HandleConnected(func(c *easyRpc.Client) {
		room.Enter(c)
	})

	svr.Handler.HandleDisconnected(func(c *easyRpc.Client) {
		room.Leave(c)
	})

	go svr.Serve(ln)

	return svr
}

func main() {
	room := NewRoom().Run()
	server := NewServer(room)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	room.Stop()
	time.Sleep(time.Second / 10)
	server.Stop()

	log.Info("server exit")
}
