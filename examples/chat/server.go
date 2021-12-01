package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/log"
	"github.com/wubbalubbaaa/easyRpcext/websocket"
)

type Message struct {
	User      uint64 `json:"user"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func NewMessage(user uint64, msg string) *Message {
	return &Message{
		User:      user,
		Message:   msg,
		Timestamp: time.Now().UnixNano(),
	}
}

type Room struct {
	users       map[*easyRpc.Client]uint64
	chEnterRoom chan *easyRpc.Client
	chLeaveRoom chan *easyRpc.Client
	chBroadcast chan *Message
	chStop      chan struct{}
}

func (room *Room) Enter(cli *easyRpc.Client) {
	room.chEnterRoom <- cli
}

func (room *Room) Leave(cli *easyRpc.Client) {
	room.chLeaveRoom <- cli
}

func (room *Room) Broadcast(msg *Message) {
	room.chBroadcast <- msg
}

func (room *Room) Run() *Room {
	go func() {
		for userCnt := uint64(10000); true; userCnt++ {
			select {
			case cli := <-room.chEnterRoom:
				room.users[cli] = userCnt
				cli.UserData = userCnt
				userid := fmt.Sprintf("%v", userCnt)
				cli.Notify("/chat/server/userid", userid, 0)
				for cli, _ := range room.users {
					cli.Notify("/chat/server/userenter", NewMessage(userCnt, ""), 0)
				}
				userCnt++
				log.Info("[user_%v] enter room", userid)
			case cli := <-room.chLeaveRoom:
				delete(room.users, cli)
				userid, _ := cli.UserData.(uint64)
				for cli, _ := range room.users {
					cli.Notify("/chat/server/userleave", NewMessage(userid, ""), 0)
				}
				log.Info("[user_%v] leave room", userid)
			case msg := <-room.chBroadcast:
				for cli, _ := range room.users {
					cli.Notify("/chat/server/broadcast", msg, 0)
				}
			case <-room.chStop:
				for cli, _ := range room.users {
					cli.Notify("/chat/server/shutdown", nil, 0)
				}
				return
			}
		}
	}()
	return room
}

func (room *Room) Stop() *Room {
	close(room.chStop)
	return room
}

func NewRoom() *Room {
	return &Room{
		users:       map[*easyRpc.Client]uint64{},
		chEnterRoom: make(chan *easyRpc.Client, 1024),
		chLeaveRoom: make(chan *easyRpc.Client, 1024),
		chBroadcast: make(chan *Message, 1024),
		chStop:      make(chan struct{}),
	}
}

func NewServer(room *Room) *easyRpc.Server {
	ln, _ := websocket.Listen(":8888", nil)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("url: %v", r.URL.String())
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "chat.html")
		} else if r.URL.Path == "/websocket.js" {
			http.ServeFile(w, r, "websocket.js")
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/ws", ln.(*websocket.Listener).Handler)
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Error("ListenAndServe: ", err)
			panic(err)
		}
	}()

	svr := easyRpc.NewServer()

	svr.Handler.Handle("/chat/user/say", func(ctx *easyRpc.Context) {
		if ctx.Client.UserData != nil {
			userid, _ := ctx.Client.UserData.(uint64)
			msg := &Message{User: userid}
			err := ctx.Bind(&msg.Message)
			if err == nil {
				room.Broadcast(msg)
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
	server.Stop()

	log.Info("server exit")
}
