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

type rpc interface {
	unregister(*service) bool
}

// service 服务
type service struct {
	name        string               // 名字 namespace.package.struct
	rpc         rpc                  // rpc
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
	return s.rpc.unregister(s)
}

// newService 创建服务
func newService(name string, i interface{}, option Options) (*service, error) {
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
		options: option,
		methods: map[string]*method{},
		name:    CombineSubject(option.namespace, name),
	}

	for _, v := range ms {
		if _, ok := s.methods[v.name]; ok {
			return nil, fmt.Errorf("service [%s] duplicate method [%s]", name, v.name)
		}
		sub := CombineSubject(s.name, v.name, s.options.id)
		s.methods[sub] = v
	}
	return s, nil
}
