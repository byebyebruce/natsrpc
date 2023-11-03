//go:generate protoc --proto_path=. --go_out=paths=source_relative:. natsrpc.proto
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
	ErrDuplicateService = errors.New("duplicate service")
	ErrNoMethod         = errors.New("no method")
)

const (
	defaultSubQueue = "_nrq" // 默认组

	headerNATSRPC = "_natsrpc_" // header method
	headerError   = "_error"    // header error
)

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

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, vPtr interface{}) error
}

var defaultServerOptions = serverOptions{
	errorHandler: func(i interface{}) {
		fmt.Fprintf(os.Stderr, "error:%v\n", i)
	},
	recoverHandler: func(i interface{}) {
		fmt.Fprintf(os.Stderr, "server panic:%v\n", i)
	},
}

var defaultServiceOptions = serviceOptions{
	namespace:  "default",
	id:         "",
	queue:      defaultSubQueue, // 默认default组，同组内只有一个service收到
	timeout:    time.Duration(3) * time.Second,
	encoder:    pb.Encoder{},
	concurrent: true,
}

var defaultClientOptions = clientOptions{
	namespace: "default",
	id:        "",
	//timeout:   time.Duration(3) * time.Second,
	encoder: pb.Encoder{},
}

type headerKey struct{}

// setHeader 填充Header
func setHeader(ctx context.Context, header map[string]string) context.Context {
	newCtx := context.WithValue(ctx, headerKey{}, header)
	return newCtx
}

// GetHeader 获得Header
func GetHeader(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}
	val := ctx.Value(headerKey{})
	if val != nil {
		return val.(map[string]string)
	}
	return nil
}
