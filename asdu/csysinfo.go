package asdu

// InterrogationCmd send a new interrogation command [C_IC_NA_1].
// subclause 7.3.4.1
// Use group 1 to 16, or 0 for the default.
// 总召唤
func InterrogationCmd(c Connect, coa CauseOfTransmission, commonAddr CommonAddr, group byte) error {
	if !(coa.Cause == Act || coa.Cause == Deact) {
		return errCmdCause
	}
	if err := checkValid(c, C_IC_NA_1, false, 1); err != nil {
		return err
	}

	u := NewASDU(c.Params(), C_IC_NA_1, false, coa, commonAddr)
	if err := u.IncVariableNumber(1); err != nil {
		return err
	}

	if err := u.AppendInfoObjAddr(InfoObjIrrelevantAddr); err != nil {
		return err
	}
	u.InfoObj = append(u.InfoObj, group+byte(Inrogen))
	return c.Send(u)
}
