package api

import (
	"context"

	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

// Example example
type Example interface {
	// Notify notify
	Notify(ctx context.Context, req *helloworld.HelloReply)
	// Request request
	Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply)
}

// Hello hello
type Hello interface {
	Notify(ctx context.Context, req *helloworld.HelloReply)
	Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply)
}