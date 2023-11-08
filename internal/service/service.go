package service

import (
	"context"
	"fmt"
	"go/ast"
	"reflect"

	"github.com/byebyebruce/natsrpc"
)

var _ natsrpc.IService = (*Service)(nil)

type Server interface {
	natsrpc.Encoder
	Remove(natsrpc.IService) bool
}

// Service 服务
type Service struct {
	name    string             // 名字
	val     interface{}        // 值
	server  Server             // rpc
	methods map[string]*method // 方法集合
	natsrpc.ServiceOptions
}

// Name 名字
func (s *Service) Name() string {
	return s.name
}

// Close 关闭
// 会取消所有订阅
func (s *Service) Close() bool {
	return s.server.Remove(s)
}

// NewService 创建服务
func NewService(server Server, name string, i interface{}, opt natsrpc.ServiceOptions) (*Service, error) {
	/*
		opt:=natsrpc.
		for _, v := range opts {
			v(&opt)
		}
	*/

	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("service must be a pointer")
	}
	typ := reflect.Indirect(val).Type()
	if !ast.IsExported(typ.Name()) {
		return nil, fmt.Errorf("service [%s] must be exported", name)
	}

	s := &Service{
		ServiceOptions: opt,
		methods:        map[string]*method{},
		name:           name,
		val:            i,
		server:         server,
	}

	ms, err := parseMethod(i)
	if nil != err {
		return nil, err
	}
	if len(ms) == 0 {
		return nil, fmt.Errorf("service [%s] has no exported method", name)
	}

	for _, v := range ms {
		if _, ok := s.methods[v.name]; ok {
			return nil, fmt.Errorf("service [%s] duplicate method [%s]", name, v.name)
		}
		s.methods[v.name] = v
	}
	return s, nil
}

func (s *Service) Call(ctx context.Context, methodName string, b []byte, interceptor natsrpc.Interceptor) ([]byte, error) {
	m, ok := s.methods[methodName]
	if !ok {
		return nil, natsrpc.ErrNoMethod
	}
	req := m.newRequest()
	if err := s.server.Decode(b, req); err != nil {
		return nil, err
	}
	resp, err := s._call(ctx, m, req, interceptor)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		b, err := s.server.Encode(resp)
		return b, err
	}
	return nil, err
}

func (s *Service) _call(ctx context.Context, m *method, req interface{}, interceptor natsrpc.Interceptor) (interface{}, error) {
	if interceptor != nil {
		next := func(ctx1 context.Context, req1 interface{}) (interface{}, error) {
			return m.handle(s.val, ctx1, req1)
		}
		return interceptor(ctx, m.name, req, next)
	} else {
		return m.handle(s.val, ctx, req)
	}
}
