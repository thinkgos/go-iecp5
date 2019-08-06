package asdu

// EndOfInitialization send a type identification M_EI_NA_1, subclass 7.3.3.1
// 初始化结束
func EndOfInitialization(c Connect, coa CauseOfTransmission, ca CommonAddr,
	ioa InfoObjAddr, coi CauseOfInitial) error {
	if err := c.Params().Valid(); err != nil {
		return err
	}

	coa.Cause = Initialized
	u := NewASDU(c.Params(), Identifier{
		M_EI_NA_1,
		VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})

	if err := u.AppendInfoObjAddr(ioa); err != nil {
		return err
	}
	u.AppendBytes(coi.Value())
	return c.Send(u)
}

// GetEndOfInitialization get GetEndOfInitialization for asud when the identification M_EI_NA_1
func (this *ASDU) GetEndOfInitialization() (InfoObjAddr, CauseOfInitial) {
	return this.DecodeInfoObjAddr(), ParseCauseOfInitial(this.infoObj[0])
}
