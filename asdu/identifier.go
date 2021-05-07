// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package asdu

import (
	"fmt"
	"strconv"
)

// about data unit identification 应用服务数据单元 - 数据单元标识符

// TypeID is the ASDU type identification.
// See companion standard 101, subclass 7.2.1.
type TypeID uint8

// The standard ASDU type identification.
// M for monitored information
// C for control information
// P for parameter
// F for file transfer.
// <0> 未用
// <1..127> 标准定义 - 兼容
// <128..135> 为路由报文保留 - 专用
// <136..255> 特殊应用 - 专用
// NOTE: 信息对象带或不带时标由标识符类型的不同序列来区别
const (
	_ TypeID = iota // 0: not defined
	// 在监视方向上的过程信息 <0..44>
	M_SP_NA_1 // 1: single-point information, 单点信息
	M_SP_TA_1 // 2: single-point information with time tag, 单点信息-带时标
	M_DP_NA_1 // 3: double-point information, 双点信息
	M_DP_TA_1 // 4: double-point information with time tag, 双点信息-带时标
	M_ST_NA_1 // 5: step position information, 步位置信息
	M_ST_TA_1 // 6: step position information with time tag, 步位置信息-带时标
	M_BO_NA_1 // 7: bitstring of 32 bit, 32位比特串
	M_BO_TA_1 // 8: bitstring of 32 bit with time tag, 32位比特串-带时标
	M_ME_NA_1 // 9: measured value, normalized value, 测量值，规一化值
	M_ME_TA_1 // 10: measured value, normalized value with time tag, 测量值，规一化值-带时标
	M_ME_NB_1 // 11: measured value, scaled value, 测量值，标度化值
	M_ME_TB_1 // 12: measured value, scaled value with time tag, 测量值带时标，标度化值-带时标
	M_ME_NC_1 // 13: measured value, short floating point number, 测量值，短浮点数
	M_ME_TC_1 // 14: measured value, short floating point number with time tag, 测量值，短浮数-带时标
	M_IT_NA_1 // 15: integrated totals, 累积量
	M_IT_TA_1 // 16: integrated totals with time tag, 累积量带时标
	M_EP_TA_1 // 17: event of protection equipment with time tag, 继电器保护设备事件-带时标
	M_EP_TB_1 // 18: packed start events of protection equipment with time tag, 继电保护设备成组启动事件-带时标
	M_EP_TC_1 // 19: packed output circuit information of protection equipment with time tag, 继电保护设备成组输出电路信息-带时标
	M_PS_NA_1 // 20: packed single-point information with status change detection, 带变位检出的成组单点信息
	M_ME_ND_1 // 21: measured value, normalized value without quality descriptor, 测量值,不带品质描述词的规一化值
	_         // 22: reserved for further compatible definitions
	_         // 23: reserved for further compatible definitions
	_         // 24: reserved for further compatible definitions
	_         // 25: reserved for further compatible definitions
	_         // 26: reserved for further compatible definitions
	_         // 27: reserved for further compatible definitions
	_         // 28: reserved for further compatible definitions
	_         // 29: reserved for further compatible definitions
	M_SP_TB_1 // 30: single-point information with time tag CP56Time2a, 单点信息-带CP56Time2a
	M_DP_TB_1 // 31: double-point information with time tag CP56Time2a, 双点信息-带CP56Time2a
	M_ST_TB_1 // 32: step position information with time tag CP56Time2a, 步位置信息-带CP56Time2a
	M_BO_TB_1 // 33: bitstring of 32 bits with time tag CP56Time2a, 32比特串-带CP56Time2a
	M_ME_TD_1 // 34: measured value, normalized value with time tag CP56Time2a, 测量值,规一化值-带CP56Time2a
	M_ME_TE_1 // 35: measured value, scaled value with time tag CP56Time2a, 测量值,标度化值-带CP56Time2a
	M_ME_TF_1 // 36: measured value, short floating point number with time tag CP56Time2a, 测量值,短浮点数-带CP56Time2a
	M_IT_TB_1 // 37: integrated totals with time tag CP56Time2a, 累积值-带CP56Time2a
	M_EP_TD_1 // 38: event of protection equipment with time tag CP56Time2a, 继电保护装置事件-带CP56Time2a
	M_EP_TE_1 // 39: packed start events of protection equipment with time tag CP56Time2a, 继电保护装置成组启动事件-带CP56Time2a
	M_EP_TF_1 // 40: packed output circuit information of protection equipment with time tag CP56Time2a, 继电保护装置成组输出电路信息-带CP56Time2a
	S_IT_TC_1 // 41: integrated totals containing time-tagged security statistics
	_         // 42: reserved for further compatible definitions
	_         // 43: reserved for further compatible definitions
	_         // 44: reserved for further compatible definitions
	// 在控制方向的过程信息 <45..69>
	C_SC_NA_1 // 45: single command 单点命令
	C_DC_NA_1 // 46: double command 双点命令
	C_RC_NA_1 // 47: regulating step command 调节步命令
	C_SE_NA_1 // 48: set-point command, normalized value 设定值命令，归一化值
	C_SE_NB_1 // 49: set-point command, scaled value 设定值命令，规度化值
	C_SE_NC_1 // 50: set-point command, short floating point number 设定值命令，短浮点数值
	C_BO_NA_1 // 51: bitstring of 32 bits 23位比特串
	_         // 52: reserved for further compatible definitions
	_         // 53: reserved for further compatible definitions
	_         // 54: reserved for further compatible definitions
	_         // 55: reserved for further compatible definitions
	_         // 56: reserved for further compatible definitions
	_         // 57: reserved for further compatible definitions
	C_SC_TA_1 // 58: single command with time tag CP56Time2a
	C_DC_TA_1 // 59: double command with time tag CP56Time2a
	C_RC_TA_1 // 60: regulating step command with time tag CP56Time2a
	C_SE_TA_1 // 61: set-point command with time tag CP56Time2a, normalized value
	C_SE_TB_1 // 62: set-point command with time tag CP56Time2a, scaled value
	C_SE_TC_1 // 63: set-point command with time tag CP56Time2a, short floating point number
	C_BO_TA_1 // 64: bitstring of 32-bit with time tag CP56Time2a
	_         // 65: reserved for further compatible definitions
	_         // 66: reserved for further compatible definitions
	_         // 67: reserved for further compatible definitions
	_         // 68: reserved for further compatible definitions
	_         // 69: reserved for further compatible definitions
	// 在监视方向的系统命令 <70..99>
	M_EI_NA_1 // 70: end of initialization  初始化结束
	_         // 71: reserved for further compatible definitions
	_         // 72: reserved for further compatible definitions
	_         // 73: reserved for further compatible definitions
	_         // 74: reserved for further compatible definitions
	_         // 75: reserved for further compatible definitions
	_         // 76: reserved for further compatible definitions
	_         // 77: reserved for further compatible definitions
	_         // 78: reserved for further compatible definitions
	_         // 79: reserved for further compatible definitions
	_         // 80: reserved for further compatible definitions
	S_CH_NA_1 // 81: authentication challenge
	S_RP_NA_1 // 82: authentication reply
	S_AR_NA_1 // 83: aggressive mode authentication request
	S_KR_NA_1 // 84: session key status request
	S_KS_NA_1 // 85: session key status
	S_KC_NA_1 // 86: session key change
	S_ER_NA_1 // 87: authentication error
	_         // 88: reserved for further compatible definitions
	_         // 89: reserved for further compatible definitions
	S_US_NA_1 // 90: user status change
	S_UQ_NA_1 // 91: update key change request
	S_UR_NA_1 // 92: update key change reply
	S_UK_NA_1 // 93: update key change — symetric
	S_UA_NA_1 // 94: update key change — asymetric
	S_UC_NA_1 // 95: update key change confirmation
	_         // 96: reserved for further compatible definitions
	_         // 97: reserved for further compatible definitions
	_         // 98: reserved for further compatible definitions
	_         // 99: reserved for further compatible definitions
	// 在控制方向的系统命令 <100..109>
	C_IC_NA_1 // 100: interrogation command 总召唤
	C_CI_NA_1 // 101: counter interrogation command 计数量召唤
	C_RD_NA_1 // 102: read command 读命令
	C_CS_NA_1 // 103: clock synchronization command 时钟同步命令
	C_TS_NA_1 // 104: test command 测试命令
	C_RP_NA_1 // 105: reset process command 复位进程命令
	C_CD_NA_1 // 106: delay acquisition command 延时获得命令
	C_TS_TA_1 // 107: test command with time tag CP56Time2a  带CP56Time2a的测试命令
	_         // 108: reserved for further compatible definitions
	_         // 109: reserved for further compatible definitions
	// 在控制方向的参数命令 <110..119>
	P_ME_NA_1 // 110: parameter of measured value, normalized value 测量值参数,规一化值
	P_ME_NB_1 // 111: parameter of measured value, scaled value 测量值参数,标度化值
	P_ME_NC_1 // 112: parameter of measured value, short floating point number 测量值参数,短浮点数
	P_AC_NA_1 // 113: parameter activation 参数激活
	_         // 114: reserved for further compatible definitions
	_         // 115: reserved for further compatible definitions
	_         // 116: reserved for further compatible definitions
	_         // 117: reserved for further compatible definitions
	_         // 118: reserved for further compatible definitions
	_         // 119: reserved for further compatible definitions
	// 文件传输 <120..127>
	F_FR_NA_1 // 120: file ready  文件准备就绪
	F_SR_NA_1 // 121: section ready 节准备就绪
	F_SC_NA_1 // 122: call directory, select file, call file, call section 如唤目录，选择文件，召唤文件，召唤节
	F_LS_NA_1 // 123: last section, last segment 最后的节，最后的段
	F_AF_NA_1 // 124: ack file, ack section 认可文件，认可节
	F_SG_NA_1 // 125: segment 段
	F_DR_TA_1 // 126: directory 目录
	F_SC_NB_1 // 127: QueryLog - request archive file (section 104) 查询日志
)

// infoObjSize maps the type identification (TypeID) to the serial octet size.
// Type extensions must register here.
var infoObjSize = map[TypeID]int{
	M_SP_NA_1: 1,
	M_SP_TA_1: 4,
	M_DP_NA_1: 1,
	M_DP_TA_1: 4,
	M_ST_NA_1: 2,
	M_ST_TA_1: 5,
	M_BO_NA_1: 5,
	M_BO_TA_1: 8,
	M_ME_NA_1: 3,
	M_ME_TA_1: 6,
	M_ME_NB_1: 3,
	M_ME_TB_1: 6,
	M_ME_NC_1: 5,
	M_ME_TC_1: 8,
	M_IT_NA_1: 5,
	M_IT_TA_1: 8,
	M_EP_TA_1: 6,
	M_EP_TB_1: 7,
	M_EP_TC_1: 7,
	M_PS_NA_1: 5,
	M_ME_ND_1: 2,

	M_SP_TB_1: 8,
	M_DP_TB_1: 8,
	M_ST_TB_1: 9,
	M_BO_TB_1: 12,
	M_ME_TD_1: 10,
	M_ME_TE_1: 10,
	M_ME_TF_1: 12,
	M_IT_TB_1: 12,
	M_EP_TD_1: 11,
	M_EP_TE_1: 11,
	M_EP_TF_1: 11,

	C_SC_NA_1: 1,
	C_DC_NA_1: 1,
	C_RC_NA_1: 1,
	C_SE_NA_1: 3,
	C_SE_NB_1: 3,
	C_SE_NC_1: 5,
	C_BO_NA_1: 4,

	M_EI_NA_1: 1,

	C_IC_NA_1: 1,
	C_CI_NA_1: 1,
	C_RD_NA_1: 0,
	C_CS_NA_1: 7,
	C_TS_NA_1: 2,
	C_RP_NA_1: 1,
	C_CD_NA_1: 2,

	P_ME_NA_1: 3,
	P_ME_NB_1: 3,
	P_ME_NC_1: 5,
	P_AC_NA_1: 1,

	F_FR_NA_1: 6,
	F_SR_NA_1: 7,
	F_SC_NA_1: 4,
	F_LS_NA_1: 5,
	F_AF_NA_1: 4,
	// F_SG_NA_1: 4 + variable,
	F_DR_TA_1: 13,
}

