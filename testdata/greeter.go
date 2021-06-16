package helloworld

import (
	"context"

	"github.com/byebyebruce/natsrpc/testdata/pb"
)

// Greeter hello
type Greeter interface {
	HiAll(ctx context.Context, req *pb.HelloRequest)
	AreYouOK(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply)
}
