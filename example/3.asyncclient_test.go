package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc/example/pb/async_client"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

type AsyncClientSvc struct{}

func (h AsyncClientSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	fmt.Println("Hello comes", req.Name)
	rp := &testdata.HelloReply{
		Message: req.Name,
	}
	return rp, nil
}
func (h AsyncClientSvc) HelloToAll(ctx context.Context, req *testdata.HelloRequest) (*testdata.Empty, error) {
	fmt.Println("HelloToAll", req.Name)
	return &testdata.Empty{}, nil
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
	svc, err := async_client.RegisterGreeterNATSRPCServer(server, ps)
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := async_client.NewGreeterAsyncClient(conn, d)
	assert.Nil(t, err)

	over := make(chan struct{})

	cli.Hello(context.Background(), &testdata.HelloRequest{Name: haha}, func(reply *testdata.HelloReply, err error) {
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
