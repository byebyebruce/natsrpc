package natsrpc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/byebyebruce/natsrpc/encode/pb"
)

var (
	ErrHeaderFormat     = errors.New("header format error")
	ErrDuplicateService = errors.New("duplicate Service")
	ErrNoMethod         = errors.New("no method")
	ErrNoMeta           = errors.New("no meta, or is not a natsrpc context")
	ErrEmptyReply       = errors.New("reply is empty")

	ErrReplyLater = errors.New("reply later") // It's not an error, when you want to reply message later, then return this.
)

const (
	defaultSubQueue = "_ns_q" // default queue group

	headerMethod = "_ns_method" // method
	headerUser   = "_ns_user"   // user header
	headerError  = "_ns_error"  // reply error
	pubSuffix    = "_pub"       // publish subject suffix
)

type ServiceRegistrar interface {
	Register(sd ServiceDesc, svc any, opt ...ServiceOption) (IService, error)
}

type IClient interface {
	Publish(method string, req interface{}, opt ...CallOption) error
	Request(ctx context.Context, method string, req interface{}, rep interface{}, opt ...CallOption) error
}

// IService 服务
type IService interface {
	Name() string
	Close() bool
}

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, vPtr interface{}) error
}

var DefaultServerOptions = ServerOptions{
	errorHandler: func(i interface{}) {
		fmt.Fprintf(os.Stderr, "error:%v\n", i)
	},
	recoverHandler: func(i interface{}) {
		fmt.Fprintf(os.Stderr, "server panic:%v\n", i)
	},
	encoder: pb.Encoder{},
}

var DefaultServiceOptions = ServiceOptions{
	timeout:         time.Duration(3) * time.Second,
	singleGoroutine: false,
	id:              "",
}

var DefaultClientOptions = ClientOptions{
	namespace: "",
	id:        "",
	encoder:   pb.Encoder{},
}
