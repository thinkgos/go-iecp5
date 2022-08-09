package cs104

import (
	"crypto/tls"
	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	type args struct {
		handler ServerHandlerInterface
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.handler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Close(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			if err := sf.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_ListenAndServer(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			if err := sf.ListenAndServer(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("ListenAndServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Params(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   *asdu.Params
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			if got := sf.Params(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Params() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Send(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	type args struct {
		a *asdu.ASDU
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			if err := sf.Send(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_SetConfig(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	type args struct {
		cfg Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			if got := sf.SetConfig(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_SetConnectionLostHandler(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	type args struct {
		f func(asdu.Connect)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			sf.SetConnectionLostHandler(tt.args.f)
		})
	}
}

func TestServer_SetInfoObjTimeZone(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	type args struct {
		zone *time.Location
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			sf.SetInfoObjTimeZone(tt.args.zone)
		})
	}
}

func TestServer_SetOnConnectionHandler(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	type args struct {
		f func(asdu.Connect)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			sf.SetOnConnectionHandler(tt.args.f)
		})
	}
}

func TestServer_SetParams(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	type args struct {
		p *asdu.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			if got := sf.SetParams(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_UnderlyingConn(t *testing.T) {
	type fields struct {
		config         Config
		params         asdu.Params
		handler        ServerHandlerInterface
		TLSConfig      *tls.Config
		mux            sync.Mutex
		sessions       map[*SrvSession]struct{}
		listen         net.Listener
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		Clog           clog.Clog
		wg             sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   net.Conn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Server{
				config:         tt.fields.config,
				params:         tt.fields.params,
				handler:        tt.fields.handler,
				TLSConfig:      tt.fields.TLSConfig,
				mux:            tt.fields.mux,
				sessions:       tt.fields.sessions,
				listen:         tt.fields.listen,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				Clog:           tt.fields.Clog,
				wg:             tt.fields.wg,
			}
			if got := sf.UnderlyingConn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnderlyingConn() = %v, want %v", got, tt.want)
			}
		})
	}
}
