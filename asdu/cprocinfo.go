package asdu

import (
	"errors"
	"math"
	"time"
)

var errCmdCause = errors.New("asdu: cause of transmission for command not act(deact)")

// SingleCmd sends a type identification C_SC_NA_1 or C_SC_TA_1. subclause 7.3.2.1
// 单命令
func SingleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	infoObjAddr InfoObjAddr, cmd bool, qoc QualifierOfCmd, Time ...time.Time) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return errCmdCause
	}
	if err := checkValid(c, typeID, false, 1); err != nil {
		return err
	}

	u := NewASDU(c.Params(), typeID, false, coa, commonAddr)
	if err := u.IncVariableNumber(1); err != nil {
		return err
	}

	if err := u.AppendInfoObjAddr(infoObjAddr); err != nil {
		return err
	}
	value := qoc.Value()
	if cmd {
		value |= 0x01
	}
	u.InfoObj = append(u.InfoObj, value)
	switch typeID {
	case C_SC_NA_1:
	case C_SC_TA_1:
		panic("TODO: append 56-bit timestamp")
	default:
		return errType
	}
	return c.Send(u)
}

// DoubleCmd sends a type identification C_DC_NA_1 or C_DC_TA_1. subclause 7.3.2.2
// 双命令
func DoubleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	infoObjAddr InfoObjAddr, p DoublePoint, qoc QualifierOfCmd, Time ...time.Time) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return errCmdCause
	}
	if err := checkValid(c, typeID, false, 1); err != nil {
		return err
	}
	u := NewASDU(c.Params(), typeID, false, coa, commonAddr)
	if err := u.IncVariableNumber(1); err != nil {
		return err
	}

	if err := u.AppendInfoObjAddr(infoObjAddr); err != nil {
		return err
	}

	u.InfoObj = append(u.InfoObj, qoc.Value()|p.Value())
	switch typeID {
	case C_DC_NA_1:
	case C_DC_TA_1:
		panic("TODO: append 56-bit timestamp")
	default:
		return errType
	}
	return c.Send(u)
}

// StepCmd sends a type C_RC_NA_1 or C_RC_TA_1. subclause 7.3.2.3
// 步调节命令
func StepCmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	infoObjAddr InfoObjAddr, p StepPosition, qoc QualifierOfCmd, Time ...time.Time) {
	panic("TODO: not implemented")
}

// SetpointCmdNormal sends a type C_SE_NA_1 or C_SE_TA_1. subclause 7.3.2.4
// 设定命令，归一化值
func SetpointCmdNormal(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	infoObjAddr InfoObjAddr, p Normalize, qos QualifierOfSetpointCmd, Time ...time.Time) {
	panic("TODO: not implemented")
}

// SetpointCmdScaled sends a type C_SE_NB_1 or C_SE_TB_1.  subclause 7.3.2.5
// 设定命令,标度化值
func SetpointCmdScaled(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	infoObjAddr InfoObjAddr, p int16, qos QualifierOfSetpointCmd, Time ...time.Time) {
	panic("TODO: not implemented")
}

// SetpointCmdFloat sends a type C_SE_NC_1 or C_SE_TC_1.  subclause 7.3.2.6
// 设定命令,短浮点数
func SetpointCmdFloat(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	infoObjAddr InfoObjAddr, p float32, qos QualifierOfSetpointCmd, Time ...time.Time) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return errCmdCause
	}
	if err := checkValid(c, typeID, false, 1); err != nil {
		return err
	}
	u := NewASDU(c.Params(), typeID, false, coa, commonAddr)
	if err := u.IncVariableNumber(1); err != nil {
		return err
	}

	bits := math.Float32bits(p)
	u.InfoObj = append(u.InfoObj, byte(bits), byte(bits>>8), byte(bits>>16), byte(bits>>24), qos.Value())

	switch typeID {
	case C_SE_NC_1:
	case C_SE_TC_1:
		panic("TODO: append 56-bit timestamp")
	default:
		return errType
	}

	return c.Send(u)
}

// BitsString32Cmd sends a type C_BO_NA_1 or C_BO_TA_1.   subclause 7.3.2.7
// 比特串
func BitsString32Cmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	infoObjAddr InfoObjAddr, p uint32, qos QualifierOfSetpointCmd, Time ...time.Time) {
	panic("TODO: not implemented")
}
