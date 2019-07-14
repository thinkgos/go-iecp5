package asdu

import (
	"encoding/binary"
	"time"
)

// 在控制方向系统信息的应用服务数据单元

// InterrogationCmd send a new interrogation command [C_IC_NA_1].
// subclass 7.3.4.1
// Use group 1 to 16, or 0 for the default.
// 总召唤命令
func InterrogationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	qoi QualifierOfInterrogation) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		C_IC_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddress(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(byte(qoi))
	return c.Send(u)
}

// QuantityInterrogationCmd send Quantity Interrogation command [C_CI_NA_1]
// subclass 7.3.4.2
// 计数量召唤命令
func QuantityInterrogationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	qcc QualifierCountCall) error {
	if coa.Cause != Act {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		C_CI_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddress(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(qcc.Value())
	return c.Send(u)
}

// [C_RD_NA_1]
// subclass 7.3.4.3
// 计数量召唤命令
func ReadCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	ioa InfoObjAddr) error {
	if !(coa.Cause == Req || coa.Cause == UnkType || coa.Cause == UnkCause ||
		coa.Cause == UnkAddr || coa.Cause == UnkInfo) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		C_RD_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddress(ioa); err != nil {
		return err
	}
	return c.Send(u)
}

// [C_CS_NA_1]
// subclass 7.3.4.4
// 时钟同步命令
func ClockSynchronizationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	t time.Time) error {
	if coa.Cause != Act {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		C_CS_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddress(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(CP56Time2a(t, u.InfoObjTimeZone)...)
	return c.Send(u)
}

// [C_TS_NA_1]
// subclass 7.3.4.5
// 测试命令
func TestCommand(c Connect, coa CauseOfTransmission, ca CommonAddr) error {
	if coa.Cause != Act {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		C_TS_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddress(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(byte(FBPTestWord&0xff), byte(FBPTestWord>>8))
	return c.Send(u)
}

// [C_RP_NA_1]
// subclass 7.3.4.6
// 复位进程命令
func ResetProcessCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	qrp QualifierOfResetProcessCmd) error {
	if coa.Cause != Act {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		C_RP_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddress(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(byte(qrp))
	return c.Send(u)
}

// [C_CD_NA_1]
// subclass 7.3.4.7
// 延时获得命令
func DelayAcquireCommand(c Connect, coa CauseOfTransmission, ca CommonAddr,
	msec uint16) error {
	if !(coa.Cause == Spont || coa.Cause == Act) {
		return ErrCmdCause
	}
	if err := c.Params().Valid(); err != nil {
		return err
	}

	u := NewASDU(c.Params(), Identifier{
		C_CD_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddress(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(CP16Time2a(msec)...)
	return c.Send(u)
}

func (this *ASDU) GetInterrogationCmd() QualifierOfInterrogation {
	// pass InfoObjAddr
	return QualifierOfInterrogation(this.infoObj[this.InfoObjAddrSize])
}

func (this *ASDU) GetQuantityInterrogationCmd() QualifierCountCall {
	// pass InfoObjAddr
	return ParseQualifierCountCall(this.infoObj[this.InfoObjAddrSize])
}

func (this *ASDU) GetReadCmd() InfoObjAddr {
	return this.DecodeInfoObjAddr()
}

func (this *ASDU) GetClockSynchronizationCmd() (time.Time, error) {
	// pass InfoObjAddr
	return ParseCP56Time2a(this.infoObj[this.InfoObjAddrSize:], this.InfoObjTimeZone)
}

func (this *ASDU) GetTestCommand() bool {
	// pass InfoObjAddr
	return binary.LittleEndian.Uint16(this.infoObj[this.InfoObjAddrSize:]) == FBPTestWord
}

func (this *ASDU) GetResetProcessCmd() QualifierOfResetProcessCmd {
	// pass InfoObjAddr
	return QualifierOfResetProcessCmd(this.infoObj[this.InfoObjAddrSize])
}

func (this *ASDU) GetDelayAcquireCommand() uint16 {
	// pass InfoObjAddr
	return binary.LittleEndian.Uint16(this.infoObj[this.InfoObjAddrSize:])
}
