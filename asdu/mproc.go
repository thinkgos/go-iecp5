package asdu

import (
	"time"
)

// 在监视方向过程信息的应用服务数据单元

type Connect interface {
	Params() *Params
	Send(a *ASDU) error
}

func checkValid(c Connect, typeID TypeID, isSequence bool, infosLen int) error {
	if infosLen == 0 {
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
		asduLen = param.IdentifierSize() + infosLen*objSize + param.InfoObjAddrSize
	} else {
		asduLen = param.IdentifierSize() + infosLen*(objSize+param.InfoObjAddrSize)
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

// single sends a type identification M_SP_NA_1, M_SP_TA_1 or M_SP_TB_1.
// subclass 7.3.1.1 - 7.3.1.2
// 单点信息
func single(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...SinglePointInfo) error {
	if err := checkValid(c, typeID, isSequence, len(infos)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.SetVariableNumber(len(infos)); err != nil {
		return err
	}
	once := false
	for _, v := range infos {
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
		u.AppendBytes(value | byte(v.Qds&0xf0))
		switch typeID {
		case M_SP_NA_1:
		case M_SP_TA_1:
			u.AppendBytes(CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_SP_TB_1:
			u.AppendBytes(CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func Single(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...SinglePointInfo) error {
	if !(coa.Cause == Back || coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return single(c, M_SP_NA_1, isSequence, coa, ca, infos...)
}

func SingleCP24Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...SinglePointInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc) {
		return ErrCmdCause
	}
	return single(c, M_SP_TA_1, false, coa, ca, infos...)
}

func SingleCP56Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...SinglePointInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc) {
		return ErrCmdCause
	}
	return single(c, M_SP_TB_1, false, coa, ca, infos...)
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

// double sends a type identification M_DP_NA_1, M_DP_TA_1 or M_DP_TB_1.
// subclass 7.3.1.3 - 7.3.1.4
// 双点信息
func double(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...DoublePointInfo) error {
	if err := checkValid(c, typeID, isSequence, len(infos)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.SetVariableNumber(len(infos)); err != nil {
		return err
	}
	once := false
	for _, v := range infos {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.AppendBytes(byte(v.Value&0x03) | byte(v.Qds&0xf0))
		switch typeID {
		case M_DP_NA_1:
		case M_DP_TA_1:
			u.AppendBytes(CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_DP_TB_1:
			u.AppendBytes(CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func Double(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...DoublePointInfo) error {
	if !(coa.Cause == Back || coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return double(c, M_DP_NA_1, isSequence, coa, ca, infos...)
}

func DoubleCP24Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...DoublePointInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc) {
		return ErrCmdCause
	}
	return double(c, M_DP_TA_1, false, coa, ca, infos...)
}

func DoubleCP56Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...DoublePointInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc) {
		return ErrCmdCause
	}
	return double(c, M_DP_TB_1, false, coa, ca, infos...)
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

// step sends a type identification M_ST_NA_1, M_ST_TA_1 or M_ST_TB_1.
// subclass 7.3.1.5 - 7.3.1.6
// 步位置信息
func step(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...StepPositionInfo) error {
	if err := checkValid(c, typeID, isSequence, len(infos)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.SetVariableNumber(len(infos)); err != nil {
		return err
	}
	once := false
	for _, v := range infos {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.AppendBytes(v.Value.Value(), byte(v.Qds))
		switch typeID {
		case M_ST_NA_1:
		case M_ST_TA_1:
			u.AppendBytes(CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_SP_TB_1:
			u.AppendBytes(CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func Step(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...StepPositionInfo) error {
	if !(coa.Cause == Back || coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return step(c, M_ST_NA_1, isSequence, coa, ca, infos...)
}

func StepCP24Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...StepPositionInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc) {
		return ErrCmdCause
	}
	return step(c, M_ST_TA_1, false, coa, ca, infos...)
}

func StepCP56Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...StepPositionInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req ||
		coa.Cause == Retrem || coa.Cause == Retloc) {
		return ErrCmdCause
	}
	return step(c, M_SP_TB_1, false, coa, ca, infos...)
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
func bitString32(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...BitString32Info) error {
	if err := checkValid(c, typeID, isSequence, len(infos)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.SetVariableNumber(len(infos)); err != nil {
		return err
	}
	once := false
	for _, v := range infos {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}
		u.AppendBitsString32(v.Value).AppendBytes(byte(v.Qds))

		switch typeID {
		case M_BO_NA_1:
		case M_BO_TA_1:
			u.AppendBytes(CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_BO_TB_1:
			u.AppendBytes(CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func BitString32(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...BitString32Info) error {
	if !(coa.Cause == Back || coa.Cause == Spont || coa.Cause == Req ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return bitString32(c, M_BO_NA_1, isSequence, coa, ca, infos...)
}

func BitString32CP24Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...BitString32Info) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return bitString32(c, M_BO_TA_1, false, coa, ca, infos...)
}

func BitString32CP56Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...BitString32Info) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return bitString32(c, M_BO_TB_1, false, coa, ca, infos...)
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

// measuredValueNormal sends a type identification M_ME_NA_1, M_ME_TA_1, M_ME_TD_1 or M_ME_ND_1.
// subclass 7.3.1.9 - 7.3.1.10
// The quality descriptor must default to info.OK for type M_ME_ND_1.
// 测量值,规一化值
func measuredValueNormal(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
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
	if err := u.SetVariableNumber(len(attrs)); err != nil {
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
		u.AppendNormalize(v.Value)
		switch typeID {
		case M_ME_NA_1:
			u.AppendBytes(byte(v.Qds))
		case M_ME_TA_1:
			u.AppendBytes(byte(v.Qds)).AppendBytes(CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_TD_1:
			u.AppendBytes(byte(v.Qds)).AppendBytes(CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_ND_1: // 不带品质
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func MeasuredValueNormal(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueNormalInfo) error {
	if !(coa.Cause == Percyc || coa.Cause == Back ||
		coa.Cause == Spont || coa.Cause == Req ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return measuredValueNormal(c, M_ME_NA_1, isSequence, coa, ca, infos...)
}

func MeasuredValueNormalCP24Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueNormalInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return measuredValueNormal(c, M_ME_TA_1, false, coa, ca, infos...)
}

func MeasuredValueNormalCP56Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueNormalInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return measuredValueNormal(c, M_ME_TD_1, false, coa, ca, infos...)
}
func MeasuredValueNormalNoQuality(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueNormalInfo) error {
	if !(coa.Cause == Percyc || coa.Cause == Back ||
		coa.Cause == Spont || coa.Cause == Req ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return measuredValueNormal(c, M_ME_ND_1, isSequence, coa, ca, infos...)
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

// measuredValueScaled sends a type identification M_ME_NB_1, M_ME_TB_1 or M_ME_TE_1.
// subclass 7.3.1.11 - 7.3.1.12
// 测量值,标度化值
func measuredValueScaled(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueScaledInfo) error {
	if err := checkValid(c, typeID, isSequence, len(infos)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.SetVariableNumber(len(infos)); err != nil {
		return err
	}
	once := false
	for _, v := range infos {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}
		u.AppendScaled(v.Value).AppendBytes(byte(v.Qds))
		switch typeID {
		case M_ME_NB_1:
		case M_ME_TB_1:
			u.AppendBytes(CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_TE_1:
			u.AppendBytes(CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func MeasuredValueScaled(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueScaledInfo) error {
	if !(coa.Cause == Percyc || coa.Cause == Back ||
		coa.Cause == Spont || coa.Cause == Req ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return measuredValueScaled(c, M_ME_NB_1, isSequence, coa, ca, infos...)
}

func MeasuredValueScaledCP24Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueScaledInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return measuredValueScaled(c, M_ME_TB_1, false, coa, ca, infos...)
}

func MeasuredValueScaledCP56Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueScaledInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return measuredValueScaled(c, M_ME_TE_1, false, coa, ca, infos...)
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

// measuredValueFloat sends a type identification M_ME_NC_1, M_ME_TC_1 or M_ME_TF_1.
// subclass 7.3.1.13 - 7.3.1.14 - 7.3.1.28
// 测量值,短浮点数
func measuredValueFloat(c Connect, typeID TypeID, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueFloatInfo) error {
	if err := checkValid(c, typeID, isSequence, len(infos)); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: isSequence},
		coa,
		0,
		ca,
	})
	if err := u.SetVariableNumber(len(infos)); err != nil {
		return err
	}
	once := false
	for _, v := range infos {
		if !isSequence || !once {
			once = true
			if err := u.AppendInfoObjAddr(v.Ioa); err != nil {
				return err
			}
		}

		u.AppendFloat32(v.Value).AppendBytes(byte(v.Qds & 0xf1))
		switch typeID {
		case M_ME_NC_1:
		case M_ME_TC_1:
			u.AppendBytes(CP24Time2a(v.Time, u.InfoObjTimeZone)...)
		case M_ME_TF_1:
			u.AppendBytes(CP56Time2a(v.Time, u.InfoObjTimeZone)...)
		default:
			return ErrTypeIDNotMatch
		}
	}
	return c.Send(u)
}

func MeasuredValueFloat(c Connect, isSequence bool, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueFloatInfo) error {
	if !(coa.Cause == Percyc || coa.Cause == Back ||
		coa.Cause == Spont || coa.Cause == Req ||
		(coa.Cause >= Inrogen && coa.Cause <= Inro16)) {
		return ErrCmdCause
	}
	return measuredValueFloat(c, M_ME_NC_1, isSequence, coa, ca, infos...)
}

func MeasuredValueFloatCP24Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueFloatInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return measuredValueFloat(c, M_ME_TC_1, false, coa, ca, infos...)
}

func MeasuredValueFloatCP56Time2a(c Connect, coa CauseOfTransmission,
	ca CommonAddr, infos ...MeasuredValueFloatInfo) error {
	if !(coa.Cause == Spont || coa.Cause == Req) {
		return ErrCmdCause
	}
	return measuredValueFloat(c, M_ME_TF_1, false, coa, ca, infos...)
}

func (this *ASDU) GetSinglePoint() ([]SinglePointInfo, error) {
	var err error

	info := make([]SinglePointInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once := 0, false; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr = this.DecodeInfoObjAddr()
		} else {
			infoObjAddr++
		}
		value := this.DecodeByte()

		var t time.Time
		switch this.Type {
		case M_SP_NA_1:
		case M_SP_TA_1:
			if t, err = this.DecodeCP24Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_SP_TB_1:
			if t, err = this.DecodeCP56Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		default:
			return nil, ErrTypeIDNotMatch
		}

		info = append(info, SinglePointInfo{
			Ioa:   infoObjAddr,
			Value: value&0x01 == 0x01,
			Qds:   QualityDescriptor(value & 0xf0),
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetDoublePoint() ([]DoublePointInfo, error) {
	var err error

	info := make([]DoublePointInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once := 0, false; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr = this.DecodeInfoObjAddr()
		} else {
			infoObjAddr++
		}
		value := this.DecodeByte()

		var t time.Time
		switch this.Type {
		case M_DP_NA_1:
		case M_DP_TA_1:
			if t, err = this.DecodeCP24Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_DP_TB_1:
			if t, err = this.DecodeCP56Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		default:
			return nil, ErrTypeIDNotMatch
		}

		info = append(info, DoublePointInfo{
			Ioa:   infoObjAddr,
			Value: DoublePoint(value & 0x03),
			Qds:   QualityDescriptor(value & 0xf0),
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetStepPosition() ([]StepPositionInfo, error) {
	var err error

	info := make([]StepPositionInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once := 0, false; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr = this.DecodeInfoObjAddr()
		} else {
			infoObjAddr++
		}
		value := ParseStepPosition(this.DecodeByte())
		qds := QualityDescriptor(this.DecodeByte())

		var t time.Time
		switch this.Type {
		case M_ST_NA_1:
		case M_ST_TA_1:
			if t, err = this.DecodeCP24Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_SP_TB_1:
			if t, err = this.DecodeCP56Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		default:
			return nil, ErrTypeIDNotMatch
		}

		info = append(info, StepPositionInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetBitString32() ([]BitString32Info, error) {
	var err error

	info := make([]BitString32Info, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once := 0, false; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr = this.DecodeInfoObjAddr()
		} else {
			infoObjAddr++
		}

		value := this.DecodeBitsString32()
		qds := QualityDescriptor(this.DecodeByte())

		var t time.Time
		switch this.Type {
		case M_BO_NA_1:
		case M_BO_TA_1:
			if t, err = this.DecodeCP24Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_BO_TB_1:
			if t, err = this.DecodeCP56Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		default:
			return nil, ErrTypeIDNotMatch
		}

		info = append(info, BitString32Info{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetMeasuredValueNormal() ([]MeasuredValueNormalInfo, error) {
	var err error

	info := make([]MeasuredValueNormalInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once := 0, false; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr = this.DecodeInfoObjAddr()
		} else {
			infoObjAddr++
		}

		value := this.DecodeNormalize()

		var t time.Time
		var qds QualityDescriptor
		switch this.Type {
		case M_ME_NA_1:
			qds = QualityDescriptor(this.DecodeByte())
		case M_ME_TA_1:
			qds = QualityDescriptor(this.DecodeByte())
			if t, err = this.DecodeCP24Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_ME_TD_1:
			qds = QualityDescriptor(this.DecodeByte())
			if t, err = this.DecodeCP56Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_ME_ND_1: // 不带品质
		default:
			return nil, ErrTypeIDNotMatch
		}

		info = append(info, MeasuredValueNormalInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetMeasuredValueScaled() ([]MeasuredValueScaledInfo, error) {
	var err error

	info := make([]MeasuredValueScaledInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once := 0, false; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr = this.DecodeInfoObjAddr()
		} else {
			infoObjAddr++
		}

		value := this.DecodeScaled()
		qds := QualityDescriptor(this.DecodeByte())

		var t time.Time
		switch this.Type {
		case M_ME_NB_1:
		case M_ME_TB_1:
			if t, err = this.DecodeCP24Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_ME_TE_1:
			if t, err = this.DecodeCP56Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		default:
			return nil, ErrTypeIDNotMatch
		}

		info = append(info, MeasuredValueScaledInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   qds,
			Time:  t})
	}
	return info, nil
}

func (this *ASDU) GetMeasuredValueFloat() ([]MeasuredValueFloatInfo, error) {
	var err error

	info := make([]MeasuredValueFloatInfo, 0, this.Variable.Number)
	infoObjAddr := InfoObjAddr(0)
	for i, once := 0, false; i < int(this.Variable.Number); i++ {
		if !this.Variable.IsSequence || !once {
			once = true
			infoObjAddr = this.DecodeInfoObjAddr()
		} else {
			infoObjAddr++
		}

		value := this.DecodeFloat()
		qua := this.DecodeByte() & 0xf1

		var t time.Time
		switch this.Type {
		case M_ME_NC_1:
		case M_ME_TC_1:
			if t, err = this.DecodeCP24Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		case M_ME_TF_1:
			if t, err = this.DecodeCP56Time2a(); err != nil {
				return nil, ErrInvalidTimeTag
			}
		default:
			return nil, ErrTypeIDNotMatch
		}
		info = append(info, MeasuredValueFloatInfo{
			Ioa:   infoObjAddr,
			Value: value,
			Qds:   QualityDescriptor(qua),
			Time:  t})
	}
	return info, nil
}
