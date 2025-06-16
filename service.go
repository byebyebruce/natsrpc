package natsrpc

import (
	"context"
	"fmt"
)

var _ ServiceInterface = (*Service)(nil)

type ServerInterface interface {
	Encoder
	Remove(string) bool
}

// Service 服务
type Service struct {
	sd      ServiceDesc           // 描述
	val     interface{}           // 值
	server  ServerInterface       // rpc
	methods map[string]MethodDesc // 方法集合
	opt     ServiceOptions
}

// Name 名字
func (s *Service) Name() string {
	return joinSubject(s.opt.namespace, s.sd.ServiceName, s.opt.id)
}

// Close 关闭
// 会取消所有订阅
func (s *Service) Close() bool {
	return s.server.Remove(s.Name())
}

// NewService 创建服务
func NewService(server ServerInterface, sd ServiceDesc, i interface{}, options ServiceOptions) (*Service, error) {
	methods := map[string]MethodDesc{}
	for _, md := range sd.Methods {
		if _, ok := methods[md.MethodName]; ok {
			return nil, fmt.Errorf("service [%s] duplicate method [%s]", sd.ServiceName, md.MethodName)
		}
		methods[md.MethodName] = md
	}
	s := &Service{
		methods: methods,
		sd:      sd,
		val:     i,
		server:  server,
		opt:     options,
	}

	return s, nil
}

func (s *Service) Call(ctx context.Context, methodName string, dec func(any) error, interceptor Interceptor) ([]byte, error) {
	m, ok := s.methods[methodName]
	if !ok {
		return nil, ErrNoMethod
	}
	//req := m.NewRequest()
	resp, err := s.call(ctx, m, dec, interceptor)
	if err != nil {
		return nil, err
	}
	if !m.IsPublish {
		if resp == nil {
			return nil, nil
		}
		return s.server.Encode(resp)
	}
	return nil, nil
}

func (s *Service) call(ctx context.Context, m MethodDesc, dec func(any) error, interceptor Interceptor) (interface{}, error) {
	if interceptor == nil {
		return m.Handler(s.val, ctx, dec)
	} else {
		invoker := func(ctx1 context.Context, req1 interface{}) (interface{}, error) {
			return m.Handler(s.val, ctx1, dec)
		}
		return interceptor(ctx, m.MethodName, dec, invoker)
	}
}
