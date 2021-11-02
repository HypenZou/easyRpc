package easyRpc

import (
	"encoding/binary"
	"errors"
	"time"
)

// Context definition
type Context struct {
	Client  *Client
	Message Message
}

// Body returns payload body
func (ctx *Context) Body(v interface{}) ([]byte, error) {
	msg := ctx.Message
	switch msg.Cmd() {
	case RPCCmdReq:
		return msg[HeadLen+msg.MethodLen():], nil
	case RPCCmdRsp, RPCCmdErr:
		return msg[HeadLen:], nil
	default:
	}
	return nil, ErrInvalidMessage
}

// Bind parses data to struct
func (ctx *Context) Bind(v interface{}) error {
	if v != nil {
		data := ctx.Message[HeadLen+ctx.Message.MethodLen():]
		switch vt := v.(type) {
		case *[]byte:
			*vt = data
		case *string:
			*vt = string(data)
		case *error:
			*vt = errors.New(bytesToStr(data))
		default:
			return ctx.Client.Codec.Unmarshal(data, v)
		}
	}
	return nil
}

func (ctx *Context) newRspMessage(cmd byte, v interface{}) Message {
	var (
		data    []byte
		msg     Message
		bodyLen int
	)

	data = valueToBytes(ctx.Client.Codec, v)

	bodyLen = len(data)
	msg = Message(make([]byte, HeadLen+bodyLen))
	binary.LittleEndian.PutUint32(msg[headerIndexBodyLenBegin:headerIndexBodyLenEnd], uint32(bodyLen))
	binary.LittleEndian.PutUint64(msg[headerIndexSeqBegin:headerIndexSeqEnd], ctx.Message.Seq())
	msg[headerIndexCmd] = cmd
	msg[headerIndexAsync] = ctx.Message.Async()
	copy(msg[HeadLen:], data)

	return msg
}

// Write responses message to client
func (ctx *Context) Write(v interface{}) error {
	msg := ctx.newRspMessage(RPCCmdRsp, v)
	return ctx.Client.PushMsg(msg, TimeForever)
}

// WriteWithTimeout responses message to client with timeout
func (ctx *Context) WriteWithTimeout(v interface{}, timeout time.Duration) error {
	msg := ctx.newRspMessage(RPCCmdRsp, v)
	return ctx.Client.PushMsg(msg, timeout)
}

// Error responses error message to client
func (ctx *Context) Error(err interface{}) error {
	msg := ctx.newRspMessage(RPCCmdErr, err)
	return ctx.Client.PushMsg(msg, TimeForever)
}
