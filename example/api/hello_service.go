package api

import (
	"context"
	"fmt"
	"github.com/byebyebruce/natsrpc/testdata/pb"
)

type HelloService struct {
}

func (a *HelloService) Notify(ctx context.Context, req *pb.HelloReply) {
	fmt.Println("begin HelloService Notify->", req.Message)
	fmt.Println("end HelloService Notify->", req.Message)
}

func (a *HelloService) Request(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply) {
	fmt.Println("begin HelloService Request", req.Name)
	repl.Message = "HelloService "+req.Name
	fmt.Println("end HelloService Request->", req.Name)
}
