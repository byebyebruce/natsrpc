package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb/request"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

type RequestSvc struct {
	idx int32
	t   *testing.T
}

func (h *RequestSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	assert.Equal(h.t, haha, req.Name)
	return &testdata.HelloReply{
		Message: haha + haha,
	}, nil
}

func (h *RequestSvc) HelloError(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return nil, fmt.Errorf(haha)
}

func TestRequest1(t *testing.T) {
	const (
		ns = "mysapce"
		id = 1234
	)
	svc, err := request.RegisterGreeterNATSRPCServer(server, &RequestSvc{t: t},
		natsrpc.WithServiceNamespace(ns),
		natsrpc.WithServiceID(id))
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := request.NewGreeterNATSRPCClient(conn,
		natsrpc.WithClientNamespace(ns),
		natsrpc.WithClientID(id))
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rep, err := cli.Hello(ctx, &testdata.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)
	assert.Equal(t, haha+haha, rep.Message)

	rep, err = cli.HelloError(ctx, &testdata.HelloRequest{
		Name: haha,
	})
	assert.NotNil(t, err)
	assert.Equal(t, haha, err.Error())

	/*
		rep, err = cli.Hello(context.Background(), &testdata.HelloRequest{
			Name: haha,
		}, natsrpc.WithCallNamespace("errornamespace"),
			natsrpc.WithCallTimeout(time.Millisecond*100))
		assert.NotNil(t, err)

	*/
}
func TestRequest(t *testing.T) {
	const (
		ns = "mysapce"
		id = 1234
	)
	svc, err := request.RegisterGreeterNATSRPCServer(server, &RequestSvc{t: t},
		natsrpc.WithServiceNamespace(ns),
		natsrpc.WithServiceID(id))
	assert.Nil(t, err)
	defer svc.Close()
	fmt.Println("svc:", svc.Name())

	cli, err := request.NewGreeterNATSRPCClient(conn,
		natsrpc.WithClientNamespace(ns),
		natsrpc.WithClientID(id))
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	rep, err := cli.Hello(ctx, &testdata.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)
	assert.Equal(t, haha+haha, rep.Message)

	rep, err = cli.HelloError(ctx, &testdata.HelloRequest{
		Name: haha,
	})
	assert.NotNil(t, err)
	assert.Equal(t, haha, err.Error())

	/*
		rep, err = cli.Hello(context.Background(), &testdata.HelloRequest{
			Name: haha,
		}, natsrpc.WithCallNamespace("errornamespace"),
			natsrpc.WithCallTimeout(time.Millisecond*100))
		assert.NotNil(t, err)

	*/
}
