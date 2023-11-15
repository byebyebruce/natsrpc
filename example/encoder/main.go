package main

import (
	"context"
	"encoding/json"
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

type Encoder struct {
	name string
}

func (e Encoder) Encode(v interface{}) ([]byte, error) {
	fmt.Println(e.name, "json encode")
	return json.Marshal(v)
}

func (e Encoder) Decode(data []byte, vPtr interface{}) error {
	fmt.Println(e.name, "json decode")
	return json.Unmarshal(data, vPtr)
}

func main() {
	conn, err := nats.Connect(*nats_url)
	example.IfNotNilPanic(err)
	defer conn.Close()

	server, err := natsrpc.NewServer(conn, natsrpc.WithServerEncoder(Encoder{"server"}))
	example.IfNotNilPanic(err)

	defer server.Close(context.Background())

	svc, err := example.RegisterGreetingNATSRPCServer(server, &HelloSvc{})
	example.IfNotNilPanic(err)
	defer svc.Close()

	cli := example.NewGreetingNATSRPCClient(conn, natsrpc.WithClientEncoder(Encoder{"client"}))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
	return &example.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
