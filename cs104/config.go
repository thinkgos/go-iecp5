package cs104

import (
	"errors"
	"time"
)

const (
	// Port is the IANA registered port number for unsecure connection.
	Port = 2404

	// PortSecure is the IANA registered port number for secure connection.
	PortSecure = 19998
)

const (
	// "t₀" 范围[1, 255]s 默认 30s
	ConnectTimeout0Min = 1 * time.Second
	ConnectTimeout0Max = 255 * time.Second
	// "t₁" 范围[1, 255]s 默认 15s. See IEC 60870-5-104, figure 18.
	SendUnackTimeout1Min = 1 * time.Second
	SendUnackTimeout1Max = 255 * time.Second
	// "t₂" 范围[1, 255]s 默认 10s, See IEC 60870-5-104, figure 10.
	RecvUnackTimeout2Min = 1 * time.Second
	RecvUnackTimeout2Max = 255 * time.Second
	// "t₃" 范围[1 second, 48 hours] 默认 20 s, See IEC 60870-5-104, subclause 5.2.
	IdleTimeout3Min = 1 * time.Second
	IdleTimeout3Max = 48 * time.Hour

	// "k" 范围[1, 32767] 默认 12. See IEC 60870-5-104, subclause 5.5.
	SendUnackLimiKtMin = 1
	SendUnackLimitKMax = 32767
	// "w" 范围 [1， 32767] 默认 8. See IEC 60870-5-104, subclause 5.5.
	RecvUnackLimitWMin = 1
	RecvUnackLimitWMax = 32767
)

// TCPConf defines an IEC 60870-5-104 configuration.
// The default is applied for each unspecified value.
type Config struct {
	// tcp连接建立的最大超时时间
	// "t₀" 范围[1, 255]s，默认 30s.
	ConnectTimeout0 time.Duration

	// I-frames 发送未收到确认的帧数上限， 一旦达到这个数，将停止传输
	// "k" 范围[1, 32767] 默认 12.
	// See IEC 60870-5-104, subclause 5.5.
	SendUnackLimitK uint16

	// 帧接收确认最长超时时间，超过此时间立即关闭连接。
	// "t₁" 范围[1, 255]s 默认 15s.
	// See IEC 60870-5-104, figure 18.
	SendUnackTimeout1 time.Duration

	// 接收端最迟在接收了w次I-frames应用规约数据单元以后发出认可。 w不超过2/3k(2/3 SendUnackMax)
	// "w" 范围 [1， 32767] 默认 8.
	// See IEC 60870-5-104, subclause 5.5.
	RecvUnackLimitW uint16

	// 发送一个接收确认的最大时间，实际上这个框架1秒内发送回复
	// "t₂" 范围[1, 255]s 默认 10s
	// See IEC 60870-5-104, figure 10.
	RecvUnackTimeout2 time.Duration

	// 触发 "TESTFR" 保活的空闲时间值，
	// "t₃" 范围[1 second, 48 hours] 默认 20 s
	// See IEC 60870-5-104, subclause 5.2.
	IdleTimeout3 time.Duration
}

// Valid applies the default (defined by IEC) for each unspecified value.
func (this *Config) Valid() error {
	if this.ConnectTimeout0 == 0 {
		this.ConnectTimeout0 = 30 * time.Second
	} else if this.ConnectTimeout0 < ConnectTimeout0Min || this.ConnectTimeout0 > ConnectTimeout0Max {
		return errors.New(`ConnectTimeout0 "t₀" not in [1, 255]s`)
	}

	if this.SendUnackLimitK == 0 {
		this.SendUnackLimitK = 12
	} else if this.SendUnackLimitK < SendUnackLimiKtMin || this.SendUnackLimitK > SendUnackLimitKMax {
		return errors.New(`SendUnackLimitK "k" not in [1, 32767]`)
	}

	if this.SendUnackTimeout1 == 0 {
		this.SendUnackTimeout1 = 15 * time.Second
	} else if this.SendUnackTimeout1 < SendUnackTimeout1Min || this.SendUnackTimeout1 > SendUnackTimeout1Max {
		return errors.New(`SendUnackTimeout1 "t₁" not in [1, 255]s`)
	}

	if this.RecvUnackLimitW == 0 {
		this.RecvUnackLimitW = 8
	} else if this.RecvUnackLimitW < RecvUnackLimitWMin || this.RecvUnackLimitW > RecvUnackLimitWMax {
		return errors.New(`RecvUnackLimitW "w" not in [1, 32767]`)
	}

	if this.RecvUnackTimeout2 == 0 {
		this.RecvUnackTimeout2 = 10 * time.Second
	} else if this.RecvUnackTimeout2 < RecvUnackTimeout2Min || this.RecvUnackTimeout2 > RecvUnackTimeout2Max {
		return errors.New(`RecvUnackTimeout2 "t₂" not in [1, 255]s`)
	}

	if this.IdleTimeout3 == 0 {
		this.IdleTimeout3 = 20 * time.Second
	} else if this.IdleTimeout3 < IdleTimeout3Min || this.IdleTimeout3 > IdleTimeout3Max {
		return errors.New(`IdleTimeout3 "t₃" not in [1 second, 48 hours]`)
	}

	return nil
}
