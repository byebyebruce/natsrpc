package example

import (
	"context"
	"fmt"
	"testing"

	"gitlab.uuzu.com/war/natsrpc"
	"gitlab.uuzu.com/war/natsrpc/example/pb"
	"gitlab.uuzu.com/war/natsrpc/example/pb/request"

	"github.com/stretchr/testify/assert"
)

type HeaderSvc struct {
	header string
}

func (h *HeaderSvc) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Hello comes", req.Name, "header:", natsrpc.Header(ctx))
	if h.header != natsrpc.Header(ctx) {
		panic("header error")
	}
	return &pb.HelloReply{
		Message: req.Name,
	}, nil
}

func (h *HeaderSvc) HelloError(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Hello comes", req.Name, "header:", natsrpc.Header(ctx))
	if h.header != natsrpc.Header(ctx) {
		panic("header error")
	}
	return &pb.HelloReply{
		Message: req.Name,
	}, nil
}

func TestHeader(t *testing.T) {
	svc, err := request.RegisterGreeter(server, &HeaderSvc{
		header: haha,
	}, natsrpc.WithServiceNamespace("header"))
	defer svc.Close()
	assert.Nil(t, err)

	cli, err := request.NewGreeterClient(enc, natsrpc.WithClientNamespace("header"))
	assert.Nil(t, err)
	rep, err := cli.Hello(natsrpc.WithHeader(context.Background(), haha), &pb.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)
	assert.Equal(t, haha, rep.GetMessage())
}
