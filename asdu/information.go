// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package asdu

// about information object 应用服务数据单元 - 信息对象

// InfoObjAddr is the information object address.
// See companion standard 101, subclass 7.2.5.
// The width is controlled by Params.InfoObjAddrSize.
// <0>: 无关的信息对象地址
// - width 1: <1..255>
// - width 2: <1..65535>
// - width 3: <1..16777215>
type InfoObjAddr uint

// InfoObjAddrIrrelevant Zero means that the information object address is irrelevant.
const InfoObjAddrIrrelevant InfoObjAddr = 0

// QualityDescriptor Quality descriptor flags attribute measured values.
// See companion standard 101, subclass 7.2.6.3.
type QualityDescriptor byte

// QualityDescriptor defined.
const (
	// QDSOverflow marks whether the value is beyond a predefined range.
	QDSOverflow QualityDescriptor = 1 << iota
	_                             // reserve
	_                             // reserve
	_                             // reserve
	// QDSBlocked flags that the value is blocked for transmission; the
	// value remains in the state that was acquired before it was blocked.
	QDSBlocked
	// QDSSubstituted flags that the value was provided by the input of
	// an operator (dispatcher) instead of an automatic source.
	QDSSubstituted
	// QDSNotTopical flags that the most recent update was unsuccessful.
	QDSNotTopical
	// QDSInvalid flags that the value was incorrectly acquired.
	QDSInvalid

	// QDSGood means no flags, no problems.
	QDSGood QualityDescriptor = 0
)

//QualityDescriptorProtection  Quality descriptor Protection Equipment flags attribute.
// See companion standard 101, subclass 7.2.6.4.
type QualityDescriptorProtection byte

// QualityDescriptorProtection defined.
const (
	_ QualityDescriptorProtection = 1 << iota // reserve
	_                                         // reserve
	_                                         // reserve
	// QDPElapsedTimeInvalid flags that the elapsed time was incorrectly acquired.
	QDPElapsedTimeInvalid
	// QDPBlocked flags that the value is blocked for transmission; the
	// value remains in the state that was acquired before it was blocked.
	QDPBlocked
	// QDPSubstituted flags that the value was provided by the input of
	// an operator (dispatcher) instead of an automatic source.
	QDPSubstituted
	// QDPNotTopical flags that the most recent update was unsuccessful.
	QDPNotTopical
	// QDPInvalid flags that the value was incorrectly acquired.
	QDPInvalid

	// QDPGood means no flags, no problems.
	QDPGood QualityDescriptorProtection = 0
)

// SinglePoint is a measured value of a switch.
// See companion standard 101, subclass 7.2.6.1.
type SinglePoint bool

// SinglePoint defined
const (
	SPIOff SinglePoint = false // 关
	SPIOn  SinglePoint = true  // 开
)

// Value single point to byte
func (sf SinglePoint) Value() bool {
	return bool(sf)
}

// ParseSinglePoint ...
func ParseSinglePoint(b byte) SinglePoint {
	return SinglePoint((b & 0x01) == 0x01)
}

// DoublePoint is a measured value of a determination aware switch.
// See companion standard 101, subclass 7.2.6.2.
type DoublePoint byte

// DoublePoint defined
const (
	DPIIndeterminateOrIntermediate DoublePoint = iota // 不确定或中间状态
	DPIDeterminedOff                                  // 确定状态开
	DPIDeterminedOn                                   // 确定状态关
	DPIIndeterminate                                  // 不确定或中间状态
)

// Value double point to byte
func (sf DoublePoint) Value() byte {
	return byte(sf & 0x03)
}

// ParseDoublePoint ...
func ParseDoublePoint(b byte) DoublePoint {
	return DoublePoint(b & 0x03)
}

// StepPosition is a measured value with transient state indication.
// 带瞬变状态指示的测量值，用于变压器步位置或其它步位置的值
// See companion standard 101, subclass 7.2.6.5.
// Val range <-64..63>
// bit[0-5]: <-64..63>
// NOTE: bit6 为符号位
// bit7: 0: 设备未在瞬变状态 1： 设备处于瞬变状态
type StepPosition struct {
	Val          int
	HasTransient bool
}

