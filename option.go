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

type (
	Handler func(svc interface{}, ctx context.Context, req interface{}) (interface{}, error)

	Invoker func(ctx context.Context, req interface{}) (interface{}, error)

	Interceptor func(ctx context.Context, method string, req interface{}, invoker Invoker) (interface{}, error)
)

// ServiceOptions Service 选项
type ServiceOptions struct {
	namespace      string        // 空间(划分隔离)
	id             string        // id
	timeout        time.Duration // handle的超时,必须要大于0
	interceptor    Interceptor   // handler's interceptor
	multiGoroutine bool          // 是否多协程
}

// ClientOptions client 选项
type ClientOptions struct {
	namespace string  // 空间(划分隔离)
	encoder   Encoder // 编码器
	//cm        callMiddleware // 调用中间件
}

// CallOptions 调用选项
type CallOptions struct {
	id     string            // id (不会覆盖clientOptions.id，只是用来标识这次调用)
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

// WithServiceSingleGoroutine 单协程,不并发handle，给那种消息需要顺序处理的情况
func WithServiceSingleGoroutine() ServiceOption {
	return func(options *ServiceOptions) {
		options.multiGoroutine = false
	}
}

// WithServerEncoder 编码
func WithServerEncoder(encoder Encoder) ServerOption {
	return func(options *ServerOptions) {
		options.encoder = encoder
	}
}

// WithServiceID id
func WithServiceID(id string) ServiceOption {
	return func(options *ServiceOptions) {
		options.id = id
	}
}

// WithServiceTimeout 超时时间
func WithServiceTimeout(timeout time.Duration) ServiceOption {
	return func(options *ServiceOptions) {
		options.timeout = timeout
	}
}

// WithServiceInterceptor handler 拦截器
func WithServiceInterceptor(i Interceptor) ServiceOption {
	return func(options *ServiceOptions) {
		options.interceptor = i
	}
}

type ClientOption func(options *ClientOptions)

// WithClientNamespace 空间集群
func WithClientNamespace(namespace string) ClientOption {
	return func(options *ClientOptions) {
		options.namespace = namespace
	}
}

// WithClientEncoder 编码
func WithClientEncoder(encoder Encoder) ClientOption {
	return func(options *ClientOptions) {
		options.encoder = encoder
	}
}

// CallOption call option
type CallOption func(options *CallOptions)

// WithCallID call id(不会覆盖clientOptions.id，只是用来标识这次调用)
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
