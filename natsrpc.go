package natsrpc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/byebyebruce/natsrpc/encode/gogopb"
)

var (
	ErrHeaderFormat     = errors.New("natsrpc: header format error")
	ErrDuplicateService = errors.New("natsrpc: duplicate service")
	ErrNoMethod         = errors.New("natsrpc: no method")
	ErrNoMeta           = errors.New("natsrpc: no meta data")
	ErrEmptyReply       = errors.New("natsrpc: reply is empty")

	// ErrReplyLater
	// It's not an error, when you want to reply message later, then return this.
	ErrReplyLater = errors.New("natsrpc: reply later")
)

const (
	pubSuffix = "_nr_pub" // publish subject suffix
)

// ServiceRegistrar 注册服务
type ServiceRegistrar interface {
	// Register 注册
	Register(sd ServiceDesc, svc any, opt ...ServiceOption) (ServiceInterface, error)
}

// ClientInterface 客户端接口
type ClientInterface interface {
	// Publish 发布
	Publish(service, method string, req interface{}, opt ...CallOption) error

	// Request 请求
	Request(ctx context.Context, service, method string, req interface{}, rep interface{}, opt ...CallOption) error
}

// ServiceInterface 服务
type ServiceInterface interface {
	// Name 名字
	Name() string

	// Close 关闭
	Close() bool
}

// Encoder 编码器
type Encoder interface {
	// Encode 编码
	Encode(v interface{}) ([]byte, error)

	// Decode 解码
	Decode(data []byte, vPtr interface{}) error
}

var defaultEncoder = gogopb.Encoder{}

// DefaultServerOptions 默认server选项
var DefaultServerOptions = ServerOptions{
	errorHandler: func(i interface{}) {
		fmt.Fprintf(os.Stderr, "error:%v\n", i)
	},
	recoverHandler: func(i interface{}) {
		fmt.Fprintf(os.Stderr, "server panic:%v\n", i)
	},
	encoder: defaultEncoder,
}

// DefaultServiceOptions 默认service选项
var DefaultServiceOptions = ServiceOptions{
	timeout:        time.Duration(5) * time.Second,
	multiGoroutine: false,
	id:             "",
}

// DefaultClientOptions 默认client选项
var DefaultClientOptions = ClientOptions{
	namespace: "",
	encoder:   defaultEncoder,
}
