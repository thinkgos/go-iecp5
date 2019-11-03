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
func (sf *Server) ListenAndServer(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		sf.Error("server run failed, %v", err)
		return
	}
	sf.mux.Lock()
	sf.listen = listen
	sf.mux.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		_ = sf.Close()
		sf.Debug("server stop")
	}()
	sf.Debug("server run")
	for {
		conn, err := listen.Accept()
		if err != nil {
			sf.Error("server run failed, %v", err)
			return
		}

		sf.wg.Add(1)
		go func() {
			sess := &SrvSession{
				Config:  sf.conf,
				params:  sf.params,
				handler: sf.handler,
				conn:    conn,
				in:      make(chan []byte, 1024),
				out:     make(chan []byte, 1024),
				recv:    make(chan []byte, 1024),
				send:    make(chan []byte, 1024), // may not block!

				Clog: sf.Clog,
			}
			sf.mux.Lock()
			sf.sessions[sess] = struct{}{}
			sf.mux.Unlock()
			sess.run(ctx)
			sf.mux.Lock()
			delete(sf.sessions, sess)
			sf.mux.Unlock()
			sf.wg.Done()
		}()
	}
}

// Close close the server
func (sf *Server) Close() error {
	var err error

	sf.mux.Lock()
	if sf.listen != nil {
		err = sf.listen.Close()
		sf.listen = nil
	}
	sf.mux.Unlock()
	sf.wg.Wait()
	return err
}

// Send imp interface Connect
func (sf *Server) Send(a *asdu.ASDU) error {
	sf.mux.Lock()
	for k := range sf.sessions {
		_ = k.Send(a.Clone())
	}
	sf.mux.Unlock()
	return nil
}

// Params imp interface Connect
func (sf *Server) Params() *asdu.Params { return sf.params }

// UnderlyingConn imp interface Connect
func (sf *Server) UnderlyingConn() net.Conn { return nil }
