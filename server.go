package natsrpc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Server server
type Server struct {
	enc      *nats.EncodedConn   // NATS Encode Conn
	mu       sync.Mutex          // lock
	services map[string]*service // 服务 name->service
}

// NewServer 构造器
func NewServer(enc *nats.EncodedConn) (*Server, error) {
	if !enc.Conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}
	d := &Server{
		enc:      enc,
		services: make(map[string]*service),
	}
	return d, nil
}

// NewServerWithConfig NewServerWithConfig
func NewServerWithConfig(cfg Config, option ...nats.Option) (*Server, error) {
	conn, err := NewNATSConn(cfg, option...)
	if nil != err {
		return nil, err
	}
	return NewServer(conn)
}

// ClearSubscription 取消所有订阅
func (s *Server) ClearSubscription() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.services {
		for _, vv := range v.subscribers {
			vv.Unsubscribe()
		}
		v.subscribers = nil
	}
}

// Close 关闭
func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.enc.FlushTimeout(3 * time.Second)
	s.enc.Close()
}

// Unregister 反注册
func (s *Server) unregister(service *service) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[service.name]; ok {
		for _, v := range service.subscribers {
			v.Unsubscribe()
		}
		service.subscribers = nil
		delete(s.services, service.name)
	}
	return false
}

// Register 注册服务
func (s *Server) Register(name string, serv interface{}, opts ...Option) (Service, error) {

	// new 一个服务
	service, err := newService(name, serv, opts...)
	if nil != err {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否重复
	if _, ok := s.services[service.name]; ok {
		return nil, fmt.Errorf("service [%s] duplicate", service.name)
	}
	service.server = s

	// TODO 如果报错了是否要unsub？
	if err := s.subscribeMethod(service); nil != err {
		return nil, err
	}
	s.services[service.name] = service
	return service, nil
}

// subscribeMethod 订阅服务的方法
func (s *Server) subscribeMethod(service *service) error {
	// 订阅
	for subject, v := range service.methods {
		m := v

		cb := func(msg *nats.Msg) {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), service.opt.timeout)
				defer cancel()

				reply, err := service.call(ctx, m, msg.Data)
				if nil != err {
					log.Printf("m.execute error[%v]", err)
					return
				}

				if "" == msg.Reply {
					return
				}

				if !s.enc.Conn.IsClosed() {
					s.enc.Publish(msg.Reply, reply)
				}
			}()
		}

		sub, subErr := s.enc.QueueSubscribe(subject, service.opt.group, cb)
		if nil != subErr {
			return subErr
		}
		service.subscribers = append(service.subscribers, sub)
	}
	return nil
}
