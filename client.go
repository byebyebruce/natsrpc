package natsrpc

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

// Client RPC client
type Client struct {
	enc         *nats.EncodedConn // NATS Encode Conn
	name        string            // 名字
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
		name:        CombineSubject(opt.namespace, serviceName, opt.id),
		opt:         opt,
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
	if c.opt.recoverHandler != nil {
		defer func() {
			if e := recover(); e != nil {
				c.opt.recoverHandler(e)
			}
		}()
	}
	// opt
	callOpt := CallOptions{
		namespace: c.opt.namespace,
		id:        c.opt.id,
		timeout:   c.opt.timeout,
	}
	for _, v := range opt {
		v(&callOpt)
	}

	isPublish := ctx == nil && rep == nil
	// ctx
	if !isPublish {
		if callOpt.timeout > 0 {
			newCtx, cancel := context.WithTimeout(ctx, callOpt.timeout)
			defer cancel()
			ctx = newCtx
		}
	}

	// subject
	subject := CombineSubject(callOpt.namespace, c.serviceName, callOpt.id, method)

	var retErr error
	next := func(ctx1 context.Context, req1 interface{}) {
		// req
		rpcReq, err := c.newRequest(subject, req1, callOpt.header)
		if err != nil {
			retErr = err
			return
		}
		if isPublish { // publish
			retErr = c.enc.Publish(subject, rpcReq)
			return
		} else { // request
			rp := &Reply{}
			// call
			err = c.enc.RequestWithContext(ctx1, subject, rpcReq, rp)
			if err != nil {
				retErr = err
				return
			}
			if len(rp.Error) > 0 {
				retErr = fmt.Errorf(rp.Error)
				return
			}
			// decode
			if err := c.enc.Enc.Decode(subject, rp.Payload, rep); err != nil {
				retErr = err
				return
			}
		}
	}
	if c.opt.cm == nil {
		next(ctx, req)
	} else {
		c.opt.cm(ctx, method, req, next)
	}
	return retErr
}

func (c *Client) newRequest(subject string, req interface{}, header map[string]string) (*Request, error) {
	payload, err := c.enc.Enc.Encode(subject, req)
	if err != nil {
		return nil, err
	}
	return &Request{
		Payload: payload,
		Header:  header,
	}, nil
}

// Request 请求
/*
func (c *Client) RequestWithOption(method string, req interface{}, rep interface{}, opt ...CallOption) error {
	callOpt := CallOptions{}
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
