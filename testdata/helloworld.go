package helloworld

import (
	"context"
	"github.com/byebyebruce/natsrpc/testdata/pb"
)

// Example example
type Example interface {
	// Notify notify
	Notify(ctx context.Context, req *pb.HelloReply)
	// Request request
	Request(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply)
}

// Hello hello
type Hello interface {
	Notify(ctx context.Context, req *pb.HelloReply)
	Request(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply)
}
