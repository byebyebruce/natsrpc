package natsrpc

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
)

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

// RunSimpleNatsServer run a simple nats server
func RunSimpleNatsServer(opts *server.Options) *server.Server {
	var defaultTestOptions = server.Options{
		Host:                  "127.0.0.1",
		Port:                  4222,
		NoLog:                 false,
		NoSigs:                true,
		MaxControlLine:        4096,
		DisableShortFirstPing: true,
	}

	if opts == nil {
		opts = &defaultTestOptions
	}
	// Optionally override for individual debugging of tests
	//opts.Trace = true
	s, err := server.NewServer(opts)
	if err != nil || s == nil {
		panic(fmt.Sprintf("No NATS Server object returned: %v", err))
	}

	s.ConfigureLogger()

	// Run server in Go routine.
	go s.Start()

	// Wait for accept loop(s) to be started
	if !s.ReadyForConnections(10 * time.Second) {
		panic("Unable to start NATS Server in Go Routine")
	}
	return s
}
