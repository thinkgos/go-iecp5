// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

package cs104

import (
	"fmt"

	"github.com/thinkgos/go-iecp5/asdu"
)

const startFrame byte = 0x68 // 启动字符

// APDU form Max size 255
//      |              APCI                   |       ASDU         |
//      | start | APDU length | control field |       ASDU         |
//                       |          APDU field size(253)           |
// bytes|    1  |    1   |        4           |                    |
const (
	APCICtlFiledSize = 4 // control filed(4)

	APDUSizeMax      = 255                                 // start(1) + length(1) + control field(4) + ASDU
	APDUFieldSizeMax = APCICtlFiledSize + asdu.ASDUSizeMax // control field(4) + ASDU
)

// U帧 控制域功能
const (
	uStartDtActive  byte = 4 << iota // 启动激活 0x04
	uStartDtConfirm                  // 启动确认 0x08
	uStopDtActive                    // 停止激活 0x10
	uStopDtConfirm                   // 停止确认 0x20
	uTestFrActive                    // 测试激活 0x40
	uTestFrConfirm                   // 测试确认 0x80
)

// I帧 含apci和asdu 信息帧.用于编号的信息传输 information
type iAPCI struct {
	sendSN, rcvSN uint16
}

func (sf iAPCI) String() string {
	return fmt.Sprintf("I[sendNO: %d, recvNO: %d]", sf.sendSN, sf.rcvSN)
}

// S帧 只含apci S帧用于主要用确认帧的正确传输,协议称是监视. supervisory
type sAPCI struct {
	rcvSN uint16
}

func (sf sAPCI) String() string {
	return fmt.Sprintf("S[recvNO: %d]", sf.rcvSN)
}

//U帧 只含apci 未编号控制信息 unnumbered
type uAPCI struct {
	function byte // bit8 测试确认
}

func (sf uAPCI) String() string {
	var s string
	switch sf.function {
	case uStartDtActive:
		s = "StartDtActive"
	case uStartDtConfirm:
		s = "StartDtConfirm"
	case uStopDtActive:
		s = "StopDtActive"
	case uStopDtConfirm:
		s = "StopDtConfirm"
	case uTestFrActive:
		s = "TestFrActive"
	case uTestFrConfirm:
		s = "TestFrConfirm"
	default:
		s = "Unknown"
	}
	return fmt.Sprintf("U[function: %s]", s)
}

// newIFrame 创建I帧 ,返回apdu
func newIFrame(sendSN, RcvSN uint16, asdus []byte) ([]byte, error) {
	if len(asdus) > asdu.ASDUSizeMax {
		return nil, fmt.Errorf("ASDU filed large than max %d", asdu.ASDUSizeMax)
	}

	b := make([]byte, len(asdus)+6)

	b[0] = startFrame
	b[1] = byte(len(asdus) + 4)
	b[2] = byte(sendSN << 1)
	b[3] = byte(sendSN >> 7)
	b[4] = byte(RcvSN << 1)
	b[5] = byte(RcvSN >> 7)
	copy(b[6:], asdus)

	return b, nil
}

// newSFrame 创建S帧,返回apdu
func newSFrame(RcvSN uint16) []byte {
	return []byte{startFrame, 4, 0x01, 0x00, byte(RcvSN << 1), byte(RcvSN >> 7)}
}

// newUFrame 创建U帧,返回apdu
func newUFrame(which byte) []byte {
	return []byte{startFrame, 4, which | 0x03, 0x00, 0x00, 0x00}
}

// APCI apci 应用规约控制信息
type APCI struct {
	start                  byte
	apduFiledLen           byte // control + asdu 的长度
	ctr1, ctr2, ctr3, ctr4 byte
}

// return frame type , APCI, remain data
func parse(apdu []byte) (interface{}, []byte) {
	apci := APCI{apdu[0], apdu[1], apdu[2], apdu[3], apdu[4], apdu[5]}
	if apci.ctr1&0x01 == 0 {
		return iAPCI{
			sendSN: uint16(apci.ctr1)>>1 + uint16(apci.ctr2)<<7,
			rcvSN:  uint16(apci.ctr3)>>1 + uint16(apci.ctr4)<<7,
		}, apdu[6:]
	}
	if apci.ctr1&0x03 == 0x01 {
		return sAPCI{
			rcvSN: uint16(apci.ctr3)>>1 + uint16(apci.ctr4)<<7,
		}, apdu[6:]
	}
	// apci.ctrl&0x03 == 0x03
	return uAPCI{
		function: apci.ctr1 & 0xfc,
	}, apdu[6:]
}
