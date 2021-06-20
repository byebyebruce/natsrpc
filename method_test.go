package natsrpc

import (
	"context"
	"testing"

	"github.com/byebyebruce/natsrpc/testdata/pb"
	"github.com/golang/protobuf/proto"
)

type MethodTest struct {
}

func (a *MethodTest) Publish(ctx context.Context, req *pb.HelloRequest) {
	_ = req.Name
}

func (a *MethodTest) Request(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	repl := &pb.HelloReply{
		Message: req.Name,
	}
	return repl, nil
}

type MethodErrorTest struct {
}

func (a *MethodErrorTest) Func1(repl *pb.HelloReply) {

}

func Test_Parse(t *testing.T) {
	m, err := parseMethod(&MethodTest{})
	if nil != err {
		t.Error(err)
	}

	if _, ok := m["Publish"]; !ok {
		t.Error("name error")
	}
	_, err = parseMethod(&MethodErrorTest{})
	if nil == err {
		t.Error(err)
	}
}

func TestMethod_Handle(t *testing.T) {
	ret, err := parseMethod(&MethodTest{})
	if nil != err {
		t.Error(err)
	}
	param := "hello"
	a := &pb.HelloRequest{Name: param}
	b, _ := proto.Marshal(a)

	var m *method
	for _, v := range ret {
		if v.name == "Request" {
			m = v
			break
		}
	}
	if m == nil {
		t.Error("m is nil")
	}
	reply, err := m.handle(context.Background(), b)
	if nil != err {
		t.Error(err)
	}

	if reply.(*pb.HelloReply).Message != param {
		t.Error("reply.Message!=param")
	}
}

func BenchmarkMethod_Handle(b *testing.B) {
	ret, err := parseMethod(&MethodTest{})
	if nil != err {
		b.Error(err)
	}
	param := "hello"
	a := &pb.HelloRequest{Name: param}
	bs, _ := proto.Marshal(a)
	b.ResetTimer()
	h := ret["Request"]
	for i := 0; i < b.N; i++ {
		_, _ = h.handle(context.Background(), bs)
	}
}
