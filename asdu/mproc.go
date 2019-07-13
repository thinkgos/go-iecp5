package asdu

import (
	"encoding/binary"
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

// SinglePointInfo are the measured value attributes.
type SinglePointInfo struct {
	Ioa InfoObjAddr
	// value of single point
	Value bool

	// Quality descriptor asdu.OK means no remarks.
	Qds QualityDescriptor

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time time.Time
}

// Single sends a type identification M_SP_NA_1, M_SP_TA_1 or M_SP_TB_1.
// subclass 7.3.1.1 - 7.3.1.2
// 单点信息
func Single(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, attrs ...SinglePointInfo) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		value := byte(0)
		if v.Value {
			value = 0x01
		}
		u.infoObj = append(u.infoObj, value|byte(v.Qds&0xf0))
		switch typeID {
		case M_SP_NA_1:
		case M_SP_TA_1:
			u.infoObj = append(u.infoObj, CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_SP_TB_1:
			u.infoObj = append(u.infoObj, CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// DoublePointInfo are the measured value attributes.
type DoublePointInfo struct {
	Ioa InfoObjAddr

	Value DoublePoint
	// Quality descriptor asdu.OK means no remarks.
	Qds QualityDescriptor

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time time.Time
}

// Double sends a type identification M_DP_NA_1, M_DP_TA_1 or M_DP_TB_1.
// subclass 7.3.1.3 - 7.3.1.4
// 双点信息
func Double(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, attrs ...DoublePointInfo) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.infoObj = append(u.infoObj, byte(v.Value&0x03)|byte(v.Qds&0xf0))
		switch typeID {
		case M_DP_NA_1:
		case M_DP_TA_1:
			u.infoObj = append(u.infoObj, CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_DP_TB_1:
			u.infoObj = append(u.infoObj, CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// StepPositionInfo are the measured value attributes.
type StepPositionInfo struct {
	Ioa InfoObjAddr

	Value StepPosition
	// Quality descriptor asdu.OK means no remarks.
	Qds QualityDescriptor

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time time.Time
}

// Step sends a type identification M_ST_NA_1, M_ST_TA_1 or M_ST_TB_1.
// subclass 7.3.1.5 - 7.3.1.6
// 步位置信息
func Step(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, attrs ...StepPositionInfo) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.infoObj = append(u.infoObj, v.Value.Value(), byte(v.Qds))
		switch typeID {
		case M_ST_NA_1:
		case M_ST_TA_1:
			u.infoObj = append(u.infoObj, CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_SP_TB_1:
			u.infoObj = append(u.infoObj, CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// BitString32Info are the measured value attributes.
type BitString32Info struct {
	Ioa InfoObjAddr

	Value uint32
	// Quality descriptor asdu.OK means no remarks.
	Qds QualityDescriptor

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time time.Time
}

// Bits sends a type identificationM_BO_NA_1, M_BO_TA_1 or M_BO_TB_1.
// subclass 7.3.1.7 - 7.3.1.8
// 比特位串
func BitString32(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, attrs ...BitString32Info) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.infoObj = append(u.infoObj, byte(v.Value), byte(v.Value>>8), byte(v.Value>>16), byte(v.Value>>24), byte(v.Qds))
		switch typeID {
		case M_BO_NA_1:
		case M_BO_TA_1:
			u.infoObj = append(u.infoObj, CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_BO_TB_1:
			u.infoObj = append(u.infoObj, CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// MeasuredValueNormalInfo are the measured value attributes.
type MeasuredValueNormalInfo struct {
	Ioa InfoObjAddr

	Value Normalize
	// Quality descriptor asdu.OK means no remarks.
	Qds QualityDescriptor

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time time.Time
}

// MeasuredValueNormal sends a type identification M_ME_NA_1, M_ME_TA_1, M_ME_TD_1 or M_ME_ND_1.
// subclass 7.3.1.9 - 7.3.1.10
// The quality descriptor must default to info.OK for type M_ME_ND_1.
// 测量值,规一化值
func MeasuredValueNormal(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, attrs ...MeasuredValueNormalInfo) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.infoObj = append(u.infoObj, byte(v.Value), byte(v.Value>>8))
		switch typeID {
		case M_ME_NA_1:
			u.infoObj = append(u.infoObj, byte(v.Qds))
		case M_ME_TA_1:
			u.infoObj = append(u.infoObj, byte(v.Qds))
			u.infoObj = append(u.infoObj, CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_TD_1:
			u.infoObj = append(u.infoObj, byte(v.Qds))
			u.infoObj = append(u.infoObj, CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_ND_1: // 不带品质
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// MeasuredValueScaledInfo are the measured value attributes.
type MeasuredValueScaledInfo struct {
	Ioa InfoObjAddr

	Value int16
	// Quality descriptor asdu.OK means no remarks.
	Qds QualityDescriptor

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time time.Time
}

// MeasuredValueScaled sends a type identification M_ME_NB_1, M_ME_TB_1 or M_ME_TE_1.
// subclass 7.3.1.11 - 7.3.1.12
// 测量值,标度化值
func MeasuredValueScaled(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, attrs ...MeasuredValueNormalInfo) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.infoObj = append(u.infoObj, byte(v.Value), byte(v.Value>>8), byte(v.Qds))
		switch typeID {
		case M_ME_NB_1:
		case M_ME_TB_1:
			u.infoObj = append(u.infoObj, CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_TE_1:
			u.infoObj = append(u.infoObj, CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

// MeasuredValueFloatInfo are the measured value attributes.
type MeasuredValueFloatInfo struct {
	Ioa InfoObjAddr

	Value float32
	// Quality descriptor asdu.OK means no remarks.
	Qds QualityDescriptor

	// The timestamp is nil when the data is invalid or
	// when the type does not include timing at all.
	Time time.Time
}

// MeasuredValueFloat sends a type identification M_ME_NC_1, M_ME_TC_1 or M_ME_TF_1.
// subclass 7.3.1.13 - 7.3.1.14 - 7.3.1.28
// 测量值,短浮点数
func MeasuredValueFloat(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, attrs ...MeasuredValueFloatInfo) error {
	if err := checkValid(c, typeID, isSequence, len(attrs)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(len(attrs)); err != nil {
		return err
	}
	once := false
	for _, v := range attrs {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		bits := math.Float32bits(v.Value)
		u.infoObj = append(u.infoObj, byte(bits), byte(bits>>8), byte(bits>>16), byte(bits>>24), byte(v.Qds&0xf1))
		switch typeID {
		case M_ME_NC_1:
		case M_ME_TC_1:
			u.infoObj = append(u.infoObj, CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_TF_1:
			u.infoObj = append(u.infoObj, CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func (this *ASDU) GetSinglePointInfo() ([]SinglePointInfo, error) {
	var err error

	info := make([]SinglePointInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.infoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}
		value := this.infoObj[offset]
		offset++

		var t time.Time
		switch this.Type {
		case M_SP_NA_1:
		case M_SP_TA_1:
			if t, err = ParseCP24Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_SP_TB_1:
			if t, err = ParseCP56Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.infoObj = this.infoObj[offset:]

		info = append(info, SinglePointInfo{
			Ioa:   infoObjAddr,
			Value: value&0x01 == 0x01,
			Qds:   QualityDescriptor(value & 0xf0),
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetDoublePointInfo() ([]DoublePointInfo, error) {
	var err error

	info := make([]DoublePointInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.infoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}
		value := this.infoObj[offset]
		offset++

		var t time.Time
		switch this.Type {
		case M_DP_NA_1:
		case M_DP_TA_1:
			if t, err = ParseCP24Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_DP_TB_1:
			if t, err = ParseCP56Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.infoObj = this.infoObj[offset:]

		info = append(info, DoublePointInfo{
			Ioa:   infoObjAddr,
			Value: DoublePoint(value & 0x03),
			Qds:   QualityDescriptor(value & 0xf0),
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetStepPositionInfo() ([]StepPositionInfo, error) {
	var err error

	info := make([]StepPositionInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.infoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}
		value := ParseStepPosition(this.infoObj[offset])
		offset++
		qds := QualityDescriptor(this.infoObj[offset])
		offset++

		var t time.Time
		switch this.Type {
		case M_ST_NA_1:
		case M_ST_TA_1:
			if t, err = ParseCP24Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_SP_TB_1:
			if t, err = ParseCP56Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.infoObj = this.infoObj[offset:]

		info = append(info, StepPositionInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetBitString32Info() ([]BitString32Info, error) {
	var err error

	info := make([]BitString32Info, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.infoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}

		value := binary.LittleEndian.Uint32(this.infoObj[offset:])
		offset += 4
		qds := QualityDescriptor(this.infoObj[offset])
		offset++

		var t time.Time
		switch this.Type {
		case M_BO_NA_1:
		case M_BO_TA_1:
			if t, err = ParseCP24Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_BO_TB_1:
			if t, err = ParseCP56Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.infoObj = this.infoObj[offset:]

		info = append(info, BitString32Info{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetMeasuredValueNormalInfo() ([]MeasuredValueNormalInfo, error) {
	var err error

	info := make([]MeasuredValueNormalInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.infoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}

		value := Normalize(binary.LittleEndian.Uint16(this.infoObj[offset:]))
		offset += 2

		var t time.Time
		var qds QualityDescriptor
		switch this.Type {
		case M_ME_NA_1:
			qds = QualityDescriptor(this.infoObj[offset])
			offset++
		case M_ME_TA_1:
			qds = QualityDescriptor(this.infoObj[offset])
			offset++
			if t, err = ParseCP24Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_ME_TD_1:
			qds = QualityDescriptor(this.infoObj[offset])
			offset++
			if t, err = ParseCP56Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		case M_ME_ND_1: // 不带品质
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.infoObj = this.infoObj[offset:]

		info = append(info, MeasuredValueNormalInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetMeasuredValueScaledInfo() ([]MeasuredValueScaledInfo, error) {
	var err error

	info := make([]MeasuredValueScaledInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.infoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}

		value := int16(binary.LittleEndian.Uint16(this.infoObj[offset:]))
		offset += 2
		qds := QualityDescriptor(this.infoObj[offset])
		offset++

		var t time.Time
		switch this.Type {
		case M_ME_NB_1:
		case M_ME_TB_1:
			if t, err = ParseCP24Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_ME_TE_1:
			if t, err = ParseCP56Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.infoObj = this.infoObj[offset:]

		info = append(info, MeasuredValueScaledInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetMeasuredValueFloatInfo() ([]MeasuredValueFloatInfo, error) {
	var err error

	info := make([]MeasuredValueFloatInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once, offset := 0, false, 0; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr, err = this.ParseInfoObjAddr(this.infoObj)
			if err != nil {
				return nil, err
			}
			offset = this.InfoObjAddrSize
		} else {
			infoObjAddr++
			offset = 0
		}

		value := math.Float32frombits(binary.LittleEndian.Uint32(this.infoObj[offset:]))
		offset += 4
		qua := this.infoObj[offset] & 0xf1
		offset++

		var t time.Time
		switch this.Type {
		case M_ME_NC_1:
		case M_ME_TC_1:
			if t, err = ParseCP24Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 3
		case M_ME_TF_1:
			if t, err = ParseCP56Time2a(this.infoObj[offset:], this.Params.InfoObjTimeZone); err != nil {
				return nil, ErrInvalidTimeTag
			}
			offset += 7
		default:
			return nil, ErrTypeIDNotMatch
		}
		this.infoObj = this.infoObj[offset:]
		info = append(info, MeasuredValueFloatInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   QualityDescriptor(qua),
			Time:  t})
	}
	return info, nil
}
