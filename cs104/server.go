package cs104

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/library/elog"
)

// TimeoutResolution is seconds according to companion standard 104,
// subclass 6.9, caption "Definition of time outs". However, thenths
// of a second make this system much more responsive i.c.w. S-frames.
const timeoutResolution = 100 * time.Millisecond

var (
	errSeqNo            = errors.New("cs104: fatal incomming sequence number disruption")
	errAckNo            = errors.New("cs104: fatal incomming acknowledge either earlier than previous or later than send")
	errAckExpire        = errors.New("cs104: fatal transmission timeout t₁")
	errStartDtAckExpire = errors.New("cs104: fatal STARTDT acknowledge timeout t₁")
	errStopDtAckExpire  = errors.New("cs104: fatal STOPDT acknowledge timeout t₁")
	errTestFrAckExpire  = errors.New("cs104: fatal TESTFR acknowledge timeout t₁")
	errAPCIIllegalFunc  = errors.New("cs104: illegal function ignored")
)

type Server struct {
	Config
	conn net.Conn

	In   chan []byte
	out  chan []byte
	recv chan []byte // for recvLoop
	send chan []byte // for sendLoop

	// see subclass 5.1 — Protection against loss and duplication of messages
	seqNoOut uint16 // sequence number of next outbound I-frame
	ackNoOut uint16 // outbound sequence number yet to be confirmed
	seqNoIn  uint16 // sequence number of next inbound I-frame
	ackNoIn  uint16 // inbound sequence number yet to be confirmed

	// maps send I-frames to their respective sequence number
	pending [1 << 15]struct {
		send time.Time
	}

	idleSince  time.Time
	cancelFunc context.CancelFunc
	ctx        context.Context
	*elog.Elog
}

// NewServer returns a cs104 server
func NewServer(conf *Config, conn net.Conn) (*Server, error) {
	err := conf.Valid()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	t := &Server{
		Config: *conf,
		conn:   conn,

		In:   make(chan []byte),
		out:  make(chan []byte),
		recv: make(chan []byte, conf.RecvUnackLimitW),
		send: make(chan []byte, conf.SendUnackLimitK), // may not block!

		idleSince:  time.Now(),
		cancelFunc: cancel,
		ctx:        ctx,
		Elog:       elog.GetElog(),
	}
	go t.recvLoop()
	go t.sendLoop()
	go t.run()
	return t, nil
}

// RecvLoop feeds t.recv.
func (t *Server) recvLoop() {
	t.Info("cs104 server: recvLoop start!")
	// 临时错误恢复，过长和过短不适合，这个需要再调试
	retryTicker := time.NewTicker(200 * time.Millisecond)

	defer func() {
		close(t.recv)
		retryTicker.Stop()
		t.cancelFunc()
		t.Info("cs104 server: recvLoop stop!")
	}()

	var deadline time.Time
	for {
		datagram := make([]byte, 255)
		length := 2
		for rdCnt := 0; rdCnt < length; {
			byteCount, err := io.ReadFull(t.conn, datagram[rdCnt:length])
			if err != nil {
				// See: https://github.com/golang/go/issues/4373
				if err != io.EOF && err != io.ErrClosedPipe || strings.Contains(err.Error(), "use of closed network connection") {
					t.Error("cs104 server: %v", err)
					return
				}

				if e, ok := err.(net.Error); ok && !e.Temporary() {
					t.Error("cs104 server: %v", err)
					return
				}

				if byteCount == 0 && err == io.EOF {
					t.Error("cs104 server: remote connect closed,%v", err)
					return
				}
				// temporary error may be recoverable
				now := <-retryTicker.C
				switch {
				case deadline.IsZero():
					deadline = now.Add(t.SendUnackTimeout1)
				case now.After(deadline):
					t.Error("%v", errAckExpire)
					return
				}
			}

			rdCnt += byteCount
			if rdCnt == 0 {
				break
			} else if rdCnt == 1 {
				if datagram[0] != startFrame {
					break
				}
			} else {
				if datagram[0] != startFrame {
					break
				}
				length = int(datagram[1]) + 2
				if length < APCICtlFiledSize+2 || length > APDUSizeMax {
					break
				}
				if rdCnt == length {
					apdu := datagram[:length]
					t.Debug("cs104 server: Raw RX [% X]", apdu)
					t.recv <- apdu // copy
				}
			}
		}
	}
}

