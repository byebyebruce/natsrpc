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

	"github.com/stretchr/testify/assert"
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
	opt := natsrpc.WithServiceMiddleware(func(ctx context.Context, method string, req interface{}, next func(context.Context, interface{})) error {
		if "HelloError" == method {
			atomic.AddInt32(&errorCount, 1)
			return fmt.Errorf(haha + haha)
		}
		start := time.Now()
		next(ctx, req)
		elapse := time.Now().Sub(start)
		fmt.Println(method, "elapse:", elapse)
		return nil
	})
	svc, err := request.RegisterGreeter(server, &ServiceMiddlewareSvc{}, opt)
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := request.NewGreeterClient(enc)
	assert.Nil(t, err)

	rep, err := cli.Hello(context.Background(), &testdata.HelloRequest{})
	assert.Nil(t, err)
	assert.Equal(t, haha, rep.Message)

	rep, err = cli.HelloError(context.Background(), &testdata.HelloRequest{})
	assert.NotNil(t, err)
	assert.EqualValues(t, errorCount, 1)
	assert.Equal(t, haha+haha, err.Error())
}
