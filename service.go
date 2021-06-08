package natsrpc

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"

	"github.com/nats-io/nats.go"
)

type service struct {
	name        string
	server      *Server
	val         reflect.Value
	subscribers []*nats.Subscription
	methods     []*method
	options     serviceOptions
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Close() {
	for _, v := range s.subscribers {
		v.Unsubscribe()
	}
	s.subscribers = nil
	s.server.Unregister(s)
}

func newService(server *Server, i interface{}, option serviceOptions) (*service, error) {
	s := &service{
		server:  server,
		options: option,
	}

	s.val = reflect.ValueOf(i)
	name := reflect.Indirect(s.val).Type().Name()
	if !ast.IsExported(name) {
		log.Fatalf("rpc server: %s is not a valid s name", name)
	}

	s.name = fmt.Sprintf("%s.%s", option.namespace, name)
	if "" != option.id {
		s.name += "." + option.id
	}
	ms, err := parseStruct(i)
	if nil != err {
		return nil, err
	}
	s.methods = ms
	return s, nil
}
