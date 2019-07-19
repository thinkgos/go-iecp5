package asdu

import (
	"net"
)

type Connect interface {
	Params() *Params
	Send(a *ASDU) error
	UnderlyingConn() net.Conn
}
