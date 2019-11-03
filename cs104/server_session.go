package cs104

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

const (
	initial uint32 = iota
	disconnected
	connected
)

// SrvSession the cs104 server session
type SrvSession struct {
	*Config
	params  *asdu.Params
	conn    net.Conn
	handler ServerHandlerInterface

	in   chan []byte // for received asdu
	out  chan []byte // for send asdu
	recv chan []byte // for recvLoop raw cs104 frame
	send chan []byte // for sendLoop raw cs104 frame

	// see subclass 5.1 — Protection against loss and duplication of messages
	seqNoOut uint16 // sequence number of next outbound I-frame
	ackNoOut uint16 // outbound sequence number yet to be confirmed
	seqNoIn  uint16 // sequence number of next inbound I-frame
	ackNoIn  uint16 // inbound sequence number yet to be confirmed
	// maps sendTime I-frames to their respective sequence number
	pending []seqPending
	//seqManage

	status uint32
	rwMux  sync.RWMutex

	*clog.Clog

	wg         sync.WaitGroup
	cancelFunc context.CancelFunc
	ctx        context.Context
}

// RecvLoop feeds t.recv.
func (this *SrvSession) recvLoop() {
	this.Debug("recvLoop start!")
	defer func() {
		this.cancelFunc()
		this.wg.Done()
		this.Debug("recvLoop stop!")
	}()

	for {
		rawData := make([]byte, APDUSizeMax)
		for rdCnt, length := 0, 2; rdCnt < length; {
			byteCount, err := io.ReadFull(this.conn, rawData[rdCnt:length])
			if err != nil {
				// See: https://github.com/golang/go/issues/4373
				if err != io.EOF && err != io.ErrClosedPipe ||
					strings.Contains(err.Error(), "use of closed network connection") {
					this.Error("receive failed, %v", err)
					return
				}

				if e, ok := err.(net.Error); ok && !e.Temporary() {
					this.Error("receive failed, %v", err)
					return
				}

				if byteCount == 0 && err == io.EOF {
					this.Error("remote connect closed, %v", err)
					return
				}
			}

			rdCnt += byteCount
			if rdCnt == 0 {
				continue
			} else if rdCnt == 1 {
				if rawData[0] != startFrame {
					rdCnt = 0
					continue
				}
			} else {
				if rawData[0] != startFrame {
					rdCnt = 0
					length = 2
					continue
				}
				length = int(rawData[1]) + 2
				if length < APCICtlFiledSize+2 || length > APDUSizeMax {
					rdCnt = 0
					length = 2
					continue
				}
				if rdCnt == length {
					apdu := rawData[:length]
					this.Debug("RX Raw[% x]", apdu)
					this.recv <- apdu
				}
			}
		}
	}
}

// sendLoop drains t.sendTime.
func (this *SrvSession) sendLoop() {
	this.Debug("sendLoop start!")
	defer func() {
		this.cancelFunc()
		this.wg.Done()
		this.Debug("sendLoop stop!")
	}()

	for {
		select {
		case <-this.ctx.Done():
			return
		case apdu := <-this.send:
			this.Debug("TX Raw[% x]", apdu)
			for wrCnt := 0; len(apdu) > wrCnt; {
				byteCount, err := this.conn.Write(apdu[wrCnt:])
				if err != nil {
					// See: https://github.com/golang/go/issues/4373
					if err != io.EOF && err != io.ErrClosedPipe ||
						strings.Contains(err.Error(), "use of closed network connection") {
						this.Error("send failed, %v", err)
						return
					}
					if e, ok := err.(net.Error); !ok || !e.Temporary() {
						this.Error("send failed, %v", err)
						return
					}
					// temporary error may be recoverable
				}
				wrCnt += byteCount
			}
		}
	}
}

