package api

import (
	"context"

	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

type ExampleService interface {
	Notify(ctx context.Context, req *helloworld.HelloReply)
	Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply)
}

type HelloService interface {
	Notify(ctx context.Context, req *helloworld.HelloReply)
	Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply)
}