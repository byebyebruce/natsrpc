package natsrpc

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ServerOptions server 选项
type ServerOptions struct {
	errorHandler   func(interface{}) // error handler
	recoverHandler func(interface{}) // recover handler
	encoder        Encoder           // 编码器
	middleware     []middleware.Middleware
	namespace      string // 空间(划分隔离)
	poolSize       int    // 协程池大小，0表示不使用池
}

type Invoker = middleware.Handler
type (
	Handler func(svc interface{}, ctx context.Context, dec func(any) error) (interface{}, error)

	//Invoker func(ctx context.Context, req interface{}) (interface{}, error)

	//Interceptor func(ctx context.Context, method string, req interface{}, invoker Invoker) (interface{}, error)
)

// ServiceOptions Service 选项
type ServiceOptions struct {
	id      string        // id
	timeout time.Duration // handle的超时,必须要大于0
	//interceptor    Interceptor   // handler's interceptor
	middleware     []middleware.Middleware
	multiGoroutine bool // 是否多协程
}

// ServerOption server option
type ServerOption func(options *ServerOptions)

func ServerMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *ServerOptions) {
		o.middleware = m
	}
}

func ServiceMiddleware(m ...middleware.Middleware) ServiceOption {
	return func(o *ServiceOptions) {
		o.middleware = m
	}
}

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

// WithServerNamespace 空间集群
func WithServerNamespace(namespace string) ServerOption {
	return func(options *ServerOptions) {
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

// WithServerPoolSize 协程池大小，0表示不使用池
func WithServerPoolSize(size int) ServerOption {
	return func(options *ServerOptions) {
		if size > 0 {
			options.poolSize = size
		}
	}
}

// WithServiceID id
func WithServiceID(id string) ServiceOption {
	return func(options *ServiceOptions) {
		options.id = id
	}
}

// WithServiceTimeout 超时时间，必须大于0
func WithServiceTimeout(timeout time.Duration) ServiceOption {
	return func(options *ServiceOptions) {
		if timeout > 0 {
			options.timeout = timeout
		}
	}
}

type ClientOption func(options *clientOptions)

// WithClientNamespace 空间集群
func WithClientNamespace(namespace string) ClientOption {
	return func(options *clientOptions) {
		options.namespace = namespace
	}
}

// WithClientMiddleware with client middleware.
func WithClientMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = m
	}
}

// WithClientEncoder 编码
func WithClientEncoder(encoder Encoder) ClientOption {
	return func(options *clientOptions) {
		options.encoder = encoder
	}
}

// WithCallID call id(不会覆盖clientOptions.id，只是用来标识这次调用)
func WithCallID(id string) CallOption {
	return func(options *callOptions) {
		options.id = id
	}
}
