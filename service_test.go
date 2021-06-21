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

func (a *A) Func2(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: req.Name,
	}, nil
}

func Test_newService(t *testing.T) {
	namespace := "test"
	serviceName := "natsrpc.A"
	id := "1"
	s, err := newService(serviceName, &A{}, WithNamespace(namespace), WithID(id))

	if nil != err {
		t.Error(err)
	}

	for k, v := range s.methods {
		if CombineSubject(namespace, serviceName, v.name, id) != k {
			t.Error("subject not match")
		}
	}
}
