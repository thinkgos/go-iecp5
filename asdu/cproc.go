// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package asdu

import (
	"time"
)

// 在控制方向过程信息的应用服务数据单元

// SingleCommandInfo 单命令 信息体
type SingleCommandInfo struct {
	Ioa   InfoObjAddr
	Value SingleCommand
	Qoc   QualifierOfCommand
	Time  time.Time
}

// SingleCmd sends a type identification [C_SC_NA_1] or [C_SC_TA_1]. 单命令, 只有单个信息对象(SQ = 0)
// [C_SC_NA_1] See companion standard 101, subclass 7.3.2.1
// [C_SC_TA_1] See companion standard 101,
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func SingleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr, cmd SingleCommandInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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

// DoubleCommandInfo 单命令 信息体
type DoubleCommandInfo struct {
	Ioa   InfoObjAddr
	value DoubleCommand
	Qoc   QualifierOfCommand
	Time  time.Time
}

// DoubleCmd sends a type identification [C_DC_NA_1] or [C_DC_TA_1]. 双命令, 只有单个信息对象(SQ = 0)
// [C_DC_NA_1] See companion standard 101, subclass 7.3.2.2
// [C_DC_TA_1] See companion standard 101,
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func DoubleCmd(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr,
	cmd DoubleCommandInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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
	u.AppendValueAndQ(cmd.value, cmd.Qoc)
	switch typeID {
	case C_DC_NA_1:
	case C_DC_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

// StepCommandInfo 步调节 信息体
type StepCommandInfo struct {
	Ioa   InfoObjAddr
	Value StepCommand
	Qoc   QualifierOfCommand
	Time  time.Time
}

// StepCmd sends a type [C_RC_NA_1] or [C_RC_TA_1]. 步调节命令, 只有单个信息对象(SQ = 0)
// [C_RC_NA_1] See companion standard 101, subclass 7.3.2.3
// [C_RC_TA_1] See companion standard 101,
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func StepCmd(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr, cmd StepCommandInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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

	u.AppendValueAndQ(cmd.Value, cmd.Qoc)
	switch typeID {
	case C_RC_NA_1:
	case C_RC_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

// SetpointCommandNormalInfo 设置命令，规一化值 信息体
type SetpointCommandNormalInfo struct {
	Ioa   InfoObjAddr
	Value NormalizedMeasurement
	Qos   QualifierOfSetpointCmd
	Time  time.Time
}

// SetpointCmdNormal sends a type [C_SE_NA_1] or [C_SE_TA_1]. 设定命令,规一化值, 只有单个信息对象(SQ = 0)
// [C_SE_NA_1] See companion standard 101, subclass 7.3.2.4
// [C_SE_TA_1] See companion standard 101,
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func SetpointCmdNormal(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr, cmd SetpointCommandNormalInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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
	u.AppendValueAndQ(cmd.Value, cmd.Qos.Value())
	switch typeID {
	case C_SE_NA_1:
	case C_SE_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

// SetpointCommandScaledInfo 设定命令,标度化值 信息体
type SetpointCommandScaledInfo struct {
	Ioa   InfoObjAddr
	Value ScaledMeasurement
	Qos   QualifierOfSetpointCmd
	Time  time.Time
}

// SetpointCmdScaled sends a type [C_SE_NB_1] or [C_SE_TB_1]. 设定命令,标度化值,只有单个信息对象(SQ = 0)
// [C_SE_NB_1] See companion standard 101, subclass 7.3.2.5
// [C_SE_TB_1] See companion standard 101,
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func SetpointCmdScaled(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr, cmd SetpointCommandScaledInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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
	u.AppendValueAndQ(cmd.Value, cmd.Qos.Value())
	switch typeID {
	case C_SE_NB_1:
	case C_SE_TB_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}
	return c.Send(u)
}

// SetpointCommandFloatInfo 设定命令, 短浮点数 信息体
type SetpointCommandFloatInfo struct {
	Ioa   InfoObjAddr
	Value ShortFloatMeasurement
	Qos   QualifierOfSetpointCmd
	Time  time.Time
}

// SetpointCmdFloat sends a type [C_SE_NC_1] or [C_SE_TC_1].设定命令,短浮点数,只有单个信息对象(SQ = 0)
// [C_SE_NC_1] See companion standard 101, subclass 7.3.2.6
// [C_SE_TC_1] See companion standard 101,
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func SetpointCmdFloat(c Connect, typeID TypeID, coa CauseOfTransmission, ca CommonAddr, cmd SetpointCommandFloatInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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

	u.AppendValueAndQ(cmd.Value, cmd.Qos.Value())

	switch typeID {
	case C_SE_NC_1:
	case C_SE_TC_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}

	return c.Send(u)
}

// BitsString32CommandInfo 比特串命令 信息体
type BitsString32CommandInfo struct {
	Ioa   InfoObjAddr
	Value BitString
	Time  time.Time
}

// BitsString32Cmd sends a type [C_BO_NA_1] or [C_BO_TA_1]. 比特串命令,只有单个信息对象(SQ = 0)
// [C_BO_NA_1] See companion standard 101, subclass 7.3.2.7
// [C_BO_TA_1] See companion standard 101,
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func BitsString32Cmd(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	cmd BitsString32CommandInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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

	u.AppendValueAndQ(cmd.Value, QOCQual(0))

	switch typeID {
	case C_BO_NA_1:
	case C_BO_TA_1:
		u.AppendBytes(CP56Time2a(cmd.Time, u.InfoObjTimeZone)...)
	default:
		return ErrTypeIDNotMatch
	}

	return c.Send(u)
}

// GetSingleCmd [C_SC_NA_1] or [C_SC_TA_1] 获取单命令信息体
func (sf *ASDU) GetSingleCmd() SingleCommandInfo {
	var s SingleCommandInfo

	s.Ioa = sf.DecodeInfoObjAddr()
	value := sf.DecodeByte()
	s.Value = value&0x01 == 0x01
	s.Qoc = ParseQualifierOfCommand(value & 0xfe)

	switch sf.Type {
	case C_SC_NA_1:
	case C_SC_TA_1:
		s.Time = sf.DecodeCP56Time2a()
	default:
		panic(ErrTypeIDNotMatch)
	}

	return s
}

// GetDoubleCmd [C_DC_NA_1] or [C_DC_TA_1] 获取双命令信息体
func (sf *ASDU) GetDoubleCmd() DoubleCommandInfo {
	var cmd DoubleCommandInfo

	cmd.Ioa = sf.DecodeInfoObjAddr()
	value := sf.DecodeByte()
	cmd.value = DoubleCommand(value & 0x03)
	cmd.Qoc = ParseQualifierOfCommand(value & 0xfc)

	switch sf.Type {
	case C_DC_NA_1:
	case C_DC_TA_1:
		cmd.Time = sf.DecodeCP56Time2a()
	default:
		panic(ErrTypeIDNotMatch)
	}

	return cmd
}

// GetStepCmd [C_RC_NA_1] or [C_RC_TA_1] 获取步调节命令信息体
func (sf *ASDU) GetStepCmd() StepCommandInfo {
	var cmd StepCommandInfo

	cmd.Ioa = sf.DecodeInfoObjAddr()
	value := sf.DecodeByte()
	cmd.Value = StepCommand(value & 0x03)
	cmd.Qoc = ParseQualifierOfCommand(value & 0xfc)

	switch sf.Type {
	case C_RC_NA_1:
	case C_RC_TA_1:
		cmd.Time = sf.DecodeCP56Time2a()
	default:
		panic(ErrTypeIDNotMatch)
	}

	return cmd
}

// GetSetpointNormalCmd [C_SE_NA_1] or [C_SE_TA_1] 获取设定命令,规一化值信息体
func (sf *ASDU) GetSetpointNormalCmd() SetpointCommandNormalInfo {
	var cmd SetpointCommandNormalInfo

	cmd.Ioa = sf.DecodeInfoObjAddr()
	cmd.Value = sf.DecodeNormalize()
	cmd.Qos = ParseQualifierOfSetpointCmd(sf.DecodeByte())

	switch sf.Type {
	case C_SE_NA_1:
	case C_SE_TA_1:
		cmd.Time = sf.DecodeCP56Time2a()
	default:
		panic(ErrTypeIDNotMatch)
	}

	return cmd
}

// GetSetpointCmdScaled [C_SE_NB_1] or [C_SE_TB_1] 获取设定命令,标度化值信息体
func (sf *ASDU) GetSetpointCmdScaled() SetpointCommandScaledInfo {
	var cmd SetpointCommandScaledInfo

	cmd.Ioa = sf.DecodeInfoObjAddr()
	cmd.Value = sf.DecodeScaled()
	cmd.Qos = ParseQualifierOfSetpointCmd(sf.DecodeByte())

	switch sf.Type {
	case C_SE_NB_1:
	case C_SE_TB_1:
		cmd.Time = sf.DecodeCP56Time2a()
	default:
		panic(ErrTypeIDNotMatch)
	}

	return cmd
}

// GetSetpointFloatCmd [C_SE_NC_1] or [C_SE_TC_1] 获取设定命令，短浮点数信息体
func (sf *ASDU) GetSetpointFloatCmd() SetpointCommandFloatInfo {
	var cmd SetpointCommandFloatInfo

	cmd.Ioa = sf.DecodeInfoObjAddr()
	cmd.Value = sf.DecodeFloat32()
	cmd.Qos = ParseQualifierOfSetpointCmd(sf.DecodeByte())

	switch sf.Type {
	case C_SE_NC_1:
	case C_SE_TC_1:
		cmd.Time = sf.DecodeCP56Time2a()
	default:
		panic(ErrTypeIDNotMatch)
	}

	return cmd
}

// GetBitsString32Cmd [C_BO_NA_1] or [C_BO_TA_1] 获取比特串命令信息体
func (sf *ASDU) GetBitsString32Cmd() BitsString32CommandInfo {
	var cmd BitsString32CommandInfo

	cmd.Ioa = sf.DecodeInfoObjAddr()
	cmd.Value = sf.DecodeBitsString32()
	switch sf.Type {
	case C_BO_NA_1:
	case C_BO_TA_1:
		cmd.Time = sf.DecodeCP56Time2a()
	default:
		panic(ErrTypeIDNotMatch)
	}

	return cmd
}
