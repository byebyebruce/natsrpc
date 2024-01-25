package natsrpc

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
)

var _ ClientInterface = (*Client)(nil)

// Client RPC client
type Client struct {
	opt  ClientOptions // 选项
	conn *nats.Conn    // nats conn
}

// NewClient 构造器
func NewClient(conn *nats.Conn, opts ...ClientOption) *Client {
	opt := DefaultClientOptions
	for _, v := range opts {
		v(&opt)
	}
	c := &Client{
		conn: conn,
		opt:  opt,
	}

	return c
}

// Publish 发布
func (c *Client) Publish(service, method string, req interface{}, opt ...CallOption) error {
	return c.call(nil, service, method, req, nil, opt...)
}

// Request 请求
func (c *Client) Request(ctx context.Context, service, method string, req interface{}, rep interface{}, opt ...CallOption) error {
	return c.call(ctx, service, method, req, rep, opt...)
}

func (c *Client) call(ctx context.Context, service, method string, req interface{}, rep interface{}, opt ...CallOption) error {
	callOpt := &CallOptions{}
	for _, v := range opt {
		v(callOpt)
	}

	payload, err := c.opt.encoder.Encode(req)
	if err != nil {
		return err
	}
	header, err := encodeHeader(method, callOpt.header)
	if err != nil {
		return err
	}
	var (
		subject   = ""
		isPublish = rep == nil
	)
	if isPublish {
		subject = joinSubject(c.opt.namespace, service, c.opt.id, pubSuffix)
	} else {
		subject = joinSubject(c.opt.namespace, service, c.opt.id)
	}

	msg := &nats.Msg{
		Subject: subject,
		Header:  header,
		Data:    payload,
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
