package natsrpc

import (
	"context"
	"fmt"
	"go/ast"
	"reflect"
)

var _ IService = (*service)(nil)

// service 服务
type service struct {
	name    string // 名字
	sub     string
	val     interface{}        // 值
	server  *Server            // rpc
	methods map[string]*method // 方法集合
	opt     serviceOptions     // 设置
}

// Name 名字
func (s *service) Name() string {
	return s.name
}

// Close 关闭
// 会取消所有订阅
func (s *service) Close() bool {
	return s.server.remove(s)
}

// newService 创建服务
func newService(name string, i interface{}, opts ...ServiceOption) (*service, error) {
	opt := defaultServiceOptions
	for _, v := range opts {
		v(&opt)
	}

	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("service must be a pointer")
	}
	typ := reflect.Indirect(val).Type()
	if !ast.IsExported(typ.Name()) {
		return nil, fmt.Errorf("service [%s] must be exported", name)
	}

	s := &service{
		opt:     opt,
		methods: map[string]*method{},
		name:    name,
		sub:     CombineSubject(opt.namespace, name, opt.id), // name = namespace.package.service.id
		val:     i,
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

func (s *service) call(ctx context.Context, methodName string, b []byte) ([]byte, error) {
	m, ok := s.methods[methodName]
	if !ok {
		return nil, ErrNoMethod
	}
	req := m.newRequest()
	if err := s.opt.encoder.Decode(b, req); err != nil {
		return nil, err
	}
	resp, err := s._call(ctx, m, req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		b, err := s.opt.encoder.Encode(resp)
		return b, err
	}
	return nil, err
}

func (s *service) _call(ctx context.Context, m *method, req interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, s.opt.timeout)
	defer cancel()

	var (
		resp interface{}
		err  error
	)
	if s.opt.mw != nil {
		next := func(ctx1 context.Context, req1 interface{}) {
			resp, err = m.handle(s.val, ctx1, req1)
		}
		if err := s.opt.mw(ctx, m.name, req, next); err != nil {
			return nil, err
		}
	} else {
		resp, err = m.handle(s.val, ctx, req)
	}

	return resp, err
}
