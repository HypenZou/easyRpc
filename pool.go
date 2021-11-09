// Copyright 2020 wubbalubbaaa. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package easyRpc

import (
	"sync"
)

var (
	// Mem Pool
	memPool = newMemPool(MaxBodyLen)

	// Context Pool
	ctxPool = sync.Pool{
		New: func() interface{} {
			return &Context{}
		},
	}

	// rpcSession Pool
	sessionPool = sync.Pool{
		New: func() interface{} {
			return &rpcSession{done: make(chan Message, 1)}
		},
	}

	// asyncHandler Pool
	// asyncHandlerPool = sync.Pool{
	// 	New: func() interface{} {
	// 		return &asyncHandler{}
	// 	},
	// }
)

func memGet(size int) []byte {
	return memPool.Get(size)
}

func memPut(b []byte) {
	memPool.Put(b)
}

func ctxGet(c *Client, msg Message) *Context {
	ctx := ctxPool.Get().(*Context)
	ctx.Client = c
	ctx.Message = msg
	return ctx
}

func ctxPut(ctx *Context) {
	ctxPool.Put(ctx)
}

func sessionGet(seq uint64) *rpcSession {
	sess := sessionPool.Get().(*rpcSession)
	sess.seq = seq
	return sess
}

func sessionPut(sess *rpcSession) {
	select {
	case msg := <-sess.done:
		memPut(msg)
	default:
	}
	sessionPool.Put(sess)
}

// func asyncHandlerGet(h HandlerFunc) *asyncHandler {
// 	handler := asyncHandlerPool.Get().(*asyncHandler)
// 	handler.h = h
// 	return handler
// }

// func asyncHandlerPut(h *asyncHandler) {
// 	asyncHandlerPool.Put(h)
// }
