package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc/example/pb"
	"github.com/byebyebruce/natsrpc/example/pb/publish"

	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := publish.NewGreeterClient(enc)
	assert.Nil(t, err)

	err = cli.HelloToAll(&pb.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, ps.name, haha)
}
