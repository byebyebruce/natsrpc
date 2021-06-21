package example

import (
	"context"
	"fmt"

	"github.com/byebyebruce/natsrpc/testdata/pb"
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
