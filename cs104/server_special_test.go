package cs104

import (
	"context"
	"github.com/thinkgos/go-iecp5/asdu"
	"testing"
)

func Test_serverSpec_SetOnConnectHandler(t *testing.T) {
	type fields struct {
		SrvSession  SrvSession
		option      ClientOption
		closeCancel context.CancelFunc
	}
	type args struct {
		f func(conn asdu.Connect)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"SetOnConnectHandler", fields{}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &serverSpec{
				SrvSession:  tt.fields.SrvSession,
				option:      tt.fields.option,
				closeCancel: tt.fields.closeCancel,
			}
			sf.SetOnConnectHandler(tt.args.f)
		})
	}
}

func Test_serverSpec_SetConnectionLostHandler(t *testing.T) {
	type fields struct {
		SrvSession  SrvSession
		option      ClientOption
		closeCancel context.CancelFunc
	}
	type args struct {
		f func(c asdu.Connect)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"SetConnectionLostHandler", fields{}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &serverSpec{
				SrvSession:  tt.fields.SrvSession,
				option:      tt.fields.option,
				closeCancel: tt.fields.closeCancel,
			}
			sf.SetConnectionLostHandler(tt.args.f)
		})
	}
}
