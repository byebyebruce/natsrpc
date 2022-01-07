package example

import (
	"context"
	"os"
	"sync"
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

const haha = "haha"

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

func (d *asyncDoer) AsyncDo(ctx context.Context, f func(cb func(ret interface{}, err error))) (interface{}, error) {
	done := make(chan struct{})
	once := sync.Once{}
	var (
		ret interface{}
		err error
	)
	cb := func(_ret interface{}, _err error) {
		once.Do(func() {
			ret, err = _ret, _err
			close(done)
		})
	}
	f1 := func() {
		f(cb)
	}

	select {
	case d.c <- f1:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-done:
			return ret, err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
