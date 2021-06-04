package xnats

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"

	"github.com/nats-io/nats.go"

	"github.com/golang/protobuf/proto"
)

const (
	callbackParameter = 3 // 回掉函数参数个数
)

var (
	errorFuncType  = errors.New("handler must be a function like [func(pb *proto.MyUser, reply string, err string)]")
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

type method struct {
	cb     func(*nats.Msg)
	replay bool
	name   string
}

// cb(*req,*resp)
func parseMethod(m reflect.Method) (*method, error) {
	const paraNum = 3

	mType := m.Type
	numArgs := mType.NumIn()
	if paraNum != numArgs {
		return nil, errorFuncType
	}

	argType := mType.In(0)
	if argType.Kind() != reflect.Ptr {
		return nil, errorFuncType
	}

	repType := mType.In(1)
	if repType.Kind() != reflect.Ptr {
		return nil, errorFuncType
	}

	oPtr := reflect.New(argType.Elem())
	_, ok := oPtr.Interface().(proto.Message)
	if !ok {
		return nil, errorFuncType
	}

	oPtr = reflect.New(repType.Elem())
	_, ok = oPtr.Interface().(proto.Message)
	if !ok {
		return nil, errorFuncType
	}

	cbValue := reflect.ValueOf(m)

	h := func(m *nats.Msg) {
		argVal := reflect.New(argType.Elem())
		pb := argVal.Interface().(proto.Message)
		if err := proto.Unmarshal(m.Data, pb); nil != err {
			//l4g.Error("[nats(%s)] cb proto.Unmarshal error=[%s]", s.name, err.Error())
		} else {
			replyVal := reflect.New(repType.Elem())
			//repPB := replyVal.Interface().(proto.Message)
			cbValue.Call([]reflect.Value{argVal, replyVal})
		}
		//l4g.Debug("[nats(%s)] sync callback sub=[%s] reply=[%s]", s.name, m.Subject, m.Reply)
	}
	return &method{cb: h, replay: true, name: m.Type.String()}, nil
}
