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

package asdu

import (
	"errors"
	"fmt"
)

// error defined
var (
	ErrTypeIdentifier = errors.New("asdu: type identification unknown")
	ErrCauseZero      = errors.New("asdu: cause of transmission 0 is not used")
	ErrCommonAddrZero = errors.New("asdu: common address 0 is not used")

	ErrParam           = errors.New("asdu: system parameter out of range")
	ErrInvalidTimeTag  = errors.New("asdu: invalid time tag")
	ErrOriginAddrFit   = errors.New("asdu: originator address not allowed with cause size 1 system parameter")
	ErrCommonAddrFit   = errors.New("asdu: common address exceeds size system parameter")
	ErrInfoObjAddrFit  = errors.New("asdu: information object address exceeds size system parameter")
	ErrInfoObjIndexFit = errors.New("asdu: information object index not in [1, 127]")
	ErrInroGroupNumFit = errors.New("asdu: interrogation group number exceeds 16")

	ErrLengthOutOfRange = fmt.Errorf("asdu: asdu filed length large than max %d", ASDUSizeMax)
	ErrNotAnyObjInfo    = errors.New("asdu: not any object information")
	ErrTypeIDNotMatch   = errors.New("asdu: type identifier doesn't match call or time tag")

	ErrCmdCause = errors.New("asdu: cause of transmission for command not standard requirement")
)
