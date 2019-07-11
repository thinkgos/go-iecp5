package asdu

import (
	"reflect"
	"testing"
	"time"
)

//
//var goldenASDUs = []struct {
//	u *ASDU
//	s string
//}{
//	{
//		&ASDU{
//			Params:     ParamsWide,
//			Identifier: Identifier{M_SP_NA_1, Percyc, 7, 1001},
//			InfoObj:       []byte{1, 2, 3, 4},
//		},
//		"M_SP_NA_1 percyc 7@1001  197121:0x04",
//	}, {
//		&ASDU{
//			Params:     ParamsNarrow,
//			Identifier: Identifier{M_DP_NA_1, Back, 0, 42},
//			InfoObj:       []byte{1, 2, 3, 4},
//		},
//		"M_DP_NA_1 back @42 1:0x02 3:0x04",
//	}, {
//		&ASDU{
//			Params:     ParamsNarrow,
//			Identifier: Identifier{M_ST_NA_1, Spont, 0, 250},
//			InfoObj:       []byte{1, 2, 3, 4, 5},
//		},
//		"M_ST_NA_1 spont @250 1:0x0203 4:0x05 <EOF>",
//	}, {
//		&ASDU{
//			Params:     ParamsNarrow,
//			Identifier: Identifier{M_ME_NC_1, Init, 0, 12},
//			InfoSeq:    true,
//			InfoObj:       []byte{99, 0, 1, 2, 3, 4, 5},
//		},
//		"M_ME_NC_1 init @12 99:0x0001020304 100:0x05 <EOF>",
//	},
//}
//
//func TestASDUStrings(t *testing.T) {
//	for _, gold := range goldenASDUs {
//		if got := gold.u.String(); got != gold.s {
//			t.Errorf("got %q, want %q", got, gold.s)
//		}
//	}
//}
//
//func TestASDUEncoding(t *testing.T) {
//	for _, gold := range goldenASDUs {
//		if strings.Contains(gold.s, " <EOF>") {
//			continue
//		}
//
//		bytes, err := gold.u.MarshalBinary()
//		if err != nil {
//			t.Error(gold.s, "marshal error:", err)
//			continue
//		}
//
//		u := NewASDU(gold.u.Params, Identifier{})
//		if err = u.UnmarshalBinary(bytes); err != nil {
//			t.Error(gold.s, "unmarshal error:", err)
//			continue
//		}
//
//		if got := u.String(); got != gold.s {
//			t.Errorf("got %q, want %q", got, gold.s)
//		}
//	}
//}

