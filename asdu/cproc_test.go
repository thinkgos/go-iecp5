package asdu

import (
	"math"
	"net"
	"reflect"
	"testing"
	"time"
)

type conn struct {
	p    *Params
	want []byte
	t    *testing.T
}

func newConn(want []byte, t *testing.T) *conn {
	return &conn{ParamsWide, want, t}
}

func (sf *conn) Params() *Params          { return sf.p }
func (sf *conn) UnderlyingConn() net.Conn { return nil }

// Send
func (sf *conn) Send(u *ASDU) error {
	data, err := u.MarshalBinary()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(sf.want, data) {
		sf.t.Errorf("Send() out = % x, want % x", data, sf.want)
	}
	return nil
}

func TestSingleCmd(t *testing.T) {
	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    SingleCommandInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid type id",
			args{
				newConn(nil, t),
				0,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SingleCommandInfo{}},
			true},
		{
			"cause not Activation and Deactivation",
			args{
				newConn(nil, t),
				C_SC_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SingleCommandInfo{}},
			true},
		{
			"C_SC_NA_1",
			args{
				newConn([]byte{byte(C_SC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x05}, t),
				C_SC_NA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SingleCommandInfo{
					0x567890,
					true,
					QualifierOfCommand{QOCShortPulseDuration, false},
					time.Time{}}},
			false},
		{
			"C_SC_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_SC_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x04}, tm0CP56Time2aBytes...), t),
				C_SC_TA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SingleCommandInfo{
					0x567890, false,
					QualifierOfCommand{QOCShortPulseDuration, false},
					tm0}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SingleCmd(tt.args.c, tt.args.typeID, tt.args.coa, tt.args.ca, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("SingleCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDoubleCmd(t *testing.T) {
	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    DoubleCommandInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid type id",
			args{
				newConn(nil, t),
				0,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				DoubleCommandInfo{}},
			true},
		{
			"cause not Activation and Deactivation",
			args{
				newConn(nil, t),
				C_DC_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				DoubleCommandInfo{}},
			true},
		{
			"C_DC_NA_1",
			args{
				newConn([]byte{byte(C_DC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x05}, t),
				C_DC_NA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				DoubleCommandInfo{
					0x567890,
					DCOOn,
					QualifierOfCommand{QOCShortPulseDuration, false},
					time.Time{}}},
			false},
		{
			"C_DC_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_DC_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...), t),
				C_DC_TA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				DoubleCommandInfo{
					0x567890,
					DCOOff,
					QualifierOfCommand{QOCShortPulseDuration, false},
					tm0}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoubleCmd(tt.args.c, tt.args.typeID, tt.args.coa, tt.args.ca, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("DoubleCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStepCmd(t *testing.T) {
	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    StepCommandInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid type id",
			args{
				newConn(nil, t),
				0,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				StepCommandInfo{}},
			true},
		{
			"cause not Activation and Deactivation", args{
				newConn(nil, t),
				C_RC_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				StepCommandInfo{}},
			true},
		{
			"C_RC_NA_1",
			args{
				newConn([]byte{byte(C_RC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x05}, t),
				C_RC_NA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				StepCommandInfo{
					0x567890,
					SCOStepDown,
					QualifierOfCommand{QOCShortPulseDuration, false},
					time.Time{}}},
			false},
		{
			"C_RC_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_RC_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...), t),
				C_RC_TA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				StepCommandInfo{
					0x567890,
					SCOStepUP,
					QualifierOfCommand{QOCShortPulseDuration, false},
					tm0}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StepCmd(tt.args.c, tt.args.typeID, tt.args.coa, tt.args.ca, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("StepCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetpointCmdNormal(t *testing.T) {
	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    SetpointCommandNormalInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid type id",
			args{
				newConn(nil, t),
				0,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandNormalInfo{}},
			true},
		{
			"cause not Activation and Deactivation",
			args{
				newConn(nil, t),
				C_SE_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SetpointCommandNormalInfo{}},
			true},
		{
			"C_SE_NA_1",
			args{
				newConn([]byte{byte(C_SE_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, t),
				C_SE_NA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandNormalInfo{
					0x567890,
					100,
					QualifierOfSetpointCmd{1, false},
					time.Time{}}},
			false},
		{
			"C_SE_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_SE_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, tm0CP56Time2aBytes...), t),
				C_SE_TA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandNormalInfo{
					0x567890, 100,
					QualifierOfSetpointCmd{1, false},
					tm0}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetpointCmdNormal(tt.args.c, tt.args.typeID, tt.args.coa, tt.args.ca, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("SetpointCmdNormal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetpointCmdScaled(t *testing.T) {
	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    SetpointCommandScaledInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid type id",
			args{
				newConn(nil, t),
				0,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandScaledInfo{}},
			true},
		{
			"cause not Activation and Deactivation",
			args{
				newConn(nil, t),
				C_SE_NB_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SetpointCommandScaledInfo{}},
			true},
		{
			"C_SE_NB_1",
			args{
				newConn([]byte{byte(C_SE_NB_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, t),
				C_SE_NB_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandScaledInfo{
					0x567890,
					100,
					QualifierOfSetpointCmd{1, false},
					time.Time{}}},
			false},
		{
			"C_SE_TB_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_SE_TB_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, tm0CP56Time2aBytes...), t),
				C_SE_TB_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandScaledInfo{
					0x567890, 100,
					QualifierOfSetpointCmd{1, false},
					tm0}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetpointCmdScaled(tt.args.c, tt.args.typeID, tt.args.coa, tt.args.ca, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("SetpointCmdScaled() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetpointCmdFloat(t *testing.T) {
	bits := math.Float32bits(100)

	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    SetpointCommandFloatInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid type id",
			args{
				newConn(nil, t),
				0,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandFloatInfo{}},
			true},
		{
			"cause not Activation and Deactivation",
			args{
				newConn(nil, t),
				C_SE_NC_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SetpointCommandFloatInfo{}},
			true},
		{
			"C_SE_NC_1",
			args{
				newConn([]byte{byte(C_SE_NC_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}, t),
				C_SE_NC_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandFloatInfo{
					0x567890,
					100,
					QualifierOfSetpointCmd{1, false},
					time.Time{}}},
			false},
		{
			"C_SE_TC_1 CP56Time2a",
			args{
				newConn(
					append([]byte{byte(C_SE_TC_1), 0x01, 0x06, 0x00, 0x34, 0x12,
						0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}, tm0CP56Time2aBytes...), t),
				C_SE_TC_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				SetpointCommandFloatInfo{
					0x567890, 100,
					QualifierOfSetpointCmd{1, false},
					tm0}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetpointCmdFloat(tt.args.c, tt.args.typeID, tt.args.coa, tt.args.ca, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("SetpointCmdFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBitsString32Cmd(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		coa        CauseOfTransmission
		commonAddr CommonAddr
		cmd        BitsString32CommandInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid type id",
			args{
				newConn(nil, t),
				0,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				BitsString32CommandInfo{}},
			true},
		{
			"cause not Activation and Deactivation",
			args{
				newConn(nil, t),
				C_BO_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				BitsString32CommandInfo{}},
			true},
		{
			"C_BO_NA_1",
			args{
				newConn([]byte{byte(C_BO_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x00, 0x00}, t),
				C_BO_NA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				BitsString32CommandInfo{
					0x567890,
					100,
					time.Time{}}},
			false},
		{
			"C_BO_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_BO_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x00, 0x00}, tm0CP56Time2aBytes...), t),
				C_BO_TA_1,
				CauseOfTransmission{Cause: Activation},
				0x1234,
				BitsString32CommandInfo{
					0x567890, 100,
					tm0}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BitsString32Cmd(tt.args.c, tt.args.typeID, tt.args.coa, tt.args.commonAddr, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("BitsString32Cmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestASDU_GetSingleCmd(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   SingleCommandInfo
	}{
		{
			"C_SC_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SC_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x05}},
			SingleCommandInfo{
				0x567890,
				true,
				QualifierOfCommand{QOCShortPulseDuration, false},
				time.Time{}},
		},
		{
			"C_SC_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SC_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x04}, tm0CP56Time2aBytes...)},
			SingleCommandInfo{
				0x567890,
				false,
				QualifierOfCommand{QOCShortPulseDuration, false},
				tm0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got := sf.GetSingleCmd()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetSingleCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetDoubleCmd(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   DoubleCommandInfo
	}{
		{
			"C_DC_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_DC_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x05}},
			DoubleCommandInfo{
				0x567890,
				DCOOn,
				QualifierOfCommand{QOCShortPulseDuration, false},
				time.Time{},
			},
		},
		{
			"C_DC_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_DC_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...)},
			DoubleCommandInfo{
				0x567890,
				DCOOff,
				QualifierOfCommand{QOCShortPulseDuration, false},
				tm0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got := sf.GetDoubleCmd()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetDoubleCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetStepCmd(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   StepCommandInfo
	}{
		{
			"C_RC_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_RC_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x05}},
			StepCommandInfo{
				0x567890,
				SCOStepDown,
				QualifierOfCommand{QOCShortPulseDuration, false},
				time.Time{}},
		},
		{
			"C_RC_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_RC_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...)},
			StepCommandInfo{
				0x567890,
				SCOStepUP,
				QualifierOfCommand{QOCShortPulseDuration, false},
				tm0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got := sf.GetStepCmd()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetStepCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetSetpointNormalCmd(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   SetpointCommandNormalInfo
	}{
		{
			"C_SE_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}},
			SetpointCommandNormalInfo{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				time.Time{}},
		},
		{
			"C_SE_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, tm0CP56Time2aBytes...)},
			SetpointCommandNormalInfo{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				tm0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got := sf.GetSetpointNormalCmd()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetSetpointNormalCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetSetpointCmdScaled(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   SetpointCommandScaledInfo
	}{
		{
			"C_SE_NB_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_NB_1},
				[]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}},
			SetpointCommandScaledInfo{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				time.Time{}},
		},
		{
			"C_SE_TB_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_TB_1},
				append([]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, tm0CP56Time2aBytes...)},
			SetpointCommandScaledInfo{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				tm0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got := sf.GetSetpointCmdScaled()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetSetpointCmdScaled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetSetpointFloatCmd(t *testing.T) {
	bits := math.Float32bits(100)

	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   SetpointCommandFloatInfo
	}{
		{
			"C_SE_NC_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_NC_1},
				[]byte{0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}},
			SetpointCommandFloatInfo{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				time.Time{}},
		},
		{
			"C_SE_TC_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_TC_1},
				append([]byte{0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}, tm0CP56Time2aBytes...)},
			SetpointCommandFloatInfo{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				tm0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got := sf.GetSetpointFloatCmd()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetSetpointFloatCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetBitsString32Cmd(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   BitsString32CommandInfo
	}{
		{
			"C_BO_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_BO_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x00, 0x00}},
			BitsString32CommandInfo{
				0x567890,
				100,
				time.Time{}},
		},
		{
			"C_BO_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_BO_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x00, 0x00}, tm0CP56Time2aBytes...)},
			BitsString32CommandInfo{
				0x567890,
				100,
				tm0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got := sf.GetBitsString32Cmd()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetBitsString32Cmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
