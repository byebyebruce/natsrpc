package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb"
	async "github.com/byebyebruce/natsrpc/example/pb/async_service"
)

type asyncDoer struct {
	c chan func()
}

func (d *asyncDoer) Do(ctx context.Context, f func()) {
	select {
	case d.c <- f:
	case <-ctx.Done():
	}
}

type AsyncServiceSvc struct{}

func (h AsyncServiceSvc) Hello(ctx context.Context, req *pb.HelloRequest, cb func(*pb.HelloReply, error)) {
	fmt.Println("Hello comes", req.Name)
	rp := &pb.HelloReply{
		Message: req.Name,
	}
	cb(rp, nil)
	cb(rp, nil) // is ok
}
func (h AsyncServiceSvc) HelloToAll(ctx context.Context, req *pb.HelloRequest) {
	fmt.Println("HelloToAll", req.Name)
}

func TestAsyncService(t *testing.T) {
	d := &asyncDoer{
		c: make(chan func()),
	}
	go func() {
		for f := range d.c {
			f()
		}
	}()
	ps := &AsyncServiceSvc{}
	svc, err := async.RegisterGreeter(server, ps, d)
	defer svc.Close()

	cli, err := async.NewGreeterClient(enc)
	natsrpc.IfNotNilPanic(err)
	const haha = "haha"
	reply, err := cli.Hello(context.Background(), &pb.HelloRequest{Name: haha})
	fmt.Println(reply, err)
	if reply.Message != haha {
		t.Error("not match")
	}

	cli.HelloToAll(&pb.HelloRequest{Name: haha})

	natsrpc.IfNotNilPanic(err)
	time.Sleep(time.Millisecond * 100)
}
