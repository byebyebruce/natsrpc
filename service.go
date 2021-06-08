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
	methods     map[string]*method
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
		methods: map[string]*method{},
	}

	s.val = reflect.ValueOf(i)
	typeName := reflect.Indirect(s.val).Type().Name()
	if !ast.IsExported(typeName) {
		log.Fatalf("rpc server: %s is not a valid s name", typeName)
	}

	s.name = fmt.Sprintf("%s.%s", option.namespace, typeName)
	if "" != option.id {
		s.name += "." + option.id
	}
	ms, err := parseStruct(i)
	if nil != err {
		return nil, err
	}

	for _, v := range ms {
		if "" == v.name {
			return nil, fmt.Errorf("method is empty %v", *v)
		}
		if _, ok := s.methods[v.name]; ok {
			return nil, fmt.Errorf("method [%s] duplicate", v.name)
		}
		s.methods[v.name] = v
	}

	return s, nil
}
