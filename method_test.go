package natsrpc

import (
	"context"
	"reflect"
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
func (a *MethodTest) AsyncRequest(ctx context.Context, req *pb.HelloRequest, f func(*pb.HelloReply, error)) {
	repl := &pb.HelloReply{
		Message: req.Name,
	}
	f(repl, nil)
}

type MethodErrorTest struct {
}

func (a *MethodErrorTest) Func1(repl *pb.HelloReply) {

}

func Test_Parse(t *testing.T) {
	m, err := parseMethod(reflect.TypeOf(&MethodTest{}))
	if nil != err {
		t.Error(err)
	}

	if _, ok := m["Publish"]; !ok {
		t.Error("name error")
	}
	_, err = parseMethod(reflect.TypeOf(&MethodErrorTest{}))
	if nil == err {
		t.Error(err)
	}
}

func TestMethod_Handle(t *testing.T) {
	s := &MethodTest{}
	ret, err := parseMethod(reflect.TypeOf(s))
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
	req, err := m.newRequest(b)
	if err != nil {
		t.Error(err)
	}
	m.handle(context.Background(), reflect.ValueOf(s), req)
	if nil != err {
		t.Error(err)
	}

	if req.reply.(*pb.HelloReply).Message != param {
		t.Error("reply.Message!=param")
	}
}

func BenchmarkMethod_Handle(b *testing.B) {
	s := &MethodTest{}
	ret, err := parseMethod(reflect.TypeOf(s))
	if nil != err {
		b.Error(err)
	}
	param := "hello"
	a := &pb.HelloRequest{Name: param}
	bs, _ := proto.Marshal(a)

	h := ret["Request"]
	val := reflect.ValueOf(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := h.newRequest(bs)
		h.handle(context.Background(), val, req)
	}
}
