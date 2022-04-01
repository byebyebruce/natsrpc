package natsrpc

import (
	"context"
	"testing"

	"github.com/byebyebruce/natsrpc/testdata"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
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

func (a *MethodTest) Publish(ctx context.Context, req *testdata.Empty) {
	_ = req
}

func (a *MethodTest) Request(ctx context.Context, req *testdata.Empty) (*testdata.Empty, error) {
	repl := &testdata.Empty{}
	return repl, nil
}

type MethodErrorTest struct {
}

func (a *MethodErrorTest) Func1(repl *testdata.Empty) {

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
	assert.NotNil(t, err)
}

func TestMethod_Handle(t *testing.T) {
	s := &MethodTest{}
	ret, err := parseMethod(s)
	assert.Nil(t, err)

	a := &testdata.Empty{}
	b, _ := pbMarshaller.Marshal(a)

	m, ok := ret["Request"]
	assert.Equal(t, true, ok)

	req := m.newRequest()
	pbMarshaller.Unmarshal(b, req)
	_, err = m.handle(s, context.Background(), req)
	assert.Nil(t, err)
}

func BenchmarkMethod_Handle(b *testing.B) {
	s := &MethodTest{}
	ret, err := parseMethod(s)
	if nil != err {
		b.Error(err)
	}

	h := ret["Request"]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := h.newRequest()
		h.handle(s, context.Background(), req)
	}
}
