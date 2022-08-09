package cs104

import (
	"context"
	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
	"net"
	"reflect"
	"sync"
	"testing"
)

func TestSrvSession_Params(t *testing.T) {
	type fields struct {
		config         *Config
		params         *asdu.Params
		conn           net.Conn
		handler        ServerHandlerInterface
		rcvASDU        chan []byte
		sendASDU       chan []byte
		rcvRaw         chan []byte
		sendRaw        chan []byte
		seqNoSend      uint16
		ackNoSend      uint16
		seqNoRcv       uint16
		ackNoRcv       uint16
		pending        []seqPending
		status         uint32
		rwMux          sync.RWMutex
		Clog           clog.Clog
		onConnection   func(asdu.Connect)
		connectionLost func(asdu.Connect)
		wg             sync.WaitGroup
		cancel         context.CancelFunc
		ctx            context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   *asdu.Params
	}{
		{"Params", fields{params: &asdu.Params{}}, &asdu.Params{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &SrvSession{
				config:         tt.fields.config,
				params:         tt.fields.params,
				conn:           tt.fields.conn,
				handler:        tt.fields.handler,
				rcvASDU:        tt.fields.rcvASDU,
				sendASDU:       tt.fields.sendASDU,
				rcvRaw:         tt.fields.rcvRaw,
				sendRaw:        tt.fields.sendRaw,
				seqNoSend:      tt.fields.seqNoSend,
				ackNoSend:      tt.fields.ackNoSend,
				seqNoRcv:       tt.fields.seqNoRcv,
				ackNoRcv:       tt.fields.ackNoRcv,
				pending:        tt.fields.pending,
				status:         tt.fields.status,
				Clog:           tt.fields.Clog,
				onConnection:   tt.fields.onConnection,
				connectionLost: tt.fields.connectionLost,
				cancel:         tt.fields.cancel,
				ctx:            tt.fields.ctx,
			}
			if got := sf.Params(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Params() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seqNoCount(t *testing.T) {
	type args struct {
		nextAckNo uint16
		nextSeqNo uint16
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{"count", args{nextAckNo: 5, nextSeqNo: 8}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seqNoCount(tt.args.nextAckNo, tt.args.nextSeqNo); got != tt.want {
				t.Errorf("seqNoCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
