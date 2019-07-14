package asdu

import (
	"math"
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

func (this *conn) Params() *Params { return this.p }

// Send
func (this *conn) Send(u *ASDU) error {
	data, err := u.MarshalBinary()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(this.want, data) {
		this.t.Errorf("Send() out = % x, want % x", data, this.want)
	}
	return nil
}

func TestSingleCmd(t *testing.T) {
	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    SingleCommandObject
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				SingleCommandObject{}},
			true},
		{
			"cause not Act and Deact",
			args{
				newConn(nil, t),
				C_SC_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SingleCommandObject{}},
			true},
		{
			"C_SC_NA_1",
			args{
				newConn([]byte{byte(C_SC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x05}, t),
				C_SC_NA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				SingleCommandObject{
					0x567890,
					true,
					QualifierOfCommand{QOCShortPulse, false},
					time.Time{}}},
			false},
		{
			"C_SC_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_SC_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x04}, tm0CP56Time2aBytes...), t),
				C_SC_TA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				SingleCommandObject{
					0x567890, false,
					QualifierOfCommand{QOCShortPulse, false},
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
		cmd    DoubleCommandObject
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				DoubleCommandObject{}},
			true},
		{
			"cause not Act and Deact",
			args{
				newConn(nil, t),
				C_DC_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				DoubleCommandObject{}},
			true},
		{
			"C_DC_NA_1",
			args{
				newConn([]byte{byte(C_DC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x05}, t),
				C_DC_NA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				DoubleCommandObject{
					0x567890,
					DCOOn,
					QualifierOfCommand{QOCShortPulse, false},
					time.Time{}}},
			false},
		{
			"C_DC_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_DC_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...), t),
				C_DC_TA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				DoubleCommandObject{
					0x567890,
					DCOOff,
					QualifierOfCommand{QOCShortPulse, false},
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
		cmd    StepCommandObject
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				StepCommandObject{}},
			true},
		{
			"cause not Act and Deact", args{
				newConn(nil, t),
				C_RC_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				StepCommandObject{}},
			true},
		{
			"C_RC_NA_1",
			args{
				newConn([]byte{byte(C_RC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x05}, t),
				C_RC_NA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				StepCommandObject{
					0x567890,
					SCOStepDown,
					QualifierOfCommand{QOCShortPulse, false},
					time.Time{}}},
			false},
		{
			"C_RC_TA_1 CP56Time2a",
			args{
				newConn(append([]byte{byte(C_RC_TA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...), t),
				C_RC_TA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				StepCommandObject{
					0x567890,
					SCOStepUP,
					QualifierOfCommand{QOCShortPulse, false},
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
		cmd    SetpointCommandNormalObject
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandNormalObject{}},
			true},
		{
			"cause not Act and Deact",
			args{
				newConn(nil, t),
				C_SE_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SetpointCommandNormalObject{}},
			true},
		{
			"C_SE_NA_1",
			args{
				newConn([]byte{byte(C_SE_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, t),
				C_SE_NA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandNormalObject{
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandNormalObject{
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
		cmd    SetpointCommandScaledObject
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandScaledObject{}},
			true},
		{
			"cause not Act and Deact",
			args{
				newConn(nil, t),
				C_SE_NB_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SetpointCommandScaledObject{}},
			true},
		{
			"C_SE_NB_1",
			args{
				newConn([]byte{byte(C_SE_NB_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, t),
				C_SE_NB_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandScaledObject{
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandScaledObject{
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
		cmd    SetpointCommandFloatObject
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandFloatObject{}},
			true},
		{
			"cause not Act and Deact",
			args{
				newConn(nil, t),
				C_SE_NC_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				SetpointCommandFloatObject{}},
			true},
		{
			"C_SE_NC_1",
			args{
				newConn([]byte{byte(C_SE_NC_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}, t),
				C_SE_NC_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandFloatObject{
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				SetpointCommandFloatObject{
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
		cmd        BitsString32CommandObject
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				BitsString32CommandObject{}},
			true},
		{
			"cause not Act and Deact",
			args{
				newConn(nil, t),
				C_BO_NA_1,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				BitsString32CommandObject{}},
			true},
		{
			"C_BO_NA_1",
			args{
				newConn([]byte{byte(C_BO_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x64, 0x00, 0x00, 0x00}, t),
				C_BO_NA_1,
				CauseOfTransmission{Cause: Act},
				0x1234,
				BitsString32CommandObject{
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
				CauseOfTransmission{Cause: Act},
				0x1234,
				BitsString32CommandObject{
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
		name    string
		fields  fields
		want    SingleCommandObject
		wantErr bool
	}{
		{
			"C_SC_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SC_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x05}},
			SingleCommandObject{
				0x567890,
				true,
				QualifierOfCommand{QOCShortPulse, false},
				time.Time{}},
			false},
		{
			"C_SC_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SC_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x04}, tm0CP56Time2aBytes...)},
			SingleCommandObject{
				0x567890,
				false,
				QualifierOfCommand{QOCShortPulse, false},
				tm0},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got, err := this.GetSingleCmd()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetSingleCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		want    DoubleCommandObject
		wantErr bool
	}{
		{
			"C_DC_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_DC_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x05}},
			DoubleCommandObject{
				0x567890,
				DCOOn,
				QualifierOfCommand{QOCShortPulse, false},
				time.Time{},
			}, false},
		{
			"C_DC_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_DC_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...)},
			DoubleCommandObject{
				0x567890,
				DCOOff,
				QualifierOfCommand{QOCShortPulse, false},
				tm0},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got, err := this.GetDoubleCmd()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetDoubleCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		want    StepCommandObject
		wantErr bool
	}{
		{
			"C_RC_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_RC_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x05}},
			StepCommandObject{
				0x567890,
				SCOStepDown,
				QualifierOfCommand{QOCShortPulse, false},
				time.Time{}},
			false},
		{
			"C_RC_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_RC_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x06}, tm0CP56Time2aBytes...)},
			StepCommandObject{
				0x567890,
				SCOStepUP,
				QualifierOfCommand{QOCShortPulse, false},
				tm0},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got, err := this.GetStepCmd()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetStepCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		want    SetpointCommandNormalObject
		wantErr bool
	}{
		{
			"C_SE_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}},
			SetpointCommandNormalObject{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				time.Time{}},
			false},
		{
			"C_SE_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, tm0CP56Time2aBytes...)},
			SetpointCommandNormalObject{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				tm0},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got, err := this.GetSetpointNormalCmd()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetSetpointNormalCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		want    SetpointCommandScaledObject
		wantErr bool
	}{
		{
			"C_SE_NB_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_NB_1},
				[]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}},
			SetpointCommandScaledObject{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				time.Time{}},
			false},
		{
			"C_SE_TB_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_TB_1},
				append([]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x01}, tm0CP56Time2aBytes...)},
			SetpointCommandScaledObject{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				tm0},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got, err := this.GetSetpointCmdScaled()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetSetpointCmdScaled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		want    SetpointCommandFloatObject
		wantErr bool
	}{
		{
			"C_SE_NC_1",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_NC_1},
				[]byte{0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}},
			SetpointCommandFloatObject{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				time.Time{}},
			false},
		{
			"C_SE_TC_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_SE_TC_1},
				append([]byte{0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}, tm0CP56Time2aBytes...)},
			SetpointCommandFloatObject{
				0x567890,
				100,
				QualifierOfSetpointCmd{1, false},
				tm0},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got, err := this.GetSetpointFloatCmd()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetSetpointFloatCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name    string
		fields  fields
		want    BitsString32CommandObject
		wantErr bool
	}{
		{
			"C_BO_NA_1",
			fields{
				ParamsWide,
				Identifier{Type: C_BO_NA_1},
				[]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x00, 0x00}},
			BitsString32CommandObject{
				0x567890,
				100,
				time.Time{}},
			false},
		{
			"C_BO_TA_1 CP56Time2a",
			fields{
				ParamsWide,
				Identifier{Type: C_BO_TA_1},
				append([]byte{0x90, 0x78, 0x56, 0x64, 0x00, 0x00, 0x00}, tm0CP56Time2aBytes...)},
			BitsString32CommandObject{
				0x567890,
				100,
				tm0},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
			}
			got, err := this.GetBitsString32Cmd()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetBitsString32Cmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetBitsString32Cmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
