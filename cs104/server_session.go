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

	in      chan []byte // for received asdu
	out     chan []byte // for send asdu
	rawRcv  chan []byte // for recvLoop raw cs104 frame
	rawSend chan []byte // for sendLoop raw cs104 frame

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

// RecvLoop feeds t.rawRcv.
func (sf *SrvSession) recvLoop() {
	sf.Debug("recvLoop started!")
	defer func() {
		sf.cancelFunc()
		sf.wg.Done()
		sf.Debug("recvLoop stopped!")
	}()

	for {
		rawData := make([]byte, APDUSizeMax)
		for rdCnt, length := 0, 2; rdCnt < length; {
			byteCount, err := io.ReadFull(sf.conn, rawData[rdCnt:length])
			if err != nil {
				// See: https://github.com/golang/go/issues/4373
				if err != io.EOF && err != io.ErrClosedPipe ||
					strings.Contains(err.Error(), "use of closed network connection") {
					sf.Error("receive failed, %v", err)
					return
				}

				if e, ok := err.(net.Error); ok && !e.Temporary() {
					sf.Error("receive failed, %v", err)
					return
				}

				if byteCount == 0 && err == io.EOF {
					sf.Error("remote connect closed, %v", err)
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
					sf.Debug("RX Raw[% x]", apdu)
					sf.rawRcv <- apdu
				}
			}
		}
	}
}

