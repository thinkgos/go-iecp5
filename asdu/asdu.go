// Package asdu provides the OSI presentation layer.
package asdu

import (
	"fmt"
	"io"
	"math/bits"
	"time"
)

const (
	ASDUSizeMax = 249 // ASDU
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
// See companion standard 101, subclause 7.1.
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
		return errParam
	}
	return nil
}

// ValidCommonAddr returns the validation result of a station address.
func (this Params) ValidCommonAddr(addr CommonAddr) error {
	if addr == InvalidCommonAddr {
		return errCommonAddrZero
	}
	if bits.Len(uint(addr)) > this.CommonAddrSize*8 {
		return errCommonAddrFit
	}
	return nil
}

// IdentifierSize the application data unit identifies size
func (this *Params) IdentifierSize() int {
	return 2 + int(this.CauseSize) + int(this.CommonAddrSize)
}

// DecodeInfoObjAddr decodes an information object address from buf.
// The function panics when the byte array is too small
// or when the address size parameter is out of bounds.
func (p *Params) DecodeInfoObjAddr(buf []byte) InfoObjAddr {
	addr := InfoObjAddr(buf[0])
	switch p.InfoObjAddrSize {
	case 1:
	case 2:
		addr |= InfoObjAddr(buf[1]) << 8
	case 3:
		addr |= InfoObjAddr(buf[1])<<8 | InfoObjAddr(buf[2])<<16
	default:
		panic(errParam)
	}

	return addr
}

// Identifier the application data unit identifies.
type Identifier struct {
	// type identification, information content
	Type TypeID
	// Variable is variable structure qualifier
	Variable Variable
	// cause of transmission submission category
	Cause Cause
	// Originator Address [1, 255] or 0 for the default.
	// The applicability is controlled by Params.CauseSize.
	OrigAddr OriginAddr
	// CommonAddr is a station address. Zero is not used.
	// The width is controlled by Params.CommonAddrSize.
	// See companion standard 101, subclause 7.2.4.
	CommonAddr CommonAddr // station address 公共地址是站地址
}

// String 返回数据单元标识符的信息like "TypeID Cause OrigAddr@CommonAddr"
func (id Identifier) String() string {
	if id.OrigAddr == 0 {
		return fmt.Sprintf("%s %s @%d", id.Type, id.Cause, id.CommonAddr)
	}
	return fmt.Sprintf("%s %s %d@%d ", id.Type, id.Cause, id.OrigAddr, id.CommonAddr)
}

// ASDU (Application Service Data Unit) is an application message.
type ASDU struct {
	*Params
	Identifier
	InfoObj   []byte            // information object serial
	bootstrap [ASDUSizeMax]byte // prevents Info malloc
}

func NewASDU(p *Params, typeID TypeID, isSequence bool, coa CauseOfTransmission, commonAddr CommonAddr) *ASDU {
	vars := Variable(0)
	if isSequence {
		vars = VariableSeq
	}
	if coa.IsTest {
		coa.Cause |= TestFlag
	}
	if coa.IsNegative {
		coa.Cause |= NegFlag
	}
	a := &ASDU{
		Params: p,
		Identifier: Identifier{
			Type:       typeID,
			Variable:   vars,
			Cause:      coa.Cause,
			CommonAddr: commonAddr,
		}}
	a.InfoObj = a.bootstrap[:(2 + a.CauseSize + a.CommonAddrSize)]
	return a
}

// AddInfoObjAddr appends an information object address to Info.
func (u *ASDU) AddInfoObjAddr(addr InfoObjAddr) error {
	switch u.InfoObjAddrSize {
	case 1:
		if addr > 255 {
			return errInfoObjAddrFit
		}
		u.InfoObj = append(u.InfoObj, byte(addr))
	case 2:
		if addr > 65535 {
			return errInfoObjAddrFit
		}
		u.InfoObj = append(u.InfoObj, byte(addr), byte(addr>>8))
	case 3:
		if addr > 16777215 {
			return errInfoObjAddrFit
		}
		u.InfoObj = append(u.InfoObj, byte(addr), byte(addr>>8), byte(addr>>16))
	default:
		return errParam
	}
	return nil
}

//if u.InfoSeq {
//	objCount = (len(u.InfoObj) - u.InfoObjAddrSize) / objSize
//} else {
//	objCount = len(u.InfoObj) / (u.InfoObjAddrSize + objSize)
//}
// See companion standard 101, subclause 7.2.2.
func (this *ASDU) IncVariableNumber(n int) error {
	objCount := this.Variable & 0x7f
	objCount += Variable(n)
	if objCount >= 128 {
		return errInfoObjIndexFit
	}
	this.Variable = (this.Variable & 0x80) | objCount
	return nil
}

