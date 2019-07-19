package cs104

import (
	"errors"
)

var (
	ErrUseClosedConnection = errors.New("use closed connection")
	ErrBufferFulled        = errors.New("buffer is full")

//errSeqNo            = errors.New("fatal incomming sequence number disruption")
//errStartDtAckExpire = errors.New("fatal STARTDT acknowledge timeout t₁")
//errStopDtAckExpire  = errors.New("fatal STOPDT acknowledge timeout t₁")
//errTestFrAckExpire  = errors.New("fatal TESTFR acknowledge timeout t₁")
)
