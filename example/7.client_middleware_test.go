package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb/request"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

type ClientMiddlewareSvc struct {
}

func (h *ClientMiddlewareSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return &testdata.HelloReply{
		Message: req.Name,
	}, nil
}

func (h *ClientMiddlewareSvc) HelloError(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return nil, fmt.Errorf(haha)
}

func TestClientMiddleware(t *testing.T) {
	cms := &ClientMiddlewareSvc{}
	svc, err := request.RegisterGreeterNATSRPCServer(server, cms)
	assert.Nil(t, err)
	defer svc.Close()

	i := 0
	cli, err := request.NewGreeterNATSRPCClient(enc,
		natsrpc.WithClientMiddleware(func(ctx context.Context, method string, req interface{}, next func(ctx context.Context, req interface{})) {
			i++
			next(ctx, req)
			i++
		}))
	assert.Nil(t, err)

	_, err = cli.Hello(context.Background(), &testdata.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)
	assert.Equal(t, i, 2)

}
