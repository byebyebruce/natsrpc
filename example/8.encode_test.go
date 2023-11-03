package example

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/pb/request"
	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/stretchr/testify/assert"
)

type testEncoder struct {
}

func (s testEncoder) Decode(b []byte, i interface{}) error {
	return json.Unmarshal(b, i)
}
func (s testEncoder) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

type EncodeSvc struct {
}

func (h *EncodeSvc) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return &testdata.HelloReply{
		Message: "ok",
	}, nil
}

func (h *EncodeSvc) HelloError(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return nil, fmt.Errorf(haha)
}

func TestEncode(t *testing.T) {
	cms1 := &EncodeSvc{}
	svc, err := request.RegisterGreeterNATSRPCServer(server, cms1, natsrpc.WithServiceEncoder(testEncoder{}))
	assert.Nil(t, err)
	defer svc.Close()

	cli, err := request.NewGreeterNATSRPCClient(conn, natsrpc.WithClientEncoder(testEncoder{}))
	assert.Nil(t, err)

	resp, err := cli.Hello(context.Background(), &testdata.HelloRequest{
		Name: haha,
	})
	assert.Nil(t, err)
	assert.Equal(t, "ok", resp.Message)

	cli, err = request.NewGreeterNATSRPCClient(conn)
	assert.Nil(t, err)

	resp, err = cli.Hello(context.Background(), &testdata.HelloRequest{
		Name: haha,
	})
	assert.NotNil(t, err)
}
