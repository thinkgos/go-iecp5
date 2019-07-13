package asdu

import (
	"reflect"
	"testing"
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
		{"cause not Act and Deact", args{
			newConn(nil, t),
			C_SC_NA_1, CauseOfTransmission{Cause: Unused}, 0x1234,
			SingleCommandObject{}},
			true},
		{"C_SC_NA_1", args{
			newConn([]byte{45, 0x01, 0x06, 0x00, 0x34, 0x12, 0x90, 0x78, 0x56, 0x05}, t),
			C_SC_NA_1, CauseOfTransmission{Cause: Act}, 0x1234,
			SingleCommandObject{
				0x567890,
				true,
				QualifierOfCommand{QOCShortPulse, false},
				tm0}},
			false},
		{"C_SC_TA_1 CP56Time2a", args{
			newConn(
				append([]byte{58, 0x01, 0x06, 0x00, 0x34, 0x12, 0x90, 0x78, 0x56, 0x05}, tm0CP56Time2aBytes...), t),
			C_SC_TA_1, CauseOfTransmission{Cause: Act}, 0x1234,
			SingleCommandObject{
				0x567890, true,
				QualifierOfCommand{QOCShortPulse, false}, tm0}},
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		cmd    SetpointNormalCommandObject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
		cmd    SetpointScaledCommandObject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
	type args struct {
		c      Connect
		typeID TypeID
		coa    CauseOfTransmission
		ca     CommonAddr
		cmd    SetpointFloatCommandObject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    SingleCommandObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
				bootstrap:  tt.fields.bootstrap,
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
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    DoubleCommandObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
				bootstrap:  tt.fields.bootstrap,
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
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    StepCommandObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
				bootstrap:  tt.fields.bootstrap,
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
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    SetpointNormalCommandObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
				bootstrap:  tt.fields.bootstrap,
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
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    SetpointScaledCommandObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
				bootstrap:  tt.fields.bootstrap,
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
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    SetpointFloatCommandObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
				bootstrap:  tt.fields.bootstrap,
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
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    BitsString32CommandObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				infoObj:    tt.fields.infoObj,
				bootstrap:  tt.fields.bootstrap,
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
