package asdu

// about information object 应用服务数据单元 - 信息对象

// InfoObjAddr is the information object address.
// The width is controlled by Params.InfoObjAddrSize.
// See companion standard 101, subclause 7.2.5.
// width 1
// <0>: 无关的信息对象地址
// <1..255>: 信息对象地址
// width 2
// <0>: 无关的信息对象地址
// <1..65535>: 信息对象地址
// width 3
// <0>: 无关的信息对象地址
// <1..16777215>: 信息对象地址
type InfoObjAddr uint

// InfoObjIrrelevantAddr Zero means that the information object address is irrelevant.
const InfoObjIrrelevantAddr InfoObjAddr = 0

// SinglePoint is a measured value of a switch.
// See companion standard 101, subclause 7.2.6.1.
type SinglePoint uint

// 单点信息
const (
	Off SinglePoint = iota
	On
)

// DoublePoint is a measured value of a determination aware switch.
// See companion standard 101, subclause 7.2.6.2.
type DoublePoint uint

// 双点信息
const (
	IndeterminateOrIntermediate DoublePoint = iota // 不确定或中间状态
	DeterminedOff                                  // 确定状态开
	DeterminedOn                                   // 确定状态关
	Indeterminate                                  // 不确定或中间状态
)

// Quality descriptor flags attribute measured values.
// See companion standard 101, subclause 7.2.6.3.
const (
	// Overflow marks whether the value is beyond a predefined range.
	Overflow = 1 << iota

	_ // reserve
	_ // reserve

	// TimeInvalid flags that the elapsed time was incorrectly acquired.
	// This attribute is only valid for events of protection equipment.
	// See companion standard 101, subclause 7.2.6.4.
	TimeInvalid

	// Blocked flags that the value is blocked for transmission; the
	// value remains in the state that was acquired before it was blocked.
	Blocked

	// Substituted flags that the value was provided by the input of
	// an operator (dispatcher) instead of an automatic source.
	Substituted

	// NotTopical flags that the most recent update was unsuccessful.
	NotTopical

	// Invalid flags that the value was incorrectly acquired.
	Invalid

	// OK means no flags, no problems.
	OK = 0
)

// StepPos is a measured value with transient state indication.
// 带瞬变状态指示的测量值，用于变压器步位置或其它步位置的值
// See companion standard 101, subclause 7.2.6.5.
type StepPos int

// NewStepPos returns a new step position.
// Values range<-64..63>
// bit[0-6]: <-64..63>
// NOTE: bit6 为符号位
// bit7: 0: 设备未在瞬变状态 1： 设备处于瞬变状态
func NewStepPos(value int, hasTransient bool) StepPos {
	p := StepPos(value & 0x7f)
	if hasTransient {
		p |= 0x80
	}
	return p
}

// ToPos 返回 value in [-64, 63] 和 hasTransient 是否瞬变状态.
func (this StepPos) ToPos() (value int, hasTransient bool) {
	u := uint(this)
	if u&0x40 == 0 {
		value = int(u & 0x3f)
	} else {
		value = int(u) | (-1 &^ 0x3f)
	}
	hasTransient = (u & 0x80) != 0
	return
}

// Normalize is a 16-bit normalized value.
// 规一化值
// See companion standard 101, subclause 7.2.6.6.
type Normalize int16

// Float64 returns the value in [-1, 1 − 2⁻¹⁵].
func (this Normalize) Float64() float64 {
	return float64(this) / 32768
}

// Qualifier Of Parameter Of Measured Values
// 测量值参数限定词
// See companion standard 101, subclause 7.2.6.24.
const (
	_          = iota // 0: not used
	Threashold        // 1: threshold value
	Smoothing         // 2: smoothing factor (filter time constant)
	LowLimit          // 3: low limit for transmission of measured values
	HighLimit         // 4: high limit for transmission of measured values

	// 5‥31: reserved for standard definitions of this companion standard (compatible range)
	// 32‥63: reserved for special use (private range)

	ChangeFlag      = 64  // bit6 marks local parameter change  当地参数改变
	InOperationFlag = 128 // bit7 marks parameter operation 参数在运行
)

// Command is a command.
// 命令限定词
// See companion standard 101, subclause 7.2.6.26.
type Command byte

// <0>: 未用
// Qual returns the qualifier of command.
//
//	0: no additional definition
//	1: short pulse duration (circuit-breaker), duration determined by a system parameter in the outstation
//	2: long pulse duration, duration determined by a system parameter in the outstation
//	3: persistent output
//	4‥8: reserved for standard definitions of this companion standard
//	9‥15: reserved for the selection of other predefined functions
//	16‥31: reserved for special use (private range)
func (this Command) Qual() uint {
	return uint((this >> 2) & 0x1f)
}

// Exec 返回命令的 executes(false) (or selects(true)).
// See section 5, subclause 6.8.
func (this Command) Exec() bool {
	return this&0x80 == 0
}

// SetpointCmd is the qualifier of a set-point command.
// 设定命令限定词
// See companion standard 101, subclause 7.2.6.39.
type SetPointCmd uint

// Qual returns the qualifier of set-point command.
//
//	0: default
//	0‥63: reserved for standard definitions of this companion standard (compatible range)
//	64‥127: reserved for special use (private range)
func (this SetPointCmd) Qual() uint {
	return uint(this & 0x7f)
}

// Exec 返回命令的 executes(false) (or selects(true)).
// See section 5, subclause 6.8.
func (this SetPointCmd) Exec() bool {
	return this&0x80 == 0
}
