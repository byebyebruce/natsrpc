package natsrpc

import (
	"context"
	"fmt"
	"go/ast"
	"reflect"
)

// Service 服务
type Service interface {
	Name() string
	Close() bool
}

// service 服务
type service struct {
	name    string             // 名字 package.struct
	server  *Server            // rpc
	methods map[string]*method // 方法集合
	opt     Options            // 设置
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

	s := &service{
		opt:     opt,
		methods: map[string]*method{},
		name:    name,
	}

	ms, err := parseMethod(i, s)
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
		// subject = namespace.package.service.method.id
		subject := CombineSubject(s.opt.namespace, s.name, v.name, s.opt.id)
		s.methods[subject] = v
	}
	return s, nil
}

func (s *service) call(ctx context.Context, m *method, b []byte) ([]byte, error) {
	if s.opt.recoverHandler != nil {
		defer func() {
			if e := recover(); e != nil {
				s.opt.recoverHandler(e)
			}
		}()
	}

	ctx, cancel := context.WithTimeout(ctx, s.opt.timeout)
	defer cancel()

	req, err := m.newRequest(b)
	if nil != err {
		return nil, err
	}

	m.handle(ctx, req)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-req.over:
		if err != nil {
			return nil, err
		}
		return s.Marshal(req.reply)
	}
}

func (s *service) Unmarshal(b []byte, i interface{}) error {
	return s.server.enc.Enc.Decode("", b, i)
}

func (s *service) Marshal(i interface{}) ([]byte, error) {
	return s.server.enc.Enc.Encode("", i)
}
