package natsrpc

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/nats-io/nats.go"
)

type Service interface {
	Name() string
	Close() bool
}

// service 服务
type service struct {
	name        string               // 名字 package.struct
	server      *Server              // rpc
	val         reflect.Value        // 值
	subscribers []*nats.Subscription // nats订阅
	methods     map[string]*method   // 方法集合
	options     Options              // 设置
}

// 名字
func (s *service) Name() string {
	return s.name
}

// Close 关闭
// 会取消所有订阅
func (s *service) Close() bool {
	return s.server.unregister(s)
}

// newService 创建服务
func newService(name string, i interface{}, opts ...Option) (*service, error) {
	opt := MakeOptions(opts...)

	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("service must be a pointer")
	}
	typ := reflect.Indirect(val).Type()
	if !ast.IsExported(typ.Name()) {
		return nil, fmt.Errorf("service [%s] must be exported", name)
	}

	ms, err := parseMethod(i)
	if nil != err {
		return nil, err
	}
	if 0 == len(ms) {
		return nil, fmt.Errorf("service [%s] has no exported method", name)
	}

	s := &service{
		val:     val,
		options: opt,
		methods: map[string]*method{},
		name:    name,
	}

	for _, v := range ms {
		if _, ok := s.methods[v.name]; ok {
			return nil, fmt.Errorf("service [%s] duplicate method [%s]", name, v.name)
		}
		// subject = namespace.package.service.method.id
		subject := CombineSubject(s.options.namespace, s.name, v.name, s.options.id)
		s.methods[subject] = v
	}
	return s, nil
}
