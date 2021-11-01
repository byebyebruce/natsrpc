package natsrpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Server server
type Server struct {
	*sync.WaitGroup
	opt      serverOptions                     // options
	enc      *nats.EncodedConn                 // NATS Encode Conn
	mu       sync.Mutex                        // lock
	services map[*service][]*nats.Subscription // 服务 name->service
}

var _ IServer = (*Server)(nil)

// NewServer 构造器
func NewServer(enc *nats.EncodedConn, option ...ServerOption) (*Server, error) {
	if !enc.Conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}

	options := defaultServerOptions
	for _, v := range option {
		v(&options)
	}

	d := &Server{
		WaitGroup: &sync.WaitGroup{},
		opt:       options,
		enc:       enc,
		services:  make(map[*service][]*nats.Subscription),
	}
	return d, nil
}

// Close 关闭
func (s *Server) Close(duration time.Duration) (err error) {
	s.ClearAllSubscription()

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	over := make(chan struct{})
	go func() {
		s.Wait()
		close(over)
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-over:
	}
	if err1 := s.enc.FlushTimeout(duration); err == nil && err1 != nil {
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
		for _, v := range sub {
			v.Unsubscribe()
		}
		delete(s.services, service)
	}
	return ok
}

// Register 注册服务
func (s *Server) Register(name string, svc interface{}, opts ...ServiceOption) (IService, error) {
	// new 一个服务
	service, err := newService(name, svc, opts...)
	if nil != err {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否重复
	if _, ok := s.services[service]; ok {
		return nil, fmt.Errorf("service [%s] duplicate", service.name)
	}
	for k := range s.services {
		if k.Name() == service.Name() {
			return nil, fmt.Errorf("service [%s] duplicate", service.name)
		}
	}
	service.server = s

	if err := s.subscribeMethod(service); nil != err {
		return nil, err
	}
	s.services[service] = make([]*nats.Subscription, 0, len(service.methods))
	return service, nil
}

// subscribeMethod 订阅服务的方法
func (s *Server) subscribeMethod(service *service) error {
	// 订阅
	for subject, v := range service.methods {
		m := v
		cb := func(msg *nats.Msg) {
			s.Add(1)
			go func() {

				defer s.Done()
				err := s.handle(context.Background(), service, m, msg.Data, msg.Reply)
				if err != nil {
					s.opt.logger.Println(err.Error())
				}
			}()
		}

		natsSub, subErr := s.enc.QueueSubscribe(subject, service.opt.group, cb)
		if nil != subErr {
			return subErr
		}
		s.services[service] = append(s.services[service], natsSub)
	}
	return nil
}

func (s *Server) handle(ctx context.Context, service *service, m *method, b []byte, r string) error {
	if s.opt.recoverHandler != nil {
		defer func() {
			if e := recover(); e != nil {
				s.opt.recoverHandler(e)
			}
		}()
	}

	reply, _ := service.handle(ctx, m, b)

	if "" == r {
		return nil
	}
	if s.enc.Conn.IsClosed() {
		return fmt.Errorf("conn colsed")
	}
	respMsg := &nats.Msg{
		Subject: r,
		Data:    reply,
	}

	// header
	/*
		if nil != err {
			respMsg.Header = nats.Header{}
			respMsg.Header.Add(headerError, err.Error())
		}
	*/
	return s.enc.Conn.PublishMsg(respMsg)
}

func (s *Server) Unmarshal(b []byte, i interface{}) error {
	return s.enc.Enc.Decode("", b, i)
}

func (s *Server) Marshal(i interface{}) ([]byte, error) {
	return s.enc.Enc.Encode("", i)
}
