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

// SinglePoint is a measured value of a switch.
// See companion standard 101, subclass 7.2.6.1.
type SinglePoint byte

// 单点信息
const (
	SPIOff SinglePoint = iota
	SPIOn
)

// Value single point to byte
func (this SinglePoint) Value() byte {
	return byte(this & 0x01)
}

// DoublePoint is a measured value of a determination aware switch.
// See companion standard 101, subclass 7.2.6.2.
type DoublePoint byte

// 双点信息
const (
	DPIIndeterminateOrIntermediate DoublePoint = iota // 不确定或中间状态
	DPIDeterminedOff                                  // 确定状态开
	DPIDeterminedOn                                   // 确定状态关
	DPIIndeterminate                                  // 不确定或中间状态
)

// Value double point to byte
func (this DoublePoint) Value() byte {
	return byte(this & 0x03)
}

// Quality descriptor flags attribute measured values.
// See companion standard 101, subclass 7.2.6.3.
type QualityDescriptor byte

// Quality descriptor flags attribute measured values.
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
	QDSGood = 0
)

// Quality descriptor Protection Equipment flags attribute.
// See companion standard 101, subclass 7.2.6.3.
type QualityDescriptorProtection byte

// Quality descriptor flags attribute Protection Equipment.
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
	QDPGood = 0
)

// StepPosition is a measured value with transient state indication.
// 带瞬变状态指示的测量值，用于变压器步位置或其它步位置的值
// See companion standard 101, subclass 7.2.6.5.
type StepPosition struct {
	Val          int
	HasTransient bool
}

// Value returns step position value.
// Values range<-64..63>
// bit[0-6]: <-64..63>
// NOTE: bit6 为符号位
// bit7: 0: 设备未在瞬变状态 1： 设备处于瞬变状态
func (this StepPosition) Value() byte {
	p := this.Val & 0x7f
	if this.HasTransient {
		p |= 0x80
	}
	return byte(p)
}

// ParseStepPosition 返回 val in [-64, 63] 和 HasTransient 是否瞬变状态.
func ParseStepPosition(b byte) StepPosition {
	step := StepPosition{HasTransient: (b & 0x80) != 0}
	if b&0x40 == 0 {
		step.Val = int(b & 0x3f)
	} else {
		step.Val = int(b) | (-1 &^ 0x3f)
	}
	return step
}

// Normalize is a 16-bit normalized value in[-1, 1 − 2⁻¹⁵]..
// 规一化值 f归一= 32768 * f真实 / 满码值
// See companion standard 101, subclass 7.2.6.6.
type Normalize int16

// Float64 returns the value in [-1, 1 − 2⁻¹⁵].
func (this Normalize) Float64() float64 {
	return float64(this) / 32768
}

// BinaryCounterReading is binary counter reading
// See companion standard 101, subclass 7.2.6.9.
type BinaryCounterReading struct {
	CounterReading   int32
	SequenceNotation byte
}

// SingleEvent is single event
// See companion standard 101, subclass 7.2.6.10.
type SingleEvent byte

// StartEvent Start event protection
type StartEvent byte

// Start event protection
// See companion standard 101, subclass 7.2.6.11.
const (
	SEPGeneralStart StartEvent = 1 << iota
	SEPStartL1
	SEPStartL2
	SEPStartL3
	SEPStartEarthCurrent
	SEPStartReverseDirection
	_
	_
)

// OutputCircuitInfo output command information
// See companion standard 101, subclass 7.2.6.12.
type OutputCircuitInfo byte

// output command information
const (
	OCIGeneralCommand = 1 << iota
	OCICommandL1
	OCICommandL2
	OCICommandL3
	// other reserved
)

// See companion standard 101, subclass 7.2.6.14.
const FBPTestWord uint16 = 0x55aa

/**************************************************/
// See companion standard 101, subclass 7.2.6.16.
type DoubleCommand byte

const (
	DCONotAllow0 DoubleCommand = iota
	DCOOn
	DCOOff
	DCONotAllow3
)

// See companion standard 101, subclass 7.2.6.17.
type StepCommand byte

const (
	SCONotAllow0 StepCommand = iota
	SCOStepDown
	SCOStepUP
	SCONotAllow3
)

// See companion standard 101, subclass 7.2.6.21.
// COICause Initialization reason
type COICause byte

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
func (this CauseOfInitial) Value() byte {
	if this.IsLocalChange {
		return byte(this.Cause | 0x80)
	}
	return byte(this.Cause)
}

// See companion standard 101, subclass 7.2.6.22.
// QualifierOfInterrogation Qualifier Of Interrogation
type QualifierOfInterrogation byte

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

// See companion standard 101, subclass 7.2.6.23.
type QCCRequest byte
type QCCFreeze byte

const (
	QCCUnused QCCRequest = iota
	QCCGroup1
	QCCGroup2
	QCCGroup3
	QCCGroup4
	QCCTotal
	// <6..31>: 为标准定义
	// <32..63>： 为特定使用保留
	QCCFrzRead          QCCFreeze = 0x00
	QCCFrzFreezeNoReset QCCFreeze = 0x40
	QCCFrzFreezeReset   QCCFreeze = 0x80
	QCCFrzReset         QCCFreeze = 0xc0
)

type QualifierCountCall struct {
	Request QCCRequest
	Freeze  QCCFreeze
}

