package natsrpc

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

// Client client
type Client struct {
	enc  *nats.EncodedConn // NATS Encode Conn
	opt  Options           // 选项
	name string            // 名字 package.service
}

// NewClient 构造器
func NewClient(enc *nats.EncodedConn, name string, opts ...Option) (*Client, error) {
	opt := MakeOptions(opts...)
	c := &Client{
		enc: enc,
		opt: opt,
	}

	c.name = name
	return c, nil
}

// NewClientWithConfig
func NewClientWithConfig(cfg Config, name string, opts ...Option) (*Client, error) {
	enc, err := NewNATSConn(cfg, nats.Name(name))
	if nil != err {
		return nil, err
	}
	return NewClient(enc, name, opts...)
}

// ID 根据ID获得client
// 不能用带ID的客户端获得
// client.ID(1000).Publish(req)
func (c *Client) ID(id interface{}) *Client {
	if "" != c.opt.id {
		return nil
	}
	ret := *c
	WithID(id)(&ret.opt)
	return &ret
}

// Publish 发布
func (c *Client) Publish(method string, req proto.Message) error {
	// subject = namespace.package.service.method.id
	subject := CombineSubject(c.opt.namespace, c.name, method, c.opt.id)
	if c.opt.isSingleThreadMode() { // 单线程模式不能同步请求
		go c.publish(subject, req)
		return nil
	}
	return c.publish(subject, req)
}

// Request 请求
func (c *Client) Request(ctx context.Context, method string, req proto.Message, rep proto.Message) error {
	if c.opt.isSingleThreadMode() { // 单线程模式不能同步请求
		panic("should call AsyncRequest in single thread mode")
	}
	subject := CombineSubject(c.opt.namespace, c.name, method, c.opt.id)
	if ctx == nil {
		ctx1, cancel := context.WithTimeout(context.Background(), c.opt.timeout)
		defer cancel()
		ctx = ctx1
	}
	return c.request(ctx, subject, req, rep)
}

// AsyncRequest 异步请求
func (c *Client) AsyncRequest(method string, req proto.Message, rep proto.Message, cb func(proto.Message, error)) {
	if !c.opt.isSingleThreadMode() { // 非单线程模式不能异步请求
		panic("call AsyncRequest only in single thread mode")
	}
	subject := CombineSubject(c.opt.namespace, c.name, method, c.opt.id)
	c.asyncRequest(subject, req, rep, c.opt, cb)
}

// Publish 发布
func (c *Client) publish(sub string, message proto.Message) error {
	return c.enc.Publish(sub, message)
}

// Request 请求
func (c *Client) request(ctx context.Context, sub string, req proto.Message, rep proto.Message) error {
	return c.enc.RequestWithContext(ctx, sub, req, rep)
}

// AsyncRequest 异步请求
func (c *Client) asyncRequest(sub string, req proto.Message, rep proto.Message, opt Options, cb func(proto.Message, error)) {
	if !opt.isSingleThreadMode() { // 非单线程模式不能异步请求
		panic("call AsyncRequest only in single thread mode")
	}
	go func() { // 不阻塞主线程
		if opt.recoverHnadler != nil {
			if e := recover(); e != nil {
				opt.recoverHnadler(e)
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), opt.timeout)
		defer cancel()
		err := c.enc.RequestWithContext(ctx, sub, req, rep)
		f := func() { // 回调
			cb(rep, err)
		}
		select {
		case opt.singleThreadCbChan <- f:
		case <-ctx.Done():
			log.Println("AsyncRequest", sub, err)
		}
	}()
}
