package cs104

import (
	"crypto/tls"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_openConnection(t *testing.T) {
	type args struct {
		uri     *url.URL
		tlsc    *tls.Config
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    net.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openConnection(tt.args.uri, tt.args.tlsc, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("openConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("openConnection() got = %v, want %v", got, tt.want)
			}
		})
	}
}
