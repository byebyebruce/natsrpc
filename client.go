package natsrpc

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

// Client client
type Client struct {
	enc  *nats.EncodedConn // NATS Conn
	opt  Options           // 选项
	name string            // 名字
}

// NewClient 构造器
// TODO 这里的service只提供了一个服务名字
func NewClient(enc *nats.EncodedConn, service interface{}, opts ...Option) (*Client, error) {
	opt := MakeOptions(opts...)
	if !enc.Conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}
	c := &Client{
		enc: enc,
		opt: opt,
	}

	ser, err := newService(service, c.opt)
	if nil != err {
		return nil, err
	}
	c.name = ser.name
	return c, nil
}

// NewClientWithConfig
func NewClientWithConfig(cfg Config, name string, s interface{}, opts ...Option) (*Client, error) {
	client, err := NewNATSConn(cfg, nats.Name(name))
	if nil != err {
		return nil, err
	}
	return NewClient(client, s, opts...)
}

// ID 根据ID获得client
// 不能用带ID的客户端获得
// client.ID(1000).Publish(req)
func (c *Client) ID(id interface{}) *Client {
	if "" != c.opt.id {
		return nil
	}
	ret := *c
	ret.opt.id = fmt.Sprintf("%v", id)
	return &ret
}

// singleThreadMode 单线程回调模式
func (c *Client) singleThreadMode() bool {
	return nil != c.opt.singleThreadCbChan
}

// Publish 发布
func (c *Client) Publish(message proto.Message) error {
	sub := combineSubject(c.name, typeName(reflect.TypeOf(message)), c.opt.id)
	if c.singleThreadMode() { // 单线程模型不阻塞
		go c.enc.Publish(sub, message)
	} else {
		return c.enc.Publish(sub, message)
	}
	return nil
}

// Request 请求
func (c *Client) Request(req proto.Message, rep proto.Message) error {
	if c.singleThreadMode() { // 单线程模式不能同步请求
		panic("should call AsyncRequest in single thread mode")
	}
	sub := combineSubject(c.name, typeName(reflect.TypeOf(req)), c.opt.id)
	return c.enc.Request(sub, req, rep, c.opt.timeout)
}

// AsyncRequest 异步请求
func (c *Client) AsyncRequest(req proto.Message, rep proto.Message, cb func(proto.Message, error)) {
	if !c.singleThreadMode() { // 非单线程模式不能异步请求
		panic("call AsyncRequest only in single thread mode")
	}
	go func() { // 不阻塞主线程
		sub := combineSubject(c.name, typeName(reflect.TypeOf(req)), c.opt.id)
		ctx, cancel := context.WithTimeout(context.Background(), c.opt.timeout)
		defer cancel()
		err := c.enc.Request(sub, req, rep, c.opt.timeout)
		f := func() { // 回调
			cb(rep, err)
		}
		select {
		case c.opt.singleThreadCbChan <- f:
		case <-ctx.Done():
			log.Println("AsyncRequest", sub, err)
		}
	}()

}
