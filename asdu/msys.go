// Copyright [2020] [thinkgos]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
