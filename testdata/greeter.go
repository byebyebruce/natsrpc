//go:generate go run ../cmd/natsrpc_codegen -s=greeter.go -cm=2
package testdata

import (
	"context"

	"github.com/byebyebruce/natsrpc/testdata/pb"
)

// Greeter hello
type Greeter interface {
	// HiAll publish to all
	HiAll(ctx context.Context, req *pb.HelloRequest)

	// AreYouOK request
	AreYouOK(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error)

	// DelayAreYouOK request
	DelayAreYouOK(ctx context.Context, req *pb.HelloRequest, reply func(*pb.HelloReply, error))
}
