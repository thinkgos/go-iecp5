package asdu

import (
	"math"
)

// 在控制方向参数的应用服务数据单元

// [P_ME_NA_1]
// subclause 7.3.5.1
// 测量参数,规一化值
// TODO: check Normalize
func ParameterNormalizedValue(c Connect, coa CauseOfTransmission, ca CommonAddr,
	ioa InfoObjAddr, value Normalize, qpm QualifierOfParameterMV) error {
	if coa.Cause != Act {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		P_ME_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}

	u.infoObj = append(u.infoObj, byte(value), byte(value>>8), qpm.Value())
	return c.Send(u)
}

// [P_ME_NB_1]
// subclause 7.3.5.2
// 测量参数,标度化值
func ParameterScaledValue(c Connect, coa CauseOfTransmission, ca CommonAddr,
	ioa InfoObjAddr, value int16, qpm QualifierOfParameterMV) error {
	if coa.Cause != Act {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		P_ME_NB_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}

	u.infoObj = append(u.infoObj, byte(value), byte(value>>8), qpm.Value())
	return c.Send(u)
}

// [P_ME_NC_1]
// subclause 7.3.5.3
// 测量参数,短浮点数
func ParameterFloatValue(c Connect, coa CauseOfTransmission, ca CommonAddr,
	ioa InfoObjAddr, value float32, qpm QualifierOfParameterMV) error {
	if coa.Cause != Act {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		P_ME_NC_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	bits := math.Float32bits(value)
	u.infoObj = append(u.infoObj, byte(bits), byte(bits>>8), byte(bits>>16), byte(bits>>24), qpm.Value())
	return c.Send(u)
}

// [P_AC_NA_1]
// subclause 7.3.5.4
// 参数激活
func ParameterActivation(c Connect, coa CauseOfTransmission, ca CommonAddr,
	ioa InfoObjAddr, qpa QualifierOfParameterAct) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		P_AC_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.infoObj = append(u.infoObj, byte(qpa))
	return c.Send(u)
}
