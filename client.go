package natsrpc

import (
	"context"

	"github.com/nats-io/nats.go"
)

// Client client
type Client struct {
	enc *nats.EncodedConn // NATS Encode Conn
}

// NewPBClient 构造器
func NewPBClient(url string, option ...nats.Option) (*Client, error) {
	enc, err := NewPBEnc(url, option...)
	if err != nil {
		return nil, err
	}
	return NewClient(enc)
}

// NewClient 构造器
func NewClient(enc *nats.EncodedConn) (*Client, error) {
	c := &Client{
		enc: enc,
	}

	return c, nil
}

// Publish 发布
func (c *Client) Publish(subject string, req interface{}) error {
	return c.enc.Publish(subject, req)
}

// Request 请求
func (c *Client) Request(ctx context.Context, subject string, req interface{}, rep interface{}) error {
	return c.enc.RequestWithContext(ctx, subject, req, rep)
}
