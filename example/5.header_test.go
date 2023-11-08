package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb/header"
	"github.com/byebyebruce/natsrpc/testdata"
	"github.com/stretchr/testify/require"
)

type HeaderSvc struct {
	header map[string]string
}

func (h *HeaderSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	fmt.Println("Hello comes", req.Name, "header:", natsrpc.GetHeader(ctx))
	hd := natsrpc.GetHeader(ctx)
	if h.header[haha] != hd[haha] {
		panic("header error")
	}
	return &testdata.HelloReply{
		Message: req.Name,
	}, nil
}

func (h *HeaderSvc) HelloPublish(ctx context.Context, req *testdata.HelloRequest) (*testdata.Empty, error) {
	fmt.Println("Hello comes", req.Name, "header:", natsrpc.GetHeader(ctx))
	hd := natsrpc.GetHeader(ctx)
	if h.header[haha] != hd[haha] {
		panic("header error")
	}
	return &testdata.Empty{}, nil
}

func TestHeader(t *testing.T) {
	hs := &HeaderSvc{
		header: map[string]string{haha: haha},
	}
	svc, err := header.RegisterGreeterNATSRPCServer(server, hs, natsrpc.WithServiceNamespace("header"))
	defer svc.Close()
	require.Nil(t, err)

	cli := header.NewGreeterNATSRPCClient(conn, natsrpc.WithClientNamespace("header"))
	require.Nil(t, err)

	rep, err := cli.Hello(context.Background(), &testdata.HelloRequest{
		Name: haha,
	}, natsrpc.WithCallHeader(hs.header))
	require.Nil(t, err)
	require.Equal(t, haha, rep.GetMessage())

	err = cli.HelloPublish(&testdata.HelloRequest{
		Name: haha,
	}, natsrpc.WithCallHeader(hs.header))
	require.Nil(t, err)
}
