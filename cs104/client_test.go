package cs104

import (
	"context"
	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_ClockSynchronizationCmd(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		coa asdu.CauseOfTransmission
		ca  asdu.CommonAddr
		t   time.Time
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.ClockSynchronizationCmd(tt.args.coa, tt.args.ca, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("ClockSynchronizationCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CounterInterrogationCmd(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		coa asdu.CauseOfTransmission
		ca  asdu.CommonAddr
		qcc asdu.QualifierCountCall
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.CounterInterrogationCmd(tt.args.coa, tt.args.ca, tt.args.qcc); (err != nil) != tt.wantErr {
				t.Errorf("CounterInterrogationCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DelayAcquireCommand(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		coa  asdu.CauseOfTransmission
		ca   asdu.CommonAddr
		msec uint16
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.DelayAcquireCommand(tt.args.coa, tt.args.ca, tt.args.msec); (err != nil) != tt.wantErr {
				t.Errorf("DelayAcquireCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_InterrogationCmd(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		coa asdu.CauseOfTransmission
		ca  asdu.CommonAddr
		qoi asdu.QualifierOfInterrogation
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.InterrogationCmd(tt.args.coa, tt.args.ca, tt.args.qoi); (err != nil) != tt.wantErr {
				t.Errorf("InterrogationCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_IsConnected(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if got := sf.IsConnected(); got != tt.want {
				t.Errorf("IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Params(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if got := sf.Params(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Params() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ReadCmd(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		coa asdu.CauseOfTransmission
		ca  asdu.CommonAddr
		ioa asdu.InfoObjAddr
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.ReadCmd(tt.args.coa, tt.args.ca, tt.args.ioa); (err != nil) != tt.wantErr {
				t.Errorf("ReadCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ResetProcessCmd(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		coa asdu.CauseOfTransmission
		ca  asdu.CommonAddr
		qrp asdu.QualifierOfResetProcessCmd
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.ResetProcessCmd(tt.args.coa, tt.args.ca, tt.args.qrp); (err != nil) != tt.wantErr {
				t.Errorf("ResetProcessCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Send(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.Send(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SendStartDt(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.SendStartDt()
		})
	}
}

func TestClient_SendStopDt(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.SendStopDt()
		})
	}
}

func TestClient_SetConnectionLostHandler(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		f func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if got := sf.SetConnectionLostHandler(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetConnectionLostHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SetOnConnectHandler(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		f func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if got := sf.SetOnConnectHandler(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetOnConnectHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Start(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UnderlyingConn(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if got := sf.UnderlyingConn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnderlyingConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_cleanUp(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.cleanUp()
		})
	}
}

func TestClient_clientHandler(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		asduPack *asdu.ASDU
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.clientHandler(tt.args.asduPack); (err != nil) != tt.wantErr {
				t.Errorf("clientHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_connectStatus(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if got := sf.connectStatus(); got != tt.want {
				t.Errorf("connectStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_handlerLoop(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.handlerLoop()
		})
	}
}

func TestClient_recvLoop(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.recvLoop()
		})
	}
}

func TestClient_run(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		ctx context.Context
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.run(tt.args.ctx)
		})
	}
}

func TestClient_running(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.running()
		})
	}
}

func TestClient_sendLoop(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.sendLoop()
		})
	}
}

func TestClient_sendUFrame(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		which byte
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.sendUFrame(tt.args.which)
		})
	}
}

func TestClient_setConnectStatus(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		status uint32
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			sf.setConnectStatus(tt.args.status)
		})
	}
}

func TestClient_updateAckNoOut(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		ackNo uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if gotOk := sf.updateAckNoOut(tt.args.ackNo); gotOk != tt.wantOk {
				t.Errorf("updateAckNoOut() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestCommand(t *testing.T) {
	type fields struct {
		option                 ClientOption
		conn                   net.Conn
		handler                ClientHandlerInterface
		rcvASDU                chan []byte
		sendASDU               chan []byte
		rcvRaw                 chan []byte
		sendRaw                chan []byte
		seqNoSend              uint16
		ackNoSend              uint16
		seqNoRcv               uint16
		ackNoRcv               uint16
		pending                []seqPending
		startDtActiveSendSince atomic.Value
		stopDtActiveSendSince  atomic.Value
		status                 uint32
		rwMux                  sync.RWMutex
		isActive               uint32
		Clog                   clog.Clog
		wg                     sync.WaitGroup
		ctx                    context.Context
		cancel                 context.CancelFunc
		closeCancel            context.CancelFunc
		onConnect              func(c *Client)
		onConnectionLost       func(c *Client)
	}
	type args struct {
		coa asdu.CauseOfTransmission
		ca  asdu.CommonAddr
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
			sf := &Client{
				option:                 tt.fields.option,
				conn:                   tt.fields.conn,
				handler:                tt.fields.handler,
				rcvASDU:                tt.fields.rcvASDU,
				sendASDU:               tt.fields.sendASDU,
				rcvRaw:                 tt.fields.rcvRaw,
				sendRaw:                tt.fields.sendRaw,
				seqNoSend:              tt.fields.seqNoSend,
				ackNoSend:              tt.fields.ackNoSend,
				seqNoRcv:               tt.fields.seqNoRcv,
				ackNoRcv:               tt.fields.ackNoRcv,
				pending:                tt.fields.pending,
				startDtActiveSendSince: tt.fields.startDtActiveSendSince,
				stopDtActiveSendSince:  tt.fields.stopDtActiveSendSince,
				status:                 tt.fields.status,
				rwMux:                  tt.fields.rwMux,
				isActive:               tt.fields.isActive,
				Clog:                   tt.fields.Clog,
				wg:                     tt.fields.wg,
				ctx:                    tt.fields.ctx,
				cancel:                 tt.fields.cancel,
				closeCancel:            tt.fields.closeCancel,
				onConnect:              tt.fields.onConnect,
				onConnectionLost:       tt.fields.onConnectionLost,
			}
			if err := sf.TestCommand(tt.args.coa, tt.args.ca); (err != nil) != tt.wantErr {
				t.Errorf("TestCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		handler ClientHandlerInterface
		o       *ClientOption
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.handler, tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
