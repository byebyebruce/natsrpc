package example

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb/request"
	"github.com/byebyebruce/natsrpc/testdata"
	"github.com/stretchr/testify/require"
)

type ServiceMiddlewareSvc struct {
}

func (h *ServiceMiddlewareSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return &testdata.HelloReply{Message: haha}, nil
}

func (h *ServiceMiddlewareSvc) HelloError(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return &testdata.HelloReply{Message: haha}, nil
}

func TestServiceMiddleware(t *testing.T) {
	var errorCount int32 = 0
	opt := natsrpc.WithServiceMiddleware(func(ctx context.Context, method string, req interface{}, next natsrpc.Invoker) (any, error) {
		if "HelloError" == method {
			atomic.AddInt32(&errorCount, 1)
			return nil, fmt.Errorf(haha + haha)
		}
		start := time.Now()
		ret, err := next(ctx, req)
		elapse := time.Now().Sub(start)
		fmt.Println(method, "elapse:", elapse)
		return ret, err
	})
	svc, err := request.RegisterGreeterNATSRPCServer(server, &ServiceMiddlewareSvc{}, opt)
	require.Nil(t, err)
	defer svc.Close()

	cli := request.NewGreeterNATSRPCClient(conn)
	require.Nil(t, err)

	rep, err := cli.Hello(context.Background(), &testdata.HelloRequest{})
	require.Nil(t, err)
	require.Equal(t, haha, rep.Message)

	rep, err = cli.HelloError(context.Background(), &testdata.HelloRequest{})
	require.NotNil(t, err)
	require.EqualValues(t, errorCount, 1)
	require.Equal(t, haha+haha, err.Error())
}
