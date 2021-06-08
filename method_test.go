package natsrpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type MethodTest struct {
}

func (a *MethodTest) Func2(ctx context.Context, req *helloworld.HelloRequest, repl *helloworld.HelloReply) {
	repl.Message = req.Name
	fmt.Println(repl.Message)
}

func TestParse(t *testing.T) {
	ret, err := parseStruct(&MethodTest{})
	if nil != err {
		t.Error(err)
	}
	a := &helloworld.HelloRequest{Name: "req"}
	b, _ := proto.Marshal(a)
	for _, v := range ret {

		fmt.Println(v.handler(context.Background(), b).(*helloworld.HelloReply))
	}

}
