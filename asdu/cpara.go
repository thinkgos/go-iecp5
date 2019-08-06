package asdu

// 在控制方向参数的应用服务数据单元

type ParameterNormalInfo struct {
	Ioa   InfoObjAddr
	Value Normalize
	Qpm   QualifierOfParameterMV
}

// [P_ME_NA_1]
// subclass 7.3.5.1
// 测量参数,规一化值
func ParameterNormal(c Connect, coa CauseOfTransmission, ca CommonAddr,
	p ParameterNormalInfo) error {
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

type ParameterScaledInfo struct {
	Ioa   InfoObjAddr
	Value int16
	Qpm   QualifierOfParameterMV
}

// [P_ME_NB_1]
// subclass 7.3.5.2
// 测量参数,标度化值
func ParameterScaled(c Connect, coa CauseOfTransmission, ca CommonAddr,
	p ParameterScaledInfo) error {
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

type ParameterFloatInfo struct {
	Ioa   InfoObjAddr
	Value float32
	Qpm   QualifierOfParameterMV
}

// [P_ME_NC_1]
// subclass 7.3.5.3
// 测量参数,短浮点数
func ParameterFloat(c Connect, coa CauseOfTransmission, ca CommonAddr,
	p ParameterFloatInfo) error {
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

type ParameterActivationInfo struct {
	Ioa InfoObjAddr
	Qpa QualifierOfParameterAct
}

// [P_AC_NA_1]
// subclass 7.3.5.4
// 参数激活
func ParameterActivation(c Connect, coa CauseOfTransmission, ca CommonAddr,
	p ParameterActivationInfo) error {
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

func (this *ASDU) GetParameterNormal() ParameterNormalInfo {
	return ParameterNormalInfo{
		this.DecodeInfoObjAddr(),
		this.DecodeNormalize(),
		ParseQualifierOfParamMV(this.infoObj[0]),
	}
}

func (this *ASDU) GetParameterScaled() ParameterScaledInfo {
	return ParameterScaledInfo{
		this.DecodeInfoObjAddr(),
		this.DecodeScaled(),
		ParseQualifierOfParamMV(this.infoObj[0]),
	}
}

func (this *ASDU) GetParameterFloat() ParameterFloatInfo {
	return ParameterFloatInfo{
		this.DecodeInfoObjAddr(),
		this.DecodeFloat(),
		ParseQualifierOfParamMV(this.infoObj[0]),
	}
}

func (this *ASDU) GetParameterActivation() ParameterActivationInfo {
	return ParameterActivationInfo{
		this.DecodeInfoObjAddr(),
		QualifierOfParameterAct(this.infoObj[0]),
	}
}
