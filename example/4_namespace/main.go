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

	const namespace = "example"

	server, err := natsrpc.NewServer(conn,
		natsrpc.WithServerNamespace(namespace))
	example.IfNotNilPanic(err)

	defer server.Close(context.Background())

	svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{namespace: namespace})
	example.IfNotNilPanic(err)
	defer svc.Close()

	client, err := natsrpc.NewClient(conn, natsrpc.WithClientNamespace(namespace))
	example.IfNotNilPanic(err)

	cli := example.NewGreetingNRClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = cli.Hello(ctx, &example.HelloRequest{
		Name: "bruce",
	})
	example.IfNotNilPanic(err)

	client, err = natsrpc.NewClient(conn, natsrpc.WithClientNamespace("wrong_namespace"))
	example.IfNotNilPanic(err)

	cli1 := example.NewGreetingNRClient(client)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = cli1.Hello(ctx, &example.HelloRequest{
		Name: "bruce",
	})
	fmt.Println("should be error", err.Error())
}

type HelloSvc struct {
	namespace string
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	return &example.HelloReply{
		Message: "hello " + req.Name + " from " + s.namespace,
	}, nil
}
