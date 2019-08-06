package asdu

import (
	"math"
	"reflect"
	"testing"
)

func TestParameterNormal(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		p   ParameterNormalInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"cause not act",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				ParameterNormalInfo{
					0x567890,
					0x3344,
					QualifierOfParameterMV{}}},
			true,
		},
		{
			"P_ME_NA_1",
			args{
				newConn([]byte{byte(P_ME_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x44, 0x33, 0x01}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				ParameterNormalInfo{
					0x567890,
					0x3344,
					QualifierOfParameterMV{
						QPMThreshold,
						false,
						false}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParameterNormal(tt.args.c, tt.args.coa, tt.args.ca, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ParameterNormal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParameterScaled(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		p   ParameterScaledInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"cause not act",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				ParameterScaledInfo{
					0x567890,
					0x3344,
					QualifierOfParameterMV{}}},
			true,
		},
		{
			"P_ME_NB_1",
			args{
				newConn([]byte{byte(P_ME_NB_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x44, 0x33, 0x01}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				ParameterScaledInfo{
					0x567890,
					0x3344,
					QualifierOfParameterMV{
						QPMThreshold,
						false,
						false}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParameterScaled(tt.args.c, tt.args.coa, tt.args.ca, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ParameterScaled() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParameterFloat(t *testing.T) {
	bits := math.Float32bits(100)

	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		p   ParameterFloatInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"cause not act",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				ParameterFloatInfo{
					0x567890,
					100,
					QualifierOfParameterMV{}}},
			true,
		},
		{
			"P_ME_NC_1",
			args{
				newConn([]byte{byte(P_ME_NC_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				ParameterFloatInfo{
					0x567890,
					100,
					QualifierOfParameterMV{
						QPMThreshold,
						false,
						false}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParameterFloat(tt.args.c, tt.args.coa, tt.args.ca, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ParameterFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParameterActivation(t *testing.T) {
	type args struct {
		c   Connect
		coa CauseOfTransmission
		ca  CommonAddr
		p   ParameterActivationInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"cause not act and deact",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				ParameterActivationInfo{
					0x567890,
					QPAUnused}},
			true,
		},
		{
			"P_AC_NA_1",
			args{
				newConn([]byte{byte(P_AC_NA_1), 0x01, 0x06, 0x00, 0x34, 0x12,
					0x90, 0x78, 0x56, 0x00}, t),
				CauseOfTransmission{Cause: Activation},
				0x1234,
				ParameterActivationInfo{
					0x567890,
					QPAUnused}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParameterActivation(tt.args.c, tt.args.coa, tt.args.ca, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ParameterActivation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestASDU_GetParameterNormal(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   ParameterNormalInfo
	}{
		{
			"P_ME_NA_1",
			fields{
				ParamsWide,
				[]byte{0x90, 0x78, 0x56, 0x44, 0x33, 0x01}},
			ParameterNormalInfo{
				0x567890,
				0x3344,
				QualifierOfParameterMV{
					QPMThreshold,
					false,
					false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			if got := this.GetParameterNormal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetParameterNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetParameterScaled(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   ParameterScaledInfo
	}{
		{
			"P_ME_NB_1",
			fields{
				ParamsWide,
				[]byte{0x90, 0x78, 0x56, 0x44, 0x33, 0x01}},
			ParameterScaledInfo{
				0x567890,
				0x3344,
				QualifierOfParameterMV{
					QPMThreshold,
					false,
					false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			if got := this.GetParameterScaled(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetParameterScaled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetParameterFloat(t *testing.T) {
	bits := math.Float32bits(100)

	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   ParameterFloatInfo
	}{
		{
			"P_ME_NC_1",
			fields{
				ParamsWide,
				[]byte{0x90, 0x78, 0x56, byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24), 0x01}},
			ParameterFloatInfo{
				0x567890,
				100,
				QualifierOfParameterMV{
					QPMThreshold,
					false,
					false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			if got := this.GetParameterFloat(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetParameterFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetParameterActivation(t *testing.T) {
	type fields struct {
		Params  *Params
		infoObj []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   ParameterActivationInfo
	}{
		{
			"P_AC_NA_1",
			fields{
				ParamsWide,
				[]byte{0x90, 0x78, 0x56, 0x00}},
			ParameterActivationInfo{
				0x567890,
				QPAUnused},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:  tt.fields.Params,
				infoObj: tt.fields.infoObj,
			}
			if got := this.GetParameterActivation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetParameterActivation() = %v, want %v", got, tt.want)
			}
		})
	}
}
