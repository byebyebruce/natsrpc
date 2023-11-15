package main

import (
	"context"
	"flag"
	"fmt"
	"time"

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

	server, err := natsrpc.NewServer(conn)
	example.IfNotNilPanic(err)

	defer server.Close(context.Background())

	const n = 10

	for i := 0; i < n; i++ {
		server, err := natsrpc.NewServer(conn)
		example.IfNotNilPanic(err)
		defer server.Close(context.Background())
		s := &HelloSvc{id: "svc" + fmt.Sprint(i)}
		svc, err := example.RegisterGreetingNATSRPCServer(server, s,
			natsrpc.WithServiceID(fmt.Sprint(i)))
		example.IfNotNilPanic(err)
		defer svc.Close()
	}

	cli := example.NewGreetingNATSRPCClient(conn)
	for i := 0; i < n; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		rep, err := cli.Hello(ctx, &example.HelloRequest{
			Name: "bruce",
		}, natsrpc.WithCallID(fmt.Sprint(i)))
		example.IfNotNilPanic(err)
		fmt.Println("call", i, rep.Message)
	}

	for i := 0; i < n; i++ {
		cli := example.NewGreetingNATSRPCClient(conn, natsrpc.WithClientID(fmt.Sprint(i)))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		rep, err := cli.Hello(ctx, &example.HelloRequest{
			Name: "bruce",
		})
		example.IfNotNilPanic(err)
		fmt.Println("client", i, rep.Message)
	}
}

type HelloSvc struct {
	id string
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	return &example.HelloReply{
		Message: "hello " + req.Name + " from " + s.id,
	}, nil
}
