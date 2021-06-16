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

type server interface {
	unregister(*service) bool
}
// service 服务
type service struct {
	name        string               // 名字
	server      server               // server
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
func newService(i interface{}, option Options) (*service, error) {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("service must be a pointer")
	}
	typeName := reflect.Indirect(val).Type().Name()
	if !ast.IsExported(typeName) {
		return nil, fmt.Errorf("service [%s] must be exported", typeName)
	}

	ms, err := parseMethod(i)
	if nil != err {
		return nil, err
	}
	if 0 == len(ms) {
		return nil, fmt.Errorf("service [%s] has no exported method", typeName)
	}

	s := &service{
		val:     val,
		options: option,
		methods: map[string]*method{},
		name:    combineSubject(option.namespace, typeName),
	}

	for _, v := range ms {
		if _, ok := s.methods[v.name]; ok {
			return nil, fmt.Errorf("service [%s] duplicate method [%s]", typeName, v.name)
		}
		sub := combineSubject(s.name, v.name, s.options.id)
		s.methods[sub] = v
	}
	return s, nil
}
