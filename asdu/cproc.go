package asdu

import (
	"math"
	"time"
)

// 在控制方向过程信息的应用服务数据单元

type SingleCommandObject struct {
	Ioa   InfoObjAddr
	Value bool
	Qoc   QualifierOfCommand
	Time  *time.Time
}

// SingleCmd sends a type identification C_SC_NA_1 or C_SC_TA_1. subclass 7.3.2.1
// 单命令
func SingleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SingleCommandObject) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})

	if err := u.AppendInfoObjAddr(cmd.Ioa); err != nil {
		return err
	}
	value := cmd.Qoc.Value()
	if cmd.Value {
		value |= 0x01
	}
	u.infoObj = append(u.infoObj, value)
	switch typeID {
	case C_SC_NA_1:
	case C_SC_TA_1:
		if cmd.Time == nil {
			return ErrInvalidTimeTag
		}
		u.infoObj = append(u.infoObj, CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type DoubleCommandObject struct {
	Ioa   InfoObjAddr
	Value DoubleCommand
	Qoc   QualifierOfCommand
	Time  *time.Time
}

// DoubleCmd sends a type identification C_DC_NA_1 or C_DC_TA_1. subclass 7.3.2.2
// 双命令
func DoubleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd DoubleCommandObject) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})

	if err := u.AppendInfoObjAddr(cmd.Ioa); err != nil {
		return err
	}

	u.infoObj = append(u.infoObj, cmd.Qoc.Value()|byte(cmd.Value&0x03))
	switch typeID {
	case C_DC_NA_1:
	case C_DC_TA_1:
		if cmd.Time == nil {
			return ErrInvalidTimeTag
		}
		u.infoObj = append(u.infoObj, CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type StepCommandObject struct {
	Ioa   InfoObjAddr
	Value StepCommand
	Qoc   QualifierOfCommand
	Time  *time.Time
}

// StepCmd sends a type C_RC_NA_1 or C_RC_TA_1. subclass 7.3.2.3
// 步调节命令
func StepCmd(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd StepCommandObject) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})

	if err := u.AppendInfoObjAddr(cmd.Ioa); err != nil {
		return err
	}

	u.infoObj = append(u.infoObj, cmd.Qoc.Value()|byte(cmd.Value&0x03))
	switch typeID {
	case C_RC_NA_1:
	case C_RC_TA_1:
		if cmd.Time == nil {
			return ErrInvalidTimeTag
		}
		u.infoObj = append(u.infoObj, CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type SetpointNormalCommandObject struct {
	Ioa   InfoObjAddr
	Value Normalize
	Qoc   QualifierOfSetpointCmd
	Time  *time.Time
}

// SetpointCmdNormal sends a type C_SE_NA_1 or C_SE_TA_1. subclass 7.3.2.4
// 设定命令，规一化值
func SetpointCmdNormal(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SetpointNormalCommandObject) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})

	if err := u.AppendInfoObjAddr(cmd.Ioa); err != nil {
		return err
	}

	u.infoObj = append(u.infoObj, byte(cmd.Value), byte(cmd.Value>>8), cmd.Qoc.Value())
	switch typeID {
	case C_SE_NA_1:
	case C_SE_TA_1:
		if cmd.Time == nil {
			return ErrInvalidTimeTag
		}
		u.infoObj = append(u.infoObj, CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type SetpointScaledCommandObject struct {
	Ioa   InfoObjAddr
	Value int16
	Qoc   QualifierOfSetpointCmd
	Time  *time.Time
}

// SetpointCmdScaled sends a type C_SE_NB_1 or C_SE_TB_1.  subclass 7.3.2.5
// 设定命令,标度化值
func SetpointCmdScaled(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SetpointScaledCommandObject) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})

	if err := u.AppendInfoObjAddr(cmd.Ioa); err != nil {
		return err
	}

	u.infoObj = append(u.infoObj, byte(cmd.Value), byte(cmd.Value>>8), cmd.Qoc.Value())
	switch typeID {
	case C_SE_NB_1:
	case C_SE_TB_1:
		if cmd.Time == nil {
			return ErrInvalidTimeTag
		}
		u.infoObj = append(u.infoObj, CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type SetpointFloatCommandObject struct {
	Ioa   InfoObjAddr
	Value float32
	Qos   QualifierOfSetpointCmd
	Time  *time.Time
}

// SetpointCmdFloat sends a type C_SE_NC_1 or C_SE_TC_1.  subclass 7.3.2.6
// 设定命令,短浮点数
func SetpointCmdFloat(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SetpointFloatCommandObject) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})

	bits := math.Float32bits(cmd.Value)
	u.infoObj = append(u.infoObj, byte(bits), byte(bits>>8), byte(bits>>16), byte(bits>>24), cmd.Qos.Value())

	switch typeID {
	case C_SE_NC_1:
	case C_SE_TC_1:
		if cmd.Time == nil {
			return ErrInvalidTimeTag
		}
		u.infoObj = append(u.infoObj, CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)

	default:
		return ErrTypeIDNotMatch
	}

	return c.Send(u)
}

