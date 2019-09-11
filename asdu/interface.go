package asdu

import (
	"net"
)

// Connect interface
type Connect interface {
	Params() *Params
	Send(a *ASDU) error
	UnderlyingConn() net.Conn
}
