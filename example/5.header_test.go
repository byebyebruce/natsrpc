package example

import (
	"context"
	"fmt"
	"testing"

	"gitlab.uuzu.com/war/natsrpc"
	"gitlab.uuzu.com/war/natsrpc/example/pb"
	"gitlab.uuzu.com/war/natsrpc/example/pb/request"
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
func TestHeader(t *testing.T) {
	const haha = "haha"
	svc, err := request.RegisterGreeter(server, &HeaderSvc{
		header: haha,
	}, natsrpc.WithServiceNamespace("header"))
	defer svc.Close()

	cli, err := request.NewGreeterClient(enc, natsrpc.WithClientNamespace("header"))
	natsrpc.IfNotNilPanic(err)
	rep, err := cli.Hello(natsrpc.WithHeader(context.Background(), haha), &pb.HelloRequest{
		Name: haha,
	})
	natsrpc.IfNotNilPanic(err)
	if rep.GetMessage() != haha {
		t.Error("not match")
	}
}
