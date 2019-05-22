package iec10x

const (
	f12startVarFrame = 0x68
)

const (
	ASDUSizeMax      = 249 // ASDU
	APDUFiledSizeMax = 253 // control(4) + ASDU
	APDUSizeMax      = 255 // start(1) +length(1) + control(4) + ASDU
)
