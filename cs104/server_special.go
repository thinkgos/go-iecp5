package cs104

import (
	"context"
	"crypto/tls"
	"errors"
	"math/rand"
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
	DefaultReconnectInterval = 1 * time.Minute
)

// OnConnectHandler when connected it will be call
type OnConnectHandler func(c ServerSpecial) error

// ConnectionLostHandler when Connection lost it will be call
type ConnectionLostHandler func(c ServerSpecial)

// ServerSpecial server special interface
type ServerSpecial interface {
	asdu.Connect

	IsConnected() bool
	IsClosed() bool
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

	server            *url.URL      // 连接的服务器端
	autoReconnect     bool          // 是否启动重连
	reconnectInterval time.Duration // 重连间隔时间
	TLSConfig         *tls.Config
	onConnect         OnConnectHandler
	onConnectionLost  ConnectionLostHandler

	cancel context.CancelFunc
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

			in:   make(chan []byte, 1024),
			out:  make(chan []byte, 1024),
			recv: make(chan []byte, 1024),
			send: make(chan []byte, 1024), // may not block!

			Clog: clog.NewWithPrefix("cs104 serverSpec => "),
		},
		autoReconnect:     true,
		reconnectInterval: DefaultReconnectInterval,
		onConnect:         func(ServerSpecial) error { return nil },
		onConnectionLost:  func(ServerSpecial) {},
	}, nil
}

// SetReconnectInterval set tcp  reconnect the host interval when connect failed after try
func (this *serverSpec) SetReconnectInterval(t time.Duration) {
	this.reconnectInterval = t
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
	this.server = remoteURL
	return nil
}

// Start start the server,and return quickly,if it nil,the server will disconnected background,other failed
func (this *serverSpec) Start() error {
	if this.server == nil {
		return errors.New("empty remote server")
	}

	go this.running()
	return nil
}

// 增加重连间隔
func (this *serverSpec) running() {
	var ctx context.Context

	this.rwMux.Lock()
	if !atomic.CompareAndSwapUint32(&this.status, initial, disconnected) {
		this.rwMux.Unlock()
		return
	}
	ctx, this.cancel = context.WithCancel(context.Background())
	this.rwMux.Unlock()
	defer this.setConnectStatus(initial)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		this.Debug("disconnected server %+v", this.server)
		conn, err := openConnection(this.server, this.TLSConfig, this.Config.ConnectTimeout0)
		if err != nil {
			this.Error("connect failed, %v", err)
			if !this.autoReconnect {
				return
			}
			time.Sleep(this.reconnectInterval)
			continue
		}
		this.Debug("connect success")
		this.conn = conn
		if err = this.onConnect(this); err != nil {
			time.Sleep(this.reconnectInterval)
			continue
		}
		this.run(ctx)
		this.onConnectionLost(this)
		select {
		case <-ctx.Done():
			return
		default:
			// 随机500ms-1s的重试，避免快速重试造成服务器许多无效连接
			time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(5)))
		}
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

func (this *serverSpec) IsClosed() bool {
	return this.connectStatus() == initial
}

func (this *serverSpec) Close() error {
	this.rwMux.Lock()
	if this.cancel != nil {
		this.cancel()
	}
	this.rwMux.Unlock()
	return nil
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
