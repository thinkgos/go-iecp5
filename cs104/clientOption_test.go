package cs104

import (
	"crypto/tls"
	"github.com/thinkgos/go-iecp5/asdu"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestClientOption_AddRemoteServer(t *testing.T) {
	type fields struct {
		config            Config
		params            asdu.Params
		server            *url.URL
		autoReconnect     bool
		reconnectInterval time.Duration
		TLSConfig         *tls.Config
	}
	type args struct {
		server string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"AddRemoteServer1", fields{}, args{}, false},
		{"AddRemoteServer2", fields{autoReconnect: true}, args{server: "2333"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ClientOption{
				config:            tt.fields.config,
				params:            tt.fields.params,
				server:            tt.fields.server,
				autoReconnect:     tt.fields.autoReconnect,
				reconnectInterval: tt.fields.reconnectInterval,
				TLSConfig:         tt.fields.TLSConfig,
			}
			if err := sf.AddRemoteServer(tt.args.server); (err != nil) != tt.wantErr {
				t.Errorf("AddRemoteServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientOption_SetAutoReconnect(t *testing.T) {
	type fields struct {
		config            Config
		params            asdu.Params
		server            *url.URL
		autoReconnect     bool
		reconnectInterval time.Duration
		TLSConfig         *tls.Config
	}
	type args struct {
		b bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ClientOption
	}{
		{"SetAutoReconnect", fields{}, args{b: false}, &ClientOption{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ClientOption{
				config:            tt.fields.config,
				params:            tt.fields.params,
				server:            tt.fields.server,
				autoReconnect:     tt.fields.autoReconnect,
				reconnectInterval: tt.fields.reconnectInterval,
				TLSConfig:         tt.fields.TLSConfig,
			}
			if got := sf.SetAutoReconnect(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAutoReconnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientOption_SetConfig(t *testing.T) {
	type fields struct {
		config            Config
		params            asdu.Params
		server            *url.URL
		autoReconnect     bool
		reconnectInterval time.Duration
		TLSConfig         *tls.Config
	}
	type args struct {
		cfg Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ClientOption
	}{
		{"SetConfig", fields{}, args{}, &ClientOption{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ClientOption{
				config:            tt.fields.config,
				params:            tt.fields.params,
				server:            tt.fields.server,
				autoReconnect:     tt.fields.autoReconnect,
				reconnectInterval: tt.fields.reconnectInterval,
				TLSConfig:         tt.fields.TLSConfig,
			}
			if got := sf.SetConfig(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientOption_SetParams(t *testing.T) {
	type fields struct {
		config            Config
		params            asdu.Params
		server            *url.URL
		autoReconnect     bool
		reconnectInterval time.Duration
		TLSConfig         *tls.Config
	}
	type args struct {
		p *asdu.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ClientOption
	}{
		{"SetParams", fields{}, args{}, &ClientOption{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ClientOption{
				config:            tt.fields.config,
				params:            tt.fields.params,
				server:            tt.fields.server,
				autoReconnect:     tt.fields.autoReconnect,
				reconnectInterval: tt.fields.reconnectInterval,
				TLSConfig:         tt.fields.TLSConfig,
			}
			if got := sf.SetParams(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientOption_SetReconnectInterval(t *testing.T) {
	type fields struct {
		config            Config
		params            asdu.Params
		server            *url.URL
		autoReconnect     bool
		reconnectInterval time.Duration
		TLSConfig         *tls.Config
	}
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ClientOption
	}{
		{"SetReconnectInterval", fields{}, args{}, &ClientOption{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ClientOption{
				config:            tt.fields.config,
				params:            tt.fields.params,
				server:            tt.fields.server,
				autoReconnect:     tt.fields.autoReconnect,
				reconnectInterval: tt.fields.reconnectInterval,
				TLSConfig:         tt.fields.TLSConfig,
			}
			if got := sf.SetReconnectInterval(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetReconnectInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientOption_SetTLSConfig(t *testing.T) {
	type fields struct {
		config            Config
		params            asdu.Params
		server            *url.URL
		autoReconnect     bool
		reconnectInterval time.Duration
		TLSConfig         *tls.Config
	}
	type args struct {
		t *tls.Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ClientOption
	}{
		{"SetTLSConfig", fields{}, args{}, &ClientOption{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ClientOption{
				config:            tt.fields.config,
				params:            tt.fields.params,
				server:            tt.fields.server,
				autoReconnect:     tt.fields.autoReconnect,
				reconnectInterval: tt.fields.reconnectInterval,
				TLSConfig:         tt.fields.TLSConfig,
			}
			if got := sf.SetTLSConfig(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetTLSConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOption(t *testing.T) {
	tests := []struct {
		name string
		want *ClientOption
	}{
		{"option", &ClientOption{config: DefaultConfig(), params: *asdu.ParamsWide, autoReconnect: true, reconnectInterval: DefaultReconnectInterval}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOption(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
