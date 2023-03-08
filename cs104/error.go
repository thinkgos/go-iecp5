// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

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
