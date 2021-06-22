package natsrpc

import (
	"context"
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
	opt         Options              // 设置
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

	ms, err := parseMethod(reflect.TypeOf(i))
	if nil != err {
		return nil, err
	}
	if len(ms) == 0 {
		return nil, fmt.Errorf("service [%s] has no exported method", name)
	}

	s := &service{
		val:     val,
		opt:     opt,
		methods: map[string]*method{},
		name:    name,
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

func (s *service) call(ctx context.Context, m *method, b []byte) (interface{}, error) {
	req, err := m.newRequest(b)
	if nil != err {
		return nil, err
	}

	fn := func() {
		if s.opt.recoverHandler != nil {
			defer func() {
				if e := recover(); e != nil {
					s.opt.recoverHandler(e)
				}
			}()
		}
		m.handle(ctx, s.val, req)
	}

	if s.opt.isSingleThreadMode() { // 单线程处理
		select {
		case <-ctx.Done():
		case s.opt.singleThreadCbChan <- fn:
		}
	} else { // 多线程处理
		go fn()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-req.over:
		return req.reply, req.err
	}
}