func ParseQualifierCountCall(b byte) QualifierCountCall {
	return QualifierCountCall{
		Request: QCCRequest(b & 0x3f),
		Freeze:  QCCFreeze(b & 0xc0),
	}
}

// Value Qualifier Count Call to byte
func (this QualifierCountCall) Value() byte {
	return byte(this.Request&0x3f) | byte(this.Freeze&0xc0)
}

// See companion standard 101, subclass 7.2.6.24.
// QPMCategory 测量参数类别
type QPMCategory byte

const (
	QPMUnused    QPMCategory = iota // 0: not used
	QPMThreshold                    // 1: threshold value
	QPMSmoothing                    // 2: smoothing factor (filter time constant)
	QPMLowLimit                     // 3: low limit for transmission of measured values
	QPMHighLimit                    // 4: high limit for transmission of measured values

	// 5‥31: reserved for standard definitions of this companion standard (compatible range)
	// 32‥63: reserved for special use (private range)

	QPMChangeFlag      QPMCategory = 0x40 // bit6 marks local parameter change  当地参数改变
	QPMInOperationFlag QPMCategory = 0x80 // bit7 marks parameter operation 参数在运行
)

// QualifierOfParameterMV Qualifier Of Parameter Of Measured Values
// 测量值参数限定词
type QualifierOfParameterMV struct {
	Category      QPMCategory
	IsChange      bool
	IsInOperation bool
}

// ParseQualifierOfParamMV
func ParseQualifierOfParamMV(b byte) QualifierOfParameterMV {
	return QualifierOfParameterMV{
		Category:      QPMCategory(b & 0x3f),
		IsChange:      b&0x40 == 0x40,
		IsInOperation: b&0x80 == 0x80,
	}
}

// Value
func (this QualifierOfParameterMV) Value() byte {
	v := this.Category & 0x3f
	if this.IsChange {
		v |= 0x40
	}
	if this.IsInOperation {
		v |= 0x80
	}
	return byte(v)
}

// Qualifier Of Parameter Activation
// 参数激活限定词
// See companion standard 101, subclass 7.2.6.25.
type QualifierOfParameterAct byte

const (
	QPAUnused QualifierOfParameterAct = iota
	// 激活/仃止激活这之前装载的参数(信息对象地址=0)
	QPADeActPrevLoadedParameter
	// 激活/仃止激活所寻址信息对象的参数
	QPADeActObjectParameter
	// 激活/仃止激活所寻址的持续循环或周期传输的信息对象
	QPADeActObjectTransmission
	// 4‥127: reserved for standard definitions of this companion standard (compatible range)
	// 128‥255: reserved for special use (private range)
)

// QOCQual is a qualifier of qual.
// See companion standard 101, subclass 7.2.6.26.
//  the qualifier of command.
type QOCQual byte

const (
	//	0: no additional definition
	QOCNoAdditionalDefinition QOCQual = iota
	//	1: short pulse duration (circuit-breaker), duration determined by a system parameter in the outstation
	QOCShortPulseDuration
	//	2: long pulse duration, duration determined by a system parameter in the outstation
	QOCLongPulseDuration
	//	3: persistant output
	QOCPersistantOutput
	//	4‥8: reserved for standard definitions of this companion standard
	//	9‥15: reserved for the selection of other predefined functions
	//	16‥31: reserved for special use (private range)
)

// QualifierOfCommand is a  qualifier of command.
// 命令限定词
type QualifierOfCommand struct {
	Qual QOCQual
	// See section 5, subclass 6.8.
	// selects(true) (or executes(false)).
	InSelect bool
}

func ParseQualifierOfCommand(b byte) QualifierOfCommand {
	return QualifierOfCommand{
		Qual:     QOCQual((b >> 2) & 0x1f),
		InSelect: b&0x80 == 0x80,
	}
}

func (this QualifierOfCommand) Value() byte {
	v := (byte(this.Qual) & 0x1f) << 2
	if this.InSelect {
		v |= 0x80
	}
	return v
}

// See companion standard 101, subclass 7.2.6.27.
// 复位进程命令限定词
type QualifierOfResetProcessCmd byte

const (
	QRPUnused QualifierOfResetProcessCmd = iota
	QPRGeneralRest
	QPRResetPendingInfoWithTimeTag
	// <3..127>: 为标准保留
	//<128..255>: 为特定使用保留
)

// QOSQual is the qualifier of a set-point command qual.
// See companion standard 101, subclass 7.2.6.39.
//	0: default
//	0‥63: reserved for standard definitions of this companion standard (compatible range)
//	64‥127: reserved for special use (private range)
type QOSQual uint

// QualifierOfSetpointCmd is a qualifier of command.
type QualifierOfSetpointCmd struct {
	Qual QOSQual
	// See section 5, subclass 6.8.
	// selects(true) (or executes(false)).
	InSelect bool
}

func ParseQualifierOfSetpointCmd(b byte) QualifierOfSetpointCmd {
	return QualifierOfSetpointCmd{
		Qual:     QOSQual(b & 0x7f),
		InSelect: b&0x80 == 0x80,
	}
}

func (this QualifierOfSetpointCmd) Value() byte {
	v := byte(this.Qual) & 0x7f
	if this.InSelect {
		v |= 0x80
	}
	return v
}

// StatusAndStatusChangeDetection
// See companion standard 101, subclass 7.2.6.40.
type StatusAndStatusChangeDetection uint32

func (this StatusAndStatusChangeDetection) Value() []byte {
	return []byte{}
}
