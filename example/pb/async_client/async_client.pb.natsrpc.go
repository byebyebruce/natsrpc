// Code generated by protoc-gen-natsrpc. DO NOT EDIT.
// source: async_client.proto

package async_client

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	natsrpc "gitlab.uuzu.com/war/natsrpc"
	pb "gitlab.uuzu.com/war/natsrpc/example/pb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// GreeterInterface
type GreeterInterface interface {
	// Hello
	Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error)
	// HelloToAll
	HelloToAll(ctx context.Context, req *pb.HelloRequest)
}

// RegisterGreeter
func RegisterGreeter(server *natsrpc.Server, s GreeterInterface, opts ...natsrpc.ServiceOption) (natsrpc.IService, error) {
	return server.Register("gitlab.uuzu.com.war.natsrpc.example.pb.async_client.Greeter", s, opts...)
}

// GreeterClient
type GreeterClient interface {
	// Hello
	Hello(ctx context.Context, req *pb.HelloRequest, cb func(*pb.HelloReply, error), opt ...natsrpc.CallOption)
	// HelloToAll
	HelloToAll(notify *pb.HelloRequest, opt ...natsrpc.CallOption) error
}
type _GreeterClient struct {
	c    *natsrpc.Client
	doer natsrpc.AsyncDoer
}

// NewGreeterClient
func NewGreeterClient(enc *nats.EncodedConn, doer natsrpc.AsyncDoer, opts ...natsrpc.ClientOption) (GreeterClient, error) {
	c, err := natsrpc.NewClient(enc, "gitlab.uuzu.com.war.natsrpc.example.pb.async_client.Greeter", opts...)
	if err != nil {
		return nil, err
	}
	ret := &_GreeterClient{
		c:    c,
		doer: doer,
	}
	return ret, nil
}
func (c *_GreeterClient) Hello(ctx context.Context, req *pb.HelloRequest, cb func(*pb.HelloReply, error), opt ...natsrpc.CallOption) {
	go func() {
		rep := &pb.HelloReply{}
		err := c.c.Request(ctx, "Hello", req, rep, opt...)
		newCb := func(_ func(interface{}, error)) {
			cb(rep, err)
		}
		c.doer.AsyncDo(ctx, newCb)
	}()
}
func (c *_GreeterClient) HelloToAll(notify *pb.HelloRequest, opt ...natsrpc.CallOption) error {
	return c.c.Publish("HelloToAll", notify, opt...)
}