// GetInfoObjSize get the serial octet size of the type identification (TypeID).
func GetInfoObjSize(id TypeID) (int, error) {
	size, exists := infoObjSize[id]
	if !exists {
		return 0, ErrTypeIdentifier
	}
	return size, nil
}

const (
	_TypeIDName0 = "M_SP_NA_1M_SP_TA_1M_DP_NA_1M_DP_TA_1M_ST_NA_1M_ST_TA_1M_BO_NA_1M_BO_TA_1M_ME_NA_1M_ME_TA_1M_ME_NB_1M_ME_TB_1M_ME_NC_1M_ME_TC_1M_IT_NA_1M_IT_TA_1M_EP_TA_1M_EP_TB_1M_EP_TC_1M_PS_NA_1M_ME_ND_1"
	_TypeIDName1 = "M_SP_TB_1M_DP_TB_1M_ST_TB_1M_BO_TB_1M_ME_TD_1M_ME_TE_1M_ME_TF_1M_IT_TB_1M_EP_TD_1M_EP_TE_1M_EP_TF_1S_IT_TC_1"
	_TypeIDName2 = "C_SC_NA_1C_DC_NA_1C_RC_NA_1C_SE_NA_1C_SE_NB_1C_SE_NC_1C_BO_NA_1"
	_TypeIDName3 = "C_SC_TA_1C_DC_TA_1C_RC_TA_1C_SE_TA_1C_SE_TB_1C_SE_TC_1C_BO_TA_1"
	_TypeIDName4 = "M_EI_NA_1"
	_TypeIDName5 = "S_CH_NA_1S_RP_NA_1S_AR_NA_1S_KR_NA_1S_KS_NA_1S_KC_NA_1S_ER_NA_1"
	_TypeIDName6 = "S_US_NA_1S_UQ_NA_1S_UR_NA_1S_UK_NA_1S_UA_NA_1S_UC_NA_1"
	_TypeIDName7 = "C_IC_NA_1C_CI_NA_1C_RD_NA_1C_CS_NA_1C_TS_NA_1C_RP_NA_1C_CD_NA_1C_TS_TA_1"
	_TypeIDName8 = "P_ME_NA_1P_ME_NB_1P_ME_NC_1P_AC_NA_1"
	_TypeIDName9 = "F_FR_NA_1F_SR_NA_1F_SC_NA_1F_LS_NA_1F_AF_NA_1F_SG_NA_1F_DR_TA_1F_SC_NB_1"
)

