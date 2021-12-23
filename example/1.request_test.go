package example

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"gitlab.uuzu.com/sanguox/natsrpc"
	"gitlab.uuzu.com/sanguox/natsrpc/example/pb"
	"gitlab.uuzu.com/sanguox/natsrpc/example/pb/request"
)

type RequestSvc struct {
	idx int32
}

func (h *RequestSvc) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Hello comes", req.Name)
	i := atomic.AddInt32(&h.idx, 1)
	if i == 1 {
		return &pb.HelloReply{
			Message: req.Name,
		}, nil
	} else {
		return nil, fmt.Errorf(req.Name)
	}
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

	rep, err = cli.Hello(context.Background(), &pb.HelloRequest{
		Name: haha,
	})
	if err == nil {
		t.Errorf("should not be nil")
	}
}
