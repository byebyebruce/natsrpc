//go:generate protoc --proto_path=. --go_out=paths=source_relative:. natsrpc.proto
//go:generate protoc --proto_path=./testdata --go_out=paths=source_relative:./testdata testdata.proto
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

type serviceMiddleware func(ctx context.Context, method string, req interface{}) error
type callMiddleware func(ctx context.Context, method string, req interface{}, options *CallOptions)

// serviceOptions service 选项
type serviceOptions struct {
	namespace string            // 空间(划分隔离)
	group     string            // sub组。空表示不分组，组内所有的sub都会收到(非空表示有分组，同组内只有一个sub会被通知到)
	id        string            // id
	timeout   time.Duration     // 请求/handle的超时
	mw        serviceMiddleware // middleware
}

// clientOptions client 选项
type clientOptions struct {
	namespace string         // 空间(划分隔离)
	id        string         // id
	timeout   time.Duration  // 请求handle的超时
	cm        callMiddleware // 调用中间件
}

// CallOptions 调用选项
type CallOptions struct {
	namespace string            // 空间(划分隔离) 会覆盖clientOptions.namesapce
	id        string            // id 会覆盖clientOptions.id
	timeout   time.Duration     // 请求handle的超时 会覆盖clientOptions.timeout
	header    map[string]string // header
}

var defaultServerOptions = serverOptions{
	errorHandler: func(i interface{}) {
		fmt.Println("error", i)
	},
	recoverHandler: func(i interface{}) {
		fmt.Println("panic", i)
	},
}

var defaultServiceOptions = serviceOptions{
	namespace: "default",
	id:        "",
	group:     "default", // 默认default组，同组内只有一个service收到
	timeout:   time.Duration(3) * time.Second,
}

var defaultClientOptions = clientOptions{
	namespace: "default",
	id:        "",
	timeout:   time.Duration(3) * time.Second,
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

// WithServiceGroup 订阅组(同组内只有个service收到)
func WithServiceGroup(group string) ServiceOption {
	return func(options *serviceOptions) {
		options.group = group
	}
}

// WithServiceNoGroup 取消订阅组(所有有个service收到)
func WithServiceNoGroup() ServiceOption {
	return func(options *serviceOptions) {
		options.group = ""
	}
}

// WithServiceNamespace 空间集群
func WithServiceNamespace(namespace string) ServiceOption {
	return func(options *serviceOptions) {
		options.namespace = namespace
	}
}

// WithServiceID id
func WithServiceID(id interface{}) ServiceOption {
	return func(options *serviceOptions) {
		options.id = fmt.Sprintf("%v", id)
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

// WithClientTimeout 默认call超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(options *clientOptions) {
		options.timeout = timeout
	}
}

// WithClientMiddleware Middleware
func WithClientMiddleware(cm callMiddleware) ClientOption {
	return func(options *clientOptions) {
		options.cm = cm
	}
}

// CallOption call option
type CallOption func(options *CallOptions)

// WithCallTimeout call 超时时间
func WithCallTimeout(timeout time.Duration) CallOption {
	return func(options *CallOptions) {
		options.timeout = timeout
	}
}

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

// WithCallNamespace 空间集群
func WithCallNamespace(namespace string) CallOption {
	return func(options *CallOptions) {
		options.namespace = namespace
	}
}

type IServer interface {
	ClearAllSubscription()
	Close(ctx context.Context) (err error)
}

type IClient interface {
	Publish(subject string, req interface{}) error
	Request(subject string, req interface{}, rep interface{}, opt ...CallOption) error
}

// IService 服务
type IService interface {
	Name() string
	Close() bool
}

type headerKey struct{}

// withHeader 填充Header
func withHeader(ctx context.Context, header map[string]string) context.Context {
	newCtx := context.WithValue(ctx, headerKey{}, header)
	return newCtx
}

// Header 获得Header
func Header(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}
	val := ctx.Value(headerKey{})
	if val != nil {
		return val.(map[string]string)
	}
	return nil
}
