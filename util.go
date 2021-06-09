package natsrpc

import (
	"go/ast"
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

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

func NewNATSConn(cfg Config, name string) (*nats.EncodedConn, error) {
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

func typeName(p reflect.Type) string {
	return strings.Trim(p.String(), "*")
}

func combineSubject(prefix string, s ...string) string {
	ret := prefix
	for _, v := range s {
		if "" == v {
			continue
		}
		ret += "." + v
	}
	return ret
}
