// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package cs101

// 采用FT1.2帧格式
const (
	startVarFrame byte = 0x68 // 长度可变帧启动字符
	startFixFrame byte = 0x10 // 长度固定帧启动字符
	endFrame      byte = 0x16
)

// 控制域定义
const (

	// 启动站到从动站特有
	FCV = 1 << 4 // 帧计数有效位
	FCB = 1 << 5 // 帧计数位
	// 从动站到启动站特有
	DFC     = 1 << 4 // 数据流控制位
	ACD_RES = 1 << 5 // 要求访问位,非平衡ACD,平衡保留
	// 启动报文位:
	// PRM = 0, 由从动站向启动站传输报文;
	// PRM = 1, 由启动站向从动站传输报文
	RPM     = 1 << 6
	RES_DIR = 1 << 7 // 非平衡保留,平衡为方向

	// 由启动站向从动站传输的报文中控制域的功能码(PRM = 1)
	FccResetRemoteLink                 = iota // 复位远方链路
	FccResetUserProcess                       // 复位用户进程
	FccBalanceTestLink                        // 链路测试功能
	FccUserDataWithConfirmed                  // 用户数据,需确认
	FccUserDataWithUnconfirmed                // 用户数据,无需确认
	_                                         // 保留
	_                                         // 制造厂和用户协商定义
	_                                         // 制造厂和用户协商定义
	FccUnbalanceWithRequestBitResponse        // 以要求访问位响应
	FccLinkStatus                             // 请求链路状态
	FccUnbalanceLevel1UserData                // 请求 1 级用户数据
	FccUnbalanceLevel2UserData                // 请求 2 级用户数据
	// 12-13: 备用
	// 14-15: 制造厂和用户协商定义

	// 从动站向启动站传输的报文中控制域的功能码(PRM = 0)
	FcsConfirmed                 = iota // 认可: 肯定认可
	FcsNConfirmed                       // 否定认可: 未收到报文,链路忙
	_                                   // 保留
	_                                   // 保留
	_                                   // 保留
	_                                   // 保留
	_                                   // 制造厂和用户协商定义
	_                                   // 制造厂和用户协商定义
	FcsUnbalanceResponse                // 用户数据
	FcsUnbalanceNegativeResponse        // 否定认哥: 无所召唤数据
	_                                   // 保留
	FcsStatus                           // 链路状态或要求访问
	// 12: 备用
	// 13: 制造厂和用户协商定义
	// 14: 链路服务未工作
	// 15: 链路服务未完成
)

// Ft12 ...
type Ft12 struct {
	start        byte
	apduFiledLen byte
	ctrl         byte
	address      uint16
	checksum     byte
	end          byte
}