func (sf TypeID) String() string {
	var s string
	switch {
	case 1 <= sf && sf <= 21:
		sf--
		s = _TypeIDName0[sf*9 : 9*(sf+1)]
	case 30 <= sf && sf <= 41:
		sf -= 30
		s = _TypeIDName1[sf*9 : 9*(sf+1)]
	case 45 <= sf && sf <= 51:
		sf -= 45
		s = _TypeIDName2[sf*9 : 9*(sf+1)]
	case 58 <= sf && sf <= 64:
		sf -= 58
		s = _TypeIDName3[sf*9 : 9*(sf+1)]
	case sf == 70:
		s = _TypeIDName4
	case 81 <= sf && sf <= 87:
		sf -= 81
		s = _TypeIDName5[sf*9 : 9*(sf+1)]
	case 90 <= sf && sf <= 95:
		sf -= 90
		s = _TypeIDName6[sf*9 : 9*(sf+1)]
	case 100 <= sf && sf <= 107:
		sf -= 100
		s = _TypeIDName7[sf*9 : 9*(sf+1)]
	case 110 <= sf && sf <= 113:
		sf -= 110
		s = _TypeIDName8[sf*9 : 9*(sf+1)]
	case 120 <= sf && sf <= 127:
		sf -= 120
		s = _TypeIDName9[sf*9 : 9*(sf+1)]
	default:
		s = strconv.FormatInt(int64(sf), 10)
	}
	return "TID<" + s + ">"
}

