package example

import (
	"context"
	"fmt"

	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

type ExampleService struct {
}

func (a *ExampleService) Notify(ctx context.Context, req *helloworld.HelloReply) {
	fmt.Println("begin ExampleService Notify->", req.Message)
	fmt.Println("end ExampleService Notify->", req.Message)
}

func (a *ExampleService) Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply) {
	fmt.Println("begin ExampleService Request", req.Name)
	repl.Message = req.Name
	fmt.Println("end ExampleService Request->", req.Name)
}
