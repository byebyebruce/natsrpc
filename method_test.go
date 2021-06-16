package natsrpc

import (
	"context"
	"testing"

	"github.com/byebyebruce/natsrpc/testdata/pb"
	"github.com/golang/protobuf/proto"
)

type MethodTest struct {
}

func (a *MethodTest) Func1(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply) {
	repl.Message = req.Name
}

type MethodErrorTest struct {
}

func (a *MethodErrorTest) Func1(repl *pb.HelloReply) {

}

func Test_Parse(t *testing.T) {
	_, err := parseMethod(&MethodTest{})
	if nil != err {
		t.Error(err)
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
	for _, v := range ret {
		reply, err := v.handle(context.Background(), b)
		if nil != err {
			t.Error(err)
		}

		if reply.(*pb.HelloReply).Message != param {
			t.Error("reply.Message!=param")
		}
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
	for i := 0; i < b.N; i++ {
		_, err := ret[0].handle(context.Background(), bs)
		if nil != err {
			b.Error(err)
		}
	}
}
