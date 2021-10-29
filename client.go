package natsrpc

import (
	"context"

	"github.com/nats-io/nats.go"
)

// Client client
type Client struct {
	enc *nats.EncodedConn // NATS Encode Conn
	opt Options
}

// NewClient 构造器
func NewClient(enc *nats.EncodedConn, opts ...Option) (*Client, error) {
	c := &Client{
		enc: enc,
		opt: MakeOptions(opts...),
	}

	return c, nil
}

// Publish 发布
func (c *Client) Publish(subject string, req interface{}) error {
	subject = CombineSubject(c.opt.namespace, subject, c.opt.id)
	return c.enc.Publish(subject, req)
}

// Request 请求
func (c *Client) Request(ctx context.Context, subject string, req interface{}, rep interface{}) error {
	subject = CombineSubject(c.opt.namespace, subject, c.opt.id)
	return c.enc.RequestWithContext(ctx, subject, req, rep)
}
