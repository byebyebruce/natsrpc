package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/byebyebruce/natsrpc/example/pb"

	"github.com/byebyebruce/natsrpc"
)

var (
	url       = flag.String("url", "nats://127.0.0.1:4222", "nats server")
	namespace = flag.String("ns", "testsapce", "namespace")
	group     = flag.String("g", "", "subscribe group")
	id        = flag.String("id", "", "service id")
)

func main() {
	flag.Parse()

	enc, err := natsrpc.NewPBEnc(*url)
	natsrpc.IfNotNilPanic(err)
	defer enc.Close()

	server, err := natsrpc.NewServer(enc)
	natsrpc.IfNotNilPanic(err)
	defer server.Close(time.Second)

	opts := []natsrpc.ServiceOption{
		natsrpc.WithNamespace(*namespace),
		natsrpc.WithServiceGroup(*group),
		natsrpc.WithID(*id),
		natsrpc.WithTimeout(time.Second),
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

type ExampleGreeter struct {
}

// HiAll publish
func (a *ExampleGreeter) HiAll(ctx context.Context, req *pb.HelloRequest) {
	fmt.Println("begin HiAll Notify->", req.Name)
	fmt.Println("end HiAll Notify->", req.Name)
}

// AreYouOK request
func (a *ExampleGreeter) AreYouOK(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("begin AreYouOK Request", req.Name)
	rep := &pb.HelloReply{
		Message: "AreYouOK " + req.Name,
	}
	fmt.Println("end AreYouOK Request->", req.Name)
	return rep, nil
}

// DelayAreYouOK async request
func (a *ExampleGreeter) DelayAreYouOK(ctx context.Context, req *pb.HelloRequest, f func(*pb.HelloReply, error)) {
	fmt.Println("begin DelayAreYouOK Request", req.Name)
	rep := &pb.HelloReply{
		Message: "DelayAreYouOK " + req.Name,
	}
	fmt.Println("end DelayAreYouOK Request->", req.Name)
	f(rep, nil)
}
