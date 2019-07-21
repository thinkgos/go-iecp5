package cs104

import (
	"reflect"
	"testing"
)

func TestIAPCI_String(t *testing.T) {
	tests := []struct {
		name string
		this iAPCI
		want string
	}{
		{"iFrame", iAPCI{sendSN: 0x02, rcvSN: 0x02}, "I[sendNO: 2, recvNO: 2]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.String(); got != tt.want {
				t.Errorf("APCI.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSAPCI_String(t *testing.T) {
	tests := []struct {
		name string
		this sAPCI
		want string
	}{
		{"sFrame", sAPCI{rcvSN: 123}, "S[recvNO: 123]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.String(); got != tt.want {
				t.Errorf("APCI.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUAPCI_String(t *testing.T) {
	tests := []struct {
		name string
		this uAPCI
		want string
	}{
		{"uFrame", uAPCI{function: uStartDtActive}, "U[function: StartDtActive]"},
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
			got, err := newIFrame(tt.args.sendSN, tt.args.RcvSN, tt.args.asdu)
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
		which byte
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
		want  interface{}
		want1 []byte
	}{
		{
			"iAPCI",
			args{[]byte{startFrame, 0x04, 0x02, 0x00, 0x03, 0x00}},
			iAPCI{sendSN: 0x01, rcvSN: 0x01},
			[]byte{},
		},
		{
			"sAPCI",
			args{[]byte{startFrame, 0x04, 0x01, 0x00, 0x02, 0x00}},
			sAPCI{rcvSN: 0x01},
			[]byte{},
		},
		{
			"uAPCI",
			args{[]byte{startFrame, 0x04, 0x07, 0x00, 0x00, 0x00}},
			uAPCI{uStartDtActive},
			[]byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parse(tt.args.apdu)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
