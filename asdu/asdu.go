// Package asdu provides the OSI presentation layer.
package asdu

import (
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"time"
)

// ASDU
// | - data unit identification- | - information object <1..n> - |
//
// | <------------  data unit identification --------->|
// | typeID | variable struct | cause | common address |
// | <------------  information object -------------->|
// | object address | element set | object time scale |

var (
	// ParamsNarrow is the smalles configuration.
	ParamsNarrow = &Params{
		CauseSize: 1, CommonAddrSize: 1,
		InfoObjAddrSize: 1, InfoObjTimeZone: time.UTC}
	// ParamsWide is the largest configuration.
	ParamsWide = &Params{
		CauseSize: 2, CommonAddrSize: 2,
		InfoObjAddrSize: 3, InfoObjTimeZone: time.UTC}
)

// Params 定义了ASDU相关特定参数
// See companion standard 101, subclause 7.1.
type Params struct {
	// cause of transmission, 传输原因字节数
	// The standard requires "b" in [1, 2].
	// Value 2 includes/activates the originator address.
	CauseSize int

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

// Valid returns the validation result.
func (this Params) Valid() error {
	switch {
	case this.CauseSize < 1 || this.CauseSize > 2:
		return errParam
	case this.CommonAddrSize < 1 || this.CommonAddrSize > 2:
		return errParam
	case this.InfoObjAddrSize < 1 || this.InfoObjAddrSize > 3:
		return errParam
	case this.InfoObjTimeZone == nil:
		return errParam
	default:
		return nil
	}
}

// ValidAddr returns the validation result of a station address.
func (this Params) ValidCommonAddr(addr CommonAddr) error {
	if addr == 0 {
		return errCommonoAddrZero
	}
	if bits.Len(uint(addr)) > this.CommonAddrSize*8 {
		return errCommonAddrFit
	}
	return nil
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
	// ASDU type identification, information content
	Type TypeID
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

// String 返回数据单元标识符的信息
func (id Identifier) String() string {
	if id.OrigAddr == 0 {
		return fmt.Sprintf("%s %s @%d",
			id.Type, id.Cause, id.CommonAddr)
	}
	return fmt.Sprintf("%s %s %d@%d ",
		id.Type, id.Cause, id.OrigAddr, id.CommonAddr)
}

// ASDU (Application Service Data Unit) is an application message.
type ASDU struct {
	*Params
	Identifier
	// marks Info as a sequence
	// number <0..127>
	// 0: 同一类型，有不同objAddress的信息元素集合
	// 1： 同一类型，相同objAddress信息元素集合
	InfoSeq   bool
	Info      []byte   // information object serial
	bootstrap [17]byte // prevents Info malloc
}

// NewASDU returns a new ASDU with the provided parameters.
func NewASDU(p *Params, id Identifier) *ASDU {
	u := &ASDU{Params: p, Identifier: id}
	u.Info = u.bootstrap[:0]
	return u
}

// MustNewInro returns a new interrogation command [C_IC_NA_1].
// Use group 1 to 16, or 0 for the default.
func MustNewInro(p *Params, commonAddr CommonAddr, origAddr OriginAddr, group uint) *ASDU {
	if group > 16 {
		group = 0
	}

	u := &ASDU{
		Params:     p,
		Identifier: Identifier{C_IC_NA_1, Act, origAddr, commonAddr},
	}

	u.Info = u.bootstrap[:p.InfoObjAddrSize+1]
	u.Info[p.InfoObjAddrSize] = byte(group + uint(Inrogen))
	return u
}

// Respond returns a new "responding" ASDU which addresses "initiating" u.
func (u *ASDU) Respond(t TypeID, c Cause) *ASDU {
	return NewASDU(u.Params, Identifier{
		CommonAddr: u.CommonAddr,
		OrigAddr:   u.OrigAddr,
		Type:       t,
		Cause:      c | u.Cause&TestFlag,
	})
}

// Reply returns a new "responding" ASDU which addresses "initiating" u with a copy of Info.
func (u *ASDU) Reply(c Cause) *ASDU {
	r := NewASDU(u.Params, u.Identifier)
	r.Cause = c | u.Cause&TestFlag
	r.InfoSeq = u.InfoSeq
	r.Info = append(r.Info, u.Info...)
	return r
}

// String returns a full description.
func (u *ASDU) String() string {
	dataSize, err := GetInfoObjSize(u.Type)
	if err != nil {
		if !u.InfoSeq {
			return fmt.Sprintf("%s: %#x", u.Identifier, u.Info)
		}
		return fmt.Sprintf("%s seq: %#x", u.Identifier, u.Info)
	}

	end := len(u.Info)
	addrSize := u.InfoObjAddrSize
	if end < addrSize {
		if !u.InfoSeq {
			return fmt.Sprintf("%s: %#x <EOF>", u.Identifier, u.Info)
		}
		return fmt.Sprintf("%s seq: %#x <EOF>", u.Identifier, u.Info)
	}
	addr := u.DecodeInfoObjAddr(u.Info)

	buf := bytes.NewBufferString(u.Identifier.String())

	for i := addrSize; ; {
		start := i
		i += dataSize
		if i > end {
			fmt.Fprintf(buf, " %d:%#x <EOF>", addr, u.Info[start:])
			break
		}
		fmt.Fprintf(buf, " %d:%#x", addr, u.Info[start:i])
		if i == end {
			break
		}

		if u.InfoSeq {
			addr++
		} else {
			start = i
			i += addrSize
			if i > end {
				fmt.Fprintf(buf, " %#x <EOF>", u.Info[start:i])
				break
			}
			addr = u.DecodeInfoObjAddr(u.Info[start:])
		}
	}

	return buf.String()
}

// MarshalBinary honors the encoding.BinaryMarshaler interface.
func (u *ASDU) MarshalBinary() (data []byte, err error) {
	switch {
	case u.Cause == Unused:
		return nil, errCauseZero
	case u.CommonAddr == InvalidCommonAddr:
		return nil, errCommonoAddrZero
	case !(u.CauseSize == 1 || u.CauseSize == 2):
		return nil, errParam
	case u.CauseSize == 1 && u.OrigAddr != 0:
		return nil, errOriginAddrFit
	case !(u.CommonAddrSize == 1 || u.CommonAddrSize == 2):
		return nil, errParam
	case u.CommonAddrSize == 1 && u.CommonAddr != GlobalCommonAddr && u.CommonAddr >= 255:
		return nil, errParam
	}

	// calculate the size declaration byte named "variable structure qualifier"
	// fixed element size
	objSize, err := GetInfoObjSize(u.Type)
	if err != nil {
		return nil, err
	}
	var vsq byte
	var objCount int
	// See companion standard 101, subclause 7.2.2.
	if u.InfoSeq {
		vsq = 0x80
		objCount = (len(u.Info) - u.InfoObjAddrSize) / objSize
	} else {
		objCount = len(u.Info) / (u.InfoObjAddrSize + objSize)
	}
	if objCount >= 128 {
		return nil, errInfoObjIndexFit
	}
	vsq |= byte(objCount)

	data = make([]byte, 0, 2+u.CauseSize+u.CommonAddrSize+len(u.Info))
	data = append(data, byte(u.Type), vsq, byte(u.Cause))
	if u.CauseSize == 2 {
		data = append(data, byte(u.OrigAddr))
	}
	if u.CommonAddrSize == 1 {
		if u.CommonAddr == GlobalCommonAddr {
			data = append(data, 255)
		} else {
			data = append(data, byte(u.CommonAddr))
		}
	} else { // 2
		data = append(data, byte(u.CommonAddr), byte(u.CommonAddr>>8))
	}

	return append(data, u.Info...), nil
}

// UnmarshalBinary honors the encoding.BinaryUnmarshaler interface.
// ASDUParams must be set in advance. All other fields are initialized.
func (u *ASDU) UnmarshalBinary(data []byte) error {
	// data unit identifier size check
	lenDUI := 2 + u.CauseSize + u.CommonAddrSize
	if lenDUI > len(data) {
		return io.EOF
	}
	u.Info = append(u.bootstrap[:0], data[lenDUI:]...)

	u.Type = TypeID(data[0])
	// fixed element size
	objSize, err := GetInfoObjSize(u.Type)
	if err != nil {
		return err
	}

	var size int
	// read the variable structure qualifier
	if vsq := data[1]; vsq > 127 {
		u.InfoSeq = true
		objCount := int(vsq & 127)
		size = u.InfoObjAddrSize + (objCount * objSize)
	} else {
		u.InfoSeq = false
		objCount := int(vsq)
		size = objCount * (u.InfoObjAddrSize + objSize)
	}

	switch {
	case size == 0:
		return errInfoObjIndexFit
	case size > len(u.Info):
		return io.EOF
	case size < len(u.Info): // not explicitly prohibited
		u.Info = u.Info[:size]
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

// AddInfoObjAddr appends an information object address to Info.
func (u *ASDU) AddInfoObjAddr(addr InfoObjAddr) error {
	switch u.InfoObjAddrSize {
	case 1:
		if addr > 255 {
			return errInfoObjAddrFit
		}
		u.Info = append(u.Info, byte(addr))

	case 2:
		if addr > 65535 {
			return errInfoObjAddrFit
		}
		u.Info = append(u.Info, byte(addr), byte(addr>>8))

	case 3:
		if addr > 16777215 {
			return errInfoObjAddrFit
		}
		u.Info = append(u.Info, byte(addr), byte(addr>>8), byte(addr>>16))
	default:
		return errParam
	}

	return nil
}
