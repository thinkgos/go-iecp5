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

// Server the common server
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

// NewServer new a server
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
		Clog:    clog.NewWithPrefix("cs104 server =>"),
	}, nil
}

// ListenAndServer run the server
func (this *Server) ListenAndServer(addr string) {
	listen, err := net.Listen("tcp", addr)
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
			sess := &SrvSession{
				Config:  this.conf,
				params:  this.params,
				handler: this.handler,
				conn:    conn,
				in:      make(chan []byte, 1024),
				out:     make(chan []byte, 1024),
				recv:    make(chan []byte, 1024),
				send:    make(chan []byte, 1024), // may not block!

				Clog: this.Clog,
			}
			sess.run(ctx)
			this.wg.Done()
		}()
	}
}

// Close close the server
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
