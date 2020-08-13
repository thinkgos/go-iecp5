// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package asdu

// 在监视方向系统信息的应用服务数据单元

// EndOfInitialization send a type identification [M_EI_NA_1],初始化结束,只有单个信息对象(SQ = 0)
// [M_EI_NA_1] See companion standard 101,subclass 7.3.3.1
// 传送原因(coa)用于
// 监视方向：
// <4> := 被初始化
func EndOfInitialization(c Connect, coa CauseOfTransmission, ca CommonAddr, ioa InfoObjAddr, coi CauseOfInitial) error {
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

// GetEndOfInitialization get GetEndOfInitialization for asdu when the identification [M_EI_NA_1]
func (sf *ASDU) GetEndOfInitialization() (InfoObjAddr, CauseOfInitial) {
	return sf.DecodeInfoObjAddr(), ParseCauseOfInitial(sf.infoObj[0])
}
