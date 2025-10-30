package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	"github.com/nats-io/nats.go"
)

var (
	nats_url = flag.String("nats_url", "nats://127.0.0.1:4222", "nats-server地址")
)

func main() {
	conn, err := nats.Connect(*nats_url)
	example.IfNotNilPanic(err)
	defer conn.Close()

	logger := log.NewStdLogger(os.Stdout)
	wm := logging.Server(logger)
	server, err := natsrpc.NewServer(conn, natsrpc.ServerMiddleware(wm))
	example.IfNotNilPanic(err)

	defer server.Close(context.Background())

	svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})
	example.IfNotNilPanic(err)
	defer svc.Close()
	select {}
}

type HelloSvc struct {
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	fmt.Println("Server Handle Hello")
	return &example.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
