package natsrpc

import (
	"fmt"
	"time"
)

// Options 设置
type Options struct {
	namespace          string        // 空间(划分隔离)
	group              string        // sub组(有分组的话，该组内只有1个sub能收到，否则全部收到
	id                 string        // id
	timeout            time.Duration // 请求/handle的超时
	singleThreadCbChan chan func()   // 单线程回调通道
}

// isSingleThreadMode 单线程模式
func (o Options) isSingleThreadMode() bool {
	return nil != o.singleThreadCbChan
}

// Namespace 空间
func (o Options) Namespace() string {
	return o.namespace
}

// ID id
func (o Options) ID() string {
	return o.id
}

// Option Option
type Option func(options *Options)

// MakeOptions 构造options
func MakeOptions(opts ...Option) Options {
	ret := Options{
		namespace: "default",
		id:        "",
		timeout:   time.Duration(3) * time.Second,
	}
	for _, v := range opts {
		v(&ret)
	}
	return ret
}

// WithGroup 订阅组
func WithGroup(group string) Option {
	return func(options *Options) {
		options.group = group
	}
}

// WithNamespace 空间集群
func WithNamespace(namespace string) Option {
	return func(options *Options) {
		options.namespace = namespace
	}
}

// WithID id
func WithID(id interface{}) Option {
	return func(options *Options) {
		options.id = fmt.Sprintf("%v", id)
	}
}

// WithTimeout 超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.timeout = timeout
	}
}

// WithSingleThreadCallback 服务单线程处理
func WithSingleThreadCallback(singleThreadCbChan chan func()) Option {
	return func(options *Options) {
		if nil == singleThreadCbChan {
			panic("singleThreadCbChan is nil")
		}
		options.singleThreadCbChan = singleThreadCbChan
	}
}
