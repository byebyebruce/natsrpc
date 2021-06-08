package natsrpc

import (
	"fmt"
	"go/ast"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"

	"github.com/golang/protobuf/proto"
)

const (
	callbackParameter = 3 // 回掉函数参数个数
)

var (
	valEmptyString = reflect.ValueOf("")
)

// Config 配置
type Config struct {
	//Cluster        string `xml:"cluster"`         // 集群名字,为了同一个nats-server各个集群下互相不影响
	Server         string `xml:"server" yaml:"server"`                   // nats://127.0.0.1:4222,nats://127.0.0.1:4223
	User           string `xml:"user" yaml:"user"`                       // 用户名
	Pwd            string `xml:"pwd" yaml:"pwd"`                         // 密码
	RequestTimeout int32  `xml:"request_timeout" yaml:"request_timeout"` // 请求超时（秒）
	ReconnectWait  int64  `xml:"reconnect_wait" yaml:"reconnect_wait"`   // 重连间隔
	MaxReconnects  int32  `xml:"max_reconnects" yaml:"max_reconnects"`   // 重连次数
}

// joinSubject 把subject用.分割组合
func joinSubject(typeName string, subjectPostfix ...interface{}) string {
	sub := typeName
	if nil != subjectPostfix {
		for _, v := range subjectPostfix {
			sub += fmt.Sprintf(".%v", v)
		}
	}
	return sub
}

// Handler 消息处理函数
// 格式：func(pb *proto.MyUser, reply string, err string)
type Handler interface{}

// msg 异步回掉消息，用reflect.Value主要为了不在主线程掉reflect.ValueOf
type msg struct {
	handler reflect.Value // 回掉函数的value
	arg     reflect.Value // 参数的value
	reply   reflect.Value // 回复字符串的value
	err     reflect.Value // 错误字符串的value
}

func checkHandler(cb interface{}) (reflect.Type, error) {
	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		return nil, errorFuncType
	}

	numArgs := cbType.NumIn()
	if callbackParameter != numArgs {
		return nil, errorFuncType
	}

	argType := cbType.In(0)
	if argType.Kind() != reflect.Ptr {
		return nil, errorFuncType
	}
	if cbType.In(1).Kind() != reflect.String {
		return nil, errorFuncType
	}
	if cbType.In(2).Kind() != reflect.String {
		return nil, errorFuncType
	}
	oPtr := reflect.New(argType.Elem())
	_, ok := oPtr.Interface().(proto.Message)
	if !ok {
		return nil, errorFuncType
	}
	return argType, nil
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

func NewNATSClient(cfg *Config, name string) (*nats.EncodedConn, error) {
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
		//l4g.Warn("[nats(%s)] Reconnected [%s]", name, nc.ConnectedUrl())
	}))
	opts = append(opts, nats.DiscoveredServersHandler(func(nc *nats.Conn) {
		//l4g.Info("[nats(%s)] DiscoveredServersHandler", name, nc.DiscoveredServers())
	}))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		//l4g.Warn("[nats(%s)] Disconnect", name)
	}))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		if nil != err {
			//l4g.Warn("[nats(%s)] DisconnectErrHandler,error=[%v]", name, err)
		}
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		//l4g.Info("[nats(%s)] ClosedHandler", name)
	}))
	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, subs *nats.Subscription, err error) {
		//l4g.Warn("[nats(%s)] ErrorHandler subs=[%s] error=[%s]", name, subs.Subject, err.Error())
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
	return enc, nil
}
