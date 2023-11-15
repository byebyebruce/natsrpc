package natsrpc

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
)

var _ IClient = (*Client)(nil)

// Client RPC client
type Client struct {
	name string        // 服务名
	opt  ClientOptions // 选项
	conn *nats.Conn
}

// NewClient 构造器
func NewClient(conn *nats.Conn, serviceName string, opts ...ClientOption) *Client {
	opt := DefaultClientOptions
	for _, v := range opts {
		v(&opt)
	}
	c := &Client{
		conn: conn,
		name: joinSubject(opt.namespace, serviceName, opt.id),
		opt:  opt,
	}

	return c
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
	callOpt := CallOptions{}
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
	var (
		subject   = ""
		isPublish = rep == nil
	)
	if isPublish {
		subject = joinSubject(c.name, callOpt.id, pubSuffix)
	} else {
		subject = joinSubject(c.name, callOpt.id)
	}

	msg := &nats.Msg{
		Subject: subject,
		Header:  h,
		Data:    b,
	}

	if isPublish {
		// publish
		return c.conn.PublishMsg(msg)
	}

	reply, err := c.conn.RequestMsgWithContext(ctx, msg)
	if err != nil {
		return err
	}

	if errStr := getErrorHeader(reply.Header); errStr != "" {
		return errors.New(errStr)
	}

	if len(reply.Data) > 0 {
		err = c.opt.encoder.Decode(reply.Data, rep)
		if err != nil {
			return err
		}
	}
	return nil
}