// sendLoop drains t.sendTime.
func (sf *SrvSession) sendLoop() {
	sf.Debug("sendLoop started!")
	defer func() {
		sf.cancelFunc()
		sf.wg.Done()
		sf.Debug("sendLoop stopped!")
	}()

	for {
		select {
		case <-sf.ctx.Done():
			return
		case apdu := <-sf.rawSend:
			sf.Debug("TX Raw[% x]", apdu)
			for wrCnt := 0; len(apdu) > wrCnt; {
				byteCount, err := sf.conn.Write(apdu[wrCnt:])
				if err != nil {
					// See: https://github.com/golang/go/issues/4373
					if err != io.EOF && err != io.ErrClosedPipe ||
						strings.Contains(err.Error(), "use of closed network connection") {
						sf.Error("rawSend failed, %v", err)
						return
					}
					if e, ok := err.(net.Error); !ok || !e.Temporary() {
						sf.Error("rawSend failed, %v", err)
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
func (sf *SrvSession) run(ctx context.Context) {
	sf.Debug("run started!")
	// before any thing make sure init
	sf.cleanUp()

	sf.ctx, sf.cancelFunc = context.WithCancel(ctx)
	sf.setConnectStatus(connected)
	sf.wg.Add(3)
	go sf.recvLoop()
	go sf.sendLoop()
	go sf.handlerLoop()

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
		sf.setConnectStatus(disconnected)
		checkTicker.Stop()
		_ = sf.conn.Close() // 连锁引发cancel
		sf.wg.Wait()
		sf.Debug("run stopped!")
	}()

	for {
		if isActive && seqNoCount(sf.ackNoOut, sf.seqNoOut) <= sf.SendUnAckLimitK {
			select {
			case o := <-sf.out:
				sf.sendIFrame(o)
				idleTimeout3Sine = time.Now()
				continue
			case <-sf.ctx.Done():
				return
			default: // make no block
			}
		}
		select {
		case <-sf.ctx.Done():
			return
		case now := <-checkTicker.C:
			// check all timeouts
			if now.Sub(testFrAliveSendSince) >= sf.SendUnAckTimeout1 {
				// now.Sub(startDtActiveSendSince) >= t.SendUnAckTimeout1 ||
				// now.Sub(stopDtActiveSendSince) >= t.SendUnAckTimeout1 ||
				sf.Error("test frame alive confirm timeout t₁")
				return
			}
			// check oldest unacknowledged outbound
			if sf.ackNoOut != sf.seqNoOut &&
				//now.Sub(sf.peek()) >= sf.SendUnAckTimeout1 {
				now.Sub(sf.pending[0].sendTime) >= sf.SendUnAckTimeout1 {
				sf.ackNoOut++
				sf.Error("fatal transmission timeout t₁")
				return
			}

			// 确定最早发送的i-Frame是否超时,超时则回复sFrame
			if sf.ackNoIn != sf.seqNoIn &&
				(now.Sub(unAckRcvSince) >= sf.RecvUnAckTimeout2 ||
					now.Sub(idleTimeout3Sine) >= timeoutResolution) {
				sf.sendSFrame(sf.seqNoIn)
				sf.ackNoIn = sf.seqNoIn
			}

			// 空闲时间到，发送TestFrActive帧,保活
			if now.Sub(idleTimeout3Sine) >= sf.IdleTimeout3 {
				sf.sendUFrame(uTestFrActive)
				testFrAliveSendSince = time.Now()
				idleTimeout3Sine = testFrAliveSendSince
			}

		case apdu := <-sf.rawRcv:
			idleTimeout3Sine = time.Now() // 每收到一个i帧,S帧,U帧, 重置空闲定时器, t3
			apci, asduVal := parse(apdu)
			switch head := apci.(type) {
			case sAPCI:
				sf.Debug("RX sFrame %v", head)
				if !sf.updateAckNoOut(head.rcvSN) {
					sf.Error("fatal incoming acknowledge either earlier than previous or later than sendTime")
					return
				}

			case iAPCI:
				sf.Debug("RX iFrame %v", head)
				if !isActive {
					sf.Warn("station not active")
					break // not active, discard apdu
				}
				if !sf.updateAckNoOut(head.rcvSN) || head.sendSN != sf.seqNoIn {
					sf.Error("fatal incoming acknowledge either earlier than previous or later than sendTime")
					return
				}

				sf.in <- asduVal
				if sf.ackNoIn == sf.seqNoIn { // first unacked
					unAckRcvSince = time.Now()
				}

				sf.seqNoIn = (sf.seqNoIn + 1) & 32767
				if seqNoCount(sf.ackNoIn, sf.seqNoIn) >= sf.RecvUnAckLimitW {
					sf.sendSFrame(sf.seqNoIn)
					sf.ackNoIn = sf.seqNoIn
				}

			case uAPCI:
				sf.Debug("RX uFrame %v", head)
				switch head.function {
				case uStartDtActive:
					sf.sendUFrame(uStartDtConfirm)
					isActive = true
				// case uStartDtConfirm:
				// 	isActive = true
				// 	startDtActiveSendSince = willNotTimeout
				case uStopDtActive:
					sf.sendUFrame(uStopDtConfirm)
					isActive = false
				// case uStopDtConfirm:
				// 	isActive = false
				// 	stopDtActiveSendSince = willNotTimeout
				case uTestFrActive:
					sf.sendUFrame(uTestFrConfirm)
				case uTestFrConfirm:
					testFrAliveSendSince = willNotTimeout
				default:
					sf.Error("illegal U-Frame functions[0x%02x] ignored", head.function)
				}
			}
		}
	}
}

// handlerLoop handler iFrame asdu
func (sf *SrvSession) handlerLoop() {
	sf.Debug("handlerLoop started")
	defer func() {
		sf.wg.Done()
		sf.Debug("handlerLoop stopped")
	}()

	for {
		select {
		case <-sf.ctx.Done():
			return
		case rawAsdu := <-sf.in:
			asduPack := asdu.NewEmptyASDU(sf.params)
			if err := asduPack.UnmarshalBinary(rawAsdu); err != nil {
				sf.Error("asdu UnmarshalBinary failed,%+v", err)
				continue
			}
			if err := sf.serverHandler(asduPack); err != nil {
				sf.Error("serverHandler falied,%+v", err)
			}
		}
	}
}

func (sf *SrvSession) setConnectStatus(status uint32) {
	sf.rwMux.Lock()
	atomic.StoreUint32(&sf.status, status)
	sf.rwMux.Unlock()
}

func (sf *SrvSession) connectStatus() uint32 {
	sf.rwMux.RLock()
	status := atomic.LoadUint32(&sf.status)
	sf.rwMux.RUnlock()
	return status
}

func (sf *SrvSession) cleanUp() {
	sf.ackNoIn = 0
	sf.ackNoOut = 0
	sf.seqNoIn = 0
	sf.seqNoOut = 0
	sf.pending = nil
	// clear sending chan buffer
loop:
	for {
		select {
		case <-sf.rawSend:
		case <-sf.rawRcv:
		case <-sf.in:
		case <-sf.out:
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

func (sf *SrvSession) sendSFrame(rcvSN uint16) {
	sf.Debug("TX sFrame %v", sAPCI{rcvSN})
	sf.rawSend <- newSFrame(rcvSN)
}

func (sf *SrvSession) sendUFrame(which byte) {
	sf.Debug("TX uFrame %v", uAPCI{which})
	sf.rawSend <- newUFrame(which)
}

func (sf *SrvSession) sendIFrame(asdu1 []byte) {
	seqNo := sf.seqNoOut

	iframe, err := newIFrame(seqNo, sf.seqNoIn, asdu1)
	if err != nil {
		return
	}
	sf.ackNoIn = sf.seqNoIn
	sf.seqNoOut = (seqNo + 1) & 32767

	//sf.push(seqPending{seqNo & 32767, time.Now()})
	sf.pending = append(sf.pending, seqPending{seqNo & 32767, time.Now()})

	sf.Debug("TX iFrame %v", iAPCI{seqNo, sf.seqNoIn})
	sf.rawSend <- iframe
}

func (sf *SrvSession) updateAckNoOut(ackNo uint16) (ok bool) {
	if ackNo == sf.ackNoOut {
		return true
	}
	// new acks validate， ack 不能在 req seq 前面,出错
	if seqNoCount(sf.ackNoOut, sf.seqNoOut) < seqNoCount(ackNo, sf.seqNoOut) {
		return false
	}

	// confirm reception
	for i, v := range sf.pending {
		if v.seq == (ackNo - 1) {
			sf.pending = sf.pending[i+1:]
			break
		}
	}
	//sf.confirmReception(ackNo)
	sf.ackNoOut = ackNo
	return true
}

func (sf *SrvSession) serverHandler(asduPack *asdu.ASDU) error {
	defer func() {
		if err := recover(); err != nil {
			sf.Critical("server handler %+v", err)
		}
	}()

	sf.Debug("ASDU %+v", asduPack)

	switch asduPack.Identifier.Type {
	case asdu.C_IC_NA_1: // InterrogationCmd
		if !(asduPack.Identifier.Coa.Cause == asdu.Activation ||
			asduPack.Identifier.Coa.Cause == asdu.Deactivation) {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCA)
		}
		ioa, qoi := asduPack.GetInterrogationCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		return sf.handler.InterrogationHandler(sf, asduPack, qoi)

	case asdu.C_CI_NA_1: // CounterInterrogationCmd
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCA)
		}
		ioa, qcc := asduPack.GetCounterInterrogationCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		return sf.handler.CounterInterrogationHandler(sf, asduPack, qcc)

	case asdu.C_RD_NA_1: // ReadCmd
		if asduPack.Identifier.Coa.Cause != asdu.Request {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCA)
		}
		return sf.handler.ReadHandler(sf, asduPack, asduPack.GetReadCmd())

	case asdu.C_CS_NA_1: // ClockSynchronizationCmd
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCA)
		}

		ioa, tm := asduPack.GetClockSynchronizationCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		return sf.handler.ClockSyncHandler(sf, asduPack, tm)

	case asdu.C_TS_NA_1: // TestCommand
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCA)
		}
		ioa, _ := asduPack.GetTestCommand()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		return asduPack.SendReplyMirror(sf, asdu.ActivationCon)

	case asdu.C_RP_NA_1: // ResetProcessCmd
		if asduPack.Identifier.Coa.Cause != asdu.Activation {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCA)
		}
		ioa, qrp := asduPack.GetResetProcessCmd()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		return sf.handler.ResetProcessHandler(sf, asduPack, qrp)
	case asdu.C_CD_NA_1: // DelayAcquireCommand
		if !(asduPack.Identifier.Coa.Cause == asdu.Activation ||
			asduPack.Identifier.Coa.Cause == asdu.Spontaneous) {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		if asduPack.CommonAddr == asdu.InvalidCommonAddr {
			return asduPack.SendReplyMirror(sf, asdu.UnknownCA)
		}
		ioa, msec := asduPack.GetDelayAcquireCommand()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return asduPack.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		return sf.handler.DelayAcquisitionHandler(sf, asduPack, msec)
	}

	if err := sf.handler.ASDUHandler(sf, asduPack); err != nil {
		return asduPack.SendReplyMirror(sf, asdu.UnknownTypeID)
	}
	return nil
}

// IsConnected get server session connected state
func (sf *SrvSession) IsConnected() bool {
	return sf.connectStatus() == connected
}

// Params get params
func (sf *SrvSession) Params() *asdu.Params {
	return sf.params
}

// Send asdu frame
func (sf *SrvSession) Send(u *asdu.ASDU) error {
	if !sf.IsConnected() {
		return ErrUseClosedConnection
	}
	data, err := u.MarshalBinary()
	if err != nil {
		return err
	}
	select {
	case sf.out <- data:
	default:
		return ErrBufferFulled
	}
	return nil
}

// UnderlyingConn got under net.conn
func (sf *SrvSession) UnderlyingConn() net.Conn {
	return sf.conn
}
