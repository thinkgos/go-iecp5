package cs104

import (
	"reflect"
	"testing"
)

func TestAPCI_parse(t *testing.T) {
	tests := []struct {
		name  string
		this  APCI
		want  interface{}
		want1 string
	}{
		{"iFrame", APCI{ctr1: 0x02, ctr3: 0x02}, iAPCI{sendSN: 0x01, rcvSN: 0x01}, iFrame},
		{"sFrame", APCI{ctr1: 0x01, ctr3: 0x02}, sAPCI{rcvSN: 0x01}, sFrame},
		{"uFrame", APCI{ctr1: 0x07}, uAPCI{function: uStartDtActive}, uFrame},
		{"uFrame", APCI{ctr1: 0x0b}, uAPCI{function: uStartDtConfirm}, uFrame},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.this.parse()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APCI.parse() got = % x, want % x", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("APCI.parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAPCI_String(t *testing.T) {
	tests := []struct {
		name string
		this APCI
		want string
	}{
		{"iFrame", APCI{ctr1: 0x02, ctr3: 0x02}, "I[send=0001, recv=0001]"},
		{"sFrame", APCI{ctr1: 0x01, ctr3: 0x02}, "S[recv=0001]"},
		{"uFrame", APCI{ctr1: 0x07}, "U[0004]"},
		{"uFrame", APCI{ctr1: 0x0b}, "U[0008]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.String(); got != tt.want {
				t.Errorf("APCI.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newIFrame(t *testing.T) {
	type args struct {
		asdu   []byte
		sendSN uint16
		RcvSN  uint16
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"asdu out of range",
			args{asdu: make([]byte, 250)},
			nil,
			true,
		},
		{
			"asdu right",
			args{[]byte{0x01, 0x02}, 0x06, 0x07},
			[]byte{startFrame, 0x06, 0x0c, 0x00, 0x0e, 0x00, 0x01, 0x02},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newIFrame(tt.args.asdu, tt.args.sendSN, tt.args.RcvSN)
			if (err != nil) != tt.wantErr {
				t.Errorf("newIFrame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newIFrame() = % x, want % x", got, tt.want)
			}
		})
	}
}

func Test_newSFrame(t *testing.T) {
	type args struct {
		RcvSN uint16
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"", args{0x06}, []byte{startFrame, 0x04, 0x01, 0x00, 0x0c, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSFrame(tt.args.RcvSN); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSFrame() = % x, want % x", got, tt.want)
			}
		})
	}
}

func Test_newUFrame(t *testing.T) {
	type args struct {
		which int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"", args{uStopDtActive}, []byte{startFrame, 0x04, 0x13, 0x00, 0x00, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUFrame(tt.args.which); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUFrame() = % x, want % x", got, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	type args struct {
		apdu []byte
	}
	tests := []struct {
		name  string
		args  args
		want  APCI
		want1 []byte
	}{
		{
			"",
			args{[]byte{startFrame, 0x04, 0x13, 0x00, 0x00, 0x00}},
			APCI{startFrame, 0x04, 0x13, 0x00, 0x00, 0x00},
			[]byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parse(tt.args.apdu)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %#v, want %#v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parse() got1 = % x, want % x", got1, tt.want1)
			}
		})
	}
}
