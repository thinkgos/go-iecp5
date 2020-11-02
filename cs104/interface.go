// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package cs104

import (
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
)

// ServerHandlerInterface is the interface of server handler
type ServerHandlerInterface interface {
	InterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfInterrogation) error
	CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error
	ReadHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr) error
	ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error
	ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error
	DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error
	ASDUHandler(asdu.Connect, *asdu.ASDU) error
}

// ClientHandlerInterface  is the interface of client handler
type ClientHandlerInterface interface {
	InterrogationHandler(asdu.Connect, *asdu.ASDU) error
	CounterInterrogationHandler(asdu.Connect, *asdu.ASDU) error
	ReadHandler(asdu.Connect, *asdu.ASDU) error
	TestCommandHandler(asdu.Connect, *asdu.ASDU) error
	ClockSyncHandler(asdu.Connect, *asdu.ASDU) error
	ResetProcessHandler(asdu.Connect, *asdu.ASDU) error
	DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU) error
	ASDUHandler(asdu.Connect, *asdu.ASDU) error
}
