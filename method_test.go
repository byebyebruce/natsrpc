package natsrpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
)

type testMarshaller struct {
}

func (s testMarshaller) Unmarshal(b []byte, i interface{}) error {
	return proto.Unmarshal(b, i.(proto.Message))
}
func (s testMarshaller) Marshal(i interface{}) ([]byte, error) {
	return proto.Marshal(i.(proto.Message))
}

var pbMarshaller = testMarshaller{}

type MethodTest struct {
}

func (a *MethodTest) Publish(ctx context.Context, req *Empty) {
	_ = req
}

func (a *MethodTest) Request(ctx context.Context, req *Empty) (*Empty, error) {
	repl := &Empty{}
	return repl, nil
}
func (a *MethodTest) AsyncRequest(ctx context.Context, req *Empty, f func(*Empty, error)) {
	repl := &Empty{}
	f(repl, nil)
}

type MethodErrorTest struct {
}

func (a *MethodErrorTest) Func1(repl *Empty) {

}

func Test_Parse(t *testing.T) {
	m, err := parseMethod(&MethodTest{}, pbMarshaller)
	if nil != err {
		t.Error(err)
	}

	if _, ok := m["Publish"]; !ok {
		t.Error("name error")
	}
	_, err = parseMethod(&MethodErrorTest{}, pbMarshaller)
	if nil == err {
		t.Error(err)
	}
}

func TestMethod_Handle(t *testing.T) {
	ret, err := parseMethod(&MethodTest{}, pbMarshaller)
	if nil != err {
		t.Error(err)
	}
	a := &Empty{}
	b, _ := pbMarshaller.Marshal(a)

	var m *method
	for _, v := range ret {
		if v.name == "Request" {
			m = v
			break
		}
	}
	if m == nil {
		t.Error("m is nil")
		return
	}
	req, err := m.newRequest(b)
	if err != nil {
		t.Error(err)
	}
	m.handle(context.Background(), req)
	if nil != err {
		t.Error(err)
	}

}

func BenchmarkMethod_Handle(b *testing.B) {
	s := &MethodTest{}
	ret, err := parseMethod(s, pbMarshaller)
	if nil != err {
		b.Error(err)
	}
	a := &Empty{}
	bs, _ := proto.Marshal(a)

	h := ret["Request"]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := h.newRequest(bs)
		h.handle(context.Background(), req)
	}
}
