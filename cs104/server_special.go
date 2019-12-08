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

// ServerSpecial server special interface
type ServerSpecial interface {
	asdu.Connect

	IsConnected() bool
	IsClosed() bool
	Start() error
	Close() error

	SetTLSConfig(t *tls.Config)
	AddRemoteServer(server string) error
	SetOnConnectHandler(f func(con net.Conn))
	SetConnectionLostHandler(f func(c ServerSpecial))

	LogMode(enable bool)
	SetLogProvider(p clog.LogProvider)
}

type serverSpec struct {
	SrvSession

	server            *url.URL      // 连接的服务器端
	autoReconnect     bool          // 是否启动重连
	reconnectInterval time.Duration // 重连间隔时间
	TLSConfig         *tls.Config
	onConnect         func(conn net.Conn)
	onConnectionLost  func(c ServerSpecial)

	closeCancel context.CancelFunc
}

// NewServerSpecial new special server
func NewServerSpecial(conf *Config, pm *asdu.Params, handler ServerHandlerInterface) (ServerSpecial, error) {
	if handler == nil {
		return nil, errors.New("invalid handler")
	}
	config := *conf
	params := *pm
	if err := config.Valid(); err != nil {
		return nil, err
	}
	if err := params.Valid(); err != nil {
		return nil, err
	}

	return &serverSpec{
		SrvSession: SrvSession{
			config:  &config,
			params:  &params,
			handler: handler,

			rcvASDU:  make(chan []byte, 1024),
			sendASDU: make(chan []byte, 1024),
			rcvRaw:   make(chan []byte, 1024),
			sendRaw:  make(chan []byte, 1024), // may not block!

			Clog: clog.NewWithPrefix("cs104 serverSpec => "),
		},
		autoReconnect:     true,
		reconnectInterval: DefaultReconnectInterval,
		onConnect:         func(conn net.Conn) {},
		onConnectionLost:  func(ServerSpecial) {},
	}, nil
}

// SetReconnectInterval set tcp  reconnect the host interval when connect failed after try
func (sf *serverSpec) SetReconnectInterval(t time.Duration) *serverSpec {
	sf.reconnectInterval = t
	return sf
}

// SetAutoReconnect enable auto reconnect
func (sf *serverSpec) SetAutoReconnect(b bool) *serverSpec {
	sf.autoReconnect = b
	return sf
}

// SetTLSConfig set tls config
func (sf *serverSpec) SetTLSConfig(t *tls.Config) *serverSpec {
	sf.TLSConfig = t
	return sf
}

// SetOnConnectHandler set on connect handler
func (sf *serverSpec) SetOnConnectHandler(f func(conn net.Conn)) *serverSpec {
	if f != nil {
		sf.onConnect = f
	}
	return sf
}

// SetConnectionLostHandler set connection lost handler
func (sf *serverSpec) SetConnectionLostHandler(f func(c ServerSpecial)) *serverSpec {
	if f != nil {
		sf.onConnectionLost = f
	}
	return sf
}

// AddRemoteServer adds a broker URI to the list of brokers to be used.
// The format should be scheme://host:port
// Default values for hostname is "127.0.0.1", for schema is "tcp://".
// An example broker URI would look like: tcp://foobar.com:1204
func (sf *serverSpec) AddRemoteServer(server string) error {
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

// Start start the server,and return quickly,if it nil,the server will disconnected background,other failed
func (sf *serverSpec) Start() error {
	if sf.server == nil {
		return errors.New("empty remote server")
	}

	go sf.running()
	return nil
}

// 增加重连间隔
func (sf *serverSpec) running() {
	var ctx context.Context

	sf.rwMux.Lock()
	if !atomic.CompareAndSwapUint32(&sf.status, initial, disconnected) {
		sf.rwMux.Unlock()
		return
	}
	ctx, sf.closeCancel = context.WithCancel(context.Background())
	sf.rwMux.Unlock()
	defer sf.setConnectStatus(initial)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		sf.Debug("connecting server %+v", sf.server)
		conn, err := openConnection(sf.server, sf.TLSConfig, sf.config.ConnectTimeout0)
		if err != nil {
			sf.Error("connect failed, %v", err)
			if !sf.autoReconnect {
				return
			}
			time.Sleep(sf.reconnectInterval)
			continue
		}
		sf.Debug("connect success")
		sf.conn = conn
		sf.onConnect(conn)
		sf.run(ctx)
		sf.onConnectionLost(sf)
		sf.Debug("disconnected server %+v", sf.server)
		select {
		case <-ctx.Done():
			return
		default:
			// 随机500ms-1s的重试，避免快速重试造成服务器许多无效连接
			time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(500)))
		}
	}
}

func (sf *serverSpec) IsClosed() bool {
	return sf.connectStatus() == initial
}

func (sf *serverSpec) Close() error {
	sf.rwMux.Lock()
	if sf.closeCancel != nil {
		sf.closeCancel()
	}
	sf.rwMux.Unlock()
	return nil
}
