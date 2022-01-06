package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gitlab.uuzu.com/war/natsrpc"
	"gitlab.uuzu.com/war/natsrpc/example/pb"
	"gitlab.uuzu.com/war/natsrpc/example/pb/request"

	"github.com/stretchr/testify/assert"
)

type RequestSvc struct {
	idx int32
	t   *testing.T
}

func (h *RequestSvc) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	assert.Equal(h.t, haha, req.Name)
	return &pb.HelloReply{
		Message: haha + haha,
	}, nil
}

func (h *RequestSvc) HelloError(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return nil, fmt.Errorf(haha)
}

func TestRequest(t *testing.T) {
	const (
		ns = "mysapce"
		id = 1234
	)
	svc, err := request.RegisterGreeter(server, &RequestSvc{t: t},
		natsrpc.WithServiceNamespace(ns),
		natsrpc.WithServiceID(id))
	assert.Nil(t, err)
	defer svc.Close()
	fmt.Println("svc:", svc.Name())

	cli, err := request.NewGreeterClient(enc,
		natsrpc.WithClientNamespace(ns),
		natsrpc.WithClientID(id))
	assert.Nil(t, err)

	rep, err := cli.Hello(context.Background(), &pb.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)
	assert.Equal(t, haha+haha, rep.Message)

	rep, err = cli.HelloError(context.Background(), &pb.HelloRequest{
		Name: haha,
	})
	assert.NotNil(t, err)
	assert.Equal(t, haha, err.Error())

	rep, err = cli.Hello(context.Background(), &pb.HelloRequest{
		Name: haha,
	}, natsrpc.WithCallNamespace("errornamespace"),
		natsrpc.WithCallTimeout(time.Millisecond*100))
	assert.NotNil(t, err)
}
