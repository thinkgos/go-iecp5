package asdu

import (
	"reflect"
	"testing"
	"time"
)

func TestInterrogationCmd(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		qoi QualifierOfInterrogation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"cause not Activation and Deactivation",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				QOIGroup1},
			true,
		},
		{
			"C_IC_NA_1",
			args{
				newConn([]byte{byte(C_IC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x00, 0x00, 0x00, 21}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				QOIGroup1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InterrogationCmd(tt.args.c, tt.args.coa, tt.args.ca, tt.args.qoi); (err != nil) != tt.wantErr {
				t.Errorf("InterrogationCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCounterInterrogationCmd(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		qcc QualifierCountCall
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"C_CI_NA_1",
			args{
				newConn([]byte{byte(C_CI_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x00, 0x00, 0x00, 0x01}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				QualifierCountCall{QCCGroup1, QCCFrzRead}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CounterInterrogationCmd(tt.args.c, tt.args.coa, tt.args.ca, tt.args.qcc); (err != nil) != tt.wantErr {
				t.Errorf("CounterInterrogationCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadCmd(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		ioa InfoObjAddr
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"C_RD_NA_1",
			args{
				newConn([]byte{byte(C_RD_NA_1), 0x01, 0x05, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56}, t),
				CauseOfTransmission{Cause: Request},
				0x1234,
				0x567890},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadCmd(tt.args.c, tt.args.coa, tt.args.ca, tt.args.ioa); (err != nil) != tt.wantErr {
				t.Errorf("ReadCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClockSynchronizationCmd(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		t   time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"C_CS_NA_1",
			args{
				newConn(append([]byte{byte(C_CS_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x00, 0x00, 0x00}, tm0CP56Time2aBytes...), t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				tm0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ClockSynchronizationCmd(tt.args.c, tt.args.coa, tt.args.ca, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("ClockSynchronizationCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTestCommand(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"C_TS_NA_1",
			args{
				newConn([]byte{byte(C_TS_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x00, 0x00, 0x00, 0xaa, 0x55}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TestCommand(tt.args.c, tt.args.coa, tt.args.ca); (err != nil) != tt.wantErr {
				t.Errorf("TestCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResetProcessCmd(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		qrp QualifierOfResetProcessCmd
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"C_RP_NA_1",
			args{
				newConn([]byte{byte(C_RP_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x00, 0x00, 0x00, 0x01}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				QPRGeneralRest},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ResetProcessCmd(tt.args.c, tt.args.coa, tt.args.ca, tt.args.qrp); (err != nil) != tt.wantErr {
				t.Errorf("ResetProcessCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDelayAcquireCommand(t *testing.T) {
	type args struct {
		c    Connect
		coa  CauseOfTransmission
		ca   CommonAddr
		msec uint16
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"cause not act and spont",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				10000},
			true,
		},
		{
			"C_CD_NA_1",
			args{
				newConn([]byte{byte(C_CD_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x00, 0x00, 0x00, 0x10, 0x27}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				10000},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DelayAcquireCommand(tt.args.c, tt.args.coa, tt.args.ca, tt.args.msec); (err != nil) != tt.wantErr {
				t.Errorf("DelayAcquireCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTestCommandCP56Time2a(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		t   time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"C_TS_TA_1",
			args{
				newConn(append([]byte{byte(C_TS_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x00, 0x00, 0x00, 0xaa, 0x55}, tm0CP56Time2aBytes...), t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				tm0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TestCommandCP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("TestCommandCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestASDU_GetInterrogationCmd(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
		want1  QualifierOfInterrogation
	}{
		{
			"C_IC_NA_1",
			fields{ParamsWide, []byte{0x00, 0x00, 0x00, 21}},
			0,
			QOIGroup1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got, got1 := this.GetInterrogationCmd()
			if got != tt.want {
				t.Errorf("ASDU.GetInterrogationCmd() QOI = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ASDU.GetInterrogationCmd() InfoObjAddr = %v, want %v", got1, tt.want1)
			}

		})
	}
}

func TestASDU_GetCounterInterrogationCmd(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
		want1  QualifierCountCall
	}{
		{
			"C_CI_NA_1",
			fields{ParamsWide, []byte{0x00, 0x00, 0x00, 0x01}},
			0,
			QualifierCountCall{QCCGroup1, QCCFrzRead},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got, got1 := this.GetCounterInterrogationCmd()
			if got != tt.want {
				t.Errorf("ASDU.GetQuantityInterrogationCmd() InfoObjAddr = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ASDU.GetQuantityInterrogationCmd() QCC = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestASDU_GetReadCmd(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
	}{
		{
			"C_RD_NA_1",
			fields{ParamsWide, []byte{0x90, 0x78, 0x56}},
			0x567890,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got := this.GetReadCmd()
			if got != tt.want {
				t.Errorf("ASDU.GetReadCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetClockSynchronizationCmd(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
		want1  time.Time
	}{
		{
			"C_CS_NA_1",
			fields{ParamsWide, append([]byte{0x00, 0x00, 0x00}, tm0CP56Time2aBytes...)},
			0,
			tm0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got, got1 := this.GetClockSynchronizationCmd()
			if got != tt.want {
				t.Errorf("ASDU.GetClockSynchronizationCmd() InfoObjAddr = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ASDU.GetClockSynchronizationCmd() time = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestASDU_GetTestCommand(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
		want1  bool
	}{
		{
			"C_CS_NA_1",
			fields{ParamsWide, []byte{0x00, 0x00, 0x00, 0xaa, 0x55}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got, got1 := this.GetTestCommand()
			if got != tt.want {
				t.Errorf("ASDU.GetTestCommand() InfoObjAddr = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ASDU.GetTestCommand() bool = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestASDU_GetResetProcessCmd(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
		want1  QualifierOfResetProcessCmd
	}{
		{
			"C_RP_NA_1",
			fields{ParamsWide, []byte{0x00, 0x00, 0x00, 0x01}},
			0,
			QPRGeneralRest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got, got1 := this.GetResetProcessCmd()
			if got != tt.want {
				t.Errorf("ASDU.GetResetProcessCmd() InfoObjAddr = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ASDU.GetResetProcessCmd() QOP = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestASDU_GetDelayAcquireCommand(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
		want1  uint16
	}{
		{
			"C_CD_NA_1",
			fields{ParamsWide, []byte{0x00, 0x00, 0x00, 0x10, 0x27}},
			0,
			10000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got, got1 := this.GetDelayAcquireCommand()
			if got != tt.want {
				t.Errorf("ASDU.GetDelayAcquireCommand() InfoObjAddr = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ASDU.GetDelayAcquireCommand() msec = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestASDU_GetTestCommandCP56Time2a(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   InfoObjAddr
		want1  bool
		want2  time.Time
	}{
		{
			"C_CS_TA_1",
			fields{ParamsWide, append([]byte{0x00, 0x00, 0x00, 0xaa, 0x55}, tm0CP56Time2aBytes...)},
			0,
			true,
			tm0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			got, got1, got2 := sf.GetTestCommandCP56Time2a()
			if got != tt.want {
				t.Errorf("GetTestCommandCP56Time2a() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetTestCommandCP56Time2a() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetTestCommandCP56Time2a() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
