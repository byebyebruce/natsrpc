package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc/example/pb/publish"
	"github.com/byebyebruce/natsrpc/testdata"
	"github.com/stretchr/testify/require"
)

type PublishSvc struct {
	name string
}

func (h *PublishSvc) HelloToAll(ctx context.Context, req *testdata.HelloRequest) (*testdata.Empty, error) {
	fmt.Println("Hello to all", req.Name)
	h.name = req.Name
	return &testdata.Empty{}, nil
}

func (h *PublishSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.Empty, error) {
	fmt.Println("Hello to all", req.Name)
	h.name = req.Name
	return &testdata.Empty{}, nil
}
func TestPublish(t *testing.T) {
	ps := &PublishSvc{}
	svc, err := publish.RegisterGreeterNATSRPCServer(server, ps)
	require.Nil(t, err)
	defer svc.Close()

	cli := publish.NewGreeterNATSRPCClient(conn)
	require.Nil(t, err)

	err = cli.HelloToAll(&testdata.HelloRequest{
		Name: haha,
	})
	require.Nil(t, err)

	time.Sleep(time.Millisecond * 100)
	require.Equal(t, ps.name, haha)
}
