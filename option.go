package natsrpc

import (
	"context"
	"fmt"
	"time"
)

// ServerOptions server 选项
type ServerOptions struct {
	errorHandler   func(interface{}) // error handler
	recoverHandler func(interface{}) // recover handler
	encoder        Encoder           // 编码器
}

type Handler func(svc interface{}, ctx context.Context, req interface{}) (interface{}, error)

type Invoker func(ctx context.Context, req interface{}) (interface{}, error)

type Interceptor func(ctx context.Context, method string, req interface{}, i Invoker) (interface{}, error)

// ServiceOptions Service 选项
type ServiceOptions struct {
	namespace   string        // 空间(划分隔离)
	queue       string        // sub组。默认只有一个sub会被通知到。空表示所有的sub都会收到
	id          string        // id
	timeout     time.Duration // 请求/handle的超时
	interceptor Interceptor   // middleware
	concurrent  bool          // 是否多线程
}

// ClientOptions client 选项
type ClientOptions struct {
	namespace string // 空间(划分隔离)
	id        string // id
	//timeout   time.Duration  // 请求handle的超时
	//cm        callMiddleware // 调用中间件
	encoder Encoder // 编码器
}

// CallOptions 调用选项
type CallOptions struct {
	id string // id 会覆盖clientOptions.id
	//timeout time.Duration     // 请求handle的超时 会覆盖clientOptions.timeout
	header map[string]string // header
}

// ServerOption server option
type ServerOption func(options *ServerOptions)

// WithErrorHandler error handler
func WithErrorHandler(h func(interface{})) ServerOption {
	return func(options *ServerOptions) {
		options.errorHandler = h
	}
}

// WithServerRecovery recover handler
func WithServerRecovery(h func(interface{})) ServerOption {
	return func(options *ServerOptions) {
		options.recoverHandler = h
	}
}

// ServiceOption Service option
type ServiceOption func(options *ServiceOptions)

// WithServiceNamespace 空间集群
func WithServiceNamespace(namespace string) ServiceOption {
	return func(options *ServiceOptions) {
		options.namespace = namespace
	}
}

// WithServerEncoder 编码
func WithServerEncoder(encoder Encoder) ServerOption {
	return func(options *ServerOptions) {
		options.encoder = encoder
	}
}

// WithServiceID id
func WithServiceID(id interface{}) ServiceOption {
	return func(options *ServiceOptions) {
		options.id = fmt.Sprintf("%v", id)
	}
}

func WithBroadcast() ServiceOption {
	return func(options *ServiceOptions) {
		options.queue = ""
	}
}

// WithServiceTimeout 超时时间
func WithServiceTimeout(timeout time.Duration) ServiceOption {
	return func(options *ServiceOptions) {
		options.timeout = timeout
	}
}

// WithServiceMiddleware 超时时间
func WithServiceMiddleware(mw Interceptor) ServiceOption {
	return func(options *ServiceOptions) {
		options.interceptor = mw
	}
}

type ClientOption func(options *ClientOptions)

// WithClientNamespace 空间集群
func WithClientNamespace(namespace string) ClientOption {
	return func(options *ClientOptions) {
		options.namespace = namespace
	}
}

// WithClientID id
func WithClientID(id string) ClientOption {
	return func(options *ClientOptions) {
		options.id = fmt.Sprintf("%v", id)
	}
}

// WithClientEncoder 编码
func WithClientEncoder(encoder Encoder) ClientOption {
	return func(options *ClientOptions) {
		options.encoder = encoder
	}
}

// WithClientTimeout 默认call超时时间
/*
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(options *ClientOptions) {
		options.timeout = timeout
	}
}
*/

// CallOption call option
type CallOption func(options *CallOptions)

// WithCallID call id
func WithCallID(id interface{}) CallOption {
	return func(options *CallOptions) {
		options.id = fmt.Sprint(id)
	}
}

// WithCallHeader header
func WithCallHeader(hd map[string]string) CallOption {
	return func(options *CallOptions) {
		options.header = hd
	}
}
