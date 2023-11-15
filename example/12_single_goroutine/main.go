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

	server, err := natsrpc.NewServer(conn)
	example.IfNotNilPanic(err)

	defer server.Close(context.Background())

	svc, err := example.RegisterGreetingNATSRPCServer(server, &HelloSvc{},
		natsrpc.WithServiceSingleGoroutine())
	example.IfNotNilPanic(err)
	defer svc.Close()

	cli := example.NewGreetingNATSRPCClient(conn)

	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(name string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			_, err := cli.Hello(ctx, &example.HelloRequest{
				Name: name,
			})
			example.IfNotNilPanic(err)
			//fmt.Println(rep.Message)
		}(strconv.Itoa(i))
	}

	wg.Wait()
}

type HelloSvc struct {
	n int
}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	s.n++
	fmt.Println(s.n) // thread safe
	time.Sleep(time.Millisecond * 10)
	return &example.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
