package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb"
	"github.com/byebyebruce/natsrpc/example/pb/request"
)

type HeaderSvc struct {
}

func (h *HeaderSvc) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Hello comes", req.Name, "header:", natsrpc.Header(ctx))
	return &pb.HelloReply{
		Message: req.Name,
	}, nil
}
func TestHeader(t *testing.T) {
	svc, err := request.RegisterGreeter(server, &HeaderSvc{})
	defer svc.Close()

	cli, err := request.NewGreeterClient(enc)
	natsrpc.IfNotNilPanic(err)
	const haha = "haha"
	rep, err := cli.Hello(natsrpc.WithHeader(context.Background(), "header"), &pb.HelloRequest{
		Name: haha,
	})
	natsrpc.IfNotNilPanic(err)
	if rep.GetMessage() != haha {
		t.Error("not match")
	}
}
