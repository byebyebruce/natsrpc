package natsrpc

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
)

// Client RPC client
type Client struct {
	sub  string        // 名字
	name string        // 服务名
	opt  clientOptions // 选项
	conn *nats.Conn
}

// NewClient 构造器
func NewClient(conn *nats.Conn, serviceName string, opts ...ClientOption) (*Client, error) {
	opt := defaultClientOptions
	for _, v := range opts {
		v(&opt)
	}
	c := &Client{
		conn: conn,
		name: serviceName,
		sub:  CombineSubject(opt.namespace, serviceName),
		opt:  opt,
	}

	return c, nil
}

// Name 名字
func (c *Client) Name() string {
	return c.name
}

// Publish 发布
func (c *Client) Publish(method string, req interface{}, opt ...CallOption) error {
	return c.call(nil, method, req, nil, opt...)
}

// Request 请求
func (c *Client) Request(ctx context.Context, method string, req interface{}, rep interface{}, opt ...CallOption) error {
	return c.call(ctx, method, req, rep, opt...)
}

func (c *Client) call(ctx context.Context, method string, req interface{}, rep interface{}, opt ...CallOption) error {
	// opt
	callOpt := CallOptions{
		id: c.opt.id,
	}
	for _, v := range opt {
		v(&callOpt)
	}

	b, err := c.opt.encoder.Encode(req)
	if err != nil {
		return err
	}
	h, err := encodeHeader(method, callOpt.header)
	if err != nil {
		return err
	}
	// subject
	subject := CombineSubject(c.sub, callOpt.id)

	msg := &nats.Msg{
		Subject: subject,
		Header:  h,
		Data:    b,
	}

	if rep == nil {
		// publish
		return c.conn.PublishMsg(msg)
	}

	reply, err := c.conn.RequestMsgWithContext(ctx, msg)
	if err != nil {
		return err
	}
	if len(reply.Data) > 0 {
		err = c.opt.encoder.Decode(reply.Data, rep)
		if err != nil {
			return err
		}
	}
	if errStr := getErrorHeader(reply.Header); errStr != "" {
		return errors.New(errStr)
	}
	return nil
}
