package natsrpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

type serviceWrapper struct {
	*Service
	subscriptions []*nats.Subscription
	ServiceOptions
}

var _ ServiceRegistrar = (*Server)(nil)

// Server RPC server
type Server struct {
	wg       sync.WaitGroup             // wait group
	mu       sync.Mutex                 // lock
	opt      ServerOptions              // options
	conn     *nats.Conn                 // NATS Encode Conn
	services map[string]*serviceWrapper // 服务 name->Service
	Encoder
}

// NewServer 构造器
func NewServer(conn *nats.Conn, option ...ServerOption) (*Server, error) {
	if !conn.IsConnected() {
		return nil, fmt.Errorf("conn is not connected")
	}

	options := DefaultServerOptions
	for _, v := range option {
		v(&options)
	}

	d := &Server{
		opt:      options,
		conn:     conn,
		services: map[string]*serviceWrapper{},
		Encoder:  options.encoder,
	}
	return d, nil
}

// Close 关闭
func (s *Server) Close(ctx context.Context) (err error) {
	s.UnSubscribeAll()

	over := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(over)
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-over:
	}
	return s.conn.FlushWithContext(ctx)
}

// UnSubscribeAll 取消所有订阅
func (s *Server) UnSubscribeAll() error {
	unsubs := make([]*nats.Subscription, 0, len(s.services))
	s.mu.Lock()
	for _, svc := range s.services {
		unsubs = append(unsubs, svc.subscriptions...)
		svc.subscriptions = nil
	}
	s.mu.Unlock()
	for _, v := range unsubs {
		v.Unsubscribe()
	}
	return s.conn.Flush()
}

// Remove 反注册
func (s *Server) Remove(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	svc, ok := s.services[name]
	if ok {
		if len(svc.subscriptions) > 0 {
			for _, subscription := range svc.subscriptions {
				subscription.Unsubscribe()
			}
			s.conn.Flush()
		}
		delete(s.services, name)
	}
	return ok
}

// Register 注册服务
func (s *Server) Register(sd ServiceDesc, val interface{}, opts ...ServiceOption) (IService, error) {
	opt := DefaultServiceOptions
	for _, v := range opts {
		v(&opt)
	}

	sd.ServiceName = joinSubject(opt.namespace, sd.ServiceName, opt.id)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[sd.ServiceName]; ok {
		return nil, ErrDuplicateService
	}
	// new 一个服务
	svc, err := NewService(s, sd, val)
	if nil != err {
		return nil, err
	}

	sw := &serviceWrapper{
		Service:        svc,
		ServiceOptions: opt,
	}
	if err := s.subscribeMethod(sw); nil != err {
		return nil, err
	}

	s.services[sd.ServiceName] = sw

	// TODO flush
	if err := s.conn.Flush(); err != nil {
		return nil, err
	}

	return svc, nil
}

// subscribeMethod 订阅服务的方法
func (s *Server) subscribeMethod(sw *serviceWrapper) error {
	cb := func(msg *nats.Msg) {
		s.wg.Add(1)
		call := func() {
			defer s.wg.Done()

			method, header, err := decodeHeader(msg.Header)
			if err != nil {
				s.opt.errorHandler(err.Error())
				return
			}
			ctx := context.Background()
			if sw.ServiceOptions.timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, sw.ServiceOptions.timeout)
				defer cancel()
			}
			meta := &metaValue{
				header: header,
				reply:  msg.Reply,
				server: s,
			}
			ctx = withMeta(ctx, meta)

			err = s.handle(ctx, sw, method, msg.Data, msg.Reply)
			if err != nil {
				s.opt.errorHandler(err.Error())
			}
		}
		if sw.singleGoroutine {
			call()
		} else {
			go call()
		}
	}

	sub := sw.Name()
	reqSub, subErr := s.conn.QueueSubscribe(sub, defaultSubQueue, cb)
	if nil != subErr {
		return subErr
	}
	sw.subscriptions = append(sw.subscriptions, reqSub)
	if len(sw.Service.sd.PublishMethods()) > 0 {
		pubSub, pubErr := s.conn.QueueSubscribe(joinSubject(sub, pubSuffix), "", cb)
		if pubErr != nil {
			go reqSub.Unsubscribe()
			return pubErr
		}
		sw.subscriptions = append(sw.subscriptions, pubSub)
	}
	return nil
}

func (s *Server) handle(ctx context.Context, sw *serviceWrapper, method string, payload []byte, replySub string) error {
	if s.opt.recoverHandler != nil {
		defer func() {
			if e := recover(); e != nil {
				s.opt.recoverHandler(e)
			}
		}()
	}

	b, err := sw.Call(ctx, method, payload, sw.ServiceOptions.interceptor)
	if errors.Is(err, ErrReplyLater) {
		// reply later
		// 用户自己回复消息
		return nil
	}
	// publish 不需要回复
	if len(replySub) == 0 {
		return nil
	}

	respMsg := &nats.Msg{
		Subject: replySub,
		Data:    b,
		Header:  makeErrorHeader(err),
	}

	return s.conn.PublishMsg(respMsg)
}
