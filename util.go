package natsrpc

import (
	"go/ast"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
)

// Config 配置
type Config struct {
	Server         string `xml:"server" yaml:"server"`                   // nats://127.0.0.1:4222,nats://127.0.0.1:4223
	User           string `xml:"user" yaml:"user"`                       // 用户名
	Pwd            string `xml:"pwd" yaml:"pwd"`                         // 密码
	RequestTimeout int32  `xml:"request_timeout" yaml:"request_timeout"` // 请求超时（秒）
	ReconnectWait  int64  `xml:"reconnect_wait" yaml:"reconnect_wait"`   // 重连间隔
	MaxReconnects  int32  `xml:"max_reconnects" yaml:"max_reconnects"`   // 重连次数
}

// isExportedOrBuiltinType 是导出或内置类型
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

// NewNATSConn 构造一个nats conn
func NewNATSConn(cfg Config, option ...nats.Option) (*nats.EncodedConn, error) {
	if cfg.ReconnectWait <= 0 {
		cfg.ReconnectWait = 1
	}
	if cfg.MaxReconnects <= 0 {
		cfg.MaxReconnects = 99999999
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 3
	}

	// 设置参数
	opts := make([]nats.Option, 0)
	if len(cfg.User) > 0 {
		opts = append(opts, nats.UserInfo(cfg.User, cfg.Pwd))
	}
	opts = append(opts, nats.ReconnectWait(time.Second*time.Duration(cfg.ReconnectWait)))
	opts = append(opts, nats.MaxReconnects(int(cfg.MaxReconnects)))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("[nats] Reconnected [%s]\n", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.DiscoveredServersHandler(func(nc *nats.Conn) {
		log.Printf("[nats] DiscoveredServersHandler [%s]\n", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		if nil != err {
			log.Printf("[nats] DisconnectErrHandler [%v]\n", err)
		}
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Printf("[nats] ClosedHandler\n")
	}))
	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, subs *nats.Subscription, err error) {
		if nil != err {
			log.Printf("[nats] ErrorHandler subs[%v] error[%v]\n", subs.Subject, err)
		}
	}))

	// 后面的可以覆盖前面的设置
	opts = append(opts, option...)

	// 创建nats enc
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

// typeName 类型名字
func typeName(p reflect.Type) string {
	return strings.Trim(p.String(), "*")
}

// CombineSubject 组合字符串成subject
func CombineSubject(prefix string, s ...string) string {
	ret := prefix
	for _, v := range s {
		if "" == v {
			continue
		}
		ret += "." + v
	}
	return ret
}
