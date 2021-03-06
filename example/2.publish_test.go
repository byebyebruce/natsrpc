package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc/example/pb/publish"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

type PublishSvc struct {
	name string
}

func (h *PublishSvc) HelloToAll(ctx context.Context, req *testdata.HelloRequest) {
	fmt.Println("Hello to all", req.Name)
	h.name = req.Name
}
func TestPublish(t *testing.T) {
	ps := &PublishSvc{}
	svc, err := publish.RegisterGreeterNATSRPCServer(server, ps)
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := publish.NewGreeterNATSRPCClient(enc)
	assert.Nil(t, err)

	err = cli.HelloToAll(&testdata.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, ps.name, haha)
}
