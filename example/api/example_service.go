package api

import (
	"context"
	"fmt"
	"github.com/byebyebruce/natsrpc/testdata/pb"
)

type ExampleService struct {
}

func (a *ExampleService) Notify(ctx context.Context, req *pb.HelloReply) {
	fmt.Println("begin ExampleService Notify->", req.Message)
	fmt.Println("end ExampleService Notify->", req.Message)
}

func (a *ExampleService) Request(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply) {
	fmt.Println("begin ExampleService Request", req.Name)
	repl.Message = "ExampleService: "+req.Name
	fmt.Println("end ExampleService Request->", req.Name)
}
