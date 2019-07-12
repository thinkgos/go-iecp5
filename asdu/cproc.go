package asdu

import (
	"math"
	"time"
)

// 在控制方向过程信息的应用服务数据单元

type SingleCommand struct {
	Ioa   InfoObjAddr
	Value bool
	Qoc   QualifierOfCommand
	Time  *time.Time
}

// SingleCmd sends a type identification C_SC_NA_1 or C_SC_TA_1. subclause 7.3.2.1
// 单命令
func SingleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SingleCommand) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := checkValid(c, typeID, false, 1); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false},
		coa,
		0,
		ca,
	})
	if err := u.IncVariableNumber(1); err != nil {
		return err
	}

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

type DoubleCommand struct {
	Ioa   InfoObjAddr
	Value DoublePoint
	Qoc   QualifierOfCommand
	Time  *time.Time
}

// DoubleCmd sends a type identification C_DC_NA_1 or C_DC_TA_1. subclause 7.3.2.2
// 双命令
func DoubleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd DoubleCommand) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := checkValid(c, typeID, false, 1); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false},
		coa,
		0,
		commonAddr,
	})
	if err := u.IncVariableNumber(1); err != nil {
		return err
	}

	if err := u.AppendInfoObjAddr(cmd.Ioa); err != nil {
		return err
	}

	u.infoObj = append(u.infoObj, cmd.Qoc.Value()|cmd.Value.Value())
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

type StepCommand struct {
	Ioa   InfoObjAddr
	Value StepPosition
	Qoc   QualifierOfCommand
	Time  *time.Time
}

// StepCmd sends a type C_RC_NA_1 or C_RC_TA_1. subclause 7.3.2.3
// 步调节命令
func StepCmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd StepCommand) {
	panic("TODO: not implemented")
}

type SetpointNormalCommand struct {
	Ioa   InfoObjAddr
	Value Normalize
	Qoc   QualifierOfSetpointCmd
	Time  *time.Time
}

// SetpointCmdNormal sends a type C_SE_NA_1 or C_SE_TA_1. subclause 7.3.2.4
// 设定命令，归一化值
func SetpointCmdNormal(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd SetpointNormalCommand) {
	panic("TODO: not implemented")
}

type SetpointScaledCommand struct {
	Ioa   InfoObjAddr
	Value int16
	Qoc   QualifierOfSetpointCmd
	Time  *time.Time
}

// SetpointCmdScaled sends a type C_SE_NB_1 or C_SE_TB_1.  subclause 7.3.2.5
// 设定命令,标度化值
func SetpointCmdScaled(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd SetpointScaledCommand) {
	panic("TODO: not implemented")
}

type SetpointFloatCommand struct {
	Ioa   InfoObjAddr
	Value float32
	Qos   QualifierOfSetpointCmd
	Time  *time.Time
}

// SetpointCmdFloat sends a type C_SE_NC_1 or C_SE_TC_1.  subclause 7.3.2.6
// 设定命令,短浮点数
func SetpointCmdFloat(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd SetpointFloatCommand) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := checkValid(c, typeID, false, 1); err != nil {
		return err
	}
	u := NewASDU(c.Params(), Identifier{
		typeID,
		VariableStruct{IsSequence: false},
		coa,
		0,
		commonAddr,
	})
	if err := u.IncVariableNumber(1); err != nil {
		return err
	}

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

type BitsString32Command struct {
	Ioa   InfoObjAddr
	Value uint32
	Qos   QualifierOfSetpointCmd
	Time  *time.Time
}

// BitsString32Cmd sends a type C_BO_NA_1 or C_BO_TA_1.   subclause 7.3.2.7
// 比特串
func BitsString32Cmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd BitsString32Command) {
	panic("TODO: not implemented")
}

func (this *ASDU) GetSingleCmd() (SingleCommand, error) {
	var err error
	var s SingleCommand

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

func (this *ASDU) GetDoubleCmd() (DoubleCommand, error) {
	var err error
	var cmd DoubleCommand

	if cmd.Ioa, err = this.ParseInfoObjAddr(this.infoObj); err != nil {
		return cmd, err
	}
	value := this.infoObj[this.InfoObjAddrSize]
	cmd.Value = DoublePoint(value & 0x03)
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

func (this *ASDU) GetStepCmd() (StepCommand, error) {
	var cmd StepCommand
	return cmd, nil
}

func (this *ASDU) GetSetpointNormalCmd() (SetpointNormalCommand, error) {
	var cmd SetpointNormalCommand
	return cmd, nil
}

func (this *ASDU) GetSetpointCmdScaled() (SetpointScaledCommand, error) {
	var cmd SetpointScaledCommand
	return cmd, nil
}
func (this *ASDU) GetSetpointFloatCmd() (SetpointFloatCommand, error) {
	var cmd SetpointFloatCommand
	return cmd, nil
}

func (this *ASDU) GetBitsString32Cmd() (BitsString32Command, error) {
	var cmd BitsString32Command
	return cmd, nil
}
