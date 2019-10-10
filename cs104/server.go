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

// timeoutResolution is seconds according to companion standard 104,
// subclass 6.9, caption "Definition of time outs". However, then
// of a second make this system much more responsive i.c.w. S-frames.
const timeoutResolution = 100 * time.Millisecond

// Server the common server
type Server struct {
	TimeoutResolution
	conf      *Config
	params    *asdu.Params
	handler   ServerHandlerInterface
	TLSConfig *tls.Config
	mux       sync.Mutex
	sessions  map[*SrvSession]struct{}
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
		conf:     conf,
		params:   params,
		handler:  handler,
		sessions: make(map[*SrvSession]struct{}),
		Clog:     clog.NewWithPrefix("cs104 server =>"),
	}, nil
}

// ListenAndServer run the server
func (this *Server) ListenAndServer(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		this.Error("server run failed, %v", err)
		return
	}
	this.mux.Lock()
	this.listen = listen
	this.mux.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		this.Close()
		this.Debug("server stop")
	}()
	this.Debug("server run")
	for {
		conn, err := listen.Accept()
		if err != nil {
			this.Error("server run failed, %v", err)
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
			this.mux.Lock()
			this.sessions[sess] = struct{}{}
			this.mux.Unlock()
			sess.run(ctx)
			this.mux.Lock()
			delete(this.sessions, sess)
			this.mux.Unlock()
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

// Send imp interface Connect
func (this *Server) Send(a *asdu.ASDU) error {
	this.mux.Lock()
	for k := range this.sessions {
		_ = k.Send(a.Clone())
	}
	this.mux.Unlock()
	return nil
}

// Params imp interface Connect
func (this *Server) Params() *asdu.Params { return this.params }

// UnderlyingConn imp interface Connect
func (this *Server) UnderlyingConn() net.Conn { return nil }