// SendLoop drains t.send.
func (t *Server) sendLoop() {
	t.Info("cs104 server: sendLoop start!")
	defer func() {
		t.cancelFunc()
		t.Info("cs104 server server: sendLoop stop!")
	}()

	for apdu := range t.send {
		t.Debug("cs104 server: Raw TX [% X]", apdu)
		for wrCnt := 0; len(apdu) > wrCnt; {
			byteCount, err := t.conn.Write(apdu[wrCnt:])
			if err != nil {
				// See: https://github.com/golang/go/issues/4373
				if err != io.EOF && err != io.ErrClosedPipe || strings.Contains(err.Error(), "use of closed network connection") {
					t.Error("cs104 server: %v", err)
					return
				}
				if e, ok := err.(net.Error); !ok || !e.Temporary() {
					t.Error("cs104 server: %v", err)
					return
				}
				// temporary error may be recoverable
			}
			wrCnt += byteCount
		}
	}
}

// Run is the big fat state machine.
func (this *Server) run() {
	this.Info("cs104 server: run start!")
	// when connected establish and not enable "data transfer" yet
	// defualt: STOPDT
	isActive := false
	checkTicker := time.NewTicker(timeoutResolution)
	defer func() {
		checkTicker.Stop()
		if this.ackNoIn != this.seqNoIn {
			select {
			case this.send <- newSFrame(this.seqNoIn):
				this.ackNoIn = this.seqNoIn
			default:
			}
		}

		close(this.send) // kill send loop
		this.conn.Close()

		// await receive loop
		for apdu := range this.recv {
			apci, _ := parse(apdu)
			switch head, f := apci.parse(); f {
			case iFrame:
				this.updateAckNoOut(head.(iAPCI).rcvSN)

			case sFrame:
				this.updateAckNoOut(head.(sAPCI).rcvSN)
			default:
				// discard
			}
		}

		close(this.In)
		this.Info("cs104 server: run stop!")
	}()

	// transmission timestamps for timeout calculation
	var willNotTimeout = time.Now().Add(time.Hour * 24 * 365 * 100)
	var unAckRecvdSince = willNotTimeout
	var testFrAliveSendSince = willNotTimeout
	// 对于server端，无需对应的U-Frame 无需判断
	// var startDtActiveSendSince = willNotTimeout
	// var stopDtActiveSendSince = willNotTimeout

	for {
		if isActive && seqNoCount(this.ackNoOut, this.seqNoOut) <= this.SendUnackLimitK {
			select {
			case o, ok := <-this.out:
				if !ok {
					return
				}
				this.submit(o)
				continue
			default:
			}
		}
		select {
		case <-this.ctx.Done():
			return

		case now := <-checkTicker.C:
			// check all timeouts
			if now.Sub(testFrAliveSendSince) >= this.SendUnackTimeout1 {
				// now.Sub(startDtActiveSendSince) >= t.SendUnackTimeout1 ||
				// now.Sub(stopDtActiveSendSince) >= t.SendUnackTimeout1 ||
				return
			}
			// check oldest unacknowledged outbound
			if this.ackNoOut != this.seqNoOut &&
				now.Sub(this.pending[this.ackNoOut].send) >= this.SendUnackTimeout1 {
				this.ackNoOut++
				return
			}

			// 确定最早发送的i-Frame是否超时
			if this.ackNoIn != this.seqNoIn &&
				(now.Sub(unAckRecvdSince) >= this.RecvUnackTimeout2 ||
					now.Sub(this.idleSince) >= timeoutResolution) {
				this.send <- newSFrame(this.seqNoIn)
				this.ackNoIn = this.seqNoIn
			}

			// 空闲时间到，发送TestFrActive帧
			if now.Sub(this.idleSince) >= this.IdleTimeout3 {
				this.send <- newUFrame(uTestFrActive)
				testFrAliveSendSince = time.Now()
				this.idleSince = testFrAliveSendSince
			}

		case apdu, ok := <-this.recv:
			if !ok {
				return
			}
			apci, asdu := parse(apdu)
			head, f := apci.parse()
			this.idleSince = time.Now() // 每收到一个i帧,S帧,U帧, 重置空闲定时器
			switch f {
			case sFrame:
				if !this.updateAckNoOut(head.(sAPCI).rcvSN) {
					return
				}

			case iFrame:
				if !isActive {
					break // not active, discard apdu
				}
				ihead := head.(iAPCI)
				if !this.updateAckNoOut(ihead.rcvSN) || ihead.sendSN != this.seqNoIn {
					return
				}

				this.In <- asdu
				if this.ackNoIn == this.seqNoIn { // first unacked
					unAckRecvdSince = time.Now()
				}

				this.seqNoIn = (this.seqNoIn + 1) & 32767
				if seqNoCount(this.ackNoIn, this.seqNoIn) >= this.RecvUnackLimitW {
					this.send <- newSFrame(this.seqNoIn)
					this.ackNoIn = this.seqNoIn
				}

			case uFrame:
				switch head.(uAPCI).function {
				case uStartDtActive:
					this.send <- newUFrame(uStartDtConfirm)
					isActive = true
				// case uStartDtConfirm:
				// 	isActive = true
				// 	startDtActiveSendSince = willNotTimeout
				case uStopDtActive:
					this.send <- newUFrame(uStopDtConfirm)
					isActive = false
				// case uStopDtConfirm:
				// 	isActive = false
				// 	stopDtActiveSendSince = willNotTimeout
				case uTestFrActive:
					this.send <- newUFrame(uTestFrConfirm)
				case uTestFrConfirm:
					testFrAliveSendSince = willNotTimeout
				default:
					this.Error("cs104 server: illegal U-Frame functions[%v]", head.(uAPCI).function)
				}
			}
		}
	}
}