// run is the big fat state machine.
func (this *SrvSession) run(ctx context.Context) {
	this.Debug("run start!")
	// before any thing make sure init
	this.cleanUp()

	this.ctx, this.cancelFunc = context.WithCancel(ctx)
	this.setConnectStatus(connected)
	this.wg.Add(3)
	go this.recvLoop()
	go this.sendLoop()
	go this.handlerLoop()

	// default: STOPDT, when connected establish and not enable "data transfer" yet
	var isActive = false
	var checkTicker = time.NewTicker(timeoutResolution)

	// transmission timestamps for timeout calculation
	var willNotTimeout = time.Now().Add(time.Hour * 24 * 365 * 100)

	var unAckRcvSince = willNotTimeout
	var idleTimeout3Sine = time.Now()         // 空闲间隔发起testFrAlive
	var testFrAliveSendSince = willNotTimeout // 当发起testFrAlive时,等待确认回复的超时间隔
	// 对于server端，无需对应的U-Frame 无需判断
	// var startDtActiveSendSince = willNotTimeout
	// var stopDtActiveSendSince = willNotTimeout

	defer func() {
		this.setConnectStatus(disconnected)
		checkTicker.Stop()
		this.conn.Close() // 连锁引发cancel
		this.wg.Wait()
		this.Debug("run stop!")
	}()

	for {
		if isActive && seqNoCount(this.ackNoOut, this.seqNoOut) <= this.SendUnAckLimitK {
			select {
			case o := <-this.out:
				this.sendIFrame(o)
				idleTimeout3Sine = time.Now()
				continue
			case <-this.ctx.Done():
				return
			default: // make no block
			}
		}
		select {
		case <-this.ctx.Done():
			return
		case now := <-checkTicker.C:
			// check all timeouts
			if now.Sub(testFrAliveSendSince) >= this.SendUnAckTimeout1 {
				// now.Sub(startDtActiveSendSince) >= t.SendUnAckTimeout1 ||
				// now.Sub(stopDtActiveSendSince) >= t.SendUnAckTimeout1 ||
				this.Error("test frame alive confirm timeout t₁")
				return
			}
			// check oldest unacknowledged outbound
			if this.ackNoOut != this.seqNoOut &&
				//now.Sub(this.peek()) >= this.SendUnAckTimeout1 {
				now.Sub(this.pending[0].sendTime) >= this.SendUnAckTimeout1 {
				this.ackNoOut++
				this.Error("fatal transmission timeout t₁")
				return
			}

			// 确定最早发送的i-Frame是否超时,超时则回复sFrame
			if this.ackNoIn != this.seqNoIn &&
				(now.Sub(unAckRcvSince) >= this.RecvUnAckTimeout2 ||
					now.Sub(idleTimeout3Sine) >= timeoutResolution) {
				this.sendSFrame(this.seqNoIn)
				this.ackNoIn = this.seqNoIn
			}

			// 空闲时间到，发送TestFrActive帧,保活
			if now.Sub(idleTimeout3Sine) >= this.IdleTimeout3 {
				this.sendUFrame(uTestFrActive)
				testFrAliveSendSince = time.Now()
				idleTimeout3Sine = testFrAliveSendSince
			}

		case apdu := <-this.recv:
			idleTimeout3Sine = time.Now() // 每收到一个i帧,S帧,U帧, 重置空闲定时器, t3
			apci, asdu := parse(apdu)
			switch head := apci.(type) {
			case sAPCI:
				this.Debug("RX sFrame %v", head)
				if !this.updateAckNoOut(head.rcvSN) {
					this.Error("fatal incoming acknowledge either earlier than previous or later than sendTime")
					return
				}

			case iAPCI:
				this.Debug("RX iFrame %v", head)
				if !isActive {
					this.Warn("station not active")
					break // not active, discard apdu
				}
				if !this.updateAckNoOut(head.rcvSN) || head.sendSN != this.seqNoIn {
					this.Error("fatal incoming acknowledge either earlier than previous or later than sendTime")
					return
				}

				this.in <- asdu
				if this.ackNoIn == this.seqNoIn { // first unacked
					unAckRcvSince = time.Now()
				}

				this.seqNoIn = (this.seqNoIn + 1) & 32767
				if seqNoCount(this.ackNoIn, this.seqNoIn) >= this.RecvUnAckLimitW {
					this.sendSFrame(this.seqNoIn)
					this.ackNoIn = this.seqNoIn
				}

			case uAPCI:
				this.Debug("RX uFrame %v", head)
				switch head.function {
				case uStartDtActive:
					this.sendUFrame(uStartDtConfirm)
					isActive = true
				// case uStartDtConfirm:
				// 	isActive = true
				// 	startDtActiveSendSince = willNotTimeout
				case uStopDtActive:
					this.sendUFrame(uStopDtConfirm)
					isActive = false
				// case uStopDtConfirm:
				// 	isActive = false
				// 	stopDtActiveSendSince = willNotTimeout
				case uTestFrActive:
					this.sendUFrame(uTestFrConfirm)
				case uTestFrConfirm:
					testFrAliveSendSince = willNotTimeout
				default:
					this.Error("illegal U-Frame functions[0x%02x] ignored", head.function)
				}
			}
		}
	}
}

