package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/nats-io/nats.go"
)

var (
	nats_url = flag.String("nats_url", "nats://127.0.0.1:4222", "nats-server地址")
)

func main() {
	conn, err := nats.Connect(*nats_url)
	example.IfNotNilPanic(err)
	defer conn.Close()

	logger := log.NewStdLogger(os.Stdout)
	mw := logging.Client(logger)
	client, err := natsrpc.NewClient(conn, natsrpc.WithClientMiddleware(mw))
	example.IfNotNilPanic(err)
	defer client.Close()

	//svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})
	//example.IfNotNilPanic(err)
	//defer svc.Close()
	//select {}

	cli := example.NewGreetingNRClient(client)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	reply, err := cli.Hello(ctx, &example.HelloRequest{
		Name: "bruce",
	})
	example.IfNotNilPanic(err)
	println("reply", reply.Message)

	// unsub
	//svc.Close()

}
