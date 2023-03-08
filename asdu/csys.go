// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package asdu

import (
	"time"
)

// 在控制方向系统信息的应用服务数据单元

// InterrogationCmd send a new interrogation command [C_IC_NA_1]. 总召唤命令, 只有单个信息对象(SQ = 0)
// [C_IC_NA_1] See companion standard 101, subclass 7.3.4.1
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
func InterrogationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr, qoi QualifierOfInterrogation) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
		return ErrCmdCause
	}

	u := NewASDU(c.Params(), Identifier{
		C_IC_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendBytes(byte(qoi))
	return c.Send(u)
}

// CounterInterrogationCmd send Counter Interrogation command [C_CI_NA_1]，计数量召唤命令，只有单个信息对象(SQ = 0)
// [C_CI_NA_1] See companion standard 101, subclass 7.3.4.2
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func CounterInterrogationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr, qcc QualifierCountCall) error {
	coa.Cause = Activation
	u := NewASDU(c.Params(), Identifier{
		C_CI_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendBytes(qcc.Value())
	return c.Send(u)
}

// ReadCmd send read command [C_RD_NA_1], 读命令, 只有单个信息对象(SQ = 0)
// [C_RD_NA_1] See companion standard 101, subclass 7.3.4.3
// 传送原因(coa)用于
// 控制方向：
// <5> := 请求
// 监视方向：
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func ReadCmd(c Connect, ca CommonAddr, ioa InfoObjAddr) error {
	u := NewASDU(c.Params(), Identifier{
		C_RD_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		ParseCauseOfTransmission(byte(Request)),
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(ioa); err != nil {
		return err
	}
	return c.Send(u)
}

// ClockSynchronizationCmd send clock sync command [C_CS_NA_1],时钟同步命令, 只有单个信息对象(SQ = 0)
// [C_CS_NA_1] See companion standard 101, subclass 7.3.4.4
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <10> := 激活终止
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func ClockSynchronizationCmd(c Connect, ca CommonAddr, t time.Time) error {
	u := NewASDU(c.Params(), Identifier{
		C_CS_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		ParseCauseOfTransmission(byte(Activation)),
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendBytes(CP56Time2a(t, u.InfoObjTimeZone)...)
	return c.Send(u)
}

// TestCommand send test command [C_TS_NA_1]，测试命令, 只有单个信息对象(SQ = 0)
// [C_TS_NA_1] See companion standard 101, subclass 7.3.4.5
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func TestCommand(c Connect, ca CommonAddr) error {
	u := NewASDU(c.Params(), Identifier{
		C_TS_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		ParseCauseOfTransmission(byte(Activation)),
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendBytes(byte(FBPTestWord&0xff), byte(FBPTestWord>>8))
	return c.Send(u)
}

// ResetProcessCmd send reset process command [C_RP_NA_1],复位进程命令, 只有单个信息对象(SQ = 0)
// [C_RP_NA_1] See companion standard 101, subclass 7.3.4.6
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func ResetProcessCmd(c Connect, ca CommonAddr, qrp QualifierOfResetProcessCmd) error {
	u := NewASDU(c.Params(), Identifier{
		C_RP_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		ParseCauseOfTransmission(byte(Activation)),
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendBytes(byte(qrp))
	return c.Send(u)
}

// DelayAcquireCommand send delay acquire command [C_CD_NA_1],延时获得命令, 只有单个信息对象(SQ = 0)
// [C_CD_NA_1] See companion standard 101, subclass 7.3.4.7
// 传送原因(coa)用于
// 控制方向：
// <3> := 突发
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func DelayAcquireCommand(c Connect, coa CauseOfTransmission, ca CommonAddr, msec uint16) error {
	if !(coa.Cause == Spontaneous || coa.Cause == Activation) {
		return ErrCmdCause
	}

	u := NewASDU(c.Params(), Identifier{
		C_CD_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendCP16Time2a(msec)
	return c.Send(u)
}

// TestCommandCP56Time2a send test command [C_TS_TA_1]，测试命令, 只有单个信息对象(SQ = 0)
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func TestCommandCP56Time2a(c Connect, ca CommonAddr, t time.Time) error {
	u := NewASDU(c.Params(), Identifier{
		C_TS_TA_1,
		VariableStruct{IsSequence: false, Number: 1},
		ParseCauseOfTransmission(byte(Activation)),
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendUint16(FBPTestWord)
	u.AppendCP56Time2a(t, u.InfoObjTimeZone)
	return c.Send(u)
}

// GetInterrogationCmd [C_IC_NA_1] 获取总召唤信息体(信息对象地址，召唤限定词)
func (sf *ASDU) GetInterrogationCmd() (InfoObjAddr, QualifierOfInterrogation) {
	return sf.DecodeInfoObjAddr(), QualifierOfInterrogation(sf.infoObj[0])
}

// GetCounterInterrogationCmd [C_CI_NA_1] 获得计量召唤信息体(信息对象地址，计量召唤限定词)
func (sf *ASDU) GetCounterInterrogationCmd() (InfoObjAddr, QualifierCountCall) {
	return sf.DecodeInfoObjAddr(), ParseQualifierCountCall(sf.infoObj[0])
}

// GetReadCmd [C_RD_NA_1] 获得读命令信息地址
func (sf *ASDU) GetReadCmd() InfoObjAddr {
	return sf.DecodeInfoObjAddr()
}

// GetClockSynchronizationCmd [C_CS_NA_1] 获得时钟同步命令信息体(信息对象地址,时间)
func (sf *ASDU) GetClockSynchronizationCmd() (InfoObjAddr, time.Time) {

	return sf.DecodeInfoObjAddr(), sf.DecodeCP56Time2a()
}

// GetTestCommand [C_TS_NA_1]，获得测试命令信息体(信息对象地址,是否是测试字)
func (sf *ASDU) GetTestCommand() (InfoObjAddr, bool) {
	return sf.DecodeInfoObjAddr(), sf.DecodeUint16() == FBPTestWord
}

// GetResetProcessCmd [C_RP_NA_1] 获得复位进程命令信息体(信息对象地址,复位进程命令限定词)
func (sf *ASDU) GetResetProcessCmd() (InfoObjAddr, QualifierOfResetProcessCmd) {
	return sf.DecodeInfoObjAddr(), QualifierOfResetProcessCmd(sf.infoObj[0])
}

// GetDelayAcquireCommand [C_CD_NA_1] 获取延时获取命令信息体(信息对象地址,延时毫秒数)
func (sf *ASDU) GetDelayAcquireCommand() (InfoObjAddr, uint16) {
	return sf.DecodeInfoObjAddr(), sf.DecodeUint16()
}

// GetTestCommandCP56Time2a [C_TS_TA_1]，获得测试命令信息体(信息对象地址,是否是测试字)
func (sf *ASDU) GetTestCommandCP56Time2a() (InfoObjAddr, bool, time.Time) {
	return sf.DecodeInfoObjAddr(), sf.DecodeUint16() == FBPTestWord, sf.DecodeCP56Time2a()
}
