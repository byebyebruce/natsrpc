package natsrpc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

// Client client
type Client struct {
	rpc  *NatsRPC
	opt  Options // 选项
	name string  // 名字
}

// NewClient 构造器
func NewClient(rpc *NatsRPC, name string, opts ...Option) (*Client, error) {
	opt := MakeOptions(opts...)
	c := &Client{
		rpc: rpc,
		opt: opt,
	}

	c.name = name
	return c, nil
}

// NewClientWithConfig
func NewClientWithConfig(cfg Config, name string, opts ...Option) (*Client, error) {
	rpc, err := NewNatsRPCWithConfig(cfg, nats.Name(name))
	if nil != err {
		return nil, err
	}
	return NewClient(rpc, name, opts...)
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
func (c *Client) Publish(method string, req proto.Message) error {
	sub := CombineSubject(c.opt.namespace, c.name, method, c.opt.id)
	return c.rpc.Publish(sub, req, c.opt)
}

// Request 请求
func (c *Client) Request(ctx context.Context, method string, req proto.Message, rep proto.Message) error {
	sub := CombineSubject(c.opt.namespace, c.name, method, c.opt.id)
	return c.rpc.Request(ctx, sub, req, rep, c.opt)
}

// AsyncRequest 异步请求
func (c *Client) AsyncRequest(method string, req proto.Message, rep proto.Message, cb func(proto.Message, error)) {
	sub := CombineSubject(c.opt.namespace, c.name, method, c.opt.id)
	c.rpc.AsyncRequest(sub, req, rep, c.opt, cb)
}
