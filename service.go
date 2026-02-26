package natsrpc

import (
	"context"
	"fmt"
)

// Service 服务
type Service struct {
	sd      ServiceDesc           // 描述
	val     interface{}           // 值
	methods map[string]MethodDesc // 方法集合
	opt     ServiceOptions
	server  *Server
}

// Name 名字
func (s *Service) Name() string {
	return joinSubject(s.server.opt.namespace, s.sd.ServiceName, s.opt.id)
}

// Close 关闭
// 会取消所有订阅
func (s *Service) Close() bool {
	return s.server.Remove(s.Name())
}

// NewService 创建服务
func NewService(server *Server, sd ServiceDesc, i interface{}, options ServiceOptions) (*Service, error) {
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

func (s *Service) call(ctx context.Context, methodName string, dec func(any) error) (any, error) {
	m, ok := s.methods[methodName]
	if !ok {
		return nil, ErrNoMethod
	}
	return m.Handler(s.val, ctx, dec)
}
