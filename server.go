package natsrpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/nats-io/nats.go"
)

const ()

type ServiceInfo struct {
	*Service
	subscriptions []*nats.Subscription
}

var _ ServiceRegistrar = (*Server)(nil)

// Server RPC server
type Server struct {
	wg       sync.WaitGroup          // wait group
	mu       sync.Mutex              // lock
	opt      ServerOptions           // options
	conn     *nats.Conn              // NATS Encode Conn
	services map[string]*ServiceInfo // 服务 name->Service
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
		services: map[string]*ServiceInfo{},
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
	for _, v := range unsubs {
		v.Unsubscribe()
	}
	s.mu.Unlock()
	if len(unsubs) > 0 {
		return s.conn.Flush()
	}
	return nil
}

// Remove 移除一个服务
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
func (s *Server) Register(sd ServiceDesc, val interface{}, opts ...ServiceOption) (ServiceInterface, error) {
	opt := DefaultServiceOptions
	for _, v := range opts {
		v(&opt)
	}

	// new 一个服务
	svc, err := NewService(s, sd, val, opt)
	if nil != err {
		return nil, err
	}

	name := svc.Name()
	s.mu.Lock()
	if _, ok := s.services[name]; ok {
		s.mu.Unlock()
		return nil, ErrDuplicateService
	}

	sw := &ServiceInfo{
		Service: svc,
	}
	if err := s.subscribeMethod(sw); nil != err {
		s.mu.Unlock()
		return nil, err
	}

	s.services[sd.ServiceName] = sw
	s.mu.Unlock()

	// TODO flush
	if err := s.conn.Flush(); err != nil {
		return nil, err
	}

	return svc, nil
}

// subscribeMethod 订阅服务的方法
func (s *Server) subscribeMethod(sw *ServiceInfo) error {
	cb := func(msg *nats.Msg) {
		s.wg.Add(1)
		call := func() {
			defer s.wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), sw.opt.timeout)
			defer cancel()
			err := s.handle(ctx, sw, msg)
			if err != nil {
				if errors.Is(err, ErrReplyLater) {
					// reply later
					// 用户自己回复消息
					return
				}
				s.opt.errorHandler(err.Error())
			}
		}
		if sw.opt.multiGoroutine {
			// TODO 自定义携程池
			go call()
		} else {
			call()
		}
	}

	sub := sw.Name()
	queue := defaultQueue
	reqSub, subErr := s.conn.QueueSubscribe(sub, queue, cb)
	if nil != subErr {
		return subErr
	}
	sw.subscriptions = append(sw.subscriptions, reqSub)
	if sw.Service.sd.hasPublishMethod() {
		pubSub, pubErr := s.conn.Subscribe(joinSubject(sub, pubSuffix), cb)
		if pubErr != nil {
			go reqSub.Unsubscribe()
			return pubErr
		}
		sw.subscriptions = append(sw.subscriptions, pubSub)
	}
	return nil
}

func (s *Server) handle(ctx context.Context, sw *ServiceInfo, msg *nats.Msg) error {
	if s.opt.recoverHandler != nil {
		defer func() {
			if e := recover(); e != nil {
				s.opt.recoverHandler(e)
			}
		}()
	}

	var replySub = msg.Reply
	method, header, err := decodeHeader(msg.Header)
	if err != nil {
		return err
	}
	payload := msg.Data

	dec := func(v any) error {
		if len(payload) > 0 {
			return s.opt.encoder.Decode(payload, v)
		}
		return nil
	}
	tr := &Transport{
		operation:    method,
		reqHeader:    Header(header),
		replySubject: replySub,
		request:      nil,
	}

	ctx = transport.NewServerContext(ctx, tr)
	if replySub != "" {
		tr.replyFunc = func(resp any, err error) error {
			return s.reply(ctx, replySub, resp, err)
		}
	}
	var h Invoker = func(ctx context.Context, _ any) (any, error) {
		return sw.call(ctx, method, dec)
	}
	if len(s.opt.middleware) > 0 || len(sw.opt.middleware) > 0 {
		mw := append(s.opt.middleware, sw.opt.middleware...)
		h = middleware.Chain(mw...)(h)
	}
	resp, err := h(ctx, nil)
	if err != nil {
		if errors.Is(err, ErrReplyLater) {
			// reply later
			// 用户自己回复消息
			return ErrReplyLater
		}
		//return err
		// err 要返回给客户端
	}

	// publish 不需要回复
	if len(replySub) == 0 {
		return nil
	}
	return s.reply(ctx, replySub, resp, err)
}

func (s *Server) reply(ctx context.Context, sub string, resp any, err error) error {
	var b []byte
	if resp != nil {
		b, err = s.opt.encoder.Encode(resp)
		if err != nil {
			return err
		}
	}
	respMsg := &nats.Msg{
		Subject: sub,
		Data:    b,
		Header:  makeErrorHeader(err),
	}

	return s.conn.PublishMsg(respMsg)
}
