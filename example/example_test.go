package example

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/byebyebruce/natsrpc"
)

type A struct {
}

func (a *A) Func2(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply) {
	repl.Message = req.Name
	fmt.Println("Func2", repl.Message)
}

func Test_Service(t *testing.T) {
	ns := RunServer(nil)
	defer ns.Shutdown()
	cfg := &natsrpc.Config{
		Server: "nats://127.0.0.1:4222",
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
	for i := 0; i <= 10; i++ {
		reply := &helloworld.HelloReply{}
		if err := client.RequestSync(&helloworld.HelloRequest{
			Name: fmt.Sprintf("hello %d", i),
		}, reply, "myspace", "A", "1"); nil != err {
			t.Error(err)
		}

		fmt.Println("reply", reply.Message)
	}

	s.Close()
}
