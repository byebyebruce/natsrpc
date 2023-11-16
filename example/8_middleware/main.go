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

	svc, err := example.RegisterGreetingNATSRPCServer(server, &HelloSvc{},
		natsrpc.WithServiceInterceptor(func(ctx context.Context, method string, req interface{}, next natsrpc.Invoker) (interface{}, error) {
			fmt.Println("middle before", method)
			defer fmt.Println("middle after", method)
			fmt.Println("method", method)
			fmt.Println("req", req)
			start := time.Now()
			rep, err := next(ctx, req)
			fmt.Println("elapse", time.Since(start).Milliseconds())
			return rep, err
		}))
	example.IfNotNilPanic(err)
	defer svc.Close()

	cli := example.NewGreetingNATSRPCClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = cli.Hello(ctx, &example.HelloRequest{
		Name: "bruce",
	})
	example.IfNotNilPanic(err)
}

type HelloSvc struct {
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	time.Sleep(time.Millisecond * 100)
	return &example.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
