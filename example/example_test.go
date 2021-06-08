package example_test

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/byebyebruce/natsrpc"
)

type A struct {
}

func (a *A) Func1(ctx context.Context, req *helloworld.HelloRequest) {
	fmt.Println(req.Name)
}

func (a *A) Func2(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply) {
	repl.Message = req.Name
	fmt.Println(repl.Message)
}

func Test_Service(t *testing.T) {
	cfg := &natsrpc.Config{
		Server: "nats://172.25.156.5:4242,nats://172.25.156.5:4252,nats://172.25.156.5:4262",
	}
	server, err := natsrpc.NewServerWithConfig(cfg, "test")
	if nil != err {
		t.Error(err)
	}
	s, err := server.Register(&A{},
		natsrpc.WithNamespace("myspace"),
		natsrpc.WithGroup("mygroup"),
		natsrpc.WithID("1"))

	if nil != err {
		t.Error(err)
	}

	client, _ := natsrpc.NewClient(cfg, "client", 0)
	reply := &helloworld.HelloReply{}
	if err := client.RequestSync(&helloworld.HelloRequest{
		Name: "hello",
	}, reply, "myspace", "group", "1"); nil != err {
		t.Error(err)
	}

	fmt.Println(*reply)
	s.Close()
}
