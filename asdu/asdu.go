// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

// Package asdu provides the OSI presentation layer.
package asdu

import (
	"fmt"
	"io"
	"math/bits"
	"time"
)

// ASDUSizeMax asdu max size
const (
	ASDUSizeMax = 249
)

// ASDU format
//       | data unit identification | information object <1..n> |
//
//       | <------------  data unit identification ------------>|
//       | typeID | variable struct | cause  |  common address  |
// bytes |    1   |      1          | [1,2]  |      [1,2]       |
//       | <------------  information object ------------------>|
//       | object address | element set  |  object time scale   |
// bytes |     [1,2,3]    |              |                      |

var (
	// ParamsNarrow is the smallest configuration.
	ParamsNarrow = &Params{CauseSize: 1, CommonAddrSize: 1, InfoObjAddrSize: 1, InfoObjTimeZone: time.UTC}
	// ParamsWide is the largest configuration.
	ParamsWide = &Params{CauseSize: 2, CommonAddrSize: 2, InfoObjAddrSize: 3, InfoObjTimeZone: time.UTC}
)

// Params 定义了ASDU相关特定参数
// See companion standard 101, subclass 7.1.
type Params struct {
	// cause of transmission, 传输原因字节数
	// The standard requires "b" in [1, 2].
	// Value 2 includes/activates the originator address.
	CauseSize int
	// Originator Address [1, 255] or 0 for the default.
	// The applicability is controlled by Params.CauseSize.
	OrigAddress OriginAddr
	// size of ASDU common address， ASDU 公共地址字节数
	// 应用服务数据单元公共地址的八位位组数目,公共地址是站地址
	// The standard requires "a" in [1, 2].
	CommonAddrSize int

	// size of ASDU information object address. 信息对象地址字节数
	// The standard requires "c" in [1, 3].
	InfoObjAddrSize int

	// InfoObjTimeZone controls the time tag interpretation.
	// The standard fails to mention this one.
	InfoObjTimeZone *time.Location
}

// Valid returns the validation result of params.
func (sf Params) Valid() error {
	if (sf.CauseSize < 1 || sf.CauseSize > 2) ||
		(sf.CommonAddrSize < 1 || sf.CommonAddrSize > 2) ||
		(sf.InfoObjAddrSize < 1 || sf.InfoObjAddrSize > 3) ||
		(sf.InfoObjTimeZone == nil) {
		return ErrParam
	}
	return nil
}

// ValidCommonAddr returns the validation result of a station common address.
func (sf Params) ValidCommonAddr(addr CommonAddr) error {
	if addr == InvalidCommonAddr {
		return ErrCommonAddrZero
	}
	if bits.Len(uint(addr)) > sf.CommonAddrSize*8 {
		return ErrCommonAddrFit
	}
	return nil
}

// IdentifierSize return the application service data unit identifies size
func (sf Params) IdentifierSize() int {
	return 2 + int(sf.CauseSize) + int(sf.CommonAddrSize)
}

// Identifier the application service data unit identifies.
type Identifier struct {
	// type identification, information content
	Type TypeID
	// Variable is variable structure qualifier
	Variable VariableStruct
	// cause of transmission submission category
	Coa CauseOfTransmission
	// Originator Address [1, 255] or 0 for the default.
	// The applicability is controlled by Params.CauseSize.
	OrigAddr OriginAddr
	// CommonAddr is a station address. Zero is not used.
	// The width is controlled by Params.CommonAddrSize.
	// See companion standard 101, subclass 7.2.4.
	CommonAddr CommonAddr // station address 公共地址是站地址
}

// String 返回数据单元标识符的信息,例： "TypeID Cause OrigAddr@CommonAddr"
func (id Identifier) String() string {
	if id.OrigAddr == 0 {
		return fmt.Sprintf("%s %s @%d", id.Type, id.Coa, id.CommonAddr)
	}
	return fmt.Sprintf("%s %s %d@%d ", id.Type, id.Coa, id.OrigAddr, id.CommonAddr)
}

