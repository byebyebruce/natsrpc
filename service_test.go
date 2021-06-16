package natsrpc

import (
	"context"
	"fmt"
	"github.com/byebyebruce/natsrpc/testdata/pb"
	"testing"
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
	s, err := newService(&A{}, MakeOptions())
	if nil != err {
		t.Error(err)
	}
	fmt.Println(s)
}
