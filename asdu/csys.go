package asdu

import (
	"time"
)

// 在控制方向系统信息的应用服务数据单元

// InterrogationCmd send a new interrogation command [C_IC_NA_1].
// subclause 7.3.4.1
// Use group 1 to 16, or 0 for the default.
// 总召唤命令
func InterrogationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr, qoi QualifierOfInterrogation) error {
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
	u.infoObj = append(u.infoObj, byte(qoi))
	return c.Send(u)
}

// subclause 7.3.4.2
// 计数量召唤命令
func QuantityInterrogationCmd(c Connect, coa CauseOfTransmission, ca CommonAddr, qcc byte) error {

	return nil
}

// subclause 7.3.4.3
// 计数量召唤命令
func ReadCommand(c Connect, coa CauseOfTransmission, ca, ioa InfoObjAddr) error {
	return nil
}

// subclause 7.3.4.4
// 时钟同步命令
func ClockSynchronizationCommand(c Connect, coa CauseOfTransmission, ca CommonAddr, t time.Time) error {
	return nil
}

// subclause 7.3.4.5
// 测试命令
func TestCommand(c Connect, coa CauseOfTransmission, ca CommonAddr, v uint16) error {
	return nil
}

// subclause 7.3.4.6
// 复位进程命令
func ResetProcessCommand(c Connect, coa CauseOfTransmission, ca CommonAddr, grp byte) error {
	return nil
}

// subclause 7.3.4.7
// 延时获得命令
func DelayAcquireCommand(c Connect, coa CauseOfTransmission, ca CommonAddr, msec uint16) error {
	return nil
}