// ASDU (Application Service Data Unit) is an application message.
type ASDU struct {
	*Params
	Identifier
	infoObj   []byte            // information object serial
	bootstrap [ASDUSizeMax]byte // prevents Info malloc
}

// NewEmptyASDU new empty asdu with special params
func NewEmptyASDU(p *Params) *ASDU {
	a := &ASDU{Params: p}
	lenDUI := a.IdentifierSize()
	a.infoObj = a.bootstrap[lenDUI:lenDUI]
	return a
}

// NewASDU new asdu with special params and identifier
func NewASDU(p *Params, identifier Identifier) *ASDU {
	a := NewEmptyASDU(p)
	a.Identifier = identifier
	return a
}

// Clone deep clone asdu
func (sf *ASDU) Clone() *ASDU {
	r := NewASDU(sf.Params, sf.Identifier)
	r.infoObj = append(r.infoObj, sf.infoObj...)
	return r
}

// SetVariableNumber See companion standard 101, subclass 7.2.2.
func (sf *ASDU) SetVariableNumber(n int) error {
	if n >= 128 {
		return ErrInfoObjIndexFit
	}
	sf.Variable.Number = byte(n)
	return nil
}

// Respond returns a new "responding" ASDU which addresses "initiating" u.
//func (u *ASDU) Respond(t TypeID, c Cause) *ASDU {
//	return NewASDU(u.Params, Identifier{
//		CommonAddr: u.CommonAddr,
//		OrigAddr:   u.OrigAddr,
//		Type:       t,
//		Cause:      c | u.Cause&TestFlag,
//	})
//}

// Reply returns a new "responding" ASDU which addresses "initiating" addr with a copy of Info.
func (sf *ASDU) Reply(c Cause, addr CommonAddr) *ASDU {
	sf.CommonAddr = addr
	r := NewASDU(sf.Params, sf.Identifier)
	r.Coa.Cause = c
	r.infoObj = append(r.infoObj, sf.infoObj...)
	return r
}

// SendReplyMirror send a reply of the mirror request but cause different
func (sf *ASDU) SendReplyMirror(c Connect, cause Cause) error {
	r := NewASDU(sf.Params, sf.Identifier)
	r.Coa.Cause = cause
	r.infoObj = append(r.infoObj, sf.infoObj...)
	return c.Send(r)
}

//// String returns a full description.
//func (u *ASDU) String() string {
//	dataSize, err := GetInfoObjSize(u.Type)
//	if err != nil {
//		if !u.InfoSeq {
//			return fmt.Sprintf("%s: %#x", u.Identifier, u.infoObj)
//		}
//		return fmt.Sprintf("%s seq: %#x", u.Identifier, u.infoObj)
//	}
//
//	end := len(u.infoObj)
//	addrSize := u.InfoObjAddrSize
//	if end < addrSize {
//		if !u.InfoSeq {
//			return fmt.Sprintf("%s: %#x <EOF>", u.Identifier, u.infoObj)
//		}
//		return fmt.Sprintf("%s seq: %#x <EOF>", u.Identifier, u.infoObj)
//	}
//	addr := u.ParseInfoObjAddr(u.infoObj)
//
//	buf := bytes.NewBufferString(u.Identifier.String())
//
//	for i := addrSize; ; {
//		start := i
//		i += dataSize
//		if i > end {
//			fmt.Fprintf(buf, " %d:%#x <EOF>", addr, u.infoObj[start:])
//			break
//		}
//		fmt.Fprintf(buf, " %d:%#x", addr, u.infoObj[start:i])
//		if i == end {
//			break
//		}
//
//		if u.InfoSeq {
//			addr++
//		} else {
//			start = i
//			i += addrSize
//			if i > end {
//				fmt.Fprintf(buf, " %#x <EOF>", u.infoObj[start:i])
//				break
//			}
//			addr = u.ParseInfoObjAddr(u.infoObj[start:])
//		}
//	}
//
//	return buf.String()
//}

