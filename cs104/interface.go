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

// ClientHandlerInterface TODO:
type ClientHandlerInterface interface {
}
