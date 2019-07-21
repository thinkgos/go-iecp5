package cs104

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

// defined default value
const (
	DefaultConnectTimeout    = 15 * time.Second
	DefaultReconnectInterval = 1 * time.Minute
)

// OnConnectHandler when connected it will be call
type OnConnectHandler func(c ServerSpecial)

// ConnectionLostHandler when Connection lost it will be call
type ConnectionLostHandler func(c ServerSpecial)

// ServerSpecial server special interface
type ServerSpecial interface {
	asdu.Connect

	IsConnected() bool
	Start() error
	Close() error

	SetTLSConfig(t *tls.Config)
	AddRemoteServer(server string) error
	SetOnConnectHandler(f OnConnectHandler)
	SetConnectionLostHandler(f ConnectionLostHandler)

	LogMode(enable bool)
	SetLogProvider(p clog.LogProvider)
}

type serverSpec struct {
	SrvSession

	Server            *url.URL      // 连接的服务器端
	connectTimeout    time.Duration // 连接超时时间
	autoReconnect     bool          // 是否启动重连
	ReconnectInterval time.Duration // 重连间隔时间
	TLSConfig         *tls.Config
	onConnect         OnConnectHandler
	onConnectionLost  ConnectionLostHandler
}

// NewServerSpecial new special server
func NewServerSpecial(conf *Config, params *asdu.Params, handler ServerHandlerInterface) (ServerSpecial, error) {
	if handler == nil {
		return nil, errors.New("invalid handler")
	}
	if err := conf.Valid(); err != nil {
		return nil, err
	}
	if err := params.Valid(); err != nil {
		return nil, err
	}

	return &serverSpec{
		SrvSession: SrvSession{
			Config:  conf,
			params:  params,
			handler: handler,

			in:   make(chan []byte, conf.RecvUnAckLimitW),
			out:  make(chan []byte, conf.SendUnAckLimitK),
			recv: make(chan []byte, conf.RecvUnAckLimitW),
			send: make(chan []byte, conf.SendUnAckLimitK), // may not block!

			Clog: clog.NewWithPrefix("cs104 serverSpec => "),
		},
		connectTimeout:    DefaultConnectTimeout,
		autoReconnect:     true,
		ReconnectInterval: DefaultReconnectInterval,
		onConnect:         func(ServerSpecial) {},
		onConnectionLost:  func(ServerSpecial) {},
	}, nil
}

// SetConnectTimeout set tcp connect the host timeout
func (this *serverSpec) SetConnectTimeout(t time.Duration) {
	this.connectTimeout = t
}

// SetReconnectInterval set tcp  reconnect the host interval when connect failed after try
func (this *serverSpec) SetReconnectInterval(t time.Duration) {
	this.ReconnectInterval = t
}

func (this *serverSpec) EnableAutoReconnect(b bool) {
	this.autoReconnect = b
}

// SetTLSConfig set tls config
func (this *serverSpec) SetTLSConfig(t *tls.Config) {
	this.TLSConfig = t
}

// AddRemoteServer adds a broker URI to the list of brokers to be used.
// The format should be scheme://host:port
// Default values for hostname is "127.0.0.1", for schema is "tcp://".
// An example broker URI would look like: tcp://foobar.com:1204
func (this *serverSpec) AddRemoteServer(server string) error {
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
	this.Server = remoteURL
	return nil
}

func (this *serverSpec) Start() error {
	if this.Server == nil {
		return errors.New("empty remote server")
	}

	go this.running()
	return nil
}

// 增加30秒 重连间隔
func (this *serverSpec) running() {
	this.rwMux.Lock()
	if !atomic.CompareAndSwapUint32(&this.status, disconnected, connecting) {
		this.rwMux.Unlock()
		return
	}
	this.rwMux.Unlock()

	for {
		this.Debug("connecting server %+v", this.Server)
		conn, err := openConnection(this.Server, this.TLSConfig, this.connectTimeout)
		if err != nil {
			this.Error("connect failed, %v", err)
			if !this.autoReconnect {
				this.setConnectStatus(disconnected)
				return
			}
			time.Sleep(this.ReconnectInterval)
			continue
		}
		this.Debug("connect success")
		this.conn = conn
		this.onConnect(this)
		this.run(context.Background())
		this.onConnectionLost(this)

	}
}

// SetOnConnectHandler set on connect handler
func (this *serverSpec) SetOnConnectHandler(f OnConnectHandler) {
	if f != nil {
		this.onConnect = f
	}
}

// SetConnectionLostHandler set connection lost handler
func (this *serverSpec) SetConnectionLostHandler(f ConnectionLostHandler) {
	if f != nil {
		this.onConnectionLost = f
	}
}

func openConnection(uri *url.URL, tlsc *tls.Config, timeout time.Duration) (net.Conn, error) {
	switch uri.Scheme {
	case "tcp":
		return net.DialTimeout("tcp", uri.Host, timeout)
	case "ssl":
		fallthrough
	case "tls":
		fallthrough
	case "tcps":
		return tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", uri.Host, tlsc)
	}
	return nil, errors.New("Unknown protocol")
}
