package asdu

// 在控制方向参数的应用服务数据单元

// subclause 7.3.5.1
// 测量参数,规一化值
func ParameterNormalizedValue(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	ioa InfoObjAddr, value Normalize, qpm QualifierOfParameterMV) error {
	return nil
}

// subclause 7.3.5.2
// 测量参数,标度化值
func ParameterScaledValue(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	ioa InfoObjAddr, value int16, qpm QualifierOfParameterMV) error {
	return nil
}

func ParameterFloatValue(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	ioa InfoObjAddr, value float32, qpm QualifierOfParameterMV) error {
	return nil
}

func ParameterActivation(c Connect, typeID TypeID, coa CauseOfTransmission, commonAddr CommonAddr,
	ioa InfoObjAddr, qpm QualifierOfParameterAct) error {
	return nil
}
