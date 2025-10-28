package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	"github.com/go-kratos/kratos/v2/transport"
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
	client, err := natsrpc.NewClient(conn)

	defer server.Close(context.Background())

	svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})
	example.IfNotNilPanic(err)
	defer svc.Close()

	cli := example.NewGreetingNRClient(client)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := cli.Hello(ctx, &example.HelloRequest{
		Name: "bruce",
	}, natsrpc.WithCallHeader(map[string][]string{"key": {"value"}}))
	example.IfNotNilPanic(err)

	println(reply.Message)
}

type HelloSvc struct {
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	tr, ok := transport.FromServerContext(ctx)
	if ok {
		fmt.Println("transport header", tr.RequestHeader())
	}
	header := tr.RequestHeader()
	fmt.Println("call header", header)

	return &example.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
