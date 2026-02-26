package natsrpc

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/nats-io/nats.go"
)

// clientOptions client 选项
type clientOptions struct {
	namespace  string  // 空间(划分隔离)
	encoder    Encoder // 编码器
	middleware []middleware.Middleware
}

// callOptions 调用选项
type callOptions struct {
	header map[string][]string // header
	id     string              // id (不会覆盖clientOptions.id，只是用来标识这次调用)
}

// CallOption call option
type CallOption func(options *callOptions)

// WithCallHeader header
func WithCallHeader(hd map[string][]string) CallOption {
	return func(options *callOptions) {
		options.header = hd
	}
}

var _ ClientInterface = (*Client)(nil)

// Client RPC client
type Client struct {
	opt  clientOptions // 选项
	conn *nats.Conn    // nats conn
}

// NewClient 构造器
func NewClient(conn *nats.Conn, opts ...ClientOption) (*Client, error) {
	if conn == nil {
		return nil, errors.New("nats connection is nil")
	}
	opt := clientOptions{
		namespace: "",
		encoder:   defaultEncoder,
	}

	for _, v := range opts {
		v(&opt)
	}
	c := &Client{
		conn: conn,
		opt:  opt,
	}

	return c, nil
}

func NewClientWithConnector(conncter func() (*nats.Conn, error), opts ...ClientOption) (*Client, error) {
	conn, err := conncter()
	if err != nil {
		return nil, err
	}
	client, err := NewClient(conn, opts...)
	if err != nil {
		conn.Close() // 清理连接，避免泄漏
		return nil, err
	}
	return client, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

// Request 请求
func (c *Client) Invoke(ctx context.Context, service, method string, req interface{}, rep interface{}, opt ...CallOption) error {
	return c.call(ctx, service, method, req, rep, opt...)
}

func (c *Client) call(ctx context.Context, service, method string, req interface{}, rep interface{}, opt ...CallOption) error {
	callOpt := &callOptions{}
	for _, v := range opt {
		v(callOpt)
	}

	return c.invoke(ctx, service, method, req, rep, callOpt)
}

func (c *Client) invoke(ctx context.Context, service, method string, req interface{}, rep interface{}, opt *callOptions) error {
	payload, err := c.opt.encoder.Encode(req)
	if err != nil {
		return err
	}
	header := encodeHeader(method, opt.header)
	var (
		subject   = ""
		isPublish = rep == nil
	)
	if isPublish {
		subject = joinSubject(c.opt.namespace, service, opt.id, pubSuffix)
	} else {
		subject = joinSubject(c.opt.namespace, service, opt.id)
	}

	ctx = transport.NewClientContext(ctx, &Transport{
		//endpoint:     client.opts.endpoint,
		reqHeader:      opt.header,
		operation:      method,
		requestSubject: subject,
		request:        req,
	})

	msg := &nats.Msg{
		Subject: subject,
		Header:  header,
		Data:    payload,
	}

	h := func(ctx context.Context, _ any) (any, error) {
		if isPublish {
			// publish
			return nil, c.conn.PublishMsg(msg)
		}

		reply, err := c.conn.RequestMsgWithContext(ctx, msg)
		if err != nil {
			return reply, err
		}

		if errStr := getErrorHeader(reply.Header); errStr != "" {
			return nil, errors.New(errStr)
		}

		if len(reply.Data) > 0 {
			err = c.opt.encoder.Decode(reply.Data, rep)
			if err != nil {
				return nil, err
			}
		}
		return reply, nil
	}
	if len(c.opt.middleware) > 0 {
		h = middleware.Chain(c.opt.middleware...)(h)
	}
	_, err = h(ctx, req)
	return err
}
