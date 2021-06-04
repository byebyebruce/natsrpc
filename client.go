package xnats

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
)

// IClient nats 客户端
type IClient interface {
	// Notify 推送(不需要回复)
	Notify(m proto.Message, subjectPostfix ...interface{})
	// Request 请求(异步)
	Request(m proto.Message, cb interface{}, subjectPostfix ...interface{})
	// RequestSync 请求(同步)
	RequestSync(req proto.Message, resp proto.Message, subjectPostfix ...interface{}) error
	// Reply 回复
	Reply(reply string, m proto.Message)
}

// Client NATS客户端
type Client struct {
	name           string
	enc            *nats.EncodedConn // NATS的Conn
	msgChan        chan *msg         // 消息通道
	argArray       []reflect.Value   // value参数数组，防止多次分配数组
	requestTimeout time.Duration
	mu             sync.Mutex
}

// NewClient 构造器
func NewClient(cfg *Config, name string, maxMsg int) (*Client, error) {

	if cfg.ReconnectWait <= 0 {
		cfg.ReconnectWait = 3
	}
	if cfg.MaxReconnects <= 0 {
		cfg.MaxReconnects = 99999999
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 3
	}

	// 设置参数
	opts := make([]nats.Option, 0)
	opts = append(opts, nats.Name(name))
	if len(cfg.User) > 0 {
		opts = append(opts, nats.UserInfo(cfg.User, cfg.Pwd))
	}
	opts = append(opts, nats.ReconnectWait(time.Second*time.Duration(cfg.ReconnectWait)))
	opts = append(opts, nats.MaxReconnects(int(cfg.MaxReconnects)))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("[nats(%s)] Reconnected [%s]", name, nc.ConnectedUrl())
	}))
	opts = append(opts, nats.DiscoveredServersHandler(func(nc *nats.Conn) {
		log.Printf("[nats(%s)] DiscoveredServersHandler %v", name, nc.DiscoveredServers())
	}))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		log.Printf("[nats(%s)] Disconnect", name)
	}))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		if nil != err {
			log.Printf("[nats(%s)] DisconnectErrHandler,error=[%v]", name, err)
		}
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Printf("[nats(%s)] ClosedHandler", name)
	}))
	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, subs *nats.Subscription, err error) {
		log.Printf("[nats(%s)] ErrorHandler subs=[%s] error=[%s]", name, subs.Subject, err.Error())
	}))

	// 创建nats client
	nc, err := nats.Connect(cfg.Server, opts...)
	if err != nil {
		return nil, err
	}
	enc, err1 := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if nil != err1 {
		return nil, err1
	}

	d := &Client{
		name:           name,
		enc:            enc,
		msgChan:        make(chan *msg, maxMsg),
		argArray:       make([]reflect.Value, callbackParameter),
		requestTimeout: time.Second * time.Duration(cfg.RequestTimeout),
	}
	return d, nil
}

// RawConn 返回conn
func (cli *Client) RawConn() *nats.EncodedConn {
	return cli.enc
}

// Notify 推送(不需要回复)
// cb格式:func(pb *proto.MyUser ,reply string, err string)
// subjectPostfix是subject的后缀 例如：m是*pb.MyUser类型的对象，subjectPostfix是1000410001，那subject是 "*pb.MyUser.1000410001"
func (cli *Client) Notify(m proto.Message, subjectPostfix ...interface{}) {
	argVal := reflect.TypeOf(m)
	sub := joinSubject(argVal.String(), subjectPostfix...)
	if err := cli.enc.Publish(sub, m); nil != err {
		log.Printf("[nats(%s)] Notify sub=[%s] error=[%s]", cli.name, sub, err.Error())
	} else {
		log.Printf("[nats(%s)] Notify sub=[%s]", cli.name, sub)
	}
}

// Request 请求
// cb格式:func(pb *proto.MyUser ,reply string, err string)
// subjectPostfix是subject的后缀 例如：m是*pb.MyUser类型的对象，subjectPostfix是1000410001，那subject是 "*pb.MyUser.1000410001"
func (cli *Client) Request(m proto.Message, cb interface{}, subjectPostfix ...interface{}) {
	argVal := reflect.TypeOf(m)
	sub := joinSubject(argVal.String(), subjectPostfix...)
	cbValue := reflect.ValueOf(cb)
	cbType := reflect.TypeOf(cb)

	// TODO 运行时linux下不检查cb，其他的平台要检查cb
	go func() {
		argType := cbType.In(0)
		oPtr := reflect.New(argType.Elem())
		err := cli.enc.Request(sub, m, oPtr.Interface(), cli.requestTimeout)
		errVal := valEmptyString
		if nil != err {
			errVal = reflect.ValueOf(err.Error())
		}
		cli.msgChan <- &msg{
			handler: cbValue,
			arg:     oPtr,
			reply:   valEmptyString,
			err:     errVal,
		}
		if nil != err {
			log.Printf("[nats(%s)] Request sub=[%s] error=[%s]", cli.name, sub, err.Error())
		} else {
			log.Printf("[nats(%s)] Request over sub=[%s]", cli.name, sub)
		}
	}()
	log.Printf("[nats(%s)] Request sub=[%s]", cli.name, sub)
}

// RequestSync 同步请求
// cb格式:func(pb *proto.MyUser ,reply string, err string)
// subjectPostfix是subject的后缀 例如：m是*pb.MyUser类型的对象，subjectPostfix是1000410001，那subject是 "*pb.MyUser.1000410001"
func (cli *Client) RequestSync(req proto.Message, resp proto.Message, subjectPostfix ...interface{}) error {
	argVal := reflect.TypeOf(req)
	sub := joinSubject(argVal.String(), subjectPostfix...)
	log.Printf("[nats(%s)] RequestSync sub=[%s]", cli.name, sub)
	return cli.enc.Request(sub, req, resp, cli.requestTimeout)
}

// Reply 回复消息
func (cli *Client) Reply(reply string, m proto.Message) {
	if err := cli.enc.Publish(reply, m); nil != err {
		log.Printf("[nats(%s)] Reply reply=[%s] error=[%s]", cli.name, reply, err.Error())
	} else {
		log.Printf("[nats(%s)] Reply reply=[%s]", cli.name, reply)
	}
}

// MsgChan 消息通道
func (cli *Client) MsgChan() <-chan *msg {
	return cli.msgChan
}

// Process 处理消息
func (cli *Client) Process(m *msg) {
	defer func() {
		if err := recover(); nil != err {
			trace := string(debug.Stack())
			fmt.Println("[nats] Process panic:", err, trace)
			log.Printf("[nats(%s)] Process panic=%v stack=%s", cli.name, err, trace)
		}
	}()
	cli.argArray[0] = m.arg
	cli.argArray[1] = m.reply
	cli.argArray[2] = m.err
	m.handler.Call(cli.argArray)
}

// Close 关闭
// 是否需要处理完通道里的消息
func (cli *Client) Close(process bool) {
	if process {
		for exit := false; exit; {
			select {
			case m := <-cli.msgChan:
				cli.Process(m)
			default:
				exit = true
			}
		}
	}
	cli.enc.FlushTimeout(time.Duration(3 * time.Second))
	cli.enc.Close()
}
