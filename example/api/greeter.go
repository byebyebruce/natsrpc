package api

import (
	"context"
	"fmt"

	"github.com/byebyebruce/natsrpc/testdata/pb"
)

type Greeter struct {
}

// HiAll publish
func (a *Greeter) HiAll(ctx context.Context, req *pb.HelloRequest) {
	fmt.Println("begin HiAll Notify->", req.Name)
	fmt.Println("end HiAll Notify->", req.Name)
}

// AreYouOK request
func (a *Greeter) AreYouOK(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply) {
	fmt.Println("begin AreYouOK Request", req.Name)
	repl.Message = "AreYouOK " + req.Name
	fmt.Println("end AreYouOK Request->", req.Name)
}