// Value returns step position value.
func (sf StepPosition) Value() byte {
	p := sf.Val & 0x7f
	if sf.HasTransient {
		p |= 0x80
	}
	return byte(p)
}

// ParseStepPosition parse byte to StepPosition.
func ParseStepPosition(b byte) StepPosition {
	step := StepPosition{HasTransient: (b & 0x80) != 0}
	if b&0x40 == 0 {
		step.Val = int(b & 0x3f)
	} else {
		step.Val = int(b) | (-1 &^ 0x3f)
	}
	return step
}

// BitString is 32 Bits String info
// 二进制状态信息
type BitString uint32

// Value ...
func (sf BitString) Value() uint32 {
	return uint32(sf)
}

// ParseBitString ...
func ParseBitString(b uint32) BitString {
	return BitString(b)
}

// NormalizedMeasurement is a 16-bit normalized value in[-1, 1 − 2⁻¹⁵]..
// 规一化值 f归一= 32768 * f真实 / 满码值
// See companion standard 101, subclass 7.2.6.6.
type NormalizedMeasurement int16

// Value ...
func (sf NormalizedMeasurement) Value() int16 {
	return int16(sf)
}

// NormalizedValue returns the value in [-1, 1 − 2⁻¹⁵].
func (sf NormalizedMeasurement) NormalizedValue() float64 {
	return float64(sf) / 32768
}

// ParseNormalizedMeasurement ...
func ParseNormalizedMeasurement(b int16) NormalizedMeasurement {
	return NormalizedMeasurement(b)
}

// ScaledMeasurement is a 16-bit scaled value in [-2¹⁵, +2¹⁵-1]
type ScaledMeasurement int16

// Value ...
func (sf ScaledMeasurement) Value() int16 {
	return int16(sf)
}

// ParseScaledMeasurement ...
func ParseScaledMeasurement(b int16) ScaledMeasurement {
	return ScaledMeasurement(b)
}

// ShortFloatMeasurement is a floa32 value
type ShortFloatMeasurement float32

// Value ...
func (sf ShortFloatMeasurement) Value() float32 {
	return float32(sf)
}

// ParseShortFloatMeasurement ...
func ParseShortFloatMeasurement(b float32) ShortFloatMeasurement {
	return ShortFloatMeasurement(b)
}

// BinaryCounterReading is binary counter reading
// See companion standard 101, subclass 7.2.6.9.
// CounterReading: 计数器读数 [bit0...bit31]
// SeqNumber: 顺序记法 [bit32...bit40]
// SQ: 顺序号 [bit32...bit36]
// CY: 进位 [bit37]
// CA: 计数量被调整
// IV: 无效
type BinaryCounterReading struct {
	CounterReading int32
	SeqNumber      byte
	HasCarry       bool
	IsAdjusted     bool
	IsInvalid      bool
}

// SingleEvent is single event
// See companion standard 101, subclass 7.2.6.10.
type SingleEvent byte

// SingleEvent dSequenceNotationefined
const (
	SEIndeterminateOrIntermediate SingleEvent = iota // 不确定或中间状态
	SEDeterminedOff                                  // 确定状态开
	SEDeterminedOn                                   // 确定状态关
	SEIndeterminate                                  // 不确定或中间状态
)

// StartEvent Start event protection
type StartEvent byte

// StartEvent defined
// See companion standard 101, subclass 7.2.6.11.
const (
	SEPGeneralStart          StartEvent = 1 << iota // 总启动
	SEPStartL1                                      // A相保护启动
	SEPStartL2                                      // B相保护启动
	SEPStartL3                                      // C相保护启动
	SEPStartEarthCurrent                            // 接地电流保护启动
	SEPStartReverseDirection                        // 反向保护启动
	// other reserved
)

// OutputCircuitInfo output command information
// See companion standard 101, subclass 7.2.6.12.
type OutputCircuitInfo byte

// OutputCircuitInfo defined
const (
	OCIGeneralCommand OutputCircuitInfo = 1 << iota // 总命令输出至输出电路
	OCICommandL1                                    // A 相保护命令输出至输出电路
	OCICommandL2                                    // B 相保护命令输出至输出电路
	OCICommandL3                                    // C 相保护命令输出至输出电路
	// other reserved
)

