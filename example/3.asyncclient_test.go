package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gitlab.uuzu.com/war/natsrpc/example/pb"
	"gitlab.uuzu.com/war/natsrpc/example/pb/async_client"

	"github.com/stretchr/testify/assert"
)

type AsyncClientSvc struct{}

func (h AsyncClientSvc) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Hello comes", req.Name)
	rp := &pb.HelloReply{
		Message: req.Name,
	}
	return rp, nil
}
func (h AsyncClientSvc) HelloToAll(ctx context.Context, req *pb.HelloRequest) {
	fmt.Println("HelloToAll", req.Name)
}

func TestAsyncClient(t *testing.T) {
	d := &asyncDoer{
		c: make(chan func()),
	}
	go func() {
		for f := range d.c {
			f()
		}
	}()
	ps := &AsyncClientSvc{}
	svc, err := async_client.RegisterGreeter(server, ps)
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := async_client.NewGreeterClient(enc, d)
	assert.Nil(t, err)

	over := make(chan struct{})

	cli.Hello(context.Background(), &pb.HelloRequest{Name: haha}, func(reply *pb.HelloReply, err error) {
		defer close(over)
		assert.Nil(t, err)
		assert.Equal(t, haha, reply.Message)
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	select {
	case <-over:
	case <-ctx.Done():
		t.Error(ctx.Err())
	}
}
