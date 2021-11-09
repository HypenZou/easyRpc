// Copyright 2020 wubbalubbaaa. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package easyRpc

import (
	"log"
	"math/rand"
	"net"
	"testing"
	"time"
)

var (
	benchAddr = "localhost:16789"

	benchServer *Server
	benchClient *Client
)

func Benchmark_Call_String_Payload_64(b *testing.B) {
	benchmarkCallStringPayload(b, randString(64))
}

func Benchmark_Call_String_Payload_128(b *testing.B) {
	benchmarkCallStringPayload(b, randString(128))
}

func Benchmark_Call_String_Payload_256(b *testing.B) {
	benchmarkCallStringPayload(b, randString(256))
}

func Benchmark_Call_String_Payload_512(b *testing.B) {
	benchmarkCallStringPayload(b, randString(512))
}

func Benchmark_Call_String_Payload_1024(b *testing.B) {
	benchmarkCallStringPayload(b, randString(1024))
}

func Benchmark_Call_String_Payload_2048(b *testing.B) {
	benchmarkCallStringPayload(b, randString(2048))
}

func Benchmark_Call_String_Payload_4096(b *testing.B) {
	benchmarkCallStringPayload(b, randString(4096))
}

func Benchmark_Call_String_Payload_8192(b *testing.B) {
	benchmarkCallStringPayload(b, randString(8192))
}

func Benchmark_Call_Bytes_Payload_64(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 64))
}

func Benchmark_Call_Bytes_Payload_128(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 128))
}

func Benchmark_Call_Bytes_Payload_256(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 256))
}

func Benchmark_Call_Bytes_Payload_512(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 512))
}

func Benchmark_Call_Bytes_Payload_1024(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 1024))
}

func Benchmark_Call_Bytes_Payload_2048(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 2048))
}

func Benchmark_Call_Bytes_Payload_4096(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 4096))
}

func Benchmark_Call_Bytes_Payload_8192(b *testing.B) {
	benchmarkCallBytesPayload(b, make([]byte, 8192))
}

func Benchmark_Call_Struct_Payload_64(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(64)})
}

func Benchmark_Call_Struct_Payload_128(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(128)})
}

func Benchmark_Call_Struct_Payload_256(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(256)})
}

func Benchmark_Call_Struct_Payload_512(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(512)})
}

func Benchmark_Call_Struct_Payload_1024(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(1024)})
}

func Benchmark_Call_Struct_Payload_2048(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(2048)})
}

func Benchmark_Call_Struct_Payload_4096(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(4096)})
}

func Benchmark_Call_Struct_Payload_8192(b *testing.B) {
	benchmarkCallStructPayload(b, &message{Payload: randString(8192)})
}

func init() {
	SetLogger(nil)
	benchServer = newBenchServer()
	benchClient = newBenchClient()
}

type message struct {
	Payload string
}

func dialer() (net.Conn, error) {
	return net.DialTimeout("tcp", benchAddr, time.Second)
}

func randString(n int) string {
	letterBytes := "/?:=&1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		ret[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(ret)
}

func newBenchServer() *Server {
	s := NewServer()
	s.Handler.Handle("/echo/string", func(ctx *Context) {
		src := ""
		err := ctx.Bind(&src)
		if err != nil {
			log.Fatalf("Bind failed: %v", err)
		}
		ctx.Write(src)
	})
	s.Handler.Handle("/echo/bytes", func(ctx *Context) {
		src := ""
		err := ctx.Bind(&src)
		if err != nil {
			log.Fatalf("Bind failed: %v", err)
		}
		ctx.Write(src)
	})
	s.Handler.Handle("/echo/struct", func(ctx *Context) {
		var src message
		err := ctx.Bind(&src)
		if err != nil {
			log.Fatalf("Bind failed: %v", err)
		}
		ctx.Write(&src)
	})
	go s.Run(benchAddr)
	time.Sleep(time.Second)
	return s
}

func newBenchClient() *Client {
	c, err := NewClient(dialer)
	if err != nil {
		log.Fatalf("NewClient() failed: %v", err)
	}
	c.Run()
	return c
}

func benchmarkCallStringPayload(b *testing.B, src string) {
	for i := 0; i < b.N; i++ {
		dst := ""
		if err := benchClient.Call("/echo/string", src, &dst, time.Second); err != nil {
			b.Fatalf("benchClient.Call() string error: %v\nsrc: %v\ndst: %v", err, src, dst)
		}
	}
}

func benchmarkCallBytesPayload(b *testing.B, src []byte) {
	for i := 0; i < b.N; i++ {
		var dst []byte
		if err := benchClient.Call("/echo/bytes", src, &dst, time.Second); err != nil {
			b.Fatalf("benchClient.Call() bytes error: %v\nsrc: %v\ndst: %v", err, src, dst)
		}
	}
}

func benchmarkCallStructPayload(b *testing.B, src *message) {
	for i := 0; i < b.N; i++ {
		var dst message
		if err := benchClient.Call("/echo/struct", src, &dst, time.Second); err != nil {
			b.Fatalf("benchClient.Call() struct error: %v\nsrc: %v\ndst: %v", err, src, dst)
		}
	}
}