// FBPTestWord test special value
// See companion standard 101, subclass 7.2.6.14.
const FBPTestWord uint16 = 0x55aa

// SingleCommand Single command
// See companion standard 101, subclass 7.2.6.15.
type SingleCommand bool

// SingleCommand defined
const (
	SCOOn  SingleCommand = true
	SCOOff SingleCommand = false
)

// Value ...
func (sf SingleCommand) Value() bool {
	return bool(sf)
}

// ParseSingleCommand ...
func ParseSingleCommand(b byte) SingleCommand {
	if (b & 0x01) == 0x01 {
		return SingleCommand(true)
	}
	return SingleCommand(false)
}

// DoubleCommand double command
// See companion standard 101, subclass 7.2.6.16.
type DoubleCommand byte

// DoubleCommand defined
const (
	DCONotAllow0 DoubleCommand = iota
	DCOOn
	DCOOff
	DCONotAllow3
)

// Value ...
func (sf DoubleCommand) Value() byte {
	return byte(sf)
}

// ParseDoubleCommand ...
func ParseDoubleCommand(b byte) DoubleCommand {
	return DoubleCommand(b & 0x03)
}

// StepCommand step command
// See companion standard 101, subclass 7.2.6.17.
type StepCommand byte

// StepCommand defined
const (
	SCONotAllow0 StepCommand = iota
	SCOStepDown
	SCOStepUP
	SCONotAllow3
)

// Value ...
func (sf StepCommand) Value() byte {
	return byte(sf)
}

// ParseDoubleCommand ...
func ParseStepCommand(b byte) StepCommand {
	return StepCommand(b & 0x03)
}

// COICause Initialization reason
// See companion standard 101, subclass 7.2.6.21.
type COICause byte

// COICause defined
// 0: 当地电源合上
// 1： 当地手动复位
// 2： 远方复位
// <3..31>: 本配讨标准备的标准定义保留
// <32...127>: 为特定使用保留
const (
	COILocalPowerOn COICause = iota
	COILocalHandReset
	COIRemoteReset
)

// CauseOfInitial cause of initial
// Cause:  see COICause
// IsLocalChange: false - 未改变当地参数的初始化
//                true - 改变当地参数后的初始化
type CauseOfInitial struct {
	Cause         COICause
	IsLocalChange bool
}

// ParseCauseOfInitial parse byte to cause of initial
func ParseCauseOfInitial(b byte) CauseOfInitial {
	return CauseOfInitial{
		Cause:         COICause(b & 0x7f),
		IsLocalChange: b&0x80 == 0x80,
	}
}

// Value CauseOfInitial to byte
func (sf CauseOfInitial) Value() byte {
	if sf.IsLocalChange {
		return byte(sf.Cause | 0x80)
	}
	return byte(sf.Cause)
}

// QualifierOfInterrogation Qualifier Of Interrogation
// See companion standard 101, subclass 7.2.6.22.
type QualifierOfInterrogation byte

// QualifierOfInterrogation defined
const (
	// <1..19>: 为标准定义保留
	QOIStation QualifierOfInterrogation = 20 + iota // interrogated by station interrogation
	QOIGroup1                                       // interrogated by group 1 interrogation
	QOIGroup2                                       // interrogated by group 2 interrogation
	QOIGroup3                                       // interrogated by group 3 interrogation
	QOIGroup4                                       // interrogated by group 4 interrogation
	QOIGroup5                                       // interrogated by group 5 interrogation
	QOIGroup6                                       // interrogated by group 6 interrogation
	QOIGroup7                                       // interrogated by group 7 interrogation
	QOIGroup8                                       // interrogated by group 8 interrogation
	QOIGroup9                                       // interrogated by group 9 interrogation
	QOIGroup10                                      // interrogated by group 10 interrogation
	QOIGroup11                                      // interrogated by group 11 interrogation
	QOIGroup12                                      // interrogated by group 12 interrogation
	QOIGroup13                                      // interrogated by group 13 interrogation
	QOIGroup14                                      // interrogated by group 14 interrogation
	QOIGroup15                                      // interrogated by group 15 interrogation
	QOIGroup16                                      // interrogated by group 16 interrogation

	// <37..63>：为标准定义保留
	// <64..255>: 为特定使用保留

	// 0:未使用
	QOIUnused QualifierOfInterrogation = 0
)

