package main

import (
	"context"
	"flag"
	"fmt"
	"sync"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	"github.com/nats-io/nats.go"
)

var (
	nats_url = flag.String("nats_url", "nats://127.0.0.1:4222", "nats-server地址")
)

var wg sync.WaitGroup

func main() {
	conn, err := nats.Connect(*nats_url)
	example.IfNotNilPanic(err)
	defer conn.Close()

	wg.Add(10)
	for i := 0; i < 10; i++ {
		server, err := natsrpc.NewServer(conn)
		example.IfNotNilPanic(err)
		defer server.Close(context.Background())
		svc, err := example.RegisterGreetingToAllNATSRPCServer(server, &HelloSvc{
			name: "svc" + fmt.Sprint(i),
		})
		example.IfNotNilPanic(err)
		defer svc.Close()
	}

	cli := example.NewGreetingToAllNATSRPCClient(conn)

	err = cli.HelloToAll(&example.HelloRequest{
		Name: "bruce",
	})
	example.IfNotNilPanic(err)

	wg.Wait()
	fmt.Println("all received")
}

type HelloSvc struct {
	name string
}

func (s *HelloSvc) HelloToAll(ctx context.Context, req *example.HelloRequest) (*example.Empty, error) {
	fmt.Println(s.name, "receive: ", req.Name)
	defer wg.Done()
	return nil, nil
}
