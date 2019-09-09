package asdu

import (
	"encoding/binary"
	"math"
	"time"
)

// AppendBytes append some bytes to info object
func (this *ASDU) AppendBytes(b ...byte) *ASDU {
	this.infoObj = append(this.infoObj, b...)
	return this
}

// DecodeByte decode a byte then the pass it
func (this *ASDU) DecodeByte() byte {
	v := this.infoObj[0]
	this.infoObj = this.infoObj[1:]
	return v
}

// AppendInfoObjAddr append information object address to information object
func (this *ASDU) AppendInfoObjAddr(addr InfoObjAddr) error {
	switch this.InfoObjAddrSize {
	case 1:
		if addr > 255 {
			return ErrInfoObjAddrFit
		}
		this.infoObj = append(this.infoObj, byte(addr))
	case 2:
		if addr > 65535 {
			return ErrInfoObjAddrFit
		}
		this.infoObj = append(this.infoObj, byte(addr), byte(addr>>8))
	case 3:
		if addr > 16777215 {
			return ErrInfoObjAddrFit
		}
		this.infoObj = append(this.infoObj, byte(addr), byte(addr>>8), byte(addr>>16))
	default:
		return ErrParam
	}
	return nil
}

// DecodeInfoObjAddr decode info object address then the pass it
func (this *ASDU) DecodeInfoObjAddr() InfoObjAddr {
	var ioa InfoObjAddr
	switch this.InfoObjAddrSize {
	case 1:
		ioa = InfoObjAddr(this.infoObj[0])
		this.infoObj = this.infoObj[1:]
	case 2:
		ioa = InfoObjAddr(this.infoObj[0]) | (InfoObjAddr(this.infoObj[1]) << 8)
		this.infoObj = this.infoObj[2:]
	case 3:
		ioa = InfoObjAddr(this.infoObj[0]) | (InfoObjAddr(this.infoObj[1]) << 8) | (InfoObjAddr(this.infoObj[2]) << 16)
		this.infoObj = this.infoObj[3:]
	default:
		panic(ErrParam)
	}
	return ioa
}

// AppendNormalize append a Normalize value to info object
func (this *ASDU) AppendNormalize(n Normalize) *ASDU {
	this.infoObj = append(this.infoObj, byte(n), byte(n>>8))
	return this
}

// DecodeNormalize decode info object byte to a Normalize value
func (this *ASDU) DecodeNormalize() Normalize {
	n := Normalize(binary.LittleEndian.Uint16(this.infoObj))
	this.infoObj = this.infoObj[2:]
	return n
}

// AppendScaled append a Scaled value to info object
// See companion standard 101, subclass 7.2.6.7.
func (this *ASDU) AppendScaled(i int16) *ASDU {
	this.infoObj = append(this.infoObj, byte(i), byte(i>>8))
	return this
}

// DecodeScaled decode info object byte to a Scaled value
func (this *ASDU) DecodeScaled() int16 {
	s := int16(binary.LittleEndian.Uint16(this.infoObj))
	this.infoObj = this.infoObj[2:]
	return s
}

// AppendFloat32 append a float32 value to info object
// See companion standard 101, subclass 7.2.6.8.
func (this *ASDU) AppendFloat32(f float32) *ASDU {
	bits := math.Float32bits(f)
	this.infoObj = append(this.infoObj, byte(bits), byte(bits>>8), byte(bits>>16), byte(bits>>24))
	return this
}

// DecodeScaled decode info object byte to a float32 value
func (this *ASDU) DecodeFloat() float32 {
	f := math.Float32frombits(binary.LittleEndian.Uint32(this.infoObj))
	this.infoObj = this.infoObj[4:]
	return f
}

// AppendBinaryCounterReading append binary couter reading value to info object
// See companion standard 101, subclass 7.2.6.9.
func (this *ASDU) AppendBinaryCounterReading(v BinaryCounterReading) *ASDU {
	this.infoObj = append(this.infoObj, byte(v.CounterReading), byte(v.CounterReading>>8),
		byte(v.CounterReading>>16), byte(v.CounterReading>>24), v.SequenceNotation)
	return this
}

// DecodeBinaryCounterReading decode info object byte to binary couter reading value
func (this *ASDU) DecodeBinaryCounterReading() BinaryCounterReading {
	v := int32(binary.LittleEndian.Uint32(this.infoObj))
	b := this.infoObj[4]
	this.infoObj = this.infoObj[5:]
	return BinaryCounterReading{v, b}
}

// AppendBitsString32 append a bits string value to info object
// See companion standard 101, subclass 7.2.6.13.
func (this *ASDU) AppendBitsString32(v uint32) *ASDU {
	this.infoObj = append(this.infoObj, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return this
}

// DecodeBitsString32 decode info object byte to a bits string value
func (this *ASDU) DecodeBitsString32() uint32 {
	v := binary.LittleEndian.Uint32(this.infoObj)
	this.infoObj = this.infoObj[4:]
	return v
}

// AppendCP56Time2a append a CP56Time2a value to info object
func (this *ASDU) AppendCP56Time2a(t time.Time, loc *time.Location) *ASDU {
	this.infoObj = append(this.infoObj, CP56Time2a(t, loc)...)
	return this
}

func (this *ASDU) DecodeCP56Time2a() time.Time {
	t := ParseCP56Time2a(this.infoObj, this.Params.InfoObjTimeZone)
	this.infoObj = this.infoObj[7:]
	return t
}

// AppendCP24Time2a append CP24Time2a to asdu info object
func (this *ASDU) AppendCP24Time2a(t time.Time, loc *time.Location) *ASDU {
	this.infoObj = append(this.infoObj, CP24Time2a(t, loc)...)
	return this
}

func (this *ASDU) DecodeCP24Time2a() time.Time {
	t := ParseCP24Time2a(this.infoObj, this.Params.InfoObjTimeZone)
	this.infoObj = this.infoObj[3:]
	return t
}

// AppendCP16Time2a append CP16Time2a to asdu info object
func (this *ASDU) AppendCP16Time2a(msec uint16) *ASDU {
	this.infoObj = append(this.infoObj, CP16Time2a(msec)...)
	return this
}

func (this *ASDU) DecodeCP16Time2a() uint16 {
	t := ParseCP16Time2a(this.infoObj)
	this.infoObj = this.infoObj[2:]
	return t
}

func (this *ASDU) DecodeStatusAndStatusChangeDetection() StatusAndStatusChangeDetection {
	// TODO
	return 0
}
