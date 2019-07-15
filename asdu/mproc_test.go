package asdu

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func Test_checkValid(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		attrsLen   int
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
			if err := checkValid(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.attrsLen); (err != nil) != tt.wantErr {
				t.Errorf("checkValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_single(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []SinglePointInfo
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
			if err := single(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("single() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingle(t *testing.T) {
	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []SinglePointInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]SinglePointInfo{}},
			true,
		},
		{
			"M_SP_NA_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_SP_NA_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x11, 0x02, 0x00, 0x00, 0x10}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]SinglePointInfo{
					{0x000001, true, QDSBlocked, time.Time{}},
					{0x000002, false, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_SP_NA_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_SP_NA_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x11, 0x10}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]SinglePointInfo{
					{0x000001, true, QDSBlocked, time.Time{}},
					{0x000002, false, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Single(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("Single() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleCP24Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []SinglePointInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]SinglePointInfo{}},
			true,
		},
		{
			"M_SP_TA_1 CP24Time2a  Number = 2",
			args{
				newConn(append(append([]byte{byte(M_SP_TA_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x11}, tm0CP24Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x10}, tm0CP24Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]SinglePointInfo{
					{0x000001, true, QDSBlocked, tm0},
					{0x000002, false, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SingleCP24Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("SingleCP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleCP56Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []SinglePointInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]SinglePointInfo{}},
			true,
		},
		{
			"M_SP_TB_1 CP56Time2a Number = 2",
			args{
				newConn(append(append([]byte{byte(M_SP_TB_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x11}, tm0CP56Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x10}, tm0CP56Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]SinglePointInfo{
					{0x000001, true, QDSBlocked, tm0},
					{0x000002, false, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SingleCP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("SingleCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_double(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []DoublePointInfo
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
			if err := double(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("double() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDouble(t *testing.T) {
	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []DoublePointInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]DoublePointInfo{}},
			true,
		},
		{
			"M_DP_NA_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_DP_NA_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x12, 0x02, 0x00, 0x00, 0x11}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]DoublePointInfo{
					{0x000001, DPIDeterminedOn, QDSBlocked, time.Time{}},
					{0x000002, DPIDeterminedOff, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_DP_NA_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_DP_NA_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x12, 0x11}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]DoublePointInfo{
					{0x000001, DPIDeterminedOn, QDSBlocked, time.Time{}},
					{0x000002, DPIDeterminedOff, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Double(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("Double() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDoubleCP24Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []DoublePointInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]DoublePointInfo{}},
			true,
		},
		{
			"M_DP_TA_1 CP24Time2a  Number = 2",
			args{
				newConn(append(append([]byte{byte(M_DP_TA_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x12}, tm0CP24Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x11}, tm0CP24Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]DoublePointInfo{
					{0x000001, DPIDeterminedOn, QDSBlocked, tm0},
					{0x000002, DPIDeterminedOff, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoubleCP24Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("DoubleCP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDoubleCP56Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []DoublePointInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]DoublePointInfo{}},
			true,
		},
		{
			"M_DP_TB_1 CP56Time2a Number = 2",
			args{
				newConn(append(append([]byte{byte(M_DP_TB_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x12}, tm0CP56Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x11}, tm0CP56Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]DoublePointInfo{
					{0x000001, DPIDeterminedOn, QDSBlocked, tm0},
					{0x000002, DPIDeterminedOff, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoubleCP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("DoubleCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_step(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []StepPositionInfo
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
			if err := step(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("step() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStep(t *testing.T) {
	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []StepPositionInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]StepPositionInfo{}},
			true,
		},
		{
			"M_ST_NA_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_ST_NA_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x10, 0x02, 0x00, 0x00, 0x02, 0x10}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]StepPositionInfo{
					{0x000001, StepPosition{Val: 0x01}, QDSBlocked, time.Time{}},
					{0x000002, StepPosition{Val: 0x02}, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_ST_NA_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_ST_NA_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x10, 0x02, 0x10}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]StepPositionInfo{
					{0x000001, StepPosition{Val: 0x01}, QDSBlocked, time.Time{}},
					{0x000002, StepPosition{Val: 0x02}, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Step(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("Step() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStepCP24Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []StepPositionInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]StepPositionInfo{}},
			true,
		},
		{
			"M_ST_TA_1 CP24Time2a  Number = 2",
			args{
				newConn(append(append([]byte{byte(M_ST_TA_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x10}, tm0CP24Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x10}, tm0CP24Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]StepPositionInfo{
					{0x000001, StepPosition{Val: 0x01}, QDSBlocked, tm0},
					{0x000002, StepPosition{Val: 0x02}, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StepCP24Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("StepCP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStepCP56Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []StepPositionInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]StepPositionInfo{}},
			true,
		},
		{
			"M_SP_TB_1 CP56Time2a Number = 2",
			args{
				newConn(append(append([]byte{byte(M_SP_TB_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x10}, tm0CP56Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x10}, tm0CP56Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]StepPositionInfo{
					{0x000001, StepPosition{Val: 0x01}, QDSBlocked, tm0},
					{0x000002, StepPosition{Val: 0x02}, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StepCP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("StepCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bitString32(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []BitString32Info
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
			if err := bitString32(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("bitString32() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBitString32(t *testing.T) {
	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []BitString32Info
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]BitString32Info{}},
			true,
		},
		{
			"M_BO_NA_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_BO_NA_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]BitString32Info{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_BO_NA_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_BO_NA_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x00, 0x10}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]BitString32Info{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BitString32(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("BitString32() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBitString32CP24Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []BitString32Info
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]BitString32Info{}},
			true,
		},
		{
			"M_BO_TA_1 CP24Time2a  Number = 2",
			args{
				newConn(append(append([]byte{byte(M_BO_TA_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10}, tm0CP24Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10}, tm0CP24Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]BitString32Info{
					{0x000001, 1, QDSBlocked, tm0},
					{0x000002, 2, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BitString32CP24Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("BitString32CP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBitString32CP56Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []BitString32Info
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]BitString32Info{}},
			true,
		},
		{
			"M_BO_TB_1 CP56Time2a Number = 2",
			args{
				newConn(append(append([]byte{byte(M_BO_TB_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10}, tm0CP56Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10}, tm0CP56Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]BitString32Info{
					{0x000001, 1, QDSBlocked, tm0},
					{0x000002, 2, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BitString32CP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("BitString32CP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_measuredValueNormal(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		attrs      []MeasuredValueNormalInfo
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
			if err := measuredValueNormal(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.attrs...); (err != nil) != tt.wantErr {
				t.Errorf("measuredValueNormal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueNormal(t *testing.T) {
	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []MeasuredValueNormalInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueNormalInfo{}},
			true,
		},
		{
			"M_ME_NA_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_ME_NA_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x10}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueNormalInfo{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_ME_NA_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_ME_NA_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x10, 0x02, 0x00, 0x10}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueNormalInfo{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueNormal(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueNormal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueNormalCP24Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []MeasuredValueNormalInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueNormalInfo{}},
			true,
		},
		{
			"M_ME_TA_1 CP24Time2a  Number = 2",
			args{
				newConn(append(append([]byte{byte(M_ME_TA_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x00, 0x10}, tm0CP24Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x00, 0x10}, tm0CP24Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]MeasuredValueNormalInfo{
					{0x000001, 1, QDSBlocked, tm0},
					{0x000002, 2, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueNormalCP24Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueNormalCP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueNormalCP56Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []MeasuredValueNormalInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueNormalInfo{}},
			true,
		},
		{
			"M_ME_TD_1 CP56Time2a Number = 2",
			args{
				newConn(append(append([]byte{byte(M_ME_TD_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x00, 0x10}, tm0CP56Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x00, 0x10}, tm0CP56Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]MeasuredValueNormalInfo{
					{0x000001, 1, QDSBlocked, tm0},
					{0x000002, 2, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueNormalCP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueNormalCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueNormalNoQuality(t *testing.T) {
	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []MeasuredValueNormalInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueNormalInfo{}},
			true,
		},
		{
			"M_ME_ND_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_ME_ND_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x02, 0x00}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueNormalInfo{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_ME_ND_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_ME_ND_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueNormalInfo{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueNormalNoQuality(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueNormalNoQuality() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_measuredValueScaled(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []MeasuredValueScaledInfo
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
			if err := measuredValueScaled(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("measuredValueScaled() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueScaled(t *testing.T) {
	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []MeasuredValueScaledInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueScaledInfo{}},
			true,
		},
		{
			"M_ME_NB_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_ME_NB_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x10}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueScaledInfo{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_ME_NB_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_ME_NB_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, 0x01, 0x00, 0x10, 0x02, 0x00, 0x10}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueScaledInfo{
					{0x000001, 1, QDSBlocked, time.Time{}},
					{0x000002, 2, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueScaled(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueScaled() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueScaledCP24Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []MeasuredValueScaledInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueScaledInfo{}},
			true,
		},
		{
			"M_ME_TB_1 CP24Time2a  Number = 2",
			args{
				newConn(append(append([]byte{byte(M_ME_TB_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x00, 0x10}, tm0CP24Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x00, 0x10}, tm0CP24Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]MeasuredValueScaledInfo{
					{0x000001, 1, QDSBlocked, tm0},
					{0x000002, 2, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueScaledCP24Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueScaledCP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueScaledCP56Time2a(t *testing.T) {
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []MeasuredValueScaledInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueScaledInfo{}},
			true,
		},
		{
			"M_ME_TE_1 CP56Time2a Number = 2",
			args{
				newConn(append(append([]byte{byte(M_ME_TE_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, 0x01, 0x00, 0x10}, tm0CP56Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, 0x02, 0x00, 0x10}, tm0CP56Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]MeasuredValueScaledInfo{
					{0x000001, 1, QDSBlocked, tm0},
					{0x000002, 2, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueScaledCP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueScaledCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_measuredValueFloat(t *testing.T) {
	type args struct {
		c          Connect
		typeID     TypeID
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []MeasuredValueFloatInfo
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
			if err := measuredValueFloat(tt.args.c, tt.args.typeID, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("measuredValueFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueFloat(t *testing.T) {
	bits1 := math.Float32bits(100)
	bits2 := math.Float32bits(101)

	type args struct {
		c          Connect
		isSequence bool
		coa        CauseOfTransmission
		ca         CommonAddr
		infos      []MeasuredValueFloatInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				false,
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueFloatInfo{}},
			true,
		},
		{
			"M_ME_NC_1 seq = false Number = 2",
			args{
				newConn([]byte{byte(M_ME_NC_1), 0x02, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, byte(bits1), byte(bits1 >> 8), byte(bits1 >> 16), byte(bits1 >> 24), 0x10,
					0x02, 0x00, 0x00, byte(bits2), byte(bits2 >> 8), byte(bits2 >> 16), byte(bits2 >> 24), 0x10}, t),
				false,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueFloatInfo{
					{0x000001, 100, QDSBlocked, time.Time{}},
					{0x000002, 101, QDSBlocked, time.Time{}},
				}},
			false,
		},
		{
			"M_ME_NC_1 seq = true Number = 2",
			args{
				newConn([]byte{byte(M_ME_NC_1), 0x82, 0x02, 0x00, 0x34, 0x12,
					0x01, 0x00, 0x00, byte(bits1), byte(bits1 >> 8), byte(bits1 >> 16), byte(bits1 >> 24), 0x10,
					byte(bits2), byte(bits2 >> 8), byte(bits2 >> 16), byte(bits2 >> 24), 0x10}, t),
				true,
				CauseOfTransmission{Cause: Back},
				0x1234,
				[]MeasuredValueFloatInfo{
					{0x000001, 100, QDSBlocked, time.Time{}},
					{0x000002, 101, QDSBlocked, time.Time{}},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueFloat(tt.args.c, tt.args.isSequence, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueFloatCP24Time2a(t *testing.T) {
	bits1 := math.Float32bits(100)
	bits2 := math.Float32bits(101)

	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []MeasuredValueFloatInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueFloatInfo{}},
			true,
		},
		{
			"M_ME_TC_1 seq = false Number = 2",
			args{
				newConn(append(append([]byte{byte(M_ME_TC_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, byte(bits1), byte(bits1 >> 8), byte(bits1 >> 16), byte(bits1 >> 24), 0x10}, tm0CP24Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, byte(bits2), byte(bits2 >> 8), byte(bits2 >> 16), byte(bits2 >> 24), 0x10}, tm0CP24Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]MeasuredValueFloatInfo{
					{0x000001, 100, QDSBlocked, tm0},
					{0x000002, 101, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueFloatCP24Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueFloatCP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasuredValueFloatCP56Time2a(t *testing.T) {
	bits1 := math.Float32bits(100)
	bits2 := math.Float32bits(101)
	type args struct {
		c     Connect
		coa   CauseOfTransmission
		ca    CommonAddr
		infos []MeasuredValueFloatInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid cause",
			args{
				newConn(nil, t),
				CauseOfTransmission{Cause: Unused},
				0x1234,
				[]MeasuredValueFloatInfo{}},
			true,
		},
		{
			"M_ME_TF_1 seq = false Number = 2",
			args{
				newConn(append(append([]byte{byte(M_ME_TF_1), 0x02, 0x03, 0x00, 0x34, 0x12},
					append([]byte{0x01, 0x00, 0x00, byte(bits1), byte(bits1 >> 8), byte(bits1 >> 16), byte(bits1 >> 24), 0x10}, tm0CP56Time2aBytes...)...),
					append([]byte{0x02, 0x00, 0x00, byte(bits2), byte(bits2 >> 8), byte(bits2 >> 16), byte(bits2 >> 24), 0x10}, tm0CP56Time2aBytes...)...), t),
				CauseOfTransmission{Cause: Spont},
				0x1234,
				[]MeasuredValueFloatInfo{
					{0x000001, 100, QDSBlocked, tm0},
					{0x000002, 101, QDSBlocked, tm0},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MeasuredValueFloatCP56Time2a(tt.args.c, tt.args.coa, tt.args.ca, tt.args.infos...); (err != nil) != tt.wantErr {
				t.Errorf("MeasuredValueFloatCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestASDU_GetSinglePoint(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []SinglePointInfo
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
			got, err := this.GetSinglePoint()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetSinglePoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetSinglePoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetDoublePoint(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []DoublePointInfo
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
			got, err := this.GetDoublePoint()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetDoublePoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetDoublePoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetStepPosition(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []StepPositionInfo
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
			got, err := this.GetStepPosition()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetStepPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetStepPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetBitString32(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []BitString32Info
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
			got, err := this.GetBitString32()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetBitString32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetBitString32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetMeasuredValueNormal(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []MeasuredValueNormalInfo
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
			got, err := this.GetMeasuredValueNormal()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetMeasuredValueNormal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetMeasuredValueNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetMeasuredValueScaled(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []MeasuredValueScaledInfo
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
			got, err := this.GetMeasuredValueScaled()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetMeasuredValueScaled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetMeasuredValueScaled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_GetMeasuredValueFloat(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		infoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []MeasuredValueFloatInfo
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
			got, err := this.GetMeasuredValueFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.GetMeasuredValueFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.GetMeasuredValueFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
