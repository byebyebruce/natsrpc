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
	conn     *nats.EncodedConn   // NATS Conn
	mu       sync.Mutex          // lock
	services map[string]*service // 服务 name->service

}

// NewServer 构造器
func NewServer(enc *nats.EncodedConn) (*Server, error) {
	if !enc.Conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}
	d := &Server{
		conn:     enc,
		services: make(map[string]*service),
	}
	return d, nil
}

// NewServerWithConfig NewServerWithConfig
func NewServerWithConfig(cfg Config, option ...nats.Option) (*Server, error) {
	client, err := NewNATSConn(cfg, option...)
	if nil != err {
		return nil, err
	}
	return NewServer(client)
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
	s.ClearSubscription()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.conn.FlushTimeout(time.Duration(3 * time.Second))
	s.conn.Close()
}

// Unregister 反注册
func (s *Server) Unregister(service *service) bool {
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
func (s *Server) Register(serv interface{}, options ...Option) (*service, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	option := defaultOption()
	for _, v := range options {
		v(&option)
	}

	// new 一个服务
	service, err := newService(serv, option)
	if nil != err {
		return nil, err
	}

	// 检查是否重复
	service.server = s
	if _, ok := s.services[service.name]; ok {
		return nil, fmt.Errorf("service [%s] duplicate", service.name)
	}
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
				ctx, cancel := context.WithTimeout(context.Background(), service.options.timeout)
				defer cancel()

				var (
					reply interface{}
					err   error
				)

				if nil != service.options.singleThreadCbChan { // 单线程处理
					over := make(chan struct{})
					fn := func() {
						defer close(over)
						reply, err = m.handle(ctx, msg.Data)
					}
					select {
					case <-ctx.Done():
						err = ctx.Err()
					case service.options.singleThreadCbChan <- fn:
						select {
						case <-ctx.Done():
							err = ctx.Err()
						case <-over:
						}
					}
				} else { // 多线程处理
					reply, err = m.handle(ctx, msg.Data)
				}
				// handle
				if nil != err {
					log.Printf("m.handle error[%v]", err)
					return
				}

				// reply
				if "" != msg.Reply && nil != reply {
					if !s.conn.Conn.IsClosed() {
						s.conn.Publish(msg.Reply, reply)
					}
				}
			}()
		}

		sub, subErr := s.conn.QueueSubscribe(subject, service.options.group, cb)
		if nil != subErr {
			return subErr
		}
		service.subscribers = append(service.subscribers, sub)
	}
	return nil
}
