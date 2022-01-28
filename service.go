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
	name    string             // 名字
	val     interface{}        // 值
	server  *Server            // rpc
	methods map[string]*method // 方法集合
	opt     serviceOptions     // 设置
}

// 名字
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
		name:    CombineSubject(opt.namespace, name, opt.id), // name = namespace.package.service.id
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
		// subject = namespace.package.service.id.method
		subject := CombineSubject(s.name, v.name)
		s.methods[subject] = v
	}
	return s, nil
}

func (s *service) handle(ctx context.Context, m *method, sub string, b []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, s.opt.timeout)
	defer cancel()

	req := m.newRequest()

	if len(b) > 0 {
		rpcReq := &Request{}
		if err := s.server.enc.Enc.Decode(sub, b, rpcReq); nil != err {
			return nil, err
		}
		if len(rpcReq.Header) > 0 {
			ctx = withHeader(ctx, rpcReq.Header)
		}
		if len(rpcReq.Payload) > 0 {
			if err := s.server.enc.Enc.Decode(sub, rpcReq.Payload, req); nil != err {
				return nil, err
			}
		}
	}

	var (
		resp interface{}
		err  error
	)
	if s.opt.mw != nil {
		if err := s.opt.mw(ctx, m.name, req); err != nil {
			return nil, err
		}
	}
	resp, err = m.handle(s.val, ctx, req)
	if err != nil {
		return nil, err
	}
	return s.server.enc.Enc.Encode(sub, resp)
}
