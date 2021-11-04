package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb"
	"github.com/byebyebruce/natsrpc/example/pb/publish"
)

type PublishSvc struct {
	name string
}

func (h *PublishSvc) HelloToAll(ctx context.Context, req *pb.HelloRequest) {
	fmt.Println("Hello to all", req.Name)
	h.name = req.Name
}
func TestPublish(t *testing.T) {
	ps := &PublishSvc{}
	svc, err := publish.RegisterGreeter(server, ps)
	defer svc.Close()

	cli, err := publish.NewGreeterClient(enc)
	natsrpc.IfNotNilPanic(err)
	const haha = "haha"
	err = cli.HelloToAll(context.Background(), &pb.HelloRequest{
		Name: haha,
	})
	natsrpc.IfNotNilPanic(err)
	time.Sleep(time.Millisecond * 100)
	if ps.name != haha {
		t.Error("not match")
	}
}
