//go:generate protoc --proto_path=. --go_out=plugins=grpc:. natsrpc.proto
package natsrpc

import "time"

const (
	headerError = "error"
)

type IServer interface {
	Close(duration time.Duration) (err error)
	ClearAllSubscription()
}

// IService 服务
type IService interface {
	Name() string
	Close() bool
}
