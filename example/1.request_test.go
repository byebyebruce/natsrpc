package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb"
	"github.com/byebyebruce/natsrpc/example/pb/request"
)

type RequestSvc struct{}

func (h RequestSvc) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Hello comes", req.Name)
	return &pb.HelloReply{
		Message: req.Name,
	}, nil
}
func TestRequest(t *testing.T) {
	svc, err := request.RegisterGreeter(server, &RequestSvc{})
	defer svc.Close()

	cli, err := request.NewGreeterClient(enc)
	natsrpc.IfNotNilPanic(err)
	const haha = "haha"
	rep, err := cli.Hello(context.Background(), &pb.HelloRequest{
		Name: haha,
	})
	natsrpc.IfNotNilPanic(err)
	if rep.GetMessage() != haha {
		t.Error("not match")
	}
}
