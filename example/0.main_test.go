package example

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.uuzu.com/war/natsrpc"
	"gitlab.uuzu.com/war/natsrpc/tool/nats_server"
)

var (
	enc    *nats.EncodedConn
	server *natsrpc.Server
)

func TestMain(m *testing.M) {
	var err error

	s := nats_server.Run(nil)
	defer s.Shutdown()

	enc, err = natsrpc.NewPBEnc(s.ClientURL())
	natsrpc.IfNotNilPanic(err)
	defer enc.Close()

	server, err = natsrpc.NewServer(enc)
	natsrpc.IfNotNilPanic(err)
	defer server.Close(time.Second)

	os.Exit(m.Run())
}

type asyncDoer struct {
	c chan func()
}

func (d *asyncDoer) Do(ctx context.Context, f func()) {
	select {
	case d.c <- f:
	case <-ctx.Done():
	}
}
