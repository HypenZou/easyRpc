package easyRpc

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// DefaultHandler .
var DefaultHandler = &handler{}

// Handler defines net message handler
type Handler interface {
	// BeforeRecv registers callback before Recv
	BeforeRecv(bh func(net.Conn) error)

	// BeforeSend registers callback before Send
	BeforeSend(bh func(net.Conn) error)

	// WrapReader wraps net.Conn to Read data with io.Reader, buffer e.g.
	WrapReader(conn net.Conn) io.Reader

	// Recv reads and returns a message from a client
	Recv(c *Client) (Message, error)

	// Send writes a message to a connection
	Send(c net.Conn, m Message) (int, error)

	// Handle registers method handler
	Handle(m string, h func(*Context))

	// OnMessage dispatches messages
	OnMessage(c *Client, m Message)
}

type handler struct {
	beforeRecv func(net.Conn) error
	beforeSend func(net.Conn) error
	routes     map[string]func(*Context)
}

func (h *handler) BeforeRecv(bh func(net.Conn) error) {
	h.beforeRecv = bh
}

func (h *handler) BeforeSend(bh func(net.Conn) error) {
	h.beforeSend = bh
}

func (h *handler) WrapReader(conn net.Conn) io.Reader {
	return bufio.NewReaderSize(conn, 1024)
}

func (h *handler) Recv(c *Client) (Message, error) {
	var (
		err     error
		message Message
	)

	if h.beforeRecv != nil {
		if err = h.beforeRecv(c.Conn); err != nil {
			return nil, err
		}
	}

	_, err = io.ReadFull(c.Reader, c.Head)
	if err != nil {
		return nil, err
	}

	message, err = c.Head.Message()
	if err == nil && len(message) > HeadLen {
		_, err = io.ReadFull(c.Reader, message[HeadLen:])
	}

	return message, err
}

func (h *handler) Send(conn net.Conn, m Message) (int, error) {
	if h.beforeSend != nil {
		if err := h.beforeSend(conn); err != nil {
			return -1, err
		}
	}
	return conn.Write(m)
}

func (h *handler) Handle(method string, cb func(*Context)) {
	if len(h.routes) == 0 {
		h.routes = map[string]func(*Context){}
	}
	if _, ok := h.routes[method]; ok {
		panic(fmt.Errorf("handler exist for method %v ", method))
	}
	h.routes[method] = cb
}

func (h *handler) OnMessage(c *Client, msg Message) {
	cmd, seq, method, body, err := msg.Parse()
	switch cmd {
	case RPCCmdReq:
		if cb, ok := h.routes[method]; ok {
			defer handlePanic()
			cb(&Context{Client: c, Message: msg})
		} else {
			DefaultLogger.Info("invalid method: [%v], %v, %v", method, body, err)
		}
	case RPCCmdRsp, RPCCmdErr:
		session, ok := c.getSession(seq)
		if ok {
			session.done <- msg
		} else {
			DefaultLogger.Info("session expired: [%v] [%v] [%v] [%v]", seq, method, string(body), err)
		}
	default:
		DefaultLogger.Info("invalid cmd: [%v]", cmd)
	}
}
