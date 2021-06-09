package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
)

var (
	server    = flag.String("server", "nats://127.0.0.1:4222", "nats server")
	namespace = flag.String("ns", "testsapce", "namespace")
	group     = flag.String("g", "", "subscribe group")
	id        = flag.String("id", "", "service id")
)

func main() {
	flag.Parse()

	cfg := natsrpc.Config{
		Server: *server,
	}

	server, err := natsrpc.NewServerWithConfig(cfg, "example_server"+*id)
	if nil != err {
		panic(err)
	}
	defer server.Close()

	opts := []natsrpc.Option{
		natsrpc.WithNamespace(*namespace),
		natsrpc.WithGroup(*group),
		natsrpc.WithID(*id),
		natsrpc.WithTimeout(time.Second)}
	s, err := server.Register(&example.MyService{}, opts...)
	if nil != err {
		panic(err)
	}
	defer s.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println(s.Name(), "start")
	<-sig
}
