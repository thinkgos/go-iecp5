// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package cs104

import (
	"context"
	"errors"
	"math/rand"
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

	SetOnConnectHandler(f func(c asdu.Connect))
	SetConnectionLostHandler(f func(c asdu.Connect))

	LogMode(enable bool)
	SetLogProvider(p clog.LogProvider)
}

type serverSpec struct {
	SrvSession
	option      ClientOption
	closeCancel context.CancelFunc
}

// NewServerSpecial new special server
func NewServerSpecial(handler ServerHandlerInterface, o *ClientOption) ServerSpecial {
	return &serverSpec{
		SrvSession: SrvSession{
			config:  &o.config,
			params:  &o.params,
			handler: handler,

			rcvASDU:  make(chan []byte, 1024),
			sendASDU: make(chan []byte, 1024),
			rcvRaw:   make(chan []byte, 1024),
			sendRaw:  make(chan []byte, 1024), // may not block!

			Clog: clog.NewLogger("cs104 serverSpec => "),
		},
		option: *o,
	}
}

// SetOnConnectHandler set on connect handler
func (sf *serverSpec) SetOnConnectHandler(f func(conn asdu.Connect)) {
	sf.onConnection = f
}

// SetConnectionLostHandler set connection lost handler
func (sf *serverSpec) SetConnectionLostHandler(f func(c asdu.Connect)) {
	sf.connectionLost = f
}

// Start start the server,and return quickly,if it nil,the server will disconnected background,other failed
func (sf *serverSpec) Start() error {
	if sf.option.server == nil {
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

		sf.Debug("connecting server %+v", sf.option.server)
		conn, err := openConnection(sf.option.server, sf.option.TLSConfig, sf.config.ConnectTimeout0)
		if err != nil {
			sf.Error("connect failed, %v", err)
			if !sf.option.autoReconnect {
				return
			}
			time.Sleep(sf.option.reconnectInterval)
			continue
		}
		sf.Debug("connect success")
		sf.conn = conn
		sf.run(ctx)
		sf.Debug("disconnected server %+v", sf.option.server)
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
