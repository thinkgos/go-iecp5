package cs104

import (
	"reflect"
	"testing"
	"time"
)

func TestConfig_Valid(t *testing.T) {
	type fields struct {
		ConnectTimeout0   time.Duration
		SendUnAckLimitK   uint16
		SendUnAckTimeout1 time.Duration
		RecvUnAckLimitW   uint16
		RecvUnAckTimeout2 time.Duration
		IdleTimeout3      time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid true", fields{}, false},
		{"valid false", fields{ConnectTimeout0: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &Config{
				ConnectTimeout0:   tt.fields.ConnectTimeout0,
				SendUnAckLimitK:   tt.fields.SendUnAckLimitK,
				SendUnAckTimeout1: tt.fields.SendUnAckTimeout1,
				RecvUnAckLimitW:   tt.fields.RecvUnAckLimitW,
				RecvUnAckTimeout2: tt.fields.RecvUnAckTimeout2,
				IdleTimeout3:      tt.fields.IdleTimeout3,
			}
			if err := sf.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{"default", DefaultConfig()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
