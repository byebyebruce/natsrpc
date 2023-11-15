package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"sync"
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

	for i := 0; i < 3; i++ {
		server, err := natsrpc.NewServer(conn)
		example.IfNotNilPanic(err)

		defer server.Close(context.Background())

		svc, err := example.RegisterGreetingNATSRPCServer(server, &HelloSvc{
			name: "svc" + strconv.Itoa(i),
		})
		example.IfNotNilPanic(err)
		defer svc.Close()
	}

	cli := example.NewGreetingNATSRPCClient(conn)

	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(name string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			rep, err := cli.Hello(ctx, &example.HelloRequest{
				Name: name,
			})
			example.IfNotNilPanic(err)
			fmt.Println(rep.Message)
		}(strconv.Itoa(i))
	}

	wg.Wait()
}

type HelloSvc struct {
	name string
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	return &example.HelloReply{
		Message: "hello " + req.Name + " from " + s.name,
	}, nil
}
