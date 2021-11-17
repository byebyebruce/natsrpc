package example

import (
	"os"
	"testing"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/extension/simpleserver"
	"github.com/nats-io/nats.go"
)

var (
	enc    *nats.EncodedConn
	server *natsrpc.Server
)

func TestMain(m *testing.M) {
	var err error

	s := simpleserver.Run(nil)
	defer s.Shutdown()

	enc, err = natsrpc.NewPBEnc(s.ClientURL())
	natsrpc.IfNotNilPanic(err)
	defer enc.Close()

	server, err = natsrpc.NewServer(enc)
	natsrpc.IfNotNilPanic(err)
	defer server.Close(time.Second)

	os.Exit(m.Run())
}
