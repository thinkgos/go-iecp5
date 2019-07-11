package asdu

import (
	"math"
	"reflect"
	"testing"
)

func TestNewStepPos(t *testing.T) {
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
	v := Normalize(-1 << 15)
	last := v.Float64()
	if last != -1 {
		t.Errorf("%#04x: got %f, want -1", uint16(v), last)
	}

	for v != 1<<15-1 {
		v++
		got := v.Float64()
		if got <= last || got >= 1 {
			t.Errorf("%#04x: got %f (%#04x was %f)", uint16(v), got, uint16(v-1), last)
		}
		last = got
	}
}

func TestNormalize_Float64(t *testing.T) {
	min := float64(-1)
	for v := math.MinInt16; v < math.MaxInt16; v++ {
		got := Normalize(v).Float64()
		if got < min || got >= 1 {
			t.Errorf("%#04x: got %f (%#04x was %f)", uint16(v), got, uint16(v-1), min)
		}
		min = got
	}
}

func TestDecodeQualifierOfCmd(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want QualifierOfCmd
	}{
		{"with selects", args{0x84}, QualifierOfCmd{1, false}},
		{"with executes", args{0x0c}, QualifierOfCmd{3, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQualifierOfCmd(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQualifierOfCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeQualifierOfSetpointCmd(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want QualifierOfSetpointCmd
	}{
		{"with selects", args{0x87}, QualifierOfSetpointCmd{7, false}},
		{"with executes", args{0x07}, QualifierOfSetpointCmd{7, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQualifierOfSetpointCmd(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQualifierOfSetpointCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSinglePoint_Value(t *testing.T) {
	tests := []struct {
		name string
		this SinglePoint
		want byte
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Value(); got != tt.want {
				t.Errorf("DoublePoint.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQualifierOfCmd_Value(t *testing.T) {
	type fields struct {
		CmdQ   CmdQualifier
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
			this := QualifierOfCmd{
				CmdQ:   tt.fields.CmdQ,
				InExec: tt.fields.InExec,
			}
			if got := this.Value(); got != tt.want {
				t.Errorf("QualifierOfCmd.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQualifierOfSetpointCmd_Value(t *testing.T) {
	type fields struct {
		CmdS   CmdSetPoint
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
				CmdS:   tt.fields.CmdS,
				InExec: tt.fields.InExec,
			}
			if got := this.Value(); got != tt.want {
				t.Errorf("QualifierOfSetpointCmd.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
