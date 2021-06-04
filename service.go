package xnats

import (
	"reflect"

	"github.com/nats-io/nats.go"
)

type Service struct {
	name        string
	s           *Server
	typ         reflect.Type
	rcvr        reflect.Value
	subscribers []*nats.Subscription
}

func (s *Service) Close() {
	for _, v := range s.subscribers {
		v.Unsubscribe()
	}
}
