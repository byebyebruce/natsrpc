package natsrpc

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type A struct {
}

func (a *A) Func1(ctx context.Context, req *helloworld.HelloRequest) {
	fmt.Println(req.Name)
}

func (a *A) Func2(ctx context.Context, req *helloworld.HelloReply, repl *helloworld.HelloReply) {
	repl.Message = req.Message
	fmt.Println(repl.Message)
}

func Test_newService(t *testing.T) {
	s, err := newService(nil, &A{}, newDefaultOption())
	if nil != err {
		t.Error(err)
	}
	fmt.Println(s)
}
