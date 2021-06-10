package natsrpc

import (
	"fmt"
	"time"
)

// options 设置
type options struct {
	namespace           string        // 空间(划分隔离)
	group               string        // sub组(有分组的话，该组内只有1个sub能收到，否则全部收到
	id                  string        // id
	timeout             time.Duration // 请求/handle的超时
	serviceSingleThread bool          // 服务单线程处理
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

// WithServiceSingleThread 服务单线程处理
func WithServiceSingleThread() Option {
	return func(options *options) {
		options.serviceSingleThread = true
	}
}
