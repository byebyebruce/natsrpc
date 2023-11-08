//go:generate protoc --proto_path=. --go_out=paths=source_relative:. types.proto
//go:generate protoc --proto_path=./testdata --go_out=paths=source_relative:./testdata testdata.proto
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
)

const (
	defaultSubQueue = "_ns_q" // 默认组

	headerMethod = "_ns_method" // header method
	headerUser   = "_ns_user"   // header method
	headerError  = "_ns_error"  // header error
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
	queue:      defaultSubQueue, // 默认default组，同组内只有一个service收到
	timeout:    time.Duration(3) * time.Second,
	concurrent: true,
	id:         "",
}

var DefaultClientOptions = ClientOptions{
	namespace: "",
	id:        "",
	//timeout:   time.Duration(3) * time.Second,
	encoder: pb.Encoder{},
}

func publishSuffix(sub string) string {
	return sub + "_pub"
}
