package asdu

import (
	"reflect"
	"testing"
)

func TestGetInfoObjSize(t *testing.T) {
	type args struct {
		id TypeID
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"defined", args{F_DR_TA_1}, 13, false},
		{"no defined", args{F_SG_NA_1}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInfoObjSize(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfoObjSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInfoObjSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeID_String(t *testing.T) {
	tests := []struct {
		name string
		this TypeID
		want string
	}{
		{"M_SP_NA_1", M_SP_NA_1, "TID<M_SP_NA_1>"},
		{"M_SP_TB_1", M_SP_TB_1, "TID<M_SP_TB_1>"},
		{"C_SC_NA_1", C_SC_NA_1, "TID<C_SC_NA_1>"},
		{"C_SC_TA_1", C_SC_TA_1, "TID<C_SC_TA_1>"},
		{"M_EI_NA_1", M_EI_NA_1, "TID<M_EI_NA_1>"},
		{"S_CH_NA_1", S_CH_NA_1, "TID<S_CH_NA_1>"},
		{"S_US_NA_1", S_US_NA_1, "TID<S_US_NA_1>"},
		{"C_IC_NA_1", C_IC_NA_1, "TID<C_IC_NA_1>"},
		{"P_ME_NA_1", P_ME_NA_1, "TID<P_ME_NA_1>"},
		{"F_FR_NA_1", F_FR_NA_1, "TID<F_FR_NA_1>"},
		{"no defined", 0, "TID<0>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.String(); got != tt.want {
				t.Errorf("TypeID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVariableStruct(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want VariableStruct
	}{
		{"no sequence", args{0x0a}, VariableStruct{Number: 0x0a}},
		{"with sequence", args{0x8a}, VariableStruct{Number: 0x0a, IsSequence: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseVariableStruct(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVariableStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariableStruct_Value(t *testing.T) {
	tests := []struct {
		name string
		this VariableStruct
		want byte
	}{
		{"no sequence", VariableStruct{Number: 0x0a}, 0x0a},
		{"with sequence", VariableStruct{Number: 0x0a, IsSequence: true}, 0x8a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Value(); got != tt.want {
				t.Errorf("VariableStruct.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariableStruct_String(t *testing.T) {
	tests := []struct {
		name string
		this VariableStruct
		want string
	}{
		{"no sequence", VariableStruct{Number: 100}, "VSQ<100>"},
		{"with sequence", VariableStruct{Number: 100, IsSequence: true}, "VSQ<sq,100>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.String(); got != tt.want {
				t.Errorf("VariableStruct.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCauseOfTransmission(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want CauseOfTransmission
	}{
		{"no test and neg", args{0x01}, CauseOfTransmission{Cause: Periodic}},
		{"with test", args{0x81}, CauseOfTransmission{Cause: Periodic, IsTest: true}},
		{"with neg", args{0x41}, CauseOfTransmission{Cause: Periodic, IsNegative: true}},
		{"with test and neg", args{0xc1}, CauseOfTransmission{Cause: Periodic, IsTest: true, IsNegative: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCauseOfTransmission(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCauseOfTransmission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCauseOfTransmission_Value(t *testing.T) {
	tests := []struct {
		name string
		this CauseOfTransmission
		want byte
	}{
		{"no test and neg", CauseOfTransmission{Cause: Periodic}, 0x01},
		{"with test", CauseOfTransmission{Cause: Periodic, IsTest: true}, 0x81},
		{"with neg", CauseOfTransmission{Cause: Periodic, IsNegative: true}, 0x41},
		{"with test and neg", CauseOfTransmission{Cause: Periodic, IsTest: true, IsNegative: true}, 0xc1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Value(); got != tt.want {
				t.Errorf("CauseOfTransmission.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCauseOfTransmission_String(t *testing.T) {
	tests := []struct {
		name string
		this CauseOfTransmission
		want string
	}{
		{"no test and neg", CauseOfTransmission{Cause: Periodic}, "COT<Periodic>"},
		{"with test", CauseOfTransmission{Cause: Periodic, IsTest: true}, "COT<Periodic,test>"},
		{"with neg", CauseOfTransmission{Cause: Periodic, IsNegative: true}, "COT<Periodic,neg>"},
		{"with test and neg", CauseOfTransmission{Cause: Periodic, IsTest: true, IsNegative: true}, "COT<Periodic,neg,test>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.String(); got != tt.want {
				t.Errorf("CauseOfTransmission.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
