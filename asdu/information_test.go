package asdu

import (
	"math"
	"testing"
)

func TestNewStepPos(t *testing.T) {
	type args struct {
		value        int
		hasTransient bool
	}
	tests := []struct {
		name string
		args args
		want StepPos
	}{
		{"值-64 处于瞬变状态", args{-64, true}, StepPos(0xc0)},
		{"值-64 未在瞬变状态", args{-64, false}, StepPos(0x40)},
		{"值7 处于瞬变状态", args{7, true}, StepPos(0x87)},
		{"值7 未在瞬变状态", args{7, false}, StepPos(0x07)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStepPos(tt.args.value, tt.args.hasTransient); got != tt.want {
				t.Errorf("NewStepPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStepPos_Pos(t *testing.T) {
	for _, HasTransient := range []bool{false, true} {
		for value := -64; value <= 63; value++ {
			gotValue, gotHasTransient := NewStepPos(value, HasTransient).ToPos()
			if gotValue != value || gotHasTransient != HasTransient {
				t.Errorf("StepPos(%d, %t) ToPos(%d, %t)", value, HasTransient, gotValue, gotHasTransient)
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

func TestCommand(t *testing.T) {
	tests := []struct {
		name     string
		this     Command
		wantQual byte
		wantExec bool
	}{
		{"with selects", Command(0x84), 1, false},
		{"with executes", Command(0x0c), 3, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Qual(); got != tt.wantQual {
				t.Errorf("Command.Qual() = %v, want %v", got, tt.wantQual)
			}
			if got := tt.this.Exec(); got != tt.wantExec {
				t.Errorf("Command.Exec() = %v, want %v", got, tt.wantExec)
			}
		})
	}
}

func TestSetPointCmd(t *testing.T) {
	tests := []struct {
		name     string
		this     SetPointCmd
		wantQual byte
		wantExec bool
	}{
		{"with selects", SetPointCmd(0x87), 7, false},
		{"with executes", SetPointCmd(0x07), 7, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Qual(); got != tt.wantQual {
				t.Errorf("SetPointCmd.Qual() = %v, want %v", got, tt.wantQual)
			}
			if got := tt.this.Exec(); got != tt.wantExec {
				t.Errorf("SetPointCmd.Exec() = %v, want %v", got, tt.wantExec)
			}
		})
	}
}
