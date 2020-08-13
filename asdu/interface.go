// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

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
