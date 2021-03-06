// Code generated by protoc-gen-natsrpc. DO NOT EDIT.
// versions:
// - protoc-gen-natsrpc v0.5.0
// source: async_service.proto

package async_service

import (
	context "context"
	fmt "fmt"
	natsrpc "github.com/byebyebruce/natsrpc"
	testdata "github.com/byebyebruce/natsrpc/testdata"
	nats_go "github.com/nats-io/nats.go"
	proto "google.golang.org/protobuf/proto"
)

var _ = new(context.Context)
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = natsrpc.Version
var _ = nats_go.Version

type GreeterNATSRPCServer interface {
	Hello(ctx context.Context, req *testdata.HelloRequest, cb func(*testdata.HelloReply, error))
	HelloToAll(ctx context.Context, req *testdata.HelloRequest)
}

// RegisterGreeterNATSRPCServer register Greeter service
func RegisterAsyncGreeter(server *natsrpc.Server, doer natsrpc.AsyncDoer, s GreeterNATSRPCServer, opts ...natsrpc.ServiceOption) (natsrpc.IService, error) {
	ss := &GreeterWrapper{
		doer: doer,
		s:    s,
	}
	return server.Register("github.com.byebyebruce.natsrpc.example.pb.async_service.Greeter", ss, opts...)
}

// GreeterWrapper DO NOT USE
type GreeterWrapper struct {
	doer natsrpc.AsyncDoer
	s    GreeterNATSRPCServer
}

// Hello DO NOT USE
func (s *GreeterWrapper) Hello(ctx context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	f := func(cb func(interface{}, error)) {
		s.s.Hello(ctx, req, func(r *testdata.HelloReply, e error) {
			cb(r, e)
		})
	}
	temp, err := s.doer.AsyncDo(ctx, f)
	if temp == nil {
		return nil, err
	}
	return temp.(*testdata.HelloReply), err
}

// HelloToAll DO NOT USE
func (s *GreeterWrapper) HelloToAll(ctx context.Context, req *testdata.HelloRequest) {
	s.doer.AsyncDo(ctx, func(_ func(interface{}, error)) {
		s.s.HelloToAll(ctx, req)
	})
}

type GreeterNATSRPCClient interface {
	Hello(ctx context.Context, req *testdata.HelloRequest, opt ...natsrpc.CallOption) (*testdata.HelloReply, error)
	HelloToAll(notify *testdata.HelloRequest, opt ...natsrpc.CallOption) error
}

type _GreeterNATSRPCClient struct {
	c *natsrpc.Client
}

// NewGreeterNATSRPCClient
func NewGreeterNATSRPCClient(enc *nats_go.EncodedConn, opts ...natsrpc.ClientOption) (GreeterNATSRPCClient, error) {
	c, err := natsrpc.NewClient(enc, "github.com.byebyebruce.natsrpc.example.pb.async_service.Greeter", opts...)
	if err != nil {
		return nil, err
	}
	ret := &_GreeterNATSRPCClient{
		c: c,
	}
	return ret, nil
}
func (c *_GreeterNATSRPCClient) Hello(ctx context.Context, req *testdata.HelloRequest, opt ...natsrpc.CallOption) (*testdata.HelloReply, error) {
	rep := &testdata.HelloReply{}
	err := c.c.Request(ctx, "Hello", req, rep, opt...)
	return rep, err
}
func (c *_GreeterNATSRPCClient) HelloToAll(notify *testdata.HelloRequest, opt ...natsrpc.CallOption) error {
	return c.c.Publish("HelloToAll", notify, opt...)
}
