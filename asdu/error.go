package asdu

import (
	"errors"
	"fmt"
)

var (
	ErrTypeIdentifier = errors.New("asdu: type identification unknown")
	ErrCauseZero      = errors.New("asdu: cause of transmission 0 is not used")
	ErrCommonAddrZero = errors.New("asdu: common address 0 is not used")

	ErrParam           = errors.New("asdu: system parameter out of range")
	ErrOriginAddrFit   = errors.New("asdu: originator address not allowed with cause size 1 system parameter")
	ErrCommonAddrFit   = errors.New("asdu: common address exceeds size system parameter")
	ErrInfoObjAddrFit  = errors.New("asdu: information object address exceeds size system parameter")
	ErrInfoObjIndexFit = errors.New("asdu: information object index not in [1, 127]")
	ErrInroGroupNumFit = errors.New("asdu: interrogation group number exceeds 16")

	ErrLengthOutOfRange = fmt.Errorf("asdu: asdu filed length large than max %d", ASDUSizeMax)
	ErrNotAnyObjInfo    = errors.New("asdu: not any object information")
	ErrTypeIDNotMatch   = errors.New("asdu: type identifier doesn't match call or time tag")
)
