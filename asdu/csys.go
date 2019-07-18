package asdu

import (
	"encoding/binary"
	"time"
)

// 在控制方向系统信息的应用服务数据单元

// InterrogationCmd send a new interrogation command [C_IC_NA_1].
// coa.Cause = Act or Deact
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
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(byte(qoi))
	return c.Send(u)
}

// CounterInterrogationCmd send Counter Interrogation command [C_CI_NA_1]
// coa.Cause always Act
// subclass 7.3.4.2
// 计数量召唤命令
func CounterInterrogationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	qcc QualifierCountCall) error {
	if err := c.Params().Valid(); err != nil {
		return err
	}
	coa.Cause = Act
	u := NewASDU(c.Params(), Identifier{
		C_CI_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(qcc.Value())
	return c.Send(u)
}

// ReadCmd  ,[C_RD_NA_1]
// coa.Cause always Req
// subclass 7.3.4.3
// 计数量召唤命令
func ReadCmd(c Connect, coa CauseOfTransmission, ca CommonAddr, ioa InfoObjAddr) error {
	if err := c.Params().Valid(); err != nil {
		return err
	}
	coa.Cause = Req
	u := NewASDU(c.Params(), Identifier{
		C_RD_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(ioa); err != nil {
		return err
	}
	return c.Send(u)
}

// ClockSynchronizationCmd [C_CS_NA_1]
// coa.Cause always Act
// subclass 7.3.4.4
// 时钟同步命令
func ClockSynchronizationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	t time.Time) error {
	if err := c.Params().Valid(); err != nil {
		return err
	}
	coa.Cause = Act
	u := NewASDU(c.Params(), Identifier{
		C_CS_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(CP56Time2a(t, u.InfoObjTimeZone)...)
	return c.Send(u)
}

// TestCommand [C_TS_NA_1]
// coa.Cause always Act(6)
// subclass 7.3.4.5
// 测试命令
func TestCommand(c Connect, coa CauseOfTransmission, ca CommonAddr) error {
	if err := c.Params().Valid(); err != nil {
		return err
	}
	coa.Cause = Req
	u := NewASDU(c.Params(), Identifier{
		C_TS_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(byte(FBPTestWord&0xff), byte(FBPTestWord>>8))
	return c.Send(u)
}

// ResetProcessCmd [C_RP_NA_1]
// coa.Cause always Act(6)
// subclass 7.3.4.6
// 复位进程命令
func ResetProcessCmd(c Connect, coa CauseOfTransmission, ca CommonAddr,
	qrp QualifierOfResetProcessCmd) error {
	if err := c.Params().Valid(); err != nil {
		return err
	}
	coa.Cause = Req
	u := NewASDU(c.Params(), Identifier{
		C_RP_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(byte(qrp))
	return c.Send(u)
}

// DelayAcquireCommand [C_CD_NA_1]
// coa.Cause = Act or Spont
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
	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.AppendBytes(CP16Time2a(msec)...)
	return c.Send(u)
}

func (this *ASDU) GetInterrogationCmd() (InfoObjAddr, QualifierOfInterrogation) {
	return this.DecodeInfoObjAddr(), QualifierOfInterrogation(this.infoObj[0])
}

func (this *ASDU) GetCounterInterrogationCmd() (InfoObjAddr, QualifierCountCall) {
	return this.DecodeInfoObjAddr(), ParseQualifierCountCall(this.infoObj[0])
}

func (this *ASDU) GetReadCmd() InfoObjAddr {
	return this.DecodeInfoObjAddr()
}

func (this *ASDU) GetClockSynchronizationCmd() (InfoObjAddr, time.Time) {
	return this.DecodeInfoObjAddr(), ParseCP56Time2a(this.infoObj, this.InfoObjTimeZone)
}

func (this *ASDU) GetTestCommand() (InfoObjAddr, bool) {
	return this.DecodeInfoObjAddr(), binary.LittleEndian.Uint16(this.infoObj) == FBPTestWord
}

func (this *ASDU) GetResetProcessCmd() (InfoObjAddr, QualifierOfResetProcessCmd) {

	return this.DecodeInfoObjAddr(), QualifierOfResetProcessCmd(this.infoObj[0])
}

func (this *ASDU) GetDelayAcquireCommand() (InfoObjAddr, uint16) {
	return this.DecodeInfoObjAddr(), binary.LittleEndian.Uint16(this.infoObj)
}
