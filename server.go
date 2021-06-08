package natsrpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Server NATS客户端
type Server struct {
	mu         sync.Mutex
	client     *nats.EncodedConn // NATS的Conn
	serviceMap map[string]*service
}

// NewServer 构造器
func NewServer(enc *nats.EncodedConn) (*Server, error) {
	d := &Server{
		client:     enc,
		serviceMap: make(map[string]*service),
	}
	return d, nil
}

func NewServerWithConfig(cfg *Config, name string) (*Server, error) {
	client, err := NewNATSClient(cfg, name)
	if nil != err {
		return nil, err
	}
	return NewServer(client)
}

// ClearSubscription 取消所有订阅
func (s *Server) ClearSubscription() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.serviceMap {
		for _, vv := range v.subscribers {
			vv.Unsubscribe()
		}
		v.subscribers = nil
	}
}

// Close 关闭
// 是否需要处理完通道里的消息
func (s *Server) Close() {
	s.ClearSubscription()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.client.FlushTimeout(time.Duration(3 * time.Second))
	s.client.Close()
}

func (s *Server) Unregister(service *service) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serviceMap[service.name]; ok {
		delete(s.serviceMap, service.name)
	}
	return false
}

func (s *Server) Register(serv interface{}, options ...ServiceOption) (*service, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	option := newDefaultOption()
	for _, v := range options {
		v(&option)
	}

	service, err := newService(s, serv, option)
	if nil != err {
		return nil, err
	}
	if _, ok := s.serviceMap[service.name]; ok {
		return nil, fmt.Errorf("service %s duplicate", service.name)
	}

	for _, v := range service.methods {
		subject := service.name + "." + v.name
		handler := func(msg *nats.Msg) {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(service.options.timeout))
				defer cancel()
				reply := v.handler(ctx, msg.Data)
				if "" != msg.Reply && nil != reply {
					if !s.client.Conn.IsClosed() {
						s.client.Publish(msg.Reply, reply)
					}
				}
			}()
		}

		sub, subErr := s.client.QueueSubscribe(subject, service.options.group, handler)
		if nil != subErr {
			return nil, subErr
		}
		service.subscribers = append(service.subscribers, sub)
	}

	s.serviceMap[service.name] = service
	return service, nil
}