func (t *Server) submit(asdu []byte) {
	seqNo := t.seqNoOut

	datagram, err := newIFrame(asdu, seqNo, t.seqNoIn)
	if err != nil {
		return
	}
	t.ackNoIn = t.seqNoIn
	t.seqNoOut = (seqNo + 1) & 32767

	p := &t.pending[seqNo&32767]
	p.send = time.Now()

	t.send <- datagram
	t.idleSince = time.Now()
}

func (t *Server) updateAckNoOut(ackNo uint16) (ok bool) {
	if ackNo == t.ackNoOut {
		return true
	}
	// new acks validate， ack 不能在 req seq 前面,出错
	if seqNoCount(t.ackNoOut, t.seqNoOut) < seqNoCount(ackNo, t.seqNoOut) {
		return false
	}

	// confirm reception
	for ackNo != t.ackNoOut {
		t.ackNoOut = (t.ackNoOut + 1) & 32767
	}

	t.ackNoOut = ackNo
	return true
}

// 回绕机制
func seqNoCount(nextAckNo, nextSeqNo uint16) uint16 {
	if nextAckNo > nextSeqNo {
		nextSeqNo += 32768
	}
	return nextSeqNo - nextAckNo
}

// SendASDU
func (this *Server) SendASDU(u *asdu.ASDU) error {
	data, err := u.MarshalBinary()
	if err != nil {
		return err
	}
	// select {
	// case this.out <- data:
	// default:
	// 	return errors.New("cs104: buffer is full")
	// }
	this.out <- data
	return nil
}

func (this *Server) Close() {
	this.cancelFunc()
}
