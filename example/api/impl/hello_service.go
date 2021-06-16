package impl

import (
	"context"
	"fmt"

	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

type HelloService struct {
}

func (a *HelloService) Notify(ctx context.Context, req *helloworld.HelloReply) {
	fmt.Println("begin HelloService Notify->", req.Message)
	fmt.Println("end HelloService Notify->", req.Message)
}

func (a *HelloService) Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply) {
	fmt.Println("begin HelloService Request", req.Name)
	repl.Message = "HelloService "+req.Name
	fmt.Println("end HelloService Request->", req.Name)
}
