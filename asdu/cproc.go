package asdu

import (
	"time"
)

// 在控制方向过程信息的应用服务数据单元

type SingleCommandObject struct {
	Ioa   InfoObjAddr
	Value bool
	Qoc   QualifierOfCommand
	Time  time.Time
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
	u.AppendBytes(value)
	switch typeID {
	case C_SC_NA_1:
	case C_SC_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type DoubleCommandObject struct {
	Ioa   InfoObjAddr
	Value DoubleCommand
	Qoc   QualifierOfCommand
	Time  time.Time
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

	u.AppendBytes(cmd.Qoc.Value() | byte(cmd.Value&0x03))
	switch typeID {
	case C_DC_NA_1:
	case C_DC_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type StepCommandObject struct {
	Ioa   InfoObjAddr
	Value StepCommand
	Qoc   QualifierOfCommand
	Time  time.Time
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

	u.AppendBytes(cmd.Qoc.Value() | byte(cmd.Value&0x03))
	switch typeID {
	case C_RC_NA_1:
	case C_RC_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type SetpointCommandNormalObject struct {
	Ioa   InfoObjAddr
	Value Normalize
	Qos   QualifierOfSetpointCmd
	Time  time.Time
}

// SetpointCmdNormal sends a type C_SE_NA_1 or C_SE_TA_1. subclass 7.3.2.4
// 设定命令，规一化值
func SetpointCmdNormal(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SetpointCommandNormalObject) error {
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
	u.AppendNormalize(cmd.Value).AppendBytes(cmd.Qos.Value())
	switch typeID {
	case C_SE_NA_1:
	case C_SE_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type SetpointCommandScaledObject struct {
	Ioa   InfoObjAddr
	Value int16
	Qos   QualifierOfSetpointCmd
	Time  time.Time
}

// SetpointCmdScaled sends a type C_SE_NB_1 or C_SE_TB_1.  subclass 7.3.2.5
// 设定命令,标度化值
func SetpointCmdScaled(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SetpointCommandScaledObject) error {
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
	u.AppendScaled(cmd.Value).AppendBytes(cmd.Qos.Value())
	switch typeID {
	case C_SE_NB_1:
	case C_SE_TB_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

type SetpointCommandFloatObject struct {
	Ioa   InfoObjAddr
	Value float32
	Qos   QualifierOfSetpointCmd
	Time  time.Time
}

// SetpointCmdFloat sends a type C_SE_NC_1 or C_SE_TC_1.  subclass 7.3.2.6
// 设定命令,短浮点数
func SetpointCmdFloat(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd SetpointCommandFloatObject) error {
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

	u.AppendFloat32(cmd.Value).AppendBytes(cmd.Qos.Value())

	switch typeID {
	case C_SE_NC_1:
	case C_SE_TC_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}

	return c.Send(u)
}

type BitsString32CommandObject struct {
	Ioa   InfoObjAddr
	Value uint32
	Time  time.Time
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
	if err := u.AppendInfoObjAddr(cmd.Ioa); err != nil {
		return err
	}

	u.AppendBitsString32(cmd.Value)

	switch typeID {
	case C_BO_NA_1:
	case C_BO_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}

	return c.Send(u)
}

func (this *ASDU) GetSingleCmd() (SingleCommandObject, error) {
	var err error
	var s SingleCommandObject

	s.Ioa = this.DecodeInfoObjAddr()
	value := this.infoObj[0]
	s.Value = value&0x01 == 0x01
	s.Qoc = ParseQualifierOfCommand(value & 0xfe)

	switch this.Type {
	case C_SC_NA_1:
	case C_SC_TA_1:
		s.Time, err = ParseCP56Time2a(this.infoObj[1:], this.InfoObjTimeZone)
		if err != nil {
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

	cmd.Ioa = this.DecodeInfoObjAddr()
	value := this.infoObj[0]
	cmd.Value = DoubleCommand(value & 0x03)
	cmd.Qoc = ParseQualifierOfCommand(value & 0xfc)

	switch this.Type {
	case C_DC_NA_1:
	case C_DC_TA_1:
		cmd.Time, err = ParseCP56Time2a(this.infoObj[1:], this.InfoObjTimeZone)
		if err != nil {
			return cmd, ErrInvalidTimeTag
		}
	default:
		return cmd, ErrTypeIDNotMatch
	}

	return cmd, nil
}

func (this *ASDU) GetStepCmd() (StepCommandObject, error) {
	var cmd StepCommandObject
	var err error

	cmd.Ioa = this.DecodeInfoObjAddr()
	value := this.infoObj[0]
	cmd.Value = StepCommand(value & 0x03)
	cmd.Qoc = ParseQualifierOfCommand(value & 0xfc)

	switch this.Type {
	case C_RC_NA_1:
	case C_RC_TA_1:
		cmd.Time, err = ParseCP56Time2a(this.infoObj[1:], this.InfoObjTimeZone)
		if err != nil {
			return cmd, ErrInvalidTimeTag
		}
	default:
		return cmd, ErrTypeIDNotMatch
	}

	return cmd, nil
}

func (this *ASDU) GetSetpointNormalCmd() (SetpointCommandNormalObject, error) {
	var cmd SetpointCommandNormalObject
	var err error

	cmd.Ioa = this.DecodeInfoObjAddr()
	cmd.Value = this.DecodeNormalize()
	cmd.Qos = ParseQualifierOfSetpointCmd(this.infoObj[0])

	switch this.Type {
	case C_SE_NA_1:
	case C_SE_TA_1:
		cmd.Time, err = ParseCP56Time2a(this.infoObj[1:], this.InfoObjTimeZone)
		if err != nil {
			return cmd, ErrInvalidTimeTag
		}
	default:
		return cmd, ErrTypeIDNotMatch
	}

	return cmd, nil
}

func (this *ASDU) GetSetpointCmdScaled() (SetpointCommandScaledObject, error) {
	var cmd SetpointCommandScaledObject
	var err error

	cmd.Ioa = this.DecodeInfoObjAddr()
	cmd.Value = this.DecodeScaled()
	cmd.Qos = ParseQualifierOfSetpointCmd(this.infoObj[0])

	switch this.Type {
	case C_SE_NB_1:
	case C_SE_TB_1:
		cmd.Time, err = ParseCP56Time2a(this.infoObj[1:], this.InfoObjTimeZone)
		if err != nil {
			return cmd, ErrInvalidTimeTag
		}
	default:
		return cmd, ErrTypeIDNotMatch
	}

	return cmd, nil
}
func (this *ASDU) GetSetpointFloatCmd() (SetpointCommandFloatObject, error) {
	var cmd SetpointCommandFloatObject
	var err error

	cmd.Ioa = this.DecodeInfoObjAddr()
	cmd.Value = this.DecodeFloat()
	cmd.Qos = ParseQualifierOfSetpointCmd(this.infoObj[0])

	switch this.Type {
	case C_SE_NC_1:
	case C_SE_TC_1:
		cmd.Time, err = ParseCP56Time2a(this.infoObj[1:], this.InfoObjTimeZone)
		if err != nil {
			return cmd, ErrInvalidTimeTag
		}
	default:
		return cmd, ErrTypeIDNotMatch
	}

	return cmd, nil
}

func (this *ASDU) GetBitsString32Cmd() (BitsString32CommandObject, error) {
	var cmd BitsString32CommandObject
	var err error

	cmd.Ioa = this.DecodeInfoObjAddr()
	cmd.Value = this.DecodeBitsString32()
	switch this.Type {
	case C_BO_NA_1:
	case C_BO_TA_1:
		cmd.Time, err = ParseCP56Time2a(this.infoObj, this.InfoObjTimeZone)
		if err != nil {
			return cmd, ErrInvalidTimeTag
		}
	default:
		return cmd, ErrTypeIDNotMatch
	}

	return cmd, nil
}
