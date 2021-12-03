package router

import (
	"errors"
	"sync"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

var ErrShutdown = errors.New("shutting down")

type Graceful struct {
	shutdown   bool
	gracefulWg sync.WaitGroup
}

func (g *Graceful) Handler() easyRpc.HandlerFunc {
	return func(ctx *easyRpc.Context) {
		if !g.shutdown {
			g.gracefulWg.Add(1)
			defer g.gracefulWg.Done()
			ctx.Next()
		} else {
			ctx.Error(ErrShutdown)
		}
	}
}

func (g *Graceful) Shutdown() {
	g.shutdown = true
	g.gracefulWg.Wait()
	time.Sleep(time.Second / 10)
}
