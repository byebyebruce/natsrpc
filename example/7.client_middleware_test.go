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
	header map[string]string
	ts     *testing.T
}

func (h *ClientMiddlewareSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	hd := natsrpc.Header(ctx)
	assert.Equal(h.ts, h.header[haha], hd[haha])
	return &testdata.HelloReply{
		Message: req.Name,
	}, nil
}

func (h *ClientMiddlewareSvc) HelloError(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return nil, fmt.Errorf(haha)
}

func TestClientMiddleware(t *testing.T) {
	cms := &ClientMiddlewareSvc{
		ts:     t,
		header: map[string]string{haha: haha},
	}
	svc, err := request.RegisterGreeter(server, cms)
	assert.Nil(t, err)
	defer svc.Close()

	ch := natsrpc.WithCallHeader(map[string]string{haha: haha})
	cli, err := request.NewGreeterClient(enc,
		natsrpc.WithClientMiddleware(func(ctx context.Context, sub string, req interface{}, options *natsrpc.CallOptions) {
			ch(options)
		}))
	assert.Nil(t, err)

	_, err = cli.Hello(context.Background(), &testdata.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)

}
