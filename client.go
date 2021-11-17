package natsrpc

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

// Client client
type Client struct {
	enc         *nats.EncodedConn // NATS Encode Conn
	serviceName string            // 服务名
	opt         clientOptions     // 选项
}

// NewClient 构造器
func NewClient(enc *nats.EncodedConn, serviceName string, opts ...ClientOption) (*Client, error) {
	opt := defaultClientOptions
	for _, v := range opts {
		v(&opt)
	}
	c := &Client{
		enc:         enc,
		serviceName: serviceName,
		opt:         opt,
	}

	return c, nil
}

// Publish 发布
func (c *Client) Publish(method string, req interface{}) error {
	subject := CombineSubject(c.opt.namespace, c.serviceName, method, c.opt.id)
	return c.enc.Publish(subject, req)
}

// Request 请求
func (c *Client) Request(ctx context.Context, method string, req interface{}, rep interface{}, opt ...CallOption) error {
	callOpt := callOptions{}
	for _, v := range opt {
		v(&callOpt)
	}
	if callOpt.timeout != nil {
		newCtx, cancel := context.WithTimeout(ctx, *callOpt.timeout)
		defer cancel()
		ctx = newCtx
	}
	subject := CombineSubject(c.opt.namespace, c.serviceName, method, c.opt.id)
	rp := &Reply{}
	err := c.enc.RequestWithContext(ctx, subject, req, rp)
	if err != nil {
		return err
	}
	if len(rp.Error) > 0 {
		return fmt.Errorf(rp.Error)
	}
	if err := c.enc.Enc.Decode(subject, rp.Data, rep); err != nil {
		return err
	}
	return nil
}

// Request 请求
/*
func (c *Client) RequestWithOption(method string, req interface{}, rep interface{}, opt ...CallOption) error {
	callOpt := callOptions{}
	for _, v := range opt {
		v(&callOpt)
	}
	if callOpt.timeout == nil {
		callOpt.timeout = &c.opt.timeout
	}
	if callOpt.id == nil {
		callOpt.id = &c.opt.id
	}

	subject := CombineSubject(c.opt.namespace, c.serviceName, method, *callOpt.id)

	if callOpt.cb == nil {
		return c.enc.Request(subject, req, rep, *callOpt.timeout)
	} else {
		go func() {
			err := c.enc.Request(subject, req, rep, *callOpt.timeout)
			callOpt.cb(rep, err)
		}()
	}
	return nil
}
*/
