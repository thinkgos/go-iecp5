package cs104

import (
	"crypto/tls"
	"net/url"
	"strings"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
)

type ClientOption struct {
	config            Config
	param             asdu.Params
	server            *url.URL      // 连接的服务器端
	autoReconnect     bool          // 是否启动重连
	reconnectInterval time.Duration // 重连间隔时间
	TLSConfig         *tls.Config   // tls配置

}

func NewOption() *ClientOption {
	return &ClientOption{
		DefaultConfig(),
		*asdu.ParamsWide,
		nil,
		true,
		DefaultReconnectInterval,
		nil,
	}
}

// SetConfig set config
func (sf *ClientOption) SetConfig(cfg Config) *ClientOption {
	if err := cfg.Valid(); err != nil {
		panic(err)
	}
	sf.config = cfg

	return sf
}

// SetParams set asdu params
func (sf *ClientOption) SetParams(p *asdu.Params) *ClientOption {
	if err := p.Valid(); err != nil {
		panic(err)
	}
	sf.param = *p

	return sf
}

// SetReconnectInterval set tcp  reconnect the host interval when connect failed after try
func (sf *ClientOption) SetReconnectInterval(t time.Duration) *ClientOption {
	sf.reconnectInterval = t
	return sf
}

// SetAutoReconnect enable auto reconnect
func (sf *ClientOption) SetAutoReconnect(b bool) *ClientOption {
	sf.autoReconnect = b
	return sf
}

// SetTLSConfig set tls config
func (sf *ClientOption) SetTLSConfig(t *tls.Config) *ClientOption {
	sf.TLSConfig = t
	return sf
}

// AddRemoteServer adds a broker URI to the list of brokers to be used.
// The format should be scheme://host:port
// Default values for hostname is "127.0.0.1", for schema is "tcp://".
// An example broker URI would look like: tcp://foobar.com:1204
func (sf *ClientOption) AddRemoteServer(server string) error {
	if len(server) > 0 && server[0] == ':' {
		server = "127.0.0.1" + server
	}
	if !strings.Contains(server, "://") {
		server = "tcp://" + server
	}
	remoteURL, err := url.Parse(server)
	if err != nil {
		return err
	}
	sf.server = remoteURL
	return nil
}
