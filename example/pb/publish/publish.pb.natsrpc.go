// Code generated by protoc-gen-natsrpc. DO NOT EDIT.
// source: publish.proto

package publish

import (
	context "context"
	fmt "fmt"
	natsrpc "github.com/byebyebruce/natsrpc"
	testdata "github.com/byebyebruce/natsrpc/testdata"
	proto "github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
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

// GreeterService Greeter service interface
type GreeterService interface {
	// HelloToAll call HelloToAll
	HelloToAll(ctx context.Context, req *testdata.HelloRequest)
}

// RegisterGreeter register Greeter service
func RegisterGreeter(server *natsrpc.Server, s GreeterService, opts ...natsrpc.ServiceOption) (natsrpc.IService, error) {
	return server.Register("github.com.byebyebruce.natsrpc.example.pb.publish.Greeter", s, opts...)
}

// GreeterClient
type GreeterClient interface {
	// HelloToAll
	HelloToAll(notify *testdata.HelloRequest, opt ...natsrpc.CallOption) error
}

type _GreeterClient struct {
	c *natsrpc.Client
}

// NewGreeterClient
func NewGreeterClient(enc *nats.EncodedConn, opts ...natsrpc.ClientOption) (GreeterClient, error) {
	c, err := natsrpc.NewClient(enc, "github.com.byebyebruce.natsrpc.example.pb.publish.Greeter", opts...)
	if err != nil {
		return nil, err
	}
	ret := &_GreeterClient{
		c: c,
	}
	return ret, nil
}
func (c *_GreeterClient) HelloToAll(notify *testdata.HelloRequest, opt ...natsrpc.CallOption) error {
	return c.c.Publish("HelloToAll", notify, opt...)
}
