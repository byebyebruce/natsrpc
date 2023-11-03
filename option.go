package natsrpc

import (
	"context"
	"fmt"
	"time"
)

// serverOptions server 选项
type serverOptions struct {
	errorHandler   func(interface{}) // error handler
	recoverHandler func(interface{}) // recover handler
}

type serviceMiddleware func(ctx context.Context, method string, req interface{}, next func(ctx context.Context, req interface{})) error
type callMiddleware func(ctx context.Context, method string, req interface{}, next func(ctx context.Context, req interface{}))

// serviceOptions service 选项
type serviceOptions struct {
	namespace  string            // 空间(划分隔离)
	queue      string            // sub组。默认只有一个sub会被通知到。空表示所有的sub都会收到
	id         string            // id
	timeout    time.Duration     // 请求/handle的超时
	mw         serviceMiddleware // middleware
	concurrent bool              // 是否多线程
	encoder    Encoder           // 编码器
}

// clientOptions client 选项
type clientOptions struct {
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
type ServerOption func(options *serverOptions)

// WithErrorHandler error handler
func WithErrorHandler(h func(interface{})) ServerOption {
	return func(options *serverOptions) {
		options.errorHandler = h
	}
}

// WithServerRecovery recover handler
func WithServerRecovery(h func(interface{})) ServerOption {
	return func(options *serverOptions) {
		options.recoverHandler = h
	}
}

// ServiceOption service option
type ServiceOption func(options *serviceOptions)

// WithServiceNamespace 空间集群
func WithServiceNamespace(namespace string) ServiceOption {
	return func(options *serviceOptions) {
		options.namespace = namespace
	}
}

// WithServiceEncoder 编码
func WithServiceEncoder(encoder Encoder) ServiceOption {
	return func(options *serviceOptions) {
		options.encoder = encoder
	}
}

// WithServiceID id
func WithServiceID(id interface{}) ServiceOption {
	return func(options *serviceOptions) {
		options.id = fmt.Sprintf("%v", id)
	}
}

func WithBroadcast() ServiceOption {
	return func(options *serviceOptions) {
		options.queue = ""
	}
}

// WithServiceTimeout 超时时间
func WithServiceTimeout(timeout time.Duration) ServiceOption {
	return func(options *serviceOptions) {
		options.timeout = timeout
	}
}

// WithServiceMiddleware 超时时间
func WithServiceMiddleware(mw serviceMiddleware) ServiceOption {
	return func(options *serviceOptions) {
		options.mw = mw
	}
}

type ClientOption func(options *clientOptions)

// WithClientNamespace 空间集群
func WithClientNamespace(namespace string) ClientOption {
	return func(options *clientOptions) {
		options.namespace = namespace
	}
}

// WithClientID id
func WithClientID(id interface{}) ClientOption {
	return func(options *clientOptions) {
		options.id = fmt.Sprintf("%v", id)
	}
}

// WithClientEncoder 编码
func WithClientEncoder(encoder Encoder) ClientOption {
	return func(options *clientOptions) {
		options.encoder = encoder
	}
}

// WithClientTimeout 默认call超时时间
/*
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(options *clientOptions) {
		options.timeout = timeout
	}
}
*/

// CallOption call option
type CallOption func(options *CallOptions)

/*
// WithCallID call id
func WithCallID(id interface{}) CallOption {
	return func(options *CallOptions) {
		options.id = fmt.Sprint(id)
	}
}

// WithCallNamespace 空间集群
func WithCallNamespace(namespace string) CallOption {
	return func(options *CallOptions) {
		options.namespace = namespace
	}
}
*/

// WithCallHeader header
func WithCallHeader(hd map[string]string) CallOption {
	return func(options *CallOptions) {
		options.header = hd
	}
}