// QCCRequest 请求 [bit0...bit5]
// See companion standard 101, subclass 7.2.6.23.
type QCCRequest byte

// QCCFreeze 冻结 [bit6,bit7]
// See companion standard 101, subclass 7.2.6.23.
type QCCFreeze byte

// QCCRequest and QCCFreeze defined
const (
	QCCUnused QCCRequest = iota
	QCCGroup1
	QCCGroup2
	QCCGroup3
	QCCGroup4
	QCCTotal
	// <6..31>: 为标准定义
	// <32..63>： 为特定使用保留
	QCCFrzRead          QCCFreeze = 0x00 // 读(无冻结或复位)
	QCCFrzFreezeNoReset QCCFreeze = 0x40 // 计数量冻结不带复位(被冻结的值为累计量)
	QCCFrzFreezeReset   QCCFreeze = 0x80 // 计数量冻结带复位(被冻结的值为增量信息)
	QCCFrzReset         QCCFreeze = 0xc0 // 计数量复位
)

// QualifierCountCall 计数量召唤命令限定词
// See companion standard 101, subclass 7.2.6.23.
type QualifierCountCall struct {
	Request QCCRequest
	Freeze  QCCFreeze
}

// ParseQualifierCountCall parse byte to QualifierCountCall
func ParseQualifierCountCall(b byte) QualifierCountCall {
	return QualifierCountCall{
		Request: QCCRequest(b & 0x3f),
		Freeze:  QCCFreeze(b & 0xc0),
	}
}

// Value QualifierCountCall to byte
func (sf QualifierCountCall) Value() byte {
	return byte(sf.Request&0x3f) | byte(sf.Freeze&0xc0)
}

// QPMCategory 测量参数类别
type QPMCategory byte

// QPMCategory defined
const (
	QPMUnused    QPMCategory = iota // 0: not used
	QPMThreshold                    // 1: threshold value
	QPMSmoothing                    // 2: smoothing factor (filter time constant)
	QPMLowLimit                     // 3: low limit for transmission of measured values
	QPMHighLimit                    // 4: high limit for transmission of measured values

	// 5‥31: reserved for standard definitions of sf companion standard (compatible range)
	// 32‥63: reserved for special use (private range)

	QPMChangeFlag      QPMCategory = 0x40 // bit6 marks local parameter change  当地参数改变
	QPMInOperationFlag QPMCategory = 0x80 // bit7 marks parameter operation 参数在运行
)

// QualifierOfParameterMV Qualifier Of Parameter Of Measured Values 测量值参数限定词
// See companion standard 101, subclass 7.2.6.24.
// QPMCategory : [bit0...bit5] 参数类型
// IsChange : [bit6]当地参数改变,false - 未改变,true - 改变
// IsInOperation : [bit7] 参数在运行,false - 运行, true - 不在运行
type QualifierOfParameterMV struct {
	Category      QPMCategory
	IsChange      bool
	IsInOperation bool
}

// ParseQualifierOfParamMV parse byte to QualifierOfParameterMV
func ParseQualifierOfParamMV(b byte) QualifierOfParameterMV {
	return QualifierOfParameterMV{
		Category:      QPMCategory(b & 0x3f),
		IsChange:      b&0x40 == 0x40,
		IsInOperation: b&0x80 == 0x80,
	}
}

// Value QualifierOfParameterMV to byte
func (sf QualifierOfParameterMV) Value() byte {
	v := byte(sf.Category) & 0x3f
	if sf.IsChange {
		v |= 0x40
	}
	if sf.IsInOperation {
		v |= 0x80
	}
	return v
}

// QualifierOfParameterAct Qualifier Of Parameter Activation 参数激活限定词
// See companion standard 101, subclass 7.2.6.25.
type QualifierOfParameterAct byte