// handlerLoop handler iFrame asdu
func (this *SrvSession) handlerLoop() {
	this.Debug("handlerLoop start")
	defer func() {
		this.wg.Done()
		this.Debug("handlerLoop stop")
	}()

	for {
		select {
		case <-this.ctx.Done():
			return
		case rawAsdu := <-this.in:
			asduPack := asdu.NewEmptyASDU(this.params)
			if err := asduPack.UnmarshalBinary(rawAsdu); err != nil {
				this.Error("asdu UnmarshalBinary failed,%+v", err)
				continue
			}
			if err := this.serverHandler(asduPack); err != nil {
				this.Error("serverHandler falied,%+v", err)
			}
		}
	}
}

func (this *SrvSession) setConnectStatus(status uint32) {
	this.rwMux.Lock()
	atomic.StoreUint32(&this.status, status)
	this.rwMux.Unlock()
}

func (this *SrvSession) connectStatus() uint32 {
	this.rwMux.RLock()
	status := atomic.LoadUint32(&this.status)
	this.rwMux.RUnlock()
	return status
}

func (this *SrvSession) cleanUp() {
	this.ackNoIn = 0
	this.ackNoOut = 0
	this.seqNoIn = 0
	this.seqNoOut = 0
	this.pending = nil
	// clear sending chan buffer
loop:
	for {
		select {
		case <-this.send:
		case <-this.recv:
		case <-this.in:
		case <-this.out:
		default:
			break loop
		}
	}
}

// 回绕机制
func seqNoCount(nextAckNo, nextSeqNo uint16) uint16 {
	if nextAckNo > nextSeqNo {
		nextSeqNo += 32768
	}
	return nextSeqNo - nextAckNo
}

func (this *SrvSession) sendSFrame(rcvSN uint16) {
	this.Debug("TX sFrame %v", sAPCI{rcvSN})
	this.send <- newSFrame(rcvSN)
}

func (this *SrvSession) sendUFrame(which byte) {
	this.Debug("TX uFrame %v", uAPCI{which})
	this.send <- newUFrame(which)
}

func (this *SrvSession) sendIFrame(asdu1 []byte) {
	seqNo := this.seqNoOut

	iframe, err := newIFrame(seqNo, this.seqNoIn, asdu1)
	if err != nil {
		return
	}
	this.ackNoIn = this.seqNoIn
	this.seqNoOut = (seqNo + 1) & 32767

	//this.push(seqPending{seqNo & 32767, time.Now()})
	this.pending = append(this.pending, seqPending{seqNo & 32767, time.Now()})

	this.Debug("TX iFrame %v", iAPCI{seqNo, this.seqNoIn})
	this.send <- iframe
}

func (this *SrvSession) updateAckNoOut(ackNo uint16) (ok bool) {
	if ackNo == this.ackNoOut {
		return true
	}
	// new acks validate， ack 不能在 req seq 前面,出错
	if seqNoCount(this.ackNoOut, this.seqNoOut) < seqNoCount(ackNo, this.seqNoOut) {
		return false
	}

	// confirm reception
	for i, v := range this.pending {
		if v.seq == (ackNo - 1) {
			this.pending = this.pending[i+1:]
			break
		}
	}
	//this.confirmReception(ackNo)
	this.ackNoOut = ackNo
	return true
}

