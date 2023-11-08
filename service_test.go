package natsrpc

import (
	"context"
	"testing"

	"github.com/byebyebruce/natsrpc/testdata"
)

type A struct {
}

func (a *A) Func1(ctx context.Context, req *testdata.Empty) {
}

func (a *A) Func2(ctx context.Context, req *testdata.Empty) (*testdata.Empty, error) {
	return &testdata.Empty{}, nil
}

func Test_Service(t *testing.T) {
	/*
		namespace := "test"
		serviceName := "natsrpc.A"
		id := "1"
		s, err := NewService(serviceName, &A{}, natsrpc.WithServiceNamespace(namespace), natsrpc.WithServiceID(id))
		require.Nil(t, err)

		for k, v := range s.methods {
			require.Equal(t, natsrpc.CombineSubject(namespace, serviceName, id, v.name), k)
		}
	*/
}
