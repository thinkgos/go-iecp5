package cs104

import (
	"errors"
)

// error defined
var (
	ErrUseClosedConnection = errors.New("Use of closed connection")
	ErrBufferFulled        = errors.New("Buffer is full")
	ErrNotActive           = errors.New("Not active")
)
