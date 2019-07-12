package asdu

import (
	"math"
	"time"
)

// 在监视方向过程信息的应用服务数据单元

type Connect interface {
	Params() *Params
	Send(a *ASDU) error
}

func checkValid(c Connect, typeID TypeID, isSequence bool, attrsLen int) error {
	if attrsLen == 0 {
		return ErrNotAnyObjInfo
	}
	objSize, err := GetInfoObjSize(typeID)
	if err != nil {
		return err
	}
	param := c.Params()
	if err := param.Valid(); err != nil {
		return err
	}

	var asduLen int
	if isSequence {
		asduLen = param.IdentifierSize() + attrsLen*objSize + param.InfoObjAddrSize
	} else {
		asduLen = param.IdentifierSize() + attrsLen*(objSize+param.InfoObjAddrSize)
	}

	if asduLen > ASDUSizeMax {
		return ErrLengthOutOfRange
	}
	return nil
}

// SinglePointInformation are the measured value attributes.
type SinglePointInformation struct {
	InfoObjAddr InfoObjAddr
	// value of single point
	Value bool

	// Quality descriptor asdu.OK means no remarks.
	QuaDesc QualityDescriptorFlag

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time *time.Time
}

// Single sends a type identification M_SP_NA_1, M_SP_TA_1 or M_SP_TB_1.
// subclause 7.3.1.1 - 7.3.1.2
// 单点信息
func Single(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	commonAddr CommonAddr, attrs ...SinglePointInformation) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		commonAddr,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.InfoObjAddr); err != nil {
				return err
			}
		}

		value := byte(0)
		if v.Value {
			value = 0x01
		}
		u.InfoObj = append(u.InfoObj, value|byte(v.QuaDesc&0xf0))
		switch typeID {
		case M_SP_NA_1:
		case M_SP_TA_1:
			panic("TODO: append 24-bit timestamp")
		case M_SP_TB_1:
			panic("TODO: append 56-bit timestamp")
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// DoublePointInformation are the measured value attributes.
type DoublePointInformation struct {
	InfoObjAddr InfoObjAddr

	Value DoublePoint
	// Quality descriptor asdu.OK means no remarks.
	QuaDesc QualityDescriptorFlag

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time *time.Time
}

// Double sends a type identification M_DP_NA_1, M_DP_TA_1 or M_DP_TB_1.
// subclause 7.3.1.3 - 7.3.1.4
// 双点信息
func Double(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	commonAddr CommonAddr, attrs ...DoublePointInformation) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		commonAddr,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.InfoObjAddr); err != nil {
				return err
			}
		}

		u.InfoObj = append(u.InfoObj, byte(v.Value&0x03)|byte(v.QuaDesc&0xf0))
		switch typeID {
		case M_DP_NA_1:
		case M_DP_TA_1:
			panic("TODO: append 24-bit timestamp")
		case M_DP_TB_1:
			panic("TODO: append 56-bit timestamp")
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// StepPositionInformation are the measured value attributes.
type StepPositionInformation struct {
	InfoObjAddr InfoObjAddr

	Value StepPosition
	// Quality descriptor asdu.OK means no remarks.
	QuaDesc QualityDescriptorFlag

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time *time.Time
}

// Step sends a type identification M_ST_NA_1, M_ST_TA_1 or M_ST_TB_1.
// subclause 7.3.1.5 - 7.3.1.6
// 步位置信息
func Step(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	commonAddr CommonAddr, attrs ...StepPositionInformation) error {
	panic("TODO: not implemented")
}

// BitString32Information are the measured value attributes.
type BitString32Information struct {
	InfoObjAddr InfoObjAddr

	Value uint32
	// Quality descriptor asdu.OK means no remarks.
	QualityDescriptor byte

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time *time.Time
}

// Bits sends a type identificationM_BO_NA_1, M_BO_TA_1 or M_BO_TB_1.
// subclause 7.3.1.7 - 7.3.1.8
// 比特位串
func BitString32(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	commonAddr CommonAddr, attrs ...BitString32Information) error {
	panic("TODO: not implement ed")
}

// MeasuredValueNormalized are the measured value attributes.
type MeasuredValueNormalized struct {
	InfoObjAddr InfoObjAddr

	Value Normalize
	// Quality descriptor asdu.OK means no remarks.
	QuaDesc QualityDescriptorFlag

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time *time.Time
}

// Normal sends a type identification M_ME_NA_1, M_ME_TA_1, M_ME_TD_1 or M_ME_ND_1.
// subclause 7.3.1.9 - 7.3.1.10
// The quality descriptor must default to info.OK for type M_ME_ND_1.
// 测量值,规一化值
func Normal(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	commonAddr CommonAddr, attrs ...MeasuredValueNormalized) error {
	panic("TODO: not implemented")
}

// MeasuredValueScaled are the measured value attributes.
type MeasuredValueScaled struct {
	InfoObjAddr InfoObjAddr

	Value int16
	// Quality descriptor asdu.OK means no remarks.
	QuaDesc QualityDescriptorFlag

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time *time.Time
}

// Scaled sends a type identification M_ME_NB_1, M_ME_TB_1 or M_ME_TE_1.
// subclause 7.3.1.11 - 7.3.1.12
// 测量值,标度化值
func Scaled(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	commonAddr CommonAddr, attrs ...MeasuredValueNormalized) error {
	panic("TODO: not implemented")
}

// MeasuredValueFloat are the measured value attributes.
type MeasuredValueFloat struct {
	InfoObjAddr InfoObjAddr

	Value float32
	// Quality descriptor asdu.OK means no remarks.
	QuaDesc QualityDescriptorFlag

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time *time.Time
}

// Float sends a type identification M_ME_NC_1, M_ME_TC_1 or M_ME_TF_1.
// subclause 7.3.1.13 - 7.3.1.14 - 7.3.1.28
// 测量值,短浮点数
func Float(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	commonAddr CommonAddr, attrs ...MeasuredValueFloat) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		commonAddr,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.InfoObjAddr); err != nil {
				return err
			}
		}

		bits := math.Float32bits(v.Value)
		u.InfoObj = append(u.InfoObj, byte(bits), byte(bits>>8), byte(bits>>16), byte(bits>>24), byte(v.QuaDesc&0xf1))
		switch typeID {
		case M_ME_NC_1:
		case M_ME_TC_1:
			panic("TODO: append 24-bit timestamp")
		case M_ME_TF_1:
			panic("TODO: append 56-bit timestamp")
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func (this *ASDU) GetSinglePointInformation() ([]SinglePointInformation, error) {
	var err error

	info := make([]SinglePointInformation, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.InfoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}
		value := this.InfoObj[offset]
		offset++

		var t *time.Time
		switch this.Type {
		case M_SP_NA_1:
		case M_SP_TA_1:
			if t = ParseCP24Time2a(this.InfoObj[offset:], this.Params.InfoObjTimeZone); t == nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_SP_TB_1:
			if t = ParseCP56Time2a(this.InfoObj[offset:], this.Params.InfoObjTimeZone); t == nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.InfoObj = this.InfoObj[offset:]
		info = append(info, SinglePointInformation{
			InfoObjAddr: infoObjAddr,
			Value:       value&0x01 == 0x01,
			QuaDesc:     QualityDescriptorFlag(value & 0xf0),
			Time:        t,
		})
	}
	return info, nil
}

func (this *ASDU) GetDoublePointInformation() ([]DoublePointInformation, error) {
	var err error

	info := make([]DoublePointInformation, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.InfoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}
		value := this.InfoObj[offset]
		offset++

		var t *time.Time
		switch this.Type {
		case M_DP_NA_1:
		case M_DP_TA_1:
			if t = ParseCP24Time2a(this.InfoObj[offset:], this.Params.InfoObjTimeZone); t == nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_DP_TB_1:
			if t = ParseCP56Time2a(this.InfoObj[offset:], this.Params.InfoObjTimeZone); t == nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.InfoObj = this.InfoObj[offset:]
		info = append(info, DoublePointInformation{
			InfoObjAddr: infoObjAddr,
			Value:       DoublePoint(value & 0x03),
			QuaDesc:     QualityDescriptorFlag(value & 0xf0),
			Time:        t,
		})
	}
	return info, nil
}

func (this *ASDU) GetStepPositionInformation() ([]StepPositionInformation, error) {
	panic("TODO: not implemented")
}

func (this *ASDU) GetBitString32Information() ([]BitString32Information, error) {
	panic("TODO: not implemented")
}

func (this *ASDU) GetMeasuredValueNormalized() ([]MeasuredValueNormalized, error) {
	panic("TODO: not implemented")
}

func (this *ASDU) GetMeasuredValueScaled() ([]MeasuredValueScaled, error) {
	panic("TODO: not implemented")
}

func (this *ASDU) GetMeasuredValueFloat() ([]MeasuredValueFloat, error) {
	var err error

	info := make([]MeasuredValueFloat, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.InfoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}

		bits := uint32(this.InfoObj[offset]) | uint32(this.InfoObj[offset+1]<<8) |
			uint32(this.InfoObj[offset+2]<<16) | uint32(this.InfoObj[offset+3]<<24)
		value := math.Float32frombits(bits)
		qua := this.InfoObj[offset+4] & 0xf1
		offset += 5

		var t *time.Time
		switch this.Type {
		case M_ME_NC_1:
		case M_ME_TC_1:
			if t = ParseCP24Time2a(this.InfoObj[offset:], this.Params.InfoObjTimeZone); t == nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_ME_TF_1:
			if t = ParseCP56Time2a(this.InfoObj[offset:], this.Params.InfoObjTimeZone); t == nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.InfoObj = this.InfoObj[offset:]
		info = append(info, MeasuredValueFloat{
			InfoObjAddr: infoObjAddr,
			Value:       value,
			QuaDesc:     QualityDescriptorFlag(qua),
			Time:        t,
		})
	}
	return info, nil
}
