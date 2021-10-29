package service_impl

import (
	"context"
	"fmt"

	"github.com/byebyebruce/natsrpc/example/pb"
)

type ExampleGreeter struct {
}

// HiAll publish
func (a *ExampleGreeter) HiAll(ctx context.Context, req *pb.HelloRequest) {
	fmt.Println("begin HiAll Notify->", req.Name)
	fmt.Println("end HiAll Notify->", req.Name)
}

// AreYouOK request
func (a *ExampleGreeter) AreYouOK(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("begin AreYouOK Request", req.Name)
	rep := &pb.HelloReply{
		Message: "AreYouOK " + req.Name,
	}
	fmt.Println("end AreYouOK Request->", req.Name)
	return rep, nil
}

// DelayAreYouOK async request
func (a *ExampleGreeter) DelayAreYouOK(ctx context.Context, req *pb.HelloRequest, f func(*pb.HelloReply, error)) {
	fmt.Println("begin DelayAreYouOK Request", req.Name)
	rep := &pb.HelloReply{
		Message: "DelayAreYouOK " + req.Name,
	}
	fmt.Println("end DelayAreYouOK Request->", req.Name)
	f(rep, nil)
}
