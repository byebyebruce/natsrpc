package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/byebyebruce/natsrpc/testdata"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	"github.com/nats-io/nats.go"
)

var (
	server              = flag.String("server", "nats://127.0.0.1:4222", "nats server")
	namespace           = flag.String("ns", "testsapce", "namespace")
	group               = flag.String("g", "", "subscribe group")
	id                  = flag.String("id", "", "service id")
	singleThreadService = flag.Bool("st", false, "service single thread handle")
)

func main() {
	flag.Parse()

	cfg := natsrpc.Config{
		Server: *server,
	}

	server, err := natsrpc.NewServerWithConfig(cfg, nats.Name("example_server"+*id))
	if nil != err {
		panic(err)
	}
	defer server.Close()

	opts := []natsrpc.Option{
		natsrpc.WithNamespace(*namespace),
		natsrpc.WithGroup(*group),
		natsrpc.WithID(*id),
		natsrpc.WithTimeout(time.Second),
		natsrpc.WithRecoveryHandler(func(e interface{}) {
			fmt.Println(e)
		})}

	if *singleThreadService {
		fnChan := make(chan func())
		go func() {
			for f := range fnChan {
				f()
			}
		}()
		opts = append(opts, natsrpc.WithSingleThreadCallback(fnChan))
	}

	s, err := testdata.RegisterGreeter(server, &example.ExampleGreeter{}, opts...)
	if nil != err {
		panic(err)
	}
	defer s.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println(s.Name(), "start")
	<-sig
}
