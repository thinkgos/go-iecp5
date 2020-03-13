package asdu

import (
	"math"
	"reflect"
	"testing"
)

func TestSinglePoint_Value(t *testing.T) {
	tests := []struct {
		name string
		this SinglePoint
		want bool
	}{
		{"off", SPIOff, false},
		{"on", SPIOn, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Value(); got != tt.want {
				t.Errorf("SinglePoint.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoublePoint_Value(t *testing.T) {
	tests := []struct {
		name string
		this DoublePoint
		want byte
	}{
		{"IndeterminateOrIntermediate", DPIIndeterminateOrIntermediate, 0x00},
		{"DeterminedOff", DPIDeterminedOff, 0x01},
		{"DeterminedOn", DPIDeterminedOn, 0x02},
		{"Indeterminate", DPIIndeterminate, 0x03},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Value(); got != tt.want {
				t.Errorf("DoublePoint.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStepPosition(t *testing.T) {
	type args struct {
		value byte
	}
	tests := []struct {
		name string
		args args
		want StepPosition
	}{
		{"值0xc0 处于瞬变状态", args{0xc0}, StepPosition{-64, true}},
		{"值0x40 未在瞬变状态", args{0x40}, StepPosition{-64, false}},
		{"值0x87 处于瞬变状态", args{0x87}, StepPosition{0x07, true}},
		{"值0x07 未在瞬变状态", args{0x07}, StepPosition{0x07, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseStepPosition(tt.args.value); got != tt.want {
				t.Errorf("NewStepPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStepPosition_Value(t *testing.T) {
	for _, HasTransient := range []bool{false, true} {
		for value := -64; value <= 63; value++ {
			got := ParseStepPosition(StepPosition{value, HasTransient}.Value())
			if got.Val != value || got.HasTransient != HasTransient {
				t.Errorf("ParseStepPosition(StepPosition(%d, %t).Value()) = StepPosition(%d, %t)", value, HasTransient, got.Val, got.HasTransient)
			}
		}
	}
}

// TestNormal tests the full value range.
func TestNormal(t *testing.T) {
	v := NormalizedMeasurement(-1 << 15)
	last := v.NormalizedValue()
	if last != -1 {
		t.Errorf("%#04x: got %f, want -1", uint16(v), last)
	}

	for v != 1<<15-1 {
		v++
		got := v.NormalizedValue()
		if got <= last || got >= 1 {
			t.Errorf("%#04x: got %f (%#04x was %f)", uint16(v), got, uint16(v-1), last)
		}
		last = got
	}
}

func TestNormalize_Float64(t *testing.T) {
	min := float64(-1)
	for v := math.MinInt16; v < math.MaxInt16; v++ {
		got := NormalizedMeasurement(v).NormalizedValue()
		if got < min || got >= 1 {
			t.Errorf("%#04x: got %f (%#04x was %f)", uint16(v), got, uint16(v-1), min)
		}
		min = got
	}
}

func TestParseQualifierOfCmd(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want QualifierOfCommand
	}{
		{"with selects", args{0x84}, QualifierOfCommand{1, true}},
		{"with executes", args{0x0c}, QualifierOfCommand{3, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQualifierOfCommand(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQualifierOfCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseQualifierOfSetpointCmd(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want QualifierOfSetpointCmd
	}{
		{"with selects", args{0x87}, QualifierOfSetpointCmd{7, true}},
		{"with executes", args{0x07}, QualifierOfSetpointCmd{7, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQualifierOfSetpointCmd(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQualifierOfSetpointCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQualifierOfCmd_Value(t *testing.T) {
	type fields struct {
		CmdQ   QOCQual
		InExec bool
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := QualifierOfCommand{
				Qual:     tt.fields.CmdQ,
				InSelect: tt.fields.InExec,
			}
			if got := this.Value(); got != tt.want {
				t.Errorf("QualifierOfCommand.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQualifierOfSetpointCmd_Value(t *testing.T) {
	type fields struct {
		CmdS   QOSQual
		InExec bool
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := QualifierOfSetpointCmd{
				Qual:     tt.fields.CmdS,
				InSelect: tt.fields.InExec,
			}
			if got := this.Value(); got != tt.want {
				t.Errorf("QualifierOfSetpointCmd.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseQualifierOfParam(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want QualifierOfParameterMV
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQualifierOfParamMV(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQualifierOfParamMV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQualifierOfParam_Value(t *testing.T) {
	type fields struct {
		ParamQ        QPMCategory
		IsChange      bool
		IsInOperation bool
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := QualifierOfParameterMV{
				Category:      tt.fields.ParamQ,
				IsChange:      tt.fields.IsChange,
				IsInOperation: tt.fields.IsInOperation,
			}
			if got := this.Value(); got != tt.want {
				t.Errorf("QualifierOfParameterMV.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
