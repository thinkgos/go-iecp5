package cs104

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

// TimeoutResolution is seconds according to companion standard 104,
// subclass 6.9, caption "Definition of time outs". However, then
// of a second make this system much more responsive i.c.w. S-frames.
const timeoutResolution = 100 * time.Millisecond

type Server struct {
	conf      *Config
	params    *asdu.Params
	handler   ServerHandlerInterface
	TLSConfig *tls.Config
	mux       sync.Mutex
	listen    net.Listener
	*clog.Clog
	wg sync.WaitGroup
}

func NewServer(conf *Config, params *asdu.Params, handler ServerHandlerInterface) (*Server, error) {
	if handler == nil {
		return nil, errors.New("invalid handler")
	}
	if err := conf.Valid(); err != nil {
		return nil, err
	}
	if err := params.Valid(); err != nil {
		return nil, err
	}

	return &Server{
		conf:    conf,
		params:  params,
		handler: handler,
		Clog:    clog.NewWithPrefix("cs104 server => "),
	}, nil
}

func (this *Server) Run() {
	listen, err := net.Listen("tcp", ":2404")
	if err != nil {
		this.Error("Server run failed, %v", err)
		return
	}
	this.mux.Lock()
	this.listen = listen
	this.mux.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		this.Close()
		this.Debug("Server stop")
	}()
	this.Debug("Server run")
	for {
		conn, err := listen.Accept()
		if err != nil {
			this.Error("Server run failed, %v", err)
			return
		}

		this.wg.Add(1)
		go func() {

			sess := &Session{
				Config:  this.conf,
				params:  this.params,
				handler: this.handler,

				in:   make(chan []byte, this.conf.RecvUnAckLimitW),
				out:  make(chan []byte, this.conf.SendUnAckLimitK),
				recv: make(chan []byte, this.conf.RecvUnAckLimitW),
				send: make(chan []byte, this.conf.SendUnAckLimitK), // may not block!

				Clog: this.Clog,
			}
			sess.run(ctx, conn)
			this.wg.Done()
		}()
	}
}

func (this *Server) Close() error {
	var err error

	this.mux.Lock()
	if this.listen != nil {
		err = this.listen.Close()
		this.listen = nil
	}
	this.mux.Unlock()
	this.wg.Wait()
	return err
}
