package natsrpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/byebyebruce/natsrpc/testdata/pb"
)

type A struct {
}

func (a *A) Func1(ctx context.Context, req *pb.HelloRequest) {
	fmt.Println(req.Name)
}

func (a *A) Func2(ctx context.Context, req *pb.HelloReply, repl *pb.HelloReply) {
	repl.Message = req.Message
	fmt.Println(repl.Message)
}

func Test_newService(t *testing.T) {
	namespace := "test"
	serviceName := "natsrpc.A"
	opt := MakeOptions()
	WithNamespace(namespace)(&opt)
	s, err := newService(serviceName, &A{}, opt)

	if nil != err {
		t.Error(err)
	}
	if s.name != fmt.Sprintf("%s.%s", namespace, serviceName) {
		t.Error("name error")
	}
}
