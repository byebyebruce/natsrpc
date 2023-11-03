package natsrpc

import (
	"context"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

// Server RPC server
type Server struct {
	wg       sync.WaitGroup                  // wait group
	mu       sync.Mutex                      // lock
	opt      serverOptions                   // options
	conn     *nats.Conn                      // NATS Encode Conn
	services map[*service]*nats.Subscription // 服务 name->service
}

var _ IServer = (*Server)(nil)

// NewServer 构造器
func NewServer(conn *nats.Conn, option ...ServerOption) (*Server, error) {
	if !conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}

	options := defaultServerOptions
	for _, v := range option {
		v(&options)
	}

	d := &Server{
		opt:      options,
		conn:     conn,
		services: make(map[*service]*nats.Subscription),
	}
	return d, nil
}

// Close 关闭
func (s *Server) Close(ctx context.Context) (err error) {
	s.ClearAllSubscription()

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
	if err1 := s.conn.Flush(); err == nil && err1 != nil {
		err = err1
	}
	return
}

// ClearAllSubscription 取消所有订阅
func (s *Server) ClearAllSubscription() {
	s.mu.Lock()
	ss := make([]*service, 0, len(s.services))
	for s := range s.services {
		ss = append(ss, s)
	}
	s.mu.Unlock()

	for _, v := range ss {
		s.remove(v)
	}
}

// Unregister 反注册
func (s *Server) remove(service *service) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	sub, ok := s.services[service]
	if ok {
		sub.Unsubscribe()
		delete(s.services, service)
	}
	return ok
}

// Register 注册服务
func (s *Server) Register(name string, handler interface{}, opts ...ServiceOption) (IService, error) {
	// new 一个服务
	svc, err := newService(name, handler, opts...)
	if nil != err {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for s := range s.services {
		if s.sub == svc.sub && s.val == handler {
			return nil, ErrDuplicateService
		}
	}

	svc.server = s

	if err := s.subscribeMethod(svc); nil != err {
		return nil, err
	}
	return svc, nil
}

// subscribeMethod 订阅服务的方法
func (s *Server) subscribeMethod(service *service) error {
	cb := func(msg *nats.Msg) {
		s.wg.Add(1)

		call := func() {
			defer s.wg.Done()

			mName, header, err := decodeHeader(msg.Header)
			if err != nil {
				s.opt.errorHandler(err.Error())
				return
			}
			ctx := context.Background()
			if header != nil {
				ctx = setHeader(ctx, header)
			}

			err = s.handle(ctx, service, mName, msg.Data, msg.Reply)
			if err != nil {
				s.opt.errorHandler(err.Error())
			}
		}
		if service.opt.concurrent {
			go call()
		} else {
			call()
		}
	}

	natsSub, subErr := s.conn.QueueSubscribe(service.sub, defaultSubQueue, cb)
	if nil != subErr {
		return subErr
	}
	s.services[service] = natsSub
	// TODO flush
	s.conn.Flush()

	return nil
}

func (s *Server) handle(ctx context.Context, svc *service, method string, payload []byte, replySub string) error {
	if s.opt.recoverHandler != nil {
		defer func() {
			if e := recover(); e != nil {
				s.opt.recoverHandler(e)
			}
		}()
	}

	b, err := svc.call(ctx, method, payload)
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
