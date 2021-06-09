package example

import (
	"context"
	"fmt"

	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

type MyService struct {
}

func (a *MyService) Notify(ctx context.Context, req *helloworld.HelloReply) {
	fmt.Println("MyService Notify->", req.Message)
}

func (a *MyService) Request(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply) {
	repl.Message = req.Name
	fmt.Println("MyService Request->", repl.Message)
}
