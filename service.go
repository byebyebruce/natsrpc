package natsrpc

import (
	"context"
	"fmt"
)

var _ IService = (*Service)(nil)

type IServer interface {
	Encoder
	Remove(string) bool
}

// Service 服务
type Service struct {
	sd      ServiceDesc           // 描述
	val     interface{}           // 值
	server  IServer               // rpc
	methods map[string]MethodDesc // 方法集合
}

// Name 名字
func (s *Service) Name() string {
	return s.sd.ServiceName
}

// Close 关闭
// 会取消所有订阅
func (s *Service) Close() bool {
	return s.server.Remove(s.Name())
}

// NewService 创建服务
func NewService(server IServer, sd ServiceDesc, i interface{}) (*Service, error) {
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
	}

	return s, nil
}

func (s *Service) Call(ctx context.Context, methodName string, b []byte, interceptor Interceptor) ([]byte, error) {
	m, ok := s.methods[methodName]
	if !ok {
		return nil, ErrNoMethod
	}
	req := m.NewRequest()
	if err := s.server.Decode(b, req); err != nil {
		return nil, err
	}
	resp, err := s.call(ctx, m, req, interceptor)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		b1, err1 := s.server.Encode(resp)
		return b1, err1
	}
	return nil, err
}

func (s *Service) call(ctx context.Context, m MethodDesc, req interface{}, interceptor Interceptor) (interface{}, error) {
	if interceptor == nil {
		return m.Handler(s.val, ctx, req)
	} else {
		invoker := func(ctx1 context.Context, req1 interface{}) (interface{}, error) {
			return m.Handler(s.val, ctx1, req1)
		}
		return interceptor(ctx, m.MethodName, req, invoker)
	}
}