func (this *SrvSession) serverHandler(asduPack *asdu.ASDU) error {
	defer func() {
		if err := recover(); err != nil {
			this.Critical("server handler %+v", err)
		}
	}()

	this.Debug("ASDU %+v", asduPack)

	switch asduPack.Identifier.Type {
	case asdu.C_IC_NA_1: // InterrogationCmd
		if !(asduPack.Identifier.Coa.Cause == asdu.Activation ||
			asduPack.Identifier.Coa.Cause == asdu.Deactivation) {
			return asduPack.SendReplyMirror(this, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(this, asdu.UnknownCA)
		}
		ioa, qoi := asduPack.GetInterrogationCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(this, asdu.UnknownIOA)
		}
		return this.handler.InterrogationHandler(this, asduPack, qoi)

	case asdu.C_CI_NA_1: // CounterInterrogationCmd
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(this, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(this, asdu.UnknownCA)
		}
		ioa, qcc := asduPack.GetCounterInterrogationCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(this, asdu.UnknownIOA)
		}
		return this.handler.CounterInterrogationHandler(this, asduPack, qcc)

	case asdu.C_RD_NA_1: // ReadCmd
		if asduPack.Identifier.Coa.Cause != asdu.Request {
			return asduPack.SendReplyMirror(this, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(this, asdu.UnknownCA)
		}
		return this.handler.ReadHandler(this, asduPack, asduPack.GetReadCmd())

	case asdu.C_CS_NA_1: // ClockSynchronizationCmd
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(this, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(this, asdu.UnknownCA)
		}

		ioa, tm := asduPack.GetClockSynchronizationCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(this, asdu.UnknownIOA)
		}
		return this.handler.ClockSyncHandler(this, asduPack, tm)

	case asdu.C_TS_NA_1: // TestCommand
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(this, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(this, asdu.UnknownCA)
		}
		ioa, _ := asduPack.GetTestCommand()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(this, asdu.UnknownIOA)
		}
		return asduPack.SendReplyMirror(this, asdu.ActivationCon)

	case asdu.C_RP_NA_1: // ResetProcessCmd
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(this, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(this, asdu.UnknownCA)
		}
		ioa, qrp := asduPack.GetResetProcessCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(this, asdu.UnknownIOA)
		}
		return this.handler.ResetProcessHandler(this, asduPack, qrp)
	case asdu.C_CD_NA_1: // DelayAcquireCommand
		if !(asduPack.Identifier.Coa.Cause == asdu.Activation ||
			asduPack.Identifier.Coa.Cause == asdu.Spontaneous) {
			return asduPack.SendReplyMirror(this, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(this, asdu.UnknownCA)
		}
		ioa, msec := asduPack.GetDelayAcquireCommand()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(this, asdu.UnknownIOA)
		}
		return this.handler.DelayAcquisitionHandler(this, asduPack, msec)
	}

	if err := this.handler.ASDUHandler(this, asduPack); err != nil {
		return asduPack.SendReplyMirror(this, asdu.UnknownTypeID)
	}
	return nil
}

// IsConnected get server session connected state
func (this *SrvSession) IsConnected() bool {
	return this.connectStatus() == connected
}

// Params get params
func (this *SrvSession) Params() *asdu.Params {
	return this.params
}

// Send asdu frame
func (this *SrvSession) Send(u *asdu.ASDU) error {
	if !this.IsConnected() {
		return ErrUseClosedConnection
	}
	data, err := u.MarshalBinary()
	if err != nil {
		return err
	}
	select {
	case this.out <- data:
	default:
		return ErrBufferFulled
	}
	return nil
}

// UnderlyingConn got under net.conn
func (this *SrvSession) UnderlyingConn() net.Conn {
	return this.conn
}
