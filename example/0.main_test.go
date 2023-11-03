package example

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/nats-io/nats.go"
)

var (
	conn   *nats.Conn
	server *natsrpc.Server
)

const haha = "haha"

func TestMain(m *testing.M) {
	var err error
	//conn, err = natsrpc.NewPBEnc(s.ClientURL())
	conn, err = nats.Connect("nats://127.0.0.1")
	natsrpc.IfNotNilPanic(err)
	defer conn.Close()

	server, err = natsrpc.NewServer(conn)
	natsrpc.IfNotNilPanic(err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	defer server.Close(ctx)

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