// VariableStruct is variable structure qualifier
// See companion standard 101, subclass 7.2.2.
// number <0..127>:  bit0 - bit6
// seq: bit7
// 0: 同一类型，有不同objAddress的信息元素集合 (地址+元素)*N
// 1： 同一类型，相同objAddress顺序信息元素集合 (一个地址,N元素*N)
type VariableStruct struct {
	Number     byte
	IsSequence bool
}

// ParseVariableStruct parse byte to variable structure qualifier
func ParseVariableStruct(b byte) VariableStruct {
	return VariableStruct{
		Number:     b & 0x7f,
		IsSequence: (b & 0x80) == 0x80,
	}
}

// Value encode variable structure to byte
func (sf VariableStruct) Value() byte {
	if sf.IsSequence {
		return sf.Number | 0x80
	}
	return sf.Number
}

// String 返回 variable structure 的格式
func (sf VariableStruct) String() string {
	if sf.IsSequence {
		return fmt.Sprintf("VSQ<sq,%d>", sf.Number)
	}
	return fmt.Sprintf("VSQ<%d>", sf.Number)
}

// CauseOfTransmission is the cause of transmission.
// See companion standard 101, subclass 7.2.3.
// | T | P/N | 5..0 cause |
// T = test, the cause of transmission for testing ,0: 未试验, 1：试验
// P/N indicates the negative (or positive) confirmation.
// Cause is the cause of transmission. bit5 - bit0
// 对由启动应用功能所请求的激活以肯定或者否定的确认 0: 肯定确认, 1: 否定确认
type CauseOfTransmission struct {
	IsTest     bool
	IsNegative bool
	Cause      Cause
}

// OriginAddr is originator address, See companion standard 101, subclass 7.2.3.
// The width is controlled by Params.CauseSize. width 2 includes/activates the originator address.
// <0>: 未用
// <1..255>: 源发地址
type OriginAddr byte

// Cause is the cause of transmission. bit5-bit0
type Cause byte

// Cause of transmission bit5-bit0
// <0> 未定义
// <1..63> 传输原因序号
// <1..47> 标准定义
// <48..63> 专用范围
// NOTE: 信息对象带或不带时标由标识符类型的不同序列来区别
const (
	Unused                  Cause = iota // unused
	Periodic                             // periodic, cyclic
	Background                           // background scan
	Spontaneous                          // spontaneous 突发
	Initialized                          // initialized
	Request                              // request or requested
	Activation                           // activation  激活
	ActivationCon                        // activation confirmation 激活确认
	Deactivation                         // deactivation 停止激活
	DeactivationCon                      // deactivation confirmation 停止激活确认
	ActivationTerm                       // activation termination 激活停止
	ReturnInfoRemote                     // return information caused by a remote command
	ReturnInfoLocal                      // return information caused by a local command
	FileTransfer                         // file transfer
	Authentication                       // authentication
	SessionKey                           // maintenance of authentication session key
	UserRoleAndUpdateKey                 // maintenance of user role and update key
	_                                    // reserved for further compatible definitions
	_                                    // reserved for further compatible definitions
	_                                    // reserved for further compatible definitions
	InterrogatedByStation                // interrogated by station interrogation
	InterrogatedByGroup1                 // interrogated by group 1 interrogation
	InterrogatedByGroup2                 // interrogated by group 2 interrogation
	InterrogatedByGroup3                 // interrogated by group 3 interrogation
	InterrogatedByGroup4                 // interrogated by group 4 interrogation
	InterrogatedByGroup5                 // interrogated by group 5 interrogation
	InterrogatedByGroup6                 // interrogated by group 6 interrogation
	InterrogatedByGroup7                 // interrogated by group 7 interrogation
	InterrogatedByGroup8                 // interrogated by group 8 interrogation
	InterrogatedByGroup9                 // interrogated by group 9 interrogation
	InterrogatedByGroup10                // interrogated by group 10 interrogation
	InterrogatedByGroup11                // interrogated by group 11 interrogation
	InterrogatedByGroup12                // interrogated by group 12 interrogation
	InterrogatedByGroup13                // interrogated by group 13 interrogation
	InterrogatedByGroup14                // interrogated by group 14 interrogation
	InterrogatedByGroup15                // interrogated by group 15 interrogation
	InterrogatedByGroup16                // interrogated by group 16 interrogation
	RequestByGeneralCounter              // requested by general counter request
	RequestByGroup1Counter               // requested by group 1 counter request
	RequestByGroup2Counter               // requested by group 2 counter request
	RequestByGroup3Counter               // requested by group 3 counter request
	RequestByGroup4Counter               // requested by group 4 counter request
	_                                    // reserved for further compatible definitions
	_                                    // reserved for further compatible definitions
	UnknownTypeID                        // unknown type identification
	UnknownCOT                           // unknown cause of transmission
	UnknownCA                            // unknown common address of ASDU
	UnknownIOA                           // unknown information object address
)

