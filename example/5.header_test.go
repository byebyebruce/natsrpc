package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb/header"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

type HeaderSvc struct {
	header map[string]string
}

func (h *HeaderSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	fmt.Println("Hello comes", req.Name, "header:", natsrpc.Header(ctx))
	hd := natsrpc.Header(ctx)
	if h.header[haha] != hd[haha] {
		panic("header error")
	}
	return &testdata.HelloReply{
		Message: req.Name,
	}, nil
}

func (h *HeaderSvc) HelloPublish(ctx context.Context, req *testdata.HelloRequest) {
	fmt.Println("Hello comes", req.Name, "header:", natsrpc.Header(ctx))
	hd := natsrpc.Header(ctx)
	if h.header[haha] != hd[haha] {
		panic("header error")
	}
}

func TestHeader(t *testing.T) {
	hs := &HeaderSvc{
		header: map[string]string{haha: haha},
	}
	svc, err := header.RegisterGreeter(server, hs, natsrpc.WithServiceNamespace("header"))
	defer svc.Close()
	assert.Nil(t, err)

	cli, err := header.NewGreeterClient(enc, natsrpc.WithClientNamespace("header"))
	assert.Nil(t, err)

	rep, err := cli.Hello(context.Background(), &testdata.HelloRequest{
		Name: haha,
	}, natsrpc.WithCallHeader(hs.header))
	assert.Nil(t, err)
	assert.Equal(t, haha, rep.GetMessage())

	err = cli.HelloPublish(&testdata.HelloRequest{
		Name: haha,
	}, natsrpc.WithCallHeader(hs.header))
	assert.Nil(t, err)
}
