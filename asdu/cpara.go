// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package asdu

// 在控制方向参数的应用服务数据单元

// ParameterNormalInfo 测量值参数,归一化值 信息体
type ParameterNormalInfo struct {
	Ioa   InfoObjAddr
	Value NormalizedMeasurement
	Qpm   QualifierOfParameterMV
}

// ParameterNormal 测量值参数,规一化值, 只有单个信息对象(SQ = 0)
// [P_ME_NA_1], See companion standard 101, subclass 7.3.5.1
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <20> := 响应站召唤
// <21> := 响应第 1 组召唤
// <22> := 响应第 2 组召唤
// 至
// <36> := 响应第 16 组召唤
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func ParameterNormal(c Connect, coa CauseOfTransmission, ca CommonAddr, p ParameterNormalInfo) error {
	if coa.Cause != Activation {
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
	if err := u.AppendInfoObjAddr(p.Ioa); err != nil {
		return err
	}
	u.AppendNormalize(p.Value)
	u.AppendBytes(p.Qpm.Value())
	return c.Send(u)
}

// ParameterScaledInfo 测量值参数,标度化值 信息体
type ParameterScaledInfo struct {
	Ioa   InfoObjAddr
	Value ScaledMeasurement
	Qpm   QualifierOfParameterMV
}

// ParameterScaled 测量值参数,标度化值, 只有单个信息对象(SQ = 0)
// [P_ME_NB_1], See companion standard 101, subclass 7.3.5.2
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <20> := 响应站召唤
// <21> := 响应第 1 组召唤
// <22> := 响应第 2 组召唤
// 至
// <36> := 响应第 16 组召唤
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func ParameterScaled(c Connect, coa CauseOfTransmission, ca CommonAddr, p ParameterScaledInfo) error {
	if coa.Cause != Activation {
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
	if err := u.AppendInfoObjAddr(p.Ioa); err != nil {
		return err
	}
	u.AppendScaled(p.Value).AppendBytes(p.Qpm.Value())
	return c.Send(u)
}

// ParameterFloatInfo 测量参数,短浮点数 信息体
type ParameterFloatInfo struct {
	Ioa   InfoObjAddr
	Value ShortFloatMeasurement
	Qpm   QualifierOfParameterMV
}

// ParameterFloat 测量值参数,短浮点数, 只有单个信息对象(SQ = 0)
// [P_ME_NC_1], See companion standard 101, subclass 7.3.5.3
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// 监视方向：
// <7> := 激活确认
// <20> := 响应站召唤
// <21> := 响应第 1 组召唤
// <22> := 响应第 2 组召唤
// 至
// <36> := 响应第 16 组召唤
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func ParameterFloat(c Connect, coa CauseOfTransmission, ca CommonAddr, p ParameterFloatInfo) error {
	if coa.Cause != Activation {
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
	if err := u.AppendInfoObjAddr(p.Ioa); err != nil {
		return err
	}
	u.AppendFloat32(p.Value).AppendBytes(p.Qpm.Value())
	return c.Send(u)
}

// ParameterActivationInfo 参数激活 信息体
type ParameterActivationInfo struct {
	Ioa InfoObjAddr
	Qpa QualifierOfParameterAct
}

// ParameterActivation 参数激活, 只有单个信息对象(SQ = 0)
// [P_AC_NA_1], See companion standard 101, subclass 7.3.5.4
// 传送原因(coa)用于
// 控制方向：
// <6> := 激活
// <8> := 停止激活
// 监视方向：
// <7> := 激活确认
// <9> := 停止激活确认
// <44> := 未知的类型标识
// <45> := 未知的传送原因
// <46> := 未知的应用服务数据单元公共地址
// <47> := 未知的信息对象地址
func ParameterActivation(c Connect, coa CauseOfTransmission, ca CommonAddr, p ParameterActivationInfo) error {
	if !(coa.Cause == Activation || coa.Cause == Deactivation) {
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
	if err := u.AppendInfoObjAddr(p.Ioa); err != nil {
		return err
	}
	u.AppendBytes(byte(p.Qpa))
	return c.Send(u)
}

// GetParameterNormal [P_ME_NA_1]，获取 测量值参数,标度化值 信息体
func (sf *ASDU) GetParameterNormal() ParameterNormalInfo {
	return ParameterNormalInfo{
		sf.DecodeInfoObjAddr(),
		sf.DecodeNormalize(),
		ParseQualifierOfParamMV(sf.infoObj[0]),
	}
}

// GetParameterScaled [P_ME_NB_1]，获取 测量值参数,归一化值 信息体
func (sf *ASDU) GetParameterScaled() ParameterScaledInfo {
	return ParameterScaledInfo{
		sf.DecodeInfoObjAddr(),
		sf.DecodeScaled(),
		ParseQualifierOfParamMV(sf.infoObj[0]),
	}
}

// GetParameterFloat [P_ME_NC_1]，获取 测量值参数,短浮点数 信息体
func (sf *ASDU) GetParameterFloat() ParameterFloatInfo {
	return ParameterFloatInfo{
		sf.DecodeInfoObjAddr(),
		sf.DecodeFloat32(),
		ParseQualifierOfParamMV(sf.infoObj[0]),
	}
}

// GetParameterActivation [P_AC_NA_1]，获取 参数激活 信息体
func (sf *ASDU) GetParameterActivation() ParameterActivationInfo {
	return ParameterActivationInfo{
		sf.DecodeInfoObjAddr(),
		QualifierOfParameterAct(sf.infoObj[0]),
	}
}
