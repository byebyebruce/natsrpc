package example

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb/request"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

var handleCnt int32 = 0

type OneSvc struct {
	id int
}

func (h *OneSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	atomic.AddInt32(&handleCnt, 1)
	// 只有1号svc回复
	if h.id != 1 {
		return nil, nil
	}
	return &testdata.HelloReply{
		Message: fmt.Sprint(h.id),
	}, nil
}

func (h *OneSvc) HelloError(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return nil, fmt.Errorf(haha)
}

func TestOneReply(t *testing.T) {
	cms1 := &OneSvc{id: 1}
	svc, err := request.RegisterGreeterNATSRPCServer(server, cms1, natsrpc.WithBroadcast())
	assert.Nil(t, err)
	defer svc.Close()

	cms2 := &OneSvc{id: 2}
	svc, err = request.RegisterGreeterNATSRPCServer(server, cms2, natsrpc.WithBroadcast())
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := request.NewGreeterNATSRPCClient(enc)
	assert.Nil(t, err)

	resp, err := cli.Hello(context.Background(), &testdata.HelloRequest{
		Name: haha,
	})
	time.Sleep(time.Millisecond * 10)
	assert.Nil(t, err)
	assert.Equal(t, "1", resp.Message)
	assert.EqualValues(t, 2, handleCnt)

}
