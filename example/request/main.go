package main

import (
	"context"
	"fmt"
	"time"

	"github.com/byebyebruce/natsrpc/example/pb"

	"github.com/byebyebruce/natsrpc"
	request "github.com/byebyebruce/natsrpc/example/request/pb"
	"github.com/byebyebruce/natsrpc/extension/simpleserver"
)

func main() {
	s := simpleserver.Run(nil)
	defer s.Shutdown()

	enc, err := natsrpc.NewPBEnc(s.ClientURL())
	natsrpc.IfNotNilPanic(err)
	defer enc.Close()

	server, err := natsrpc.NewServer(enc)
	natsrpc.IfNotNilPanic(err)
	defer server.Close(time.Second)

	svc, err := request.RegisterGreeter(server, &HelloSvc{})
	defer svc.Close()

	cli, err := request.NewGreeterClient(enc)
	natsrpc.IfNotNilPanic(err)

	err = cli.ToAll(context.Background(), &natsrpc.Empty{})
	natsrpc.IfNotNilPanic(err)

	const haha = "haha"
	rep, err := cli.Hello(context.Background(), &request.HelloRequest{
		Name: haha,
	})
	natsrpc.IfNotNilPanic(err)
	if rep.GetMessage() != haha {
		panic("not match")
	}
}

type HelloSvc struct {
}

func (h HelloSvc) Hello(ctx context.Context, req *request.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Hello comes", req.Name)
	return &pb.HelloReply{
		Message: req.Name,
	}, nil
}

func (h HelloSvc) ToAll(ctx context.Context, req *natsrpc.Empty) {
	fmt.Println("To all", req.String())
}