// Causal semantics description
var causeSemantics = []string{
	"Unused0",
	"Periodic",
	"Background",
	"Spontaneous",
	"Initialized",
	"Request",
	"Activation",
	"ActivationCon",
	"Deactivation",
	"DeactivationCon",
	"ActivationTerm",
	"ReturnInfoRemote",
	"ReturnInfoLocal",
	"FileTransfer",
	"Authentication",
	"SessionKey",
	"UserRoleAndUpdateKey",
	"Reserved17",
	"Reserved18",
	"Reserved19",
	"InterrogatedByStation",
	"InterrogatedByGroup1",
	"InterrogatedByGroup2",
	"InterrogatedByGroup3",
	"InterrogatedByGroup4",
	"InterrogatedByGroup5",
	"InterrogatedByGroup6",
	"InterrogatedByGroup7",
	"InterrogatedByGroup8",
	"InterrogatedByGroup9",
	"InterrogatedByGroup10",
	"InterrogatedByGroup11",
	"InterrogatedByGroup12",
	"InterrogatedByGroup13",
	"InterrogatedByGroup14",
	"InterrogatedByGroup15",
	"InterrogatedByGroup16",
	"RequestByGeneralCounter",
	"RequestByGroup1Counter",
	"RequestByGroup2Counter",
	"RequestByGroup3Counter",
	"RequestByGroup4Counter",
	"Reserved42",
	"Reserved43",
	"UnknownTypeID",
	"UnknownCOT",
	"UnknownCA",
	"UnknownIOA",
	"Special48",
	"Special49",
	"Special50",
	"Special51",
	"Special52",
	"Special53",
	"Special54",
	"Special55",
	"Special56",
	"Special57",
	"Special58",
	"Special59",
	"Special60",
	"Special61",
	"Special62",
	"Special63",
}

// ParseCauseOfTransmission parse byte to cause of transmission
func ParseCauseOfTransmission(b byte) CauseOfTransmission {
	return CauseOfTransmission{
		IsNegative: (b & 0x40) == 0x40,
		IsTest:     (b & 0x80) == 0x80,
		Cause:      Cause(b & 0x3f),
	}
}

// Value encode cause of transmission to byte
func (sf CauseOfTransmission) Value() byte {
	v := sf.Cause
	if sf.IsNegative {
		v |= 0x40
	}
	if sf.IsTest {
		v |= 0x80
	}
	return byte(v)
}

// String 返回Cause的字符串,包含相应应用的",neg" and ",test"
func (sf CauseOfTransmission) String() string {
	s := "COT<" + causeSemantics[sf.Cause]
	switch {
	case sf.IsNegative && sf.IsTest:
		s += ",neg,test"
	case sf.IsNegative:
		s += ",neg"
	case sf.IsTest:
		s += ",test"
	}
	return s + ">"
}

// CommonAddr is a station address.
// The width is controlled by Params.CommonAddrSize.
// width 1:
//      <0>: 未用
//      <1..254>: 站地址
//      <255>: 全局地址
// width 2:
//      <0>: 未用
//      <1..65534>: 站地址
//      <65535>: 全局地址
type CommonAddr uint16

// special commonAddr
const (
	// InvalidCommonAddr is the invalid common address.
	InvalidCommonAddr CommonAddr = 0
	// GlobalCommonAddr is the broadcast address. Use is restricted
	// to C_IC_NA_1, C_CI_NA_1, C_CS_NA_1 and C_RP_NA_1.
	// When in 8-bit mode 255 is mapped to this value on the fly.
	GlobalCommonAddr CommonAddr = 65535
)
