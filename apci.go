package iec10x

import "errors"

// U帧 ctr1
const (
	uStartDtActive  = iota // 启动激活
	uStartDtConfirm        // 启动确认
	uStopDtActive          // 停止激活
	uStopDtConfirm         // 停止确认
	uTestFrActive          // 测试激活
	uTestFrConfirm         // 测试确认

)

type APCI struct {
	start                  byte
	apduFiled              byte
	ctr1, ctr2, ctr3, ctr4 byte
}

type iFrame struct {
	sendSN, rcvSN uint16
}

type sFrame struct {
	rcvSN uint16
}

type uFrame struct {
	testFrConfirm  bool // bit8 测试确认
	testFrActive   bool // bit7 测试激活
	stopDtConfirm  bool // bit6 停止确认
	stopDtActive   bool // bit5 停止激活
	startDtConfirm bool // bit4 启动确认
	startDtActive  bool // bit3 启动激活
}

// 序列化apic
func (this APCI) bytes() []byte {
	return []byte{
		this.start,
		this.apduFiled,
		this.ctr1,
		this.ctr2,
		this.ctr3,
		this.ctr4,
	}
}

// 解析到I,S,U帧
func (this APCI) parse() interface{} {
	if this.ctr1&0x01 == 0 {
		return iFrame{
			sendSN: uint16(this.ctr1)>>1 + uint16(this.ctr2)<<7,
			rcvSN:  uint16(this.ctr3)>>1 + uint16(this.ctr4)<<7,
		}
	}

	if this.ctr1&0x03 == 0x01 {
		return sFrame{
			rcvSN: uint16(this.ctr3)>>1 + uint16(this.ctr4)<<7,
		}
	}
	return uFrame{
		startDtConfirm: this.ctr1&0x04 != 0,
		startDtActive:  this.ctr1&0x08 != 0,
		stopDtConfirm:  this.ctr1&0x10 != 0,
		stopDtActive:   this.ctr1&0x20 != 0,
		testFrConfirm:  this.ctr1&0x40 != 0,
		testFrActive:   this.ctr1&0x80 != 0,
	}
}

// NewIFrameAPCI 创建I帧
func NewIFrameAPCI(asduLen byte, sendSN, RcvSN uint16) (*APCI, error) {
	if asduLen+4 > APDUFiledSizeMax {
		return nil, errors.New("apdu filed large than 253")
	}
	return &APCI{
		start: f12startVarFrame,
		ctr1:  byte(sendSN << 1),
		ctr2:  byte(sendSN >> 7),
		ctr3:  byte(RcvSN << 1),
		ctr4:  byte(RcvSN >> 7),
	}, nil
}

// NewSFrameAPCI 新建S帧
func NewSFrameAPCI(asduLen byte, RcvSN uint16) (*APCI, error) {
	if asduLen+4 > APDUFiledSizeMax {
		return nil, errors.New("apdu filed large than 253")
	}
	return &APCI{
		start: f12startVarFrame,
		ctr1:  0x01,
		ctr2:  0,
		ctr3:  byte(RcvSN << 1),
		ctr4:  byte(RcvSN >> 7),
	}, nil
}

// NewUFrameAPCI 新建U帧
func NewUFrameAPCI(asduLen byte, which int) (*APCI, error) {
	if asduLen+4 > APDUFiledSizeMax {
		return nil, errors.New("apdu filed large than 253")
	}
	apci := &APCI{start: f12startVarFrame, ctr1: 0x03}
	switch which {
	case uStartDtActive:
		apci.ctr1 |= 0x04
	case uStartDtConfirm:
		apci.ctr1 |= 0x08
	case uStopDtActive:
		apci.ctr1 |= 0x10
	case uStopDtConfirm:
		apci.ctr1 |= 0x20
	case uTestFrActive:
		apci.ctr1 |= 0x40
	case uTestFrConfirm:
		apci.ctr1 |= 0x80
	default:
		return nil, errors.New("unknow control filed type")
	}
	return apci, nil
}