func TestParams_Valid(t *testing.T) {
	tests := []struct {
		name    string
		this    *Params
		wantErr bool
	}{
		{"invalid", &Params{}, true},
		{"ParamsNarrow", ParamsNarrow, false},
		{"ParamsWide", ParamsWide, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Params.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParams_ValidCommonAddr(t *testing.T) {
	type args struct {
		addr CommonAddr
	}
	tests := []struct {
		name    string
		this    *Params
		args    args
		wantErr bool
	}{
		{"common address zero", ParamsNarrow, args{InvalidCommonAddr}, true},
		{"common address size(1),invalid", ParamsNarrow, args{256}, true},
		{"common address size(1),valid", ParamsNarrow, args{255}, false},
		{"common address size(2),valid", ParamsWide, args{65535}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.ValidCommonAddr(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("Params.ValidCommonAddr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParams_IdentifierSize(t *testing.T) {
	tests := []struct {
		name string
		this *Params
		want int
	}{
		{"ParamsNarrow(4)", ParamsNarrow, 4},
		{"ParamsWide(6)", ParamsWide, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.IdentifierSize(); got != tt.want {
				t.Errorf("Params.IdentifierSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_AppendInfoObjAddr(t *testing.T) {
	type args struct {
		addr InfoObjAddr
	}
	tests := []struct {
		name    string
		this    *ASDU
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.AppendInfoObjAddr(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("ASDU.AppendInfoObjAddr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestASDU_ParseInfoObjAddr(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		InfoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    InfoObjAddr
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				InfoObj:    tt.fields.InfoObj,
				bootstrap:  tt.fields.bootstrap,
			}
			got, err := this.ParseInfoObjAddr(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.ParseInfoObjAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ASDU.ParseInfoObjAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_IncVariableNumber(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		InfoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				InfoObj:    tt.fields.InfoObj,
				bootstrap:  tt.fields.bootstrap,
			}
			if err := this.IncVariableNumber(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("ASDU.IncVariableNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestASDU_Reply(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		InfoObj    []byte
		bootstrap  [ASDUSizeMax]byte
	}
	type args struct {
		c    Cause
		addr CommonAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ASDU
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ASDU{
				Params:     tt.fields.Params,
				Identifier: tt.fields.Identifier,
				InfoObj:    tt.fields.InfoObj,
				bootstrap:  tt.fields.bootstrap,
			}
			if got := this.Reply(tt.args.c, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ASDU.Reply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASDU_MarshalBinary(t *testing.T) {
	type fields struct {
		Params     *Params
		Identifier Identifier
		InfoObj    []byte
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		{"unused cause", fields{
			ParamsNarrow,
			Identifier{M_SP_NA_1, VariableStruct{}, CauseOfTransmission{Cause: Unused}, 0, 0x80},
			nil,
		}, nil, true},
		{"invalid cause size", fields{
			&Params{CauseSize: 0, CommonAddrSize: 1, InfoObjAddrSize: 1, InfoObjTimeZone: time.UTC},
			Identifier{M_SP_NA_1, VariableStruct{}, CauseOfTransmission{Cause: Act}, 0, 0x80},
			nil,
		}, nil, true},
		{"cause size(1),but origAddress not equal zero", fields{
			&Params{CauseSize: 1, CommonAddrSize: 1, InfoObjAddrSize: 1, InfoObjTimeZone: time.UTC},
			Identifier{M_SP_NA_1, VariableStruct{}, CauseOfTransmission{Cause: Act}, 1, 0x80},
			nil,
		}, nil, true},
		{"invalid common address", fields{
			&Params{CauseSize: 1, CommonAddrSize: 1, InfoObjAddrSize: 1, InfoObjTimeZone: time.UTC},
			Identifier{M_SP_NA_1, VariableStruct{}, CauseOfTransmission{Cause: Act}, 0, InvalidCommonAddr},
			nil,
		}, nil, true},
		{"invalid common address size", fields{
			&Params{CauseSize: 1, CommonAddrSize: 0, InfoObjAddrSize: 1, InfoObjTimeZone: time.UTC},
			Identifier{M_SP_NA_1, VariableStruct{}, CauseOfTransmission{Cause: Act}, 0, 0x80},
			nil,
		}, nil, true},
		{"common size(1),but common address equal 255", fields{
			&Params{CauseSize: 1, CommonAddrSize: 1, InfoObjAddrSize: 1, InfoObjTimeZone: time.UTC},
			Identifier{M_SP_NA_1, VariableStruct{}, CauseOfTransmission{Cause: Act}, 0, 255},
			nil,
		}, nil, true},
		{"ParamsNarrow", fields{
			ParamsNarrow,
			Identifier{M_SP_NA_1, VariableStruct{Number: 1}, CauseOfTransmission{Cause: Act}, 0, 0x80},
			[]byte{0x00, 0x01, 0x02, 0x03},
		}, []byte{0x01, 0x01, 0x06, 0x80, 0x00, 0x01, 0x02, 0x03}, false},
		{"ParamsNarrow global address", fields{
			ParamsNarrow,
			Identifier{M_SP_NA_1, VariableStruct{Number: 1}, CauseOfTransmission{Cause: Act}, 0, GlobalCommonAddr},
			[]byte{0x00, 0x01, 0x02, 0x03},
		}, []byte{0x01, 0x01, 0x06, 0xff, 0x00, 0x01, 0x02, 0x03}, false},
		{"ParamsWide", fields{
			ParamsWide,
			Identifier{M_SP_NA_1, VariableStruct{Number: 1}, CauseOfTransmission{Cause: Act}, 0, 0x6080},
			[]byte{0x00, 0x01, 0x02, 0x03},
		}, []byte{0x01, 0x01, 0x06, 0x00, 0x80, 0x60, 0x00, 0x01, 0x02, 0x03}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := NewASDU(tt.fields.Params, tt.fields.Identifier)
			this.InfoObj = append(this.InfoObj, tt.fields.InfoObj...)

			gotData, err := this.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("ASDU.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ASDU.MarshalBinary() = % x, want % x", gotData, tt.wantData)
			}
		})
	}
}

func TestASDU_UnmarshalBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		Params  *Params
		args    args
		want    []byte
		wantErr bool
	}{
		{"invalid param", &Params{}, args{}, // 125
			[]byte{}, true},
		{"less than data unit identifier size", ParamsWide, args{[]byte{0x0b, 0x01, 0x06, 0x80}},
			[]byte{}, true},
		{"type id fix size error", ParamsWide, args{[]byte{0x07d, 0x01, 0x06, 0x00, 0x80, 0x60}},
			[]byte{}, true},

		{"ParamsNarrow global address", ParamsNarrow, args{[]byte{0x0b, 0x01, 0x06, 0x80, 0x00, 0x01, 0x02, 0x03}},
			[]byte{0x00, 0x01, 0x02, 0x03}, false},
		{"ParamsNarrow", ParamsNarrow, args{[]byte{0x0b, 0x01, 0x06, 0xff, 0x00, 0x01, 0x02, 0x03}},
			[]byte{0x00, 0x01, 0x02, 0x03}, false},
		{"ParamsWide", ParamsWide, args{[]byte{0x01, 0x01, 0x06, 0x00, 0x80, 0x60, 0x00, 0x01, 0x02, 0x03}},
			[]byte{0x00, 0x01, 0x02, 0x03}, false},
		{"ParamsWide sequence", ParamsWide, args{[]byte{0x01, 0x81, 0x06, 0x00, 0x80, 0x60, 0x00, 0x01, 0x02, 0x03}},
			[]byte{0x00, 0x01, 0x02, 0x03}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := NewEmptyASDU(tt.Params)
			if err := this.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ASDU.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(this.InfoObj, tt.want) {
				t.Errorf("ASDU.UnmarshalBinary() got % x, want % x", this.InfoObj, tt.want)
			}
		})
	}
}
