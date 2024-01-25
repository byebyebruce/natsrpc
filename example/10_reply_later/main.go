package main

import (
	"context"
	"flag"
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

	client := natsrpc.NewClient(conn)

	svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})
	example.IfNotNilPanic(err)
	defer svc.Close()

	cli := example.NewGreetingNRClient(client)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	reply, err := cli.Hello(ctx, &example.HelloRequest{
		Name: "bruce",
	})
	example.IfNotNilPanic(err)

	println(reply.Message)
}

type HelloSvc struct {
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	go func() {
		time.Sleep(time.Second)
		natsrpc.Reply(ctx, &example.HelloReply{
			Message: "hello " + req.Name + " reply later",
		}, nil)
	}()
	return nil, natsrpc.ErrReplyLater
}
