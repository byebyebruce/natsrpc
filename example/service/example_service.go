package service

import (
	"context"
	"fmt"

	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

type ExampleService struct {
}

func (a *ExampleService) Notify(ctx context.Context, req *helloworld.HelloReply) {
	fmt.Println("ExampleService Notify->", req.Message)
}

func (a *ExampleService) Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply) {
	repl.Message = req.Name
	fmt.Println("ExampleService Request->", repl.Message)
}
