package asdu

import "errors"

var (
	errTypeIdentifier  = errors.New("asdu: type identification unknown")
	errCauseZero       = errors.New("asdu: cause of transmission 0 is not used")
	errCommonoAddrZero = errors.New("asdu: common address 0 is not used")

	errParam           = errors.New("asdu: fixed system parameter out of range")
	errOriginAddrFit   = errors.New("asdu: originator address not allowed with cause size 1 system parameter")
	errCommonAddrFit   = errors.New("asdu: common address exceeds size system parameter")
	errInfoObjAddrFit  = errors.New("asdu: information object address exceeds size system parameter")
	errInfoObjIndexFit = errors.New("asdu: information object index not in [1, 127]")
	errInroGroupNumFit = errors.New("asdu: interrogation group number exceeds 16")
)
