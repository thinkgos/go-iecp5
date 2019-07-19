package cs104

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

const (
	DefaultConnectTimeout       = 30 * time.Second
	DefaultReconnectInterval    = 1 * time.Minute
	DefaultMaxReconnectInterval = 10 * time.Minute
)

type OnConnectHandler func(c *ServerSpec)
type ConnectionLostHandler func(c *ServerSpec)
type ErrorHandler func(err error)

type ServerSpecial interface {
	asdu.Connect
	IsConnected() bool
	Connect() error
	Close() error

	SetTLSConfig(t *tls.Config)
	AddRemoteServer(server string) error
	SetOnConnectHandler(f OnConnectHandler)
	SetConnectionLostHandler(f ConnectionLostHandler)
}

type ServerSpec struct {
	Session

	Server               *url.URL      // 连接的服务器端
	connectTimeout       time.Duration // 连接超时时间
	autoReconnect        bool          // 是否启动重连
	MaxReconnectInterval time.Duration // 重连间隔时间
	rw                   sync.RWMutex
	status               bool
	TLSConfig            *tls.Config
	onConnect            OnConnectHandler
	onConnectionLost     ConnectionLostHandler
}

func NewServerSpecial(conf *Config, params *asdu.Params, handler ServerHandlerInterface) *ServerSpec {
	return &ServerSpec{
		Session: Session{
			Config:  conf,
			params:  params,
			handler: handler,

			in:   make(chan []byte, conf.RecvUnAckLimitW),
			out:  make(chan []byte, conf.SendUnAckLimitK),
			recv: make(chan []byte, conf.RecvUnAckLimitW),
			send: make(chan []byte, conf.SendUnAckLimitK), // may not block!

			Clog: clog.NewWithPrefix("cs104 ServerSpec => "),
		},
		connectTimeout:       DefaultConnectTimeout,
		autoReconnect:        true,
		MaxReconnectInterval: DefaultReconnectInterval,
		onConnect:            func(*ServerSpec) {},
		onConnectionLost:     func(*ServerSpec) {},
	}
}

func (this *ServerSpec) SetTLSConfig(t *tls.Config) {
	this.TLSConfig = t
}

// AddRemoteServer adds a broker URI to the list of brokers to be used.
// The format should be scheme://host:port
// Default values for hostname is "127.0.0.1", for schema is "tcp://".
// An example broker URI would look like: tcp://foobar.com:1204
func (this *ServerSpec) AddRemoteServer(server string) error {
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

func (this *ServerSpec) IsConnected() bool {
	this.rw.RLock()
	b := this.status
	this.rw.RUnlock()
	return b
}

func (this *ServerSpec) Connect() error {
	if this.Server == nil {
		return errors.New("empty remote server")
	}
	if this.handler == nil {
		return errors.New("invalid handler")
	}
	if err := this.Config.Valid(); err != nil {
		return err
	}
	if err := this.params.Valid(); err != nil {
		return err
	}
	go this.connect()
	return nil
}

// 增加30秒 重连间隔
func (this *ServerSpec) connect() {
	this.rw.Lock()
	conn, err := openConnection(this.Server, this.TLSConfig, this.connectTimeout)
	if err != nil {
		this.rw.Unlock()
		if this.autoReconnect {
			go this.connect()
		}
		return
	}
	this.status = true
	this.ctx, this.cancelFunc = context.WithCancel(context.Background())
	this.conn = conn
	this.in = make(chan []byte, this.Config.RecvUnAckLimitW)
	this.out = make(chan []byte, this.Config.SendUnAckLimitK)
	this.recv = make(chan []byte, this.Config.RecvUnAckLimitW)
	this.send = make(chan []byte, this.Config.SendUnAckLimitK) // may not block!
	this.rw.Unlock()
	this.onConnect(this)
	this.wg.Add(3)
	go this.recvLoop()
	go this.sendLoop()
	go this.runHandler()
	this.runMonitor()
	this.wg.Wait()
	this.onConnectionLost(this)
	if this.autoReconnect {
		go this.connect()
	}
}
func (this *ServerSpec) Close() error {
	return nil
}

func (this *ServerSpec) SetOnConnectHandler(f OnConnectHandler) {
	this.onConnect = f
}
func (this *ServerSpec) SetConnectionLostHandler(f ConnectionLostHandler) {
	this.onConnectionLost = f
}

func openConnection(uri *url.URL, tlsc *tls.Config, timeout time.Duration) (net.Conn, error) {
	switch uri.Scheme {
	case "tcp":
		conn, err := net.DialTimeout("tcp", uri.Host, timeout)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case "ssl":
		fallthrough
	case "tls":
		fallthrough
	case "tcps":
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", uri.Host, tlsc)
		if err != nil {
			return nil, err
		}
		return conn, nil

	}
	return nil, errors.New("Unknown protocol")
}
