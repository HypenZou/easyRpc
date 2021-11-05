// Copyright 2020 lesismal. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package easyRpc

import "errors"

// client error
var (
	// ErrClientTimeout
	ErrClientTimeout = errors.New("timeout")
	// ErrClientQueueIsFull
	ErrClientQueueIsFull = errors.New("timeout: rpc client's send queue is full")
	// ErrClientReconnecting
	ErrClientReconnecting = errors.New("client reconnecting")
	// ErrClientStopped
	ErrClientStopped = errors.New("client stopped")
)

// message error
var (
	// ErrInvalidBodyLen
	ErrInvalidBodyLen = errors.New("invalid body length")
	// ErrInvalidMessage
	ErrInvalidMessage = errors.New("invalid message")
)

// other error
var (
	// ErrTimeout
	ErrTimeout = errors.New("timeout")

	// ErrUnexpected
	ErrUnexpected = errors.New("unexpected error")
)
