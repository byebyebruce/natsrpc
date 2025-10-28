package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	"github.com/go-kratos/kratos/v2/middleware"
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

	mw := func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			tr, _ := transport.FromServerContext(ctx)
			method := tr.Operation()

			fmt.Println("middle before", method)
			defer fmt.Println("middle after", method)
			fmt.Println("method", method)
			fmt.Println("req", req)
			start := time.Now()
			rep, err := handler(ctx, req)
			fmt.Println("elapse", time.Since(start).Milliseconds())
			return rep, err
		}
	}
	server, err := natsrpc.NewServer(conn, natsrpc.ServerMiddleware(mw))
	example.IfNotNilPanic(err)
	defer server.Close(context.Background())
	client, err := natsrpc.NewClient(conn)

	svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})

	example.IfNotNilPanic(err)
	defer svc.Close()

	cli := example.NewGreetingNRClient(client)

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
