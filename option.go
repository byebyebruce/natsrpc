package natsrpc

import (
	"fmt"
	"time"
)

// options 设置
type options struct {
	group     string
	namespace string
	id        string
	timeout   time.Duration
}

// defaultOption 默认设置
func defaultOption() options {
	return options{
		namespace: "default",
		id:        "",
		timeout:   time.Duration(3) * time.Second,
	}
}

// Option Option
type Option func(options *options)

// WithGroup 订阅组
func WithGroup(group string) Option {
	return func(options *options) {
		options.group = group
	}
}

// WithNamespace 空间集群
func WithNamespace(namespace string) Option {
	return func(options *options) {
		options.namespace = namespace
	}
}

// WithID id
func WithID(id interface{}) Option {
	return func(options *options) {
		options.id = fmt.Sprintf("%v", id)
	}
}

// WithTimeout 超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(options *options) {
		options.timeout = timeout
	}
}