// MustNewInro returns a new interrogation command [C_IC_NA_1].
// Use group 1 to 16, or 0 for the default.
//func MustNewInro(p *Params, commonAddr CommonAddr, origAddr OriginAddr, group uint) *ASDU {
//	if group > 16 {
//		group = 0
//	}
//
//	u := &ASDU{
//		Params:     p,
//		Identifier: Identifier{C_IC_NA_1, Act, origAddr, commonAddr},
//	}
//
//	u.InfoObj = u.bootstrap[:p.InfoObjAddrSize+1]
//	u.InfoObj[p.InfoObjAddrSize] = byte(group + uint(Inrogen))
//	return u
//}
//
//// Respond returns a new "responding" ASDU which addresses "initiating" u.
//func (u *ASDU) Respond(t TypeID, c Cause) *ASDU {
//	return NewASDU(u.Params, Identifier{
//		CommonAddr: u.CommonAddr,
//		OrigAddr:   u.OrigAddr,
//		Type:       t,
//		Cause:      c | u.Cause&TestFlag,
//	})
//}
//
//// Reply returns a new "responding" ASDU which addresses "initiating" u with a copy of Info.
//func (u *ASDU) Reply(c Cause) *ASDU {
//	r := NewASDU(u.Params, u.Identifier)
//	r.Cause = c | u.Cause&TestFlag
//	r.InfoSeq = u.InfoSeq
//	r.InfoObj = append(r.InfoObj, u.InfoObj...)
//	return r
//}
//
//// String returns a full description.
//func (u *ASDU) String() string {
//	dataSize, err := GetInfoObjSize(u.Type)
//	if err != nil {
//		if !u.InfoSeq {
//			return fmt.Sprintf("%s: %#x", u.Identifier, u.InfoObj)
//		}
//		return fmt.Sprintf("%s seq: %#x", u.Identifier, u.InfoObj)
//	}
//
//	end := len(u.InfoObj)
//	addrSize := u.InfoObjAddrSize
//	if end < addrSize {
//		if !u.InfoSeq {
//			return fmt.Sprintf("%s: %#x <EOF>", u.Identifier, u.InfoObj)
//		}
//		return fmt.Sprintf("%s seq: %#x <EOF>", u.Identifier, u.InfoObj)
//	}
//	addr := u.DecodeInfoObjAddr(u.InfoObj)
//
//	buf := bytes.NewBufferString(u.Identifier.String())
//
//	for i := addrSize; ; {
//		start := i
//		i += dataSize
//		if i > end {
//			fmt.Fprintf(buf, " %d:%#x <EOF>", addr, u.InfoObj[start:])
//			break
//		}
//		fmt.Fprintf(buf, " %d:%#x", addr, u.InfoObj[start:i])
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
//				fmt.Fprintf(buf, " %#x <EOF>", u.InfoObj[start:i])
//				break
//			}
//			addr = u.DecodeInfoObjAddr(u.InfoObj[start:])
//		}
//	}
//
//	return buf.String()
//}

// MarshalBinary honors the encoding.BinaryMarshaler interface.
func (u *ASDU) MarshalBinary() (data []byte, err error) {
	switch {
	case u.Cause == Unused:
		return nil, errCauseZero
	case !(u.CauseSize == 1 || u.CauseSize == 2):
		return nil, errParam
	case u.CauseSize == 1 && u.OrigAddr != 0:
		return nil, errOriginAddrFit
	case u.CommonAddr == InvalidCommonAddr:
		return nil, errCommonAddrZero
	case !(u.CommonAddrSize == 1 || u.CommonAddrSize == 2):
		return nil, errParam
	case u.CommonAddrSize == 1 && u.CommonAddr != GlobalCommonAddr && u.CommonAddr >= 255:
		return nil, errParam
	}

	offset := 0
	u.InfoObj[offset] = byte(u.Type)
	offset++
	u.InfoObj[offset] = byte(u.Variable)
	offset++
	u.InfoObj[offset] = byte(u.Cause)
	offset++
	if u.CauseSize == 2 {
		u.InfoObj[offset] = byte(u.OrigAddr)
		offset++
	}
	if u.CommonAddrSize == 1 {
		if u.CommonAddr == GlobalCommonAddr {
			u.InfoObj[offset] = 255
		} else {
			u.InfoObj[offset] = byte(u.CommonAddr)
		}
	} else { // 2
		u.InfoObj[offset] = byte(u.CommonAddr)
		offset++
		u.InfoObj[offset] = byte(u.CommonAddr >> 8)
	}

	return u.InfoObj, nil
}

// UnmarshalBinary honors the encoding.BinaryUnmarshaler interface.
// ASDUParams must be set in advance. All other fields are initialized.
func (u *ASDU) UnmarshalBinary(data []byte) error {
	// data unit identifier size check
	lenDUI := u.IdentifierSize()
	if lenDUI > len(data) {
		return io.EOF
	}
	u.InfoObj = append(u.bootstrap[:0], data[lenDUI:]...)

	u.Type = TypeID(data[0])
	// fixed element size
	objSize, err := GetInfoObjSize(u.Type)
	if err != nil {
		return err
	}

	u.Variable = Variable(data[1])
	var size int
	// read the variable structure qualifier
	if u.Variable&VariableSeq == VariableSeq {
		size = u.InfoObjAddrSize + (int(u.Variable&0x7f) * objSize)
	} else {
		size = int(u.Variable) * (u.InfoObjAddrSize + objSize)
	}

	switch {
	case size == 0:
		return errInfoObjIndexFit
	case size > len(u.InfoObj):
		return io.EOF
	case size < len(u.InfoObj): // not explicitly prohibited
		u.InfoObj = u.InfoObj[:size]
	}

	u.Cause = Cause(data[2])
	switch u.CauseSize {
	case 1:
		u.OrigAddr = 0
	case 2:
		u.OrigAddr = OriginAddr(data[3])
	default:
		return errParam
	}

	switch u.CommonAddrSize {
	case 1:
		addr := CommonAddr(data[lenDUI-1])
		if addr == 255 { // map 8-bit variant to 16-bit equivalent
			addr = GlobalCommonAddr
		}
		u.CommonAddr = addr
	case 2:
		u.CommonAddr = CommonAddr(data[lenDUI-2]) | CommonAddr(data[lenDUI-1])<<8
	default:
		return errParam
	}
	return nil
}
