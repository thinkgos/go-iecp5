// Copyright [2020] [thinkgos]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
