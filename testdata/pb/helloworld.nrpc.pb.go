

// Greeter
type Greeter interface {
	// AreYouOKAsync
	AreYouOKAsync(ctx context.Context, req *HelloRequest, reply func(*HelloReply, error))
}


// RegisterGreeter
func RegisterGreeter(server *natsrpc.Server, s Greeter, opts ...natsrpc.Option) (natsrpc.Service, error) {
	return server.Register("xxx.Greeter", s, opts...)
}


// GreeterClient
type GreeterClient struct {
	c *natsrpc.Client
}

// NewGreeterClient
func NewGreeterClient(enc *nats.EncodedConn, opts ...natsrpc.Option) (*GreeterClient, error) {
	c, err := natsrpc.NewClient(enc, "xxx.Greeter", opts...)
	if err != nil {
		return nil, err
	}
	ret := &GreeterClient{
		c:c,
	}
	return ret, nil
}

// ID 根据ID获得client
func (c *GreeterClient) ID(id interface{}) *GreeterClient {
	return &GreeterClient{
		c : c.c.ID(id),
	}
}

// AreYouOKAsync
func (c *GreeterClient) AreYouOKAsync(req *HelloRequest, cb func(*HelloReply, error)){
	rep := &HelloReply{}
	f := func(_ proto.Message, err error) {
		cb(rep, err)
	}
	c.c.AsyncRequest("AreYouOK", req, rep, f)
}