type BitsString32CommandObject struct {
	Ioa   InfoObjAddr
	Value uint32
	Qos   QualifierOfSetpointCmd
	Time  *time.Time
}

// BitsString32Cmd sends a type C_BO_NA_1 or C_BO_TA_1. subclass 7.3.2.7
// 比特串
func BitsString32Cmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd BitsString32CommandObject) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		commonAddr,
	})

	u.infoObj = append(u.infoObj, byte(cmd.Value), byte(cmd.Value>>8), byte(cmd.Value>>16), byte(cmd.Value>>24))

	switch typeID {
	case C_BO_NA_1:
	case C_BO_TA_1:
		if cmd.Time == nil {
			return ErrInvalidTimeTag
		}
		u.infoObj = append(u.infoObj, CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)

	default:
		return ErrTypeIDNotMatch
	}

	return c.Send(u)
}

func (this *ASDU) GetSingleCmd() (SingleCommandObject, error) {
	var err error
	var s SingleCommandObject

	if s.Ioa, err = this.ParseInfoObjAddr(this.infoObj); err != nil {
		return s, err
	}
	value := this.infoObj[this.InfoObjAddrSize]
	s.Value = value&0x01 == 0x01
	s.Qoc = ParseQualifierOfCommand(value & 0xfe)

	switch this.Type {
	case C_SC_NA_1:
	case C_SC_TA_1:
		s.Time = ParseCP56Time2a(this.infoObj[this.InfoObjAddrSize+1:], this.InfoObjTimeZone)
		if s.Time == nil {
			return s, ErrInvalidTimeTag
		}
	default:
		return s, ErrTypeIDNotMatch
	}

	return s, nil
}

func (this *ASDU) GetDoubleCmd() (DoubleCommandObject, error) {
	var err error
	var cmd DoubleCommandObject

	if cmd.Ioa, err = this.ParseInfoObjAddr(this.infoObj); err != nil {
		return cmd, err
	}
	value := this.infoObj[this.InfoObjAddrSize]
	cmd.Value = DoubleCommand(value & 0x03)
	cmd.Qoc = ParseQualifierOfCommand(value & 0xfc)

	switch this.Type {
	case C_SC_NA_1:
	case C_SC_TA_1:
		cmd.Time = ParseCP56Time2a(this.infoObj[this.InfoObjAddrSize+1:], this.InfoObjTimeZone)
		if cmd.Time == nil {
			return cmd, ErrInvalidTimeTag
		}
	default:
		return cmd, ErrTypeIDNotMatch
	}

	return cmd, nil
}

func (this *ASDU) GetStepCmd() (StepCommandObject, error) {
	var cmd StepCommandObject
	return cmd, nil
}

func (this *ASDU) GetSetpointNormalCmd() (SetpointNormalCommandObject, error) {
	var cmd SetpointNormalCommandObject
	return cmd, nil
}

func (this *ASDU) GetSetpointCmdScaled() (SetpointScaledCommandObject, error) {
	var cmd SetpointScaledCommandObject
	return cmd, nil
}
func (this *ASDU) GetSetpointFloatCmd() (SetpointFloatCommandObject, error) {
	var cmd SetpointFloatCommandObject
	return cmd, nil
}

func (this *ASDU) GetBitsString32Cmd() (BitsString32CommandObject, error) {
	var cmd BitsString32CommandObject
	return cmd, nil
}
