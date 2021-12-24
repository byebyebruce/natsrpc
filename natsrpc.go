//go:generate protoc --proto_path=. --go_out=paths=source_relative:. natsrpc.proto
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

// serviceOptions service 选项
type serviceOptions struct {
	namespace string        // 空间(划分隔离)
	group     string        // sub组。空表示不分组，组内所有的sub都会收到(非空表示有分组，同组内只有一个sub会被通知到)
	id        string        // id
	timeout   time.Duration // 请求/handle的超时
}

// clientOptions client 选项
type clientOptions struct {
	namespace string        // 空间(划分隔离)
	id        string        // id
	timeout   time.Duration // 请求handle的超时
}

// callOptions 调用选项
type callOptions struct {
	namespace string        // 空间(划分隔离) 会覆盖clientOptions.namesapce
	id        string        // id 会覆盖clientOptions.id
	timeout   time.Duration // 请求handle的超时 会覆盖clientOptions.timeout
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
	group:     "",
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

// WithServiceGroup 订阅组
func WithServiceGroup(group string) ServiceOption {
	return func(options *serviceOptions) {
		options.group = group
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

// CallOption call option
type CallOption func(options *callOptions)

// WithCallTimeout call 超时时间
func WithCallTimeout(timeout time.Duration) CallOption {
	return func(options *callOptions) {
		options.timeout = timeout
	}
}

// WithCallID call id
func WithCallID(id interface{}) CallOption {
	return func(options *callOptions) {
		options.id = fmt.Sprint(id)
	}
}

// WithCallNamespace 空间集群
func WithCallNamespace(namespace string) CallOption {
	return func(options *callOptions) {
		options.namespace = namespace
	}
}

type IServer interface {
	ClearAllSubscription()
	Close(duration time.Duration) (err error)
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

// WithHeader 填充Header
func WithHeader(ctx context.Context, header string) context.Context {
	newCtx := context.WithValue(ctx, headerKey{}, header)
	return newCtx
}

// Header 获得Header
func Header(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	val := ctx.Value(headerKey{})
	if val != nil {
		return val.(string)
	}
	return ""
}
