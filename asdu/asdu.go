// Package asdu provides the OSI presentation layer.
package asdu

import (
	"fmt"
	"io"
	"math/bits"
	"time"
)

const (
	ASDUSizeMax = 249 // ASDU max size
)

// ASDU form
// | data unit identification | information object <1..n> |
//
//      | <------------  data unit identification --------->|
//      | typeID | variable struct | cause | common address |
// bytes|    1   |      1          | [1,2] |     [1,2]      |
//      | <------------  information object -------------->|
//      | object address | element set | object time scale |
// bytes|     [1,2,3]    |             |                   |

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
func (this Params) Valid() error {
	if (this.CauseSize < 1 || this.CauseSize > 2) ||
		(this.CommonAddrSize < 1 || this.CommonAddrSize > 2) ||
		(this.InfoObjAddrSize < 1 || this.InfoObjAddrSize > 3) ||
		(this.InfoObjTimeZone == nil) {
		return ErrParam
	}
	return nil
}

// ValidCommonAddr returns the validation result of a station address.
func (this Params) ValidCommonAddr(addr CommonAddr) error {
	if addr == InvalidCommonAddr {
		return ErrCommonAddrZero
	}
	if bits.Len(uint(addr)) > this.CommonAddrSize*8 {
		return ErrCommonAddrFit
	}
	return nil
}

// IdentifierSize the application data unit identifies size
func (this Params) IdentifierSize() int {
	return 2 + int(this.CauseSize) + int(this.CommonAddrSize)
}

// Identifier the application data unit identifies.
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

// String 返回数据单元标识符的信息like "TypeID Cause OrigAddr@CommonAddr"
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

func NewEmptyASDU(p *Params) *ASDU {
	a := &ASDU{Params: p}
	lenDUI := a.IdentifierSize()
	a.infoObj = a.bootstrap[lenDUI:lenDUI]
	return a
}

func NewASDU(p *Params, identifier Identifier) *ASDU {
	a := NewEmptyASDU(p)
	a.Identifier = identifier
	return a
}

// SetVariableNumber See companion standard 101, subclass 7.2.2.
func (this *ASDU) SetVariableNumber(n int) error {
	if n >= 128 {
		return ErrInfoObjIndexFit
	}
	this.Variable.Number = byte(n)
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

// Reply returns a new "responding" ASDU which addresses "initiating" u with a copy of Info.
func (this *ASDU) Reply(c Cause, addr CommonAddr) *ASDU {
	this.CommonAddr = addr
	r := NewASDU(this.Params, this.Identifier)
	r.Coa.Cause = c
	r.infoObj = append(r.infoObj, this.infoObj...)
	return r
}

func (this *ASDU) ReplyMirror(c Cause) *ASDU {
	this.Coa.Cause = c
	return this
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
func (this *ASDU) MarshalBinary() (data []byte, err error) {
	switch {
	case this.Coa.Cause == Unused:
		return nil, ErrCauseZero
	case !(this.CauseSize == 1 || this.CauseSize == 2):
		return nil, ErrParam
	case this.CauseSize == 1 && this.OrigAddr != 0:
		return nil, ErrOriginAddrFit
	case this.CommonAddr == InvalidCommonAddr:
		return nil, ErrCommonAddrZero
	case !(this.CommonAddrSize == 1 || this.CommonAddrSize == 2):
		return nil, ErrParam
	case this.CommonAddrSize == 1 && this.CommonAddr != GlobalCommonAddr && this.CommonAddr >= 255:
		return nil, ErrParam
	}

	raw := this.bootstrap[:(this.IdentifierSize() + len(this.infoObj))]
	raw[0] = byte(this.Type)
	raw[1] = this.Variable.Value()
	raw[2] = byte(this.Coa.Value())
	offset := 3
	if this.CauseSize == 2 {
		raw[offset] = byte(this.OrigAddr)
		offset++
	}
	if this.CommonAddrSize == 1 {
		if this.CommonAddr == GlobalCommonAddr {
			raw[offset] = 255
		} else {
			raw[offset] = byte(this.CommonAddr)
		}
	} else { // 2
		raw[offset] = byte(this.CommonAddr)
		offset++
		raw[offset] = byte(this.CommonAddr >> 8)
	}
	return raw, nil
}

// UnmarshalBinary honors the encoding.BinaryUnmarshaler interface.
// ASDUParams must be set in advance. All other fields are initialized.
func (this *ASDU) UnmarshalBinary(rawAsdu []byte) error {
	if !(this.CauseSize == 1 || this.CauseSize == 2) ||
		!(this.CommonAddrSize == 1 || this.CommonAddrSize == 2) {
		return ErrParam
	}

	// rawAsdu unit identifier size check
	lenDUI := this.IdentifierSize()
	if lenDUI > len(rawAsdu) {
		return io.EOF
	}

	// parse rawAsdu unit identifier
	this.Type = TypeID(rawAsdu[0])
	this.Variable = ParseVariableStruct(rawAsdu[1])
	this.Coa = ParseCauseOfTransmission(rawAsdu[2])
	if this.CauseSize == 1 {
		this.OrigAddr = 0
	} else {
		this.OrigAddr = OriginAddr(rawAsdu[3])
	}
	if this.CommonAddrSize == 1 {
		this.CommonAddr = CommonAddr(rawAsdu[lenDUI-1])
		if this.CommonAddr == 255 { // map 8-bit variant to 16-bit equivalent
			this.CommonAddr = GlobalCommonAddr
		}
	} else { // 2
		this.CommonAddr = CommonAddr(rawAsdu[lenDUI-2]) | CommonAddr(rawAsdu[lenDUI-1])<<8
	}
	// information object
	this.infoObj = append(this.bootstrap[lenDUI:lenDUI], rawAsdu[lenDUI:]...)
	return this.fixInfoObjSize()
}

func (this *ASDU) fixInfoObjSize() error {
	// fixed element size
	objSize, err := GetInfoObjSize(this.Type)
	if err != nil {
		return err
	}

	var size int
	// read the variable structure qualifier
	if this.Variable.IsSequence {
		size = this.InfoObjAddrSize + int(this.Variable.Number)*objSize
	} else {
		size = int(this.Variable.Number) * (this.InfoObjAddrSize + objSize)
	}

	switch {
	case size == 0:
		return ErrInfoObjIndexFit
	case size > len(this.infoObj):
		return io.EOF
	case size < len(this.infoObj): // not explicitly prohibited
		this.infoObj = this.infoObj[:size]
	}

	return nil
}
