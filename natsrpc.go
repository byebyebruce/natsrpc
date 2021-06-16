package natsrpc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/nats-io/nats.go"
)

// NatsRPC server
type NatsRPC struct {
	enc      *nats.EncodedConn   // NATS Encode Conn
	mu       sync.Mutex          // lock
	services map[string]*service // 服务 name->service
}

// NewNatsRPC 构造器
func NewNatsRPC(enc *nats.EncodedConn) (*NatsRPC, error) {
	if !enc.Conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}
	d := &NatsRPC{
		enc:      enc,
		services: make(map[string]*service),
	}
	return d, nil
}

// NewNatsRPCWithConfig NewServerWithConfig
func NewNatsRPCWithConfig(cfg Config, option ...nats.Option) (*NatsRPC, error) {
	conn, err := NewNATSConn(cfg, option...)
	if nil != err {
		return nil, err
	}
	return NewNatsRPC(conn)
}

// ClearSubscription 取消所有订阅
func (rpc *NatsRPC) ClearSubscription() {
	rpc.mu.Lock()
	defer rpc.mu.Unlock()
	for _, v := range rpc.services {
		for _, vv := range v.subscribers {
			vv.Unsubscribe()
		}
		v.subscribers = nil
	}
}

// Close 关闭
func (rpc *NatsRPC) Close() {
	rpc.mu.Lock()
	defer rpc.mu.Unlock()
	_ = rpc.enc.FlushTimeout(3 * time.Second)
	rpc.enc.Close()
}

// Unregister 反注册
func (rpc *NatsRPC) unregister(service *service) bool {
	rpc.mu.Lock()
	defer rpc.mu.Unlock()
	if _, ok := rpc.services[service.name]; ok {
		for _, v := range service.subscribers {
			v.Unsubscribe()
		}
		service.subscribers = nil
		delete(rpc.services, service.name)
	}
	return false
}

// Register 注册服务
func (rpc *NatsRPC) Register(serv interface{}, opts ...Option) (Service, error) {
	opt := MakeOptions(opts...)

	// new 一个服务
	service, err := newService(serv, opt)
	if nil != err {
		return nil, err
	}
	rpc.mu.Lock()
	defer rpc.mu.Unlock()

	// 检查是否重复
	if _, ok := rpc.services[service.name]; ok {
		return nil, fmt.Errorf("service [%s] duplicate", service.name)
	}
	service.rpc = rpc

	// TODO 如果报错了是否要unsub？
	if err := rpc.subscribeMethod(service); nil != err {
		return nil, err
	}
	rpc.services[service.name] = service
	return service, nil
}

// subscribeMethod 订阅服务的方法
func (rpc *NatsRPC) subscribeMethod(service *service) error {
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
					if !rpc.enc.Conn.IsClosed() {
						rpc.enc.Publish(msg.Reply, reply)
					}
				}
			}()
		}

		sub, subErr := rpc.enc.QueueSubscribe(subject, service.options.group, cb)
		if nil != subErr {
			return subErr
		}
		service.subscribers = append(service.subscribers, sub)
	}
	return nil
}

// Publish 发布
func (rpc *NatsRPC) Publish(sub string, message proto.Message, opt Options) error {
	if opt.isSingleThreadMode() { // 单线程模型不阻塞
		go rpc.enc.Publish(sub, message)
	} else {
		return rpc.enc.Publish(sub, message)
	}
	return nil
}

// Request 请求
func (rpc *NatsRPC) Request(ctx context.Context, sub string, req proto.Message, rep proto.Message, opt Options) error {
	if opt.isSingleThreadMode() { // 单线程模式不能同步请求
		panic("should call AsyncRequest in single thread mode")
	}
	return rpc.enc.RequestWithContext(ctx, sub, req, rep)
}

// AsyncRequest 异步请求
func (rpc *NatsRPC) AsyncRequest(ctx context.Context, sub string, req proto.Message, rep proto.Message, opt Options, cb func(proto.Message, error)) {
	if !opt.isSingleThreadMode() { // 非单线程模式不能异步请求
		panic("call AsyncRequest only in single thread mode")
	}
	go func() { // 不阻塞主线程
		err := rpc.enc.RequestWithContext(ctx, sub, req, rep)
		f := func() { // 回调
			cb(rep, err)
		}
		select {
		case opt.singleThreadCbChan <- f:
		case <-ctx.Done():
			log.Println("AsyncRequest", sub, err)
		}
	}()

}
