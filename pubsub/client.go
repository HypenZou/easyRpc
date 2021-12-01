// Copyright 2020 wubbalubbaaa. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"net"
	"sync"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/log"
	"github.com/wubbalubbaaa/easyRpc/util"
)

// Client .
type Client struct {
	*easyRpc.Client

	Password string

	mux sync.RWMutex

	topicHandlerMap map[string]TopicHandler

	onPublishHandler TopicHandler
}

// Authenticate .
func (c *Client) Authenticate() error {
	if c.Password == "" {
		return nil
	}
	err := c.Call(routeAuthenticate, c.Password, nil, time.Second*5)
	if err == nil {
		log.Info("%v [Authenticate] success from\t%v", c.Handler.LogTag(), c.Conn.RemoteAddr())
	} else {
		log.Error("%v [Authenticate] failed: %v, from\t%v", c.Handler.LogTag(), err, c.Conn.RemoteAddr())
	}
	return err
}

// Subscribe .
func (c *Client) Subscribe(topicName string, h TopicHandler, timeout time.Duration) error {
	topic, err := newTopic(topicName, nil)
	if err != nil {
		return err
	}
	bs, err := topic.toBytes()
	if err != nil {
		return err
	}

	c.mux.Lock()
	// if _, ok := c.topicHandlerMap[topicName]; ok {
	// 	panic(fmt.Errorf("handler exist for topic [%v]", topicName))
	// }
	c.topicHandlerMap[topicName] = h
	c.mux.Unlock()

	err = c.Call(routeSubscribe, bs, nil, timeout)
	if err == nil {
		log.Info("%v [Subscribe] [topic: '%v'] success from\t%v", c.Handler.LogTag(), topicName, c.Conn.RemoteAddr())
	} else {
		c.mux.Lock()
		delete(c.topicHandlerMap, topicName)
		c.mux.Unlock()
		log.Error("%v [Subscribe] [topic: '%v'] failed: %v, from\t%v", c.Handler.LogTag(), topicName, err, c.Conn.RemoteAddr())
	}
	return err
}

// Unsubscribe .
func (c *Client) Unsubscribe(topicName string, timeout time.Duration) error {
	topic, err := newTopic(topicName, nil)
	if err != nil {
		return err
	}
	bs, err := topic.toBytes()
	if err != nil {
		return err
	}
	err = c.Call(routeUnsubscribe, bs, nil, timeout)
	if err == nil {
		c.mux.Lock()
		delete(c.topicHandlerMap, topic.Name)
		c.mux.Unlock()
		log.Info("%v[Unsubscribe] [topic: '%v'] success from\t%v", c.Handler.LogTag(), topicName, c.Conn.RemoteAddr())
	} else {
		log.Error("%v[Unsubscribe] [topic: '%v'] failed: %v, from\t%v", c.Handler.LogTag(), topicName, err, c.Conn.RemoteAddr())
	}
	return err
}

// Publish .
func (c *Client) Publish(topicName string, v interface{}, timeout time.Duration) error {
	topic, err := newTopic(topicName, util.ValueToBytes(c.Codec, v))
	if err != nil {
		return err
	}
	bs, err := topic.toBytes()
	if err != nil {
		return err
	}

	err = c.Call(routePublish, bs, nil, timeout)
	if err != nil {
		log.Error("%v [Publish] [topic: '%v'] failed: %v, from\t%v", c.Handler.LogTag(), topicName, err, c.Conn.RemoteAddr())
	}
	return err
}

// PublishToOne .
func (c *Client) PublishToOne(topicName string, v interface{}, timeout time.Duration) error {
	topic, err := newTopic(topicName, util.ValueToBytes(c.Codec, v))
	if err != nil {
		return err
	}
	bs, err := topic.toBytes()
	if err != nil {
		return err
	}

	err = c.Call(routePublishToOne, bs, nil, timeout)
	if err != nil {
		log.Error("%v [PublishToOne] [topic: '%v'] failed: %v, from\t%v", c.Handler.LogTag(), topicName, err, c.Conn.RemoteAddr())
	}
	return err
}

// OnPublish .
func (c *Client) OnPublish(h TopicHandler) {
	c.onPublishHandler = h
}

// func (c *Client) invalidTopic(topic string) bool {
// 	c.mux.RLock()
// 	_, ok := c.topicHandlerMap[topic]
// 	c.mux.RUnlock()

// 	return !ok
// }

func (c *Client) initTopics() {
	c.mux.RLock()
	for name := range c.topicHandlerMap {
		topicName := name
		go util.Safe(func() {
			for i := 0; i < 10; i++ {
				topic, _ := newTopic(topicName, util.ValueToBytes(c.Codec, nil))
				bs, _ := topic.toBytes()
				err := c.Call(routeSubscribe, bs, nil, time.Second*10)
				if err == nil {
					log.Info("%v [Subscribe] [topic: '%v'] success from\t%v", c.Handler.LogTag(), topicName, c.Conn.RemoteAddr())
					break
				} else {
					log.Error("%v [Subscribe] [topic: '%v'] %v times failed: %v, from\t%v", c.Handler.LogTag(), topicName, i+1, err, c.Conn.RemoteAddr())
				}
				time.Sleep(time.Second)
			}
		})
	}
	c.mux.RUnlock()
}

func (c *Client) onPublish(ctx *easyRpc.Context) {
	defer util.Recover()

	topic := &Topic{}
	msg := ctx.Message
	if msg.IsError() {
		log.Warn("%v [Publish IN] failed [%v], to\t%v", c.Handler.LogTag(), msg.Error(), ctx.Client.Conn.RemoteAddr())
		return
	}
	err := topic.fromBytes(ctx.Body())
	if err != nil {
		log.Warn("%v [Publish IN] failed [%v], to\t%v", c.Handler.LogTag(), err, ctx.Client.Conn.RemoteAddr())
		return
	}

	if c.onPublishHandler == nil {
		c.mux.RLock()
		if h, ok := c.topicHandlerMap[topic.Name]; ok {
			h(topic)
			c.mux.RUnlock()
		} else {
			c.mux.RUnlock()
		}
	} else {
		c.onPublishHandler(topic)
	}
}

// NewClient .
func NewClient(dialer func() (net.Conn, error)) (*Client, error) {
	c, err := easyRpc.NewClient(dialer)
	if err != nil {
		return nil, err
	}
	c.Handler.SetLogTag("[APS CLI]")
	cli := &Client{
		Client:          c,
		topicHandlerMap: map[string]TopicHandler{},
	}
	cli.Handler = cli.Handler.Clone()
	cli.Handler.Handle(routePublish, cli.onPublish)
	reconnected := false
	cli.Handler.HandleConnected(func(c *easyRpc.Client) {
		if reconnected == false {
			reconnected = true
		} else {
			if cli.Authenticate() == nil {
				cli.initTopics()
			}
		}
	})
	return cli, nil
}
