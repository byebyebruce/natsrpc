package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc/example/pb/async_service"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

type AsyncServiceSvc struct{}

func (h AsyncServiceSvc) Hello(ctx context.Context, req *testdata.HelloRequest, cb func(*testdata.HelloReply, error)) {
	fmt.Println("Hello comes", req.Name)
	rp := &testdata.HelloReply{
		Message: req.Name,
	}
	cb(rp, nil)
	cb(rp, nil) // is ok
}
func (h AsyncServiceSvc) HelloToAll(ctx context.Context, req *testdata.HelloRequest) {
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
	svc, err := async_service.RegisterAsyncGreeter(server, d, ps)
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := async_service.NewGreeterNATSRPCClient(enc)
	assert.Nil(t, err)

	reply, err := cli.Hello(context.Background(), &testdata.HelloRequest{Name: haha})
	assert.Nil(t, err)
	fmt.Println(reply, err)
	assert.Equal(t, haha, reply.Message)

	cli.HelloToAll(&testdata.HelloRequest{Name: haha})

	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 100)
}