// MarshalBinary honors the encoding.BinaryMarshaler interface.
func (sf *ASDU) MarshalBinary() (data []byte, err error) {
	switch {
	case sf.Coa.Cause == Unused:
		return nil, ErrCauseZero
	case !(sf.CauseSize == 1 || sf.CauseSize == 2):
		return nil, ErrParam
	case sf.CauseSize == 1 && sf.OrigAddr != 0:
		return nil, ErrOriginAddrFit
	case sf.CommonAddr == InvalidCommonAddr:
		return nil, ErrCommonAddrZero
	case !(sf.CommonAddrSize == 1 || sf.CommonAddrSize == 2):
		return nil, ErrParam
	case sf.CommonAddrSize == 1 && sf.CommonAddr != GlobalCommonAddr && sf.CommonAddr >= 255:
		return nil, ErrParam
	}

	raw := sf.bootstrap[:(sf.IdentifierSize() + len(sf.infoObj))]
	raw[0] = byte(sf.Type)
	raw[1] = sf.Variable.Value()
	raw[2] = sf.Coa.Value()
	offset := 3
	if sf.CauseSize == 2 {
		raw[offset] = byte(sf.OrigAddr)
		offset++
	}
	if sf.CommonAddrSize == 1 {
		raw[offset] = byte(sf.CommonAddr)
	} else { // 2
		raw[offset] = byte(sf.CommonAddr)
		offset++
		raw[offset] = byte(sf.CommonAddr >> 8)
	}
	return raw, nil
}

// UnmarshalBinary honors the encoding.BinaryUnmarshaler interface.
// ASDUParams must be set in advance. All other fields are initialized.
func (sf *ASDU) UnmarshalBinary(rawAsdu []byte) error {
	if !(sf.CauseSize == 1 || sf.CauseSize == 2) ||
		!(sf.CommonAddrSize == 1 || sf.CommonAddrSize == 2) {
		return ErrParam
	}

	// rawAsdu unit identifier size check
	lenDUI := sf.IdentifierSize()
	if lenDUI > len(rawAsdu) {
		return io.EOF
	}

	// parse rawAsdu unit identifier
	sf.Type = TypeID(rawAsdu[0])
	sf.Variable = ParseVariableStruct(rawAsdu[1])
	sf.Coa = ParseCauseOfTransmission(rawAsdu[2])
	if sf.CauseSize == 1 {
		sf.OrigAddr = 0
	} else {
		sf.OrigAddr = OriginAddr(rawAsdu[3])
	}
	if sf.CommonAddrSize == 1 {
		sf.CommonAddr = CommonAddr(rawAsdu[lenDUI-1])
		if sf.CommonAddr == 255 { // map 8-bit variant to 16-bit equivalent
			sf.CommonAddr = GlobalCommonAddr
		}
	} else { // 2
		sf.CommonAddr = CommonAddr(rawAsdu[lenDUI-2]) | CommonAddr(rawAsdu[lenDUI-1])<<8
	}
	// information object
	sf.infoObj = append(sf.bootstrap[lenDUI:lenDUI], rawAsdu[lenDUI:]...)
	return sf.fixInfoObjSize()
}

// fixInfoObjSize fix information object size
func (sf *ASDU) fixInfoObjSize() error {
	// fixed element size
	objSize, err := GetInfoObjSize(sf.Type)
	if err != nil {
		return err
	}

	var size int
	// read the variable structure qualifier
	if sf.Variable.IsSequence {
		size = sf.InfoObjAddrSize + int(sf.Variable.Number)*objSize
	} else {
		size = int(sf.Variable.Number) * (sf.InfoObjAddrSize + objSize)
	}

	switch {
	case size == 0:
		return ErrInfoObjIndexFit
	case size > len(sf.infoObj):
		return io.EOF
	case size < len(sf.infoObj): // not explicitly prohibited
		sf.infoObj = sf.infoObj[:size]
	}

	return nil
}