// QualifierOfParameterAct defined
const (
	QPAUnused QualifierOfParameterAct = iota
	// 激活/停止激活这之前装载的参数(信息对象地址=0)
	QPADeActPrevLoadedParameter
	// 激活/停止激活所寻址信息对象的参数
	QPADeActObjectParameter
	// 激活/停止激活所寻址的持续循环或周期传输的信息对象
	QPADeActObjectTransmission
	// 4‥127: reserved for standard definitions of sf companion standard (compatible range)
	// 128‥255: reserved for special use (private range)
)

// QOCQual the qualifier of qual.
// See companion standard 101, subclass 7.2.6.26.
type QOCQual byte

// QOCQual defined
const (
	// 0: no additional definition
	// 无另外的定义
	QOCNoAdditionalDefinition QOCQual = iota
	// 1: short pulse duration (circuit-breaker), duration determined by a system parameter in the outstation
	// 短脉冲持续时间(断路器),持续时间由被控站内的系统参数所确定
	QOCShortPulseDuration
	// 2: long pulse duration, duration determined by a system parameter in the outstation
	// 长脉冲持续时间,持续时间由被控站内的系统参数所确定
	QOCLongPulseDuration
	// 3: persistent output
	// 持续输出
	QOCPersistentOutput
	//	4‥8: reserved for standard definitions of sf companion standard
	//	9‥15: reserved for the selection of other predefined functions
	//	16‥31: reserved for special use (private range)
)

// QualifierOfCommand is a  qualifier of command. 命令限定词
// See companion standard 101, subclass 7.2.6.26.
// See section 5, subclass 6.8.
// InSelect: true - selects, false - executes.
type QualifierOfCommand struct {
	Qual     QOCQual
	InSelect bool
}

// ParseQualifierOfCommand parse byte to QualifierOfCommand
func ParseQualifierOfCommand(b byte) QualifierOfCommand {
	return QualifierOfCommand{
		Qual:     QOCQual((b >> 2) & 0x1f),
		InSelect: b&0x80 == 0x80,
	}
}

// Value QualifierOfCommand to byte
func (sf QualifierOfCommand) Value() byte {
	v := (byte(sf.Qual) & 0x1f) << 2
	if sf.InSelect {
		v |= 0x80
	}
	return v
}

// QualifierOfResetProcessCmd 复位进程命令限定词
// See companion standard 101, subclass 7.2.6.27.
type QualifierOfResetProcessCmd byte

// QualifierOfResetProcessCmd defined
const (
	// 未采用
	QRPUnused QualifierOfResetProcessCmd = iota
	// 进程的总复位
	QPRGeneralRest
	// 复位事件缓冲区等待处理的带时标的信息
	QPRResetPendingInfoWithTimeTag
	// <3..127>: 为标准保留
	//<128..255>: 为特定使用保留
)

/*
TODO: file 文件相关未定义
*/

// QOSQual is the qualifier of a set-point command qual.
// See companion standard 101, subclass 7.2.6.39.
//	0: default
//	0‥63: reserved for standard definitions of sf companion standard (compatible range)
//	64‥127: reserved for special use (private range)
type QOSQual byte

// QualifierOfSetpointCmd is a qualifier of command. 设定命令限定词
// See section 5, subclass 6.8.
// InSelect: true - selects, false - executes.
type QualifierOfSetpointCmd struct {
	Qual     QOSQual
	InSelect bool
}

// ParseQualifierOfSetpointCmd parse byte to QualifierOfSetpointCmd
func ParseQualifierOfSetpointCmd(b byte) QualifierOfSetpointCmd {
	return QualifierOfSetpointCmd{
		Qual:     QOSQual(b & 0x7f),
		InSelect: b&0x80 == 0x80,
	}
}

// Value QualifierOfSetpointCmd to byte
func (sf QualifierOfSetpointCmd) Value() byte {
	v := byte(sf.Qual) & 0x7f
	if sf.InSelect {
		v |= 0x80
	}
	return v
}

// StatusAndStatusChangeDetection 状态和状态变位检出
// See companion standard 101, subclass 7.2.6.40.
type StatusAndStatusChangeDetection uint32
