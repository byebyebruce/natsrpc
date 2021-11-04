//go:generate protoc --proto_path=. --go_out=paths=source_relative:. natsrpc.proto
package natsrpc

import (
	"log"
	"time"
)

const (
//headerError = "error"
)

var defaultServerOptions = serverOptions{
	logger: &log.Logger{},
	//recoverHandler: func(i interface{}) {
	//	fmt.Println("panic", i)
	//},
}

var defaultServiceOptions = serviceOptions{
	namespace: "default",
	id:        "",
	group:     "", // 空表示不分组(同组内只有一个sub会被通知到)
	timeout:   time.Duration(3) * time.Second,
}

var defaultClientOptions = clientOptions{
	namespace: "default",
	id:        "",
	timeout:   time.Duration(3) * time.Second,
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
