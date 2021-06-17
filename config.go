package natsrpc

// Config 配置
type Config struct {
	Server         string `xml:"server" yaml:"server" json:"server"`                            // nats://127.0.0.1:4222,nats://127.0.0.1:4223
	User           string `xml:"user" yaml:"user" json:"user"`                                  // 用户名
	Pwd            string `xml:"pwd" yaml:"pwd" json:"pwd"`                                     // 密码
	RequestTimeout int32  `xml:"request_timeout" yaml:"request_timeout" json:"request_timeout"` // 请求超时（秒）
	ReconnectWait  int64  `xml:"reconnect_wait" yaml:"reconnect_wait" json:"reconnect_wait"`    // 重连间隔
	MaxReconnects  int32  `xml:"max_reconnects" yaml:"max_reconnects" json:"max_reconnects"`    // 重连次数
}
