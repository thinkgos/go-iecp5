package asdu

import (
	"reflect"
	"testing"
	"time"
)

var tm0 = time.Date(2019, 6, 5, 4, 3, 0, 513000000, time.UTC)
var tm0Serial = []byte{0x01, 0x02, 0x03, 0x04, 0x65, 0x06, 0x13}
var tm1 = time.Date(2019, 12, 15, 14, 13, 3, 83000000, time.UTC)
var tm1Serial = []byte{0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x0c, 0x13}

func TestCP56Time2a(t *testing.T) {
	type args struct {
		t   time.Time
		loc *time.Location
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"20190605", args{tm0, nil}, tm0Serial},
		{"20191215", args{tm1, time.UTC}, tm1Serial},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CP56Time2a(tt.args.t, tt.args.loc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CP56Time2a() = % x, want % x", got, tt.want)
			}
		})
	}
}

func TestParseCP56Time2a(t *testing.T) {
	type args struct {
		bytes []byte
		loc   *time.Location
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"invalid flag", args{
			[]byte{0x01, 0x02, 0x83, 0x04, 0x65, 0x06, 0x13},
			nil},
			time.Time{},
			true},
		{"20190605", args{tm0Serial, nil}, tm0, false},
		{"20191215", args{tm1Serial, time.UTC}, tm1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCP56Time2a(tt.args.bytes, tt.args.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCP56Time2a() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCP56Time2a() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCP24Time2a(t *testing.T) {
	type args struct {
		t   time.Time
		loc *time.Location
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"3 Minutes 513 Milliseconds", args{tm0, nil}, tm0Serial[:3]},
		{"13 Minutes 3083 Milliseconds", args{tm1, time.UTC}, tm1Serial[:3]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CP24Time2a(tt.args.t, tt.args.loc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CP24Time2a() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCP24Time2a(t *testing.T) {
	type args struct {
		bytes []byte
		loc   *time.Location
	}
	tests := []struct {
		name     string
		args     args
		wantMsec int
		wantMin  int
		wantErr  bool
	}{
		{"invalid flag", args{[]byte{0x01, 0x02, 0x83}, nil},
			0, 0, true},
		{"3 Minutes 513 Milliseconds", args{tm0Serial[:3], nil},
			513, 3, false},
		{"13 Minutes 3083 Milliseconds", args{tm1Serial[:3], time.UTC},
			3083, 13, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCP24Time2a(tt.args.bytes, tt.args.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCP24Time2a() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			msec := (got.Nanosecond()/int(time.Millisecond) + got.Second()*1000)
			if msec != tt.wantMsec {
				t.Errorf("ParseCP24Time2a() go Millisecond = %v, want %v", msec, tt.wantMsec)
			}
			if got.Minute() != tt.wantMin {
				t.Errorf("ParseCP24Time2a() got Minute = %v, want %v", got.Minute(), tt.wantMin)
			}
		})
	}
}

func TestCP16Time2a(t *testing.T) {
	type args struct {
		msec uint16
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"513 Milliseconds", args{513}, []byte{0x01, 0x02}},
		{"3083 Milliseconds", args{3083}, []byte{0x0b, 0x0c}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CP16Time2a(tt.args.msec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CP16Time2a() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCP16Time2a(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{"513 Milliseconds", args{[]byte{0x01, 0x02}}, 513},
		{"3083 Milliseconds", args{[]byte{0x0b, 0x0c}}, 3083},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCP16Time2a(tt.args.b); got != tt.want {
				t.Errorf("ParseCP16Time2a() = %v, want %v", got, tt.want)
			}
		})
	}
}
