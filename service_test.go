package natsrpc

import (
	"context"
	"testing"
)

type A struct {
}

func (a *A) Func1(ctx context.Context, req *Empty) {
}

func (a *A) Func2(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, nil
}

func Test_Service(t *testing.T) {
	namespace := "test"
	serviceName := "natsrpc.A"
	id := "1"
	s, err := newService(serviceName, &A{}, WithServiceNamespace(namespace), WithServiceID(id))

	if nil != err {
		t.Error(err)
	}

	for k, v := range s.methods {
		if CombineSubject(namespace, serviceName, v.name, id) != k {
			t.Error("subject not match")
		}
	}
}
