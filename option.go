package natsrpc

import (
	"fmt"
	"log"
	"time"
)

type serverOptions struct {
	logger         *log.Logger       // logger
	recoverHandler func(interface{}) // recover handler
}

// ServerOption
type ServerOption func(options *serverOptions)

// WithLogger logger
func WithServerLogger(logger *log.Logger) ServerOption {
	return func(options *serverOptions) {
		options.logger = logger
	}
}

// WithServerRecovery recover handler
func WithServerRecovery(h func(interface{})) ServerOption {
	return func(options *serverOptions) {
		options.recoverHandler = h
	}
}

// serviceOptions 设置
type serviceOptions struct {
	namespace string        // 空间(划分隔离)
	group     string        // sub组(有分组的话，该组内只有1个sub能收到，否则全部收到
	id        string        // id
	timeout   time.Duration // 请求/handle的超时
}

// ServiceOption ServiceOption
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

// clientOptions 设置
type clientOptions struct {
	namespace string        // 空间(划分隔离)
	id        string        // id
	timeout   time.Duration // 请求/handle的超时
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

// WithClientTimeout 超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(options *clientOptions) {
		options.timeout = timeout
	}
}

type callOptions struct {
	timeout *time.Duration
}

type CallOption func(options *callOptions)

// WithCallTimeout
func WithCallTimeout(timeout time.Duration) CallOption {
	return func(options *callOptions) {
		options.timeout = &timeout
	}
}
