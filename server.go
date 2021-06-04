package xnats

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

// Server NATS客户端
type Server struct {
	name           string
	enc            *nats.EncodedConn  // NATS的Conn
	msgChan        chan *msg          // 消息通道
	argArray       []reflect.Value    // value参数数组，防止多次分配数组
	handleMap      map[string]Handler // 消息处理函数map
	requestTimeout time.Duration
	mu             sync.Mutex
	subscribers    []*nats.Subscription
}

// NewServer 构造器
func NewServer(enc *nats.EncodedConn, name string, maxMsg int) (*Server, error) {

	d := &Server{
		name:           name,
		enc:            enc,
		handleMap:      make(map[string]Handler),
		msgChan:        make(chan *msg, maxMsg),
		argArray:       make([]reflect.Value, callbackParameter),
		requestTimeout: time.Second * time.Duration(1),
		subscribers:    make([]*nats.Subscription, 0),
	}
	return d, nil
}

// ClearSubscription 取消所有订阅
func (s *Server) ClearSubscription() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.subscribers {
		v.Unsubscribe()
	}
	s.subscribers = make([]*nats.Subscription, 0)
	log.Printf("[nats(%s)] ClearSubscription", s.name)
}

// RawConn 返回conn
func (s *Server) RawConn() *nats.EncodedConn {
	return s.enc
}

// Reply 回复消息
func (s *Server) Reply(reply string, m proto.Message) {
	if err := s.enc.Publish(reply, m); nil != err {
		log.Printf("[nats(%s)] Reply reply=[%s] error=[%s]", s.name, reply, err.Error())
	} else {
		log.Printf("[nats(%s)] Reply reply=[%s]", s.name, reply)
	}
}

func (s *Server) parseHandler(cb interface{}, subjectPostfix ...interface{}) (func(*nats.Msg), string) {
	// 检查回掉函数的格式
	argType, err := checkHandler(cb)
	if nil != err {
		panic(err)
	}
	// 算subject
	sub := joinSubject(argType.String(), subjectPostfix...)

	cbValue := reflect.ValueOf(cb)

	// 回掉函数
	h := func(m *nats.Msg) {
		argVal := reflect.New(argType.Elem())
		pb := argVal.Interface().(proto.Message)
		if err := proto.Unmarshal(m.Data, pb); nil != err {
			log.Printf("[nats(%s)] cb proto.Unmarshal error=[%s]", s.name, err.Error())
		} else {
			s.msgChan <- &msg{
				handler: cbValue,
				arg:     argVal,
				reply:   reflect.ValueOf(m.Reply),
				err:     valEmptyString,
			}
		}
		log.Printf("[nats(%s)] callback sub=[%s] reply=[%s]", s.name, m.Subject, m.Reply)
	}
	return h, sub
}

// RegisterHandler 注册异步回掉函数，消息会发送消息通道msgChan
// grouped:是否分组(分组的只有一个能收到)
// cb格式:func(pb *proto.MyUser ,reply string, err string)
// subjectPostfix是subject的后缀 例如：m是*pb.MyUser类型的对象，subjectPostfix是1000410001，那subject是 "*pb.MyUser.1000410001"
func (s *Server) RegisterHandler(cb interface{}, grouped bool, subjectPostfix ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	h, sub := s.parseHandler(cb, subjectPostfix...)
	if _, ok := s.handleMap[sub]; ok {
		panic("handler exists")
	}
	// 保存下，防止重复
	s.handleMap[sub] = cb

	var err error
	var subscription *nats.Subscription
	if grouped {
		subscription, err = s.enc.QueueSubscribe(sub, "group", h)
	} else {
		subscription, err = s.enc.Subscribe(sub, h)
	}
	if nil != err {
		panic(err)
	} else {
		s.subscribers = append(s.subscribers, subscription)
	}
	log.Printf("[nats(%s)] RegisterHandler=[%s] grouped[%v]", s.name, sub, grouped)
}

func (s *Server) parseSyncHandler(cb interface{}, subjectPostfix ...interface{}) (func(*nats.Msg), string) {
	// 检查回掉函数的格式
	argType, err := checkHandler(cb)
	if nil != err {
		panic(err)
	}

	cbValue := reflect.ValueOf(cb)

	sub := joinSubject(argType.String(), subjectPostfix...)
	h := func(m *nats.Msg) {
		argVal := reflect.New(argType.Elem())
		pb := argVal.Interface().(proto.Message)
		if err := proto.Unmarshal(m.Data, pb); nil != err {
			log.Printf("[nats(%s)] cb proto.Unmarshal error=[%s]", s.name, err.Error())
		} else {
			cbValue.Call([]reflect.Value{argVal, reflect.ValueOf(m.Reply), valEmptyString})
		}
		log.Printf("[nats(%s)] sync callback sub=[%s] reply=[%s]", s.name, m.Subject, m.Reply)
	}
	return h, sub
}

// RegisterSyncHandler 注册同步回掉函数
// grouped:是否分组(分组的只有一个能收到)
// cb格式:func(pb *proto.MyUser ,reply string, err string)
// subjectPostfix是subject的后缀 例如：m是*pb.MyUser类型的对象，subjectPostfix是1000410001，那subject是 "*pb.MyUser.1000410001"
func (s *Server) RegisterSyncHandler(cb interface{}, grouped bool, subjectPostfix ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	h, sub := s.parseSyncHandler(cb, subjectPostfix...)

	var err error
	var subscription *nats.Subscription
	if grouped {
		subscription, err = s.enc.QueueSubscribe(sub, "group", h)
	} else {
		subscription, err = s.enc.Subscribe(sub, h)
	}
	if nil != err {
		panic(err)
	} else {
		s.subscribers = append(s.subscribers, subscription)
	}

	log.Printf("[nats(%s)] RegisterSyncHandler=[%s] grouped[%v]", s.name, sub, grouped)
}

// MsgChan 消息通道
func (s *Server) MsgChan() <-chan *msg {
	return s.msgChan
}

// Process 处理消息
func (s *Server) Process(m *msg) {
	defer func() {
		if err := recover(); nil != err {
			trace := string(debug.Stack())
			fmt.Println("[nats] Process panic:", err, trace)
			log.Printf("[nats(%s)] Process panic=%v stack=%s", s.name, err, trace)
		}
	}()
	s.argArray[0] = m.arg
	s.argArray[1] = m.reply
	s.argArray[2] = m.err
	m.handler.Call(s.argArray)
}

// Close 关闭
// 是否需要处理完通道里的消息
func (s *Server) Close(process bool) {
	if process {
		for exit := false; exit; {
			select {
			case m := <-s.msgChan:
				s.Process(m)
			default:
				exit = true
			}
		}
	}
	s.enc.FlushTimeout(time.Duration(3 * time.Second))
	s.enc.Close()
}

func (s *Server) Register(handler interface{}, options ...ServiceOption) (*Service, error) {
	service := &Service{
		s: s,
	}
	service.rcvr = reflect.ValueOf(handler)
	s.name = reflect.Indirect(service.rcvr).Type().Name()
	service.typ = reflect.TypeOf(service.rcvr)
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}

	option := newDefaultOption()
	for _, v := range options {
		v(&option)
	}

	subPrefix := fmt.Sprintf("%s.%s.%s", option.namespace, s.name, option.id)

	//s.method = make(map[string]*methodType)
	for i := 0; i < service.typ.NumMethod(); i++ {
		method := service.typ.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		sub := subPrefix + "." + method.Name

		log.Printf("rpc server: register %v\n", sub)
	}

	return service, nil
}
