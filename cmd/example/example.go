//go:generate go run ../gencode.go -s=example.go -d=a.go -p=aa
package service

import (
	"context"

	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

type ExampleService interface {
	Notify(ctx context.Context, req *helloworld.HelloReply)
	Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply)
}
