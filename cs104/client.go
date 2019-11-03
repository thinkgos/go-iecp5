package cs104

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

// Client is an IEC104 master
type Client struct {
	conf    Config
	param   asdu.Params
	conn    net.Conn
	handler ClientHandlerInterface

	// channel
	rcvASDU  chan []byte // for received asdu
	sendASDU chan []byte // for send asdu
	rcvRaw   chan []byte // for recvLoop raw cs104 frame
	sendRaw  chan []byte // for sendLoop raw cs104 frame

	// I帧的发送与接收序号
	seqNoSend uint16 // sequence number of next outbound I-frame
	ackNoSend uint16 // outbound sequence number yet to be confirmed
	seqNoRcv  uint16 // sequence number of next inbound I-frame
	ackNoRcv  uint16 // inbound sequence number yet to be confirmed

	// maps sendTime I-frames to their respective sequence number
	pending []seqPending

	startDtActiveSendSince atomic.Value // 当发送startDtActive时,等待确认回复的超时间隔
	stopDtActiveSendSince  atomic.Value // 当发起stopDtActive时,等待确认回复的超时

	// 连接状态
	status uint32
	rwMux  sync.RWMutex

	// 其他
	*clog.Clog

	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	closeCancel context.CancelFunc

	server            *url.URL      // 连接的服务器端
	autoReconnect     bool          // 是否启动重连
	reconnectInterval time.Duration // 重连间隔时间
	TLSConfig         *tls.Config   // tls配置
	onConnect         func(c *Client) error
	onConnectionLost  func(c *Client)
}

// NewClient returns an IEC104 master
func NewClient(conf *Config, params *asdu.Params, handler ClientHandlerInterface) (*Client, error) {
	if handler == nil {
		return nil, errors.New("invalid handler")
	}
	if err := conf.Valid(); err != nil {
		return nil, err
	}
	if err := params.Valid(); err != nil {
		return nil, err
	}
	return &Client{
		conf:     *conf,
		param:    *params,
		handler:  handler,
		rcvASDU:  make(chan []byte, 1024),
		sendASDU: make(chan []byte, 1024),
		rcvRaw:   make(chan []byte, 1024),
		sendRaw:  make(chan []byte, 1024), // may not block!
		Clog:     clog.NewWithPrefix("cs104 client =>"),

		autoReconnect:     true,
		reconnectInterval: DefaultReconnectInterval,
		onConnect:         func(*Client) error { return nil },
		onConnectionLost:  func(*Client) {},
	}, nil
}

// SetReconnectInterval set tcp  reconnect the host interval when connect failed after try
func (sf *Client) SetReconnectInterval(t time.Duration) {
	sf.reconnectInterval = t
}

// SetAutoReconnect enable auto reconnect
func (sf *Client) SetAutoReconnect(b bool) {
	sf.autoReconnect = b
}

// SetTLSConfig set tls config
func (sf *Client) SetTLSConfig(t *tls.Config) {
	sf.TLSConfig = t
}

// AddRemoteServer adds a broker URI to the list of brokers to be used.
// The format should be scheme://host:port
// Default values for hostname is "127.0.0.1", for schema is "tcp://".
// An example broker URI would look like: tcp://foobar.com:1204
func (sf *Client) AddRemoteServer(server string) error {
	if len(server) > 0 && server[0] == ':' {
		server = "127.0.0.1" + server
	}
	if !strings.Contains(server, "://") {
		server = "tcp://" + server
	}
	remoteURL, err := url.Parse(server)
	if err != nil {
		return err
	}
	sf.server = remoteURL
	return nil
}

// SetOnConnectHandler set on connect handler
func (sf *Client) SetOnConnectHandler(f func(c *Client) error) {
	if f != nil {
		sf.onConnect = f
	}
}

// SetConnectionLostHandler set connection lost handler
func (sf *Client) SetConnectionLostHandler(f func(c *Client)) {
	if f != nil {
		sf.onConnectionLost = f
	}
}

// Start start the server,and return quickly,if it nil,the server will disconnected background,other failed
func (sf *Client) Start() error {
	if sf.server == nil {
		return errors.New("empty remote server")
	}

	go sf.running()
	return nil
}

// Connect is
func (sf *Client) running() {
	var ctx context.Context

	sf.rwMux.Lock()
	if !atomic.CompareAndSwapUint32(&sf.status, initial, disconnected) {
		sf.rwMux.Unlock()
		return
	}
	ctx, sf.closeCancel = context.WithCancel(context.Background())
	sf.rwMux.Unlock()
	defer sf.setConnectStatus(initial)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		sf.Debug("connecting server %+v", sf.server)
		conn, err := openConnection(sf.server, sf.TLSConfig, sf.conf.ConnectTimeout0)
		if err != nil {
			sf.Error("connect failed, %v", err)
			if !sf.autoReconnect {
				return
			}
			time.Sleep(sf.reconnectInterval)
			continue
		}
		sf.Debug("connect success")
		sf.conn = conn
		if err = sf.onConnect(sf); err != nil {
			time.Sleep(sf.reconnectInterval)
			continue
		}
		sf.run(ctx)
		sf.onConnectionLost(sf)
		sf.Debug("disconnected server %+v", sf.server)
		select {
		case <-ctx.Done():
			return
		default:
			// 随机500ms-1s的重试，避免快速重试造成服务器许多无效连接
			time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(500)))
		}
	}
}

func (sf *Client) recvLoop() {
	sf.Debug("recvLoop started")
	defer func() {
		sf.cancel()
		sf.wg.Done()
		sf.Debug("recvLoop stopped")
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
				if rdCnt == 0 && err == io.EOF {
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
					sf.rcvRaw <- apdu
				}
			}
		}
	}
}

func (sf *Client) sendLoop() {
	sf.Debug("sendLoop started")
	defer func() {
		sf.cancel()
		sf.wg.Done()
		sf.Debug("sendLoop stopped")
	}()
	for {
		select {
		case <-sf.ctx.Done():
			return
		case apdu := <-sf.sendRaw:
			sf.Debug("TX Raw[% x]", apdu)
			for wrCnt := 0; len(apdu) > wrCnt; {
				byteCount, err := sf.conn.Write(apdu[wrCnt:])
				if err != nil {
					// See: https://github.com/golang/go/issues/4373
					if err != io.EOF && err != io.ErrClosedPipe ||
						strings.Contains(err.Error(), "use of closed network connection") {
						sf.Error("sendRaw failed, %v", err)
						return
					}
					if e, ok := err.(net.Error); !ok || !e.Temporary() {
						sf.Error("sendRaw failed, %v", err)
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
func (sf *Client) run(ctx context.Context) {
	sf.Debug("run started!")
	// before any thing make sure init
	sf.cleanUp()

	sf.ctx, sf.cancel = context.WithCancel(ctx)
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

	sf.startDtActiveSendSince.Store(willNotTimeout)
	sf.stopDtActiveSendSince.Store(willNotTimeout)

	sendSFrame := func(rcvSN uint16) {
		sf.Debug("TX sFrame %v", sAPCI{rcvSN})
		sf.sendRaw <- newSFrame(rcvSN)
	}

	sendIFrame := func(asdu1 []byte) {
		seqNo := sf.seqNoSend

		iframe, err := newIFrame(seqNo, sf.seqNoRcv, asdu1)
		if err != nil {
			return
		}
		sf.ackNoRcv = sf.seqNoRcv
		sf.seqNoSend = (seqNo + 1) & 32767
		sf.pending = append(sf.pending, seqPending{seqNo & 32767, time.Now()})

		sf.Debug("TX iFrame %v", iAPCI{seqNo, sf.seqNoRcv})
		sf.sendRaw <- iframe
	}

	defer func() {
		sf.setConnectStatus(disconnected)
		checkTicker.Stop()
		_ = sf.conn.Close() // 连锁引发cancel
		sf.wg.Wait()
		sf.Debug("run stopped!")
	}()

	sf.SendStartDt() // 发送startDt激活指令
	for {
		if isActive && seqNoCount(sf.ackNoSend, sf.seqNoSend) <= sf.conf.SendUnAckLimitK {
			select {
			case o := <-sf.sendASDU:
				sendIFrame(o)
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
			if now.Sub(testFrAliveSendSince) >= sf.conf.SendUnAckTimeout1 ||
				now.Sub(sf.startDtActiveSendSince.Load().(time.Time)) >= sf.conf.SendUnAckTimeout1 ||
				now.Sub(sf.stopDtActiveSendSince.Load().(time.Time)) >= sf.conf.SendUnAckTimeout1 {
				sf.Error("test frame alive confirm timeout t₁")
				return
			}
			// check oldest unacknowledged outbound
			if sf.ackNoSend != sf.seqNoSend &&
				//now.Sub(sf.peek()) >= sf.SendUnAckTimeout1 {
				now.Sub(sf.pending[0].sendTime) >= sf.conf.SendUnAckTimeout1 {
				sf.ackNoSend++
				sf.Error("fatal transmission timeout t₁")
				return
			}

			// 确定最早发送的i-Frame是否超时,超时则回复sFrame
			if sf.ackNoRcv != sf.seqNoRcv &&
				(now.Sub(unAckRcvSince) >= sf.conf.RecvUnAckTimeout2 ||
					now.Sub(idleTimeout3Sine) >= timeoutResolution) {
				sendSFrame(sf.seqNoRcv)
				sf.ackNoRcv = sf.seqNoRcv
			}

			// 空闲时间到，发送TestFrActive帧,保活
			if now.Sub(idleTimeout3Sine) >= sf.conf.IdleTimeout3 {
				sf.sendUFrame(uTestFrActive)
				testFrAliveSendSince = time.Now()
				idleTimeout3Sine = testFrAliveSendSince
			}

		case apdu := <-sf.rcvRaw:
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
				if !sf.updateAckNoOut(head.rcvSN) || head.sendSN != sf.seqNoRcv {
					sf.Error("fatal incoming acknowledge either earlier than previous or later than sendTime")
					return
				}

				sf.rcvASDU <- asduVal
				if sf.ackNoRcv == sf.seqNoRcv { // first unacked
					unAckRcvSince = time.Now()
				}

				sf.seqNoRcv = (sf.seqNoRcv + 1) & 32767
				if seqNoCount(sf.ackNoRcv, sf.seqNoRcv) >= sf.conf.RecvUnAckLimitW {
					sendSFrame(sf.seqNoRcv)
					sf.ackNoRcv = sf.seqNoRcv
				}

			case uAPCI:
				sf.Debug("RX uFrame %v", head)
				switch head.function {
				//case uStartDtActive:
				//	sf.sendUFrame(uStartDtConfirm)
				//	isActive = true
				case uStartDtConfirm:
					isActive = true
					sf.startDtActiveSendSince.Store(willNotTimeout)
				//case uStopDtActive:
				//	sf.sendUFrame(uStopDtConfirm)
				//	isActive = false
				case uStopDtConfirm:
					isActive = false
					sf.stopDtActiveSendSince.Store(willNotTimeout)
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

func (sf *Client) handlerLoop() {
	sf.Debug("handlerLoop started")
	defer func() {
		sf.wg.Done()
		sf.Debug("handlerLoop stopped")
	}()

	for {
		select {
		case <-sf.ctx.Done():
			return
		case rawAsdu := <-sf.rcvASDU:
			asduPack := asdu.NewEmptyASDU(&sf.param)
			if err := asduPack.UnmarshalBinary(rawAsdu); err != nil {
				sf.Warn("asdu UnmarshalBinary failed,%+v", err)
				continue
			}
			if err := sf.clientHandler(asduPack); err != nil {
				sf.Warn("Falied handling I frame, error: %v", err)
			}
		}
	}
}

func (sf *Client) setConnectStatus(status uint32) {
	sf.rwMux.Lock()
	atomic.StoreUint32(&sf.status, status)
	sf.rwMux.Unlock()
}

func (sf *Client) connectStatus() uint32 {
	sf.rwMux.RLock()
	status := atomic.LoadUint32(&sf.status)
	sf.rwMux.RUnlock()
	return status
}

func (sf *Client) cleanUp() {
	sf.ackNoRcv = 0
	sf.ackNoSend = 0
	sf.seqNoRcv = 0
	sf.seqNoSend = 0
	sf.pending = nil
	// clear sending chan buffer
loop:
	for {
		select {
		case <-sf.sendRaw:
		case <-sf.rcvRaw:
		case <-sf.rcvASDU:
		case <-sf.sendASDU:
		default:
			break loop
		}
	}
}

func (sf *Client) sendUFrame(which byte) {
	sf.Debug("TX uFrame %v", uAPCI{which})
	sf.sendRaw <- newUFrame(which)
}

func (sf *Client) updateAckNoOut(ackNo uint16) (ok bool) {
	if ackNo == sf.ackNoSend {
		return true
	}
	// new acks validate， ack 不能在 req seq 前面,出错
	if seqNoCount(sf.ackNoSend, sf.seqNoSend) < seqNoCount(ackNo, sf.seqNoSend) {
		return false
	}

	// confirm reception
	for i, v := range sf.pending {
		if v.seq == (ackNo - 1) {
			sf.pending = sf.pending[i+1:]
			break
		}
	}

	sf.ackNoSend = ackNo
	return true
}

// IsConnected get server session connected state
func (sf *Client) IsConnected() bool {
	return sf.connectStatus() == connected
}

// 接收到I帧后根据TYPEID进行不同的处理,分别调用对应的接口函数
// TODO: fix response handler
func (sf *Client) clientHandler(asduPack *asdu.ASDU) error {
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

// Params returns params of client
func (sf *Client) Params() *asdu.Params {
	return &sf.param
}

// Send send asdu
func (sf *Client) Send(a *asdu.ASDU) error {
	if !sf.IsConnected() {
		return ErrUseClosedConnection
	}
	//if !sf.isServerActive {
	//	return fmt.Errorf("ErrorUnactive")
	//}
	data, err := a.MarshalBinary()
	if err != nil {
		return err
	}
	select {
	case sf.sendASDU <- data:
	default:
		return ErrBufferFulled
	}
	return nil
}

// UnderlyingConn returns underlying conn of client
func (sf *Client) UnderlyingConn() net.Conn {
	return sf.conn
}

// Close close all
func (sf *Client) Close() error {
	sf.rwMux.Lock()
	if sf.closeCancel != nil {
		sf.closeCancel()
	}
	sf.rwMux.Unlock()
	return nil
}

// SendStartDt start data transmission on this connection
func (sf *Client) SendStartDt() {
	sf.startDtActiveSendSince.Store(time.Now())
	sf.sendUFrame(uStartDtActive)
}

// SendStopDt stop data transmission on this connection
func (sf *Client) SendStopDt() {
	sf.stopDtActiveSendSince.Store(time.Now())
	sf.sendUFrame(uStopDtActive)
}

//InterrogationCmd wrap asdu.InterrogationCmd
func (sf *Client) InterrogationCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, qoi asdu.QualifierOfInterrogation) error {
	return asdu.InterrogationCmd(sf, coa, ca, qoi)
}

// CounterInterrogationCmd wrap asdu.CounterInterrogationCmd
func (sf *Client) CounterInterrogationCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, qcc asdu.QualifierCountCall) error {
	return asdu.CounterInterrogationCmd(sf, coa, ca, qcc)
}

// ClockSynchronizationCmd wrap asdu.ClockSynchronizationCmd
func (sf *Client) ClockSynchronizationCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, t time.Time) error {
	return asdu.ClockSynchronizationCmd(sf, coa, ca, t)
}

// ResetProcessCmd wrap asdu.ResetProcessCmd
func (sf *Client) ResetProcessCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, qrp asdu.QualifierOfResetProcessCmd) error {
	return asdu.ResetProcessCmd(sf, coa, ca, qrp)
}

// TestCommand  wrap asdu.TestCommand
func (sf *Client) TestCommand(coa asdu.CauseOfTransmission, ca asdu.CommonAddr) error {
	return asdu.TestCommand(sf, coa, ca)
}

// TestCommandCP56Time2a send test command [C_TS_TA_1]，测试命令, 只有单个信息对象(SQ = 0)
func (sf *Client) TestCommandCP56Time2a(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, t time.Time) error {
	if err := sf.Params().Valid(); err != nil {
		return err
	}
	u := asdu.NewASDU(sf.Params(), asdu.Identifier{
		asdu.C_TS_TA_1,
		asdu.VariableStruct{IsSequence: false, Number: 1},
		coa,
		0,
		ca,
	})
	if err := u.AppendInfoObjAddr(asdu.InfoObjAddrIrrelevant); err != nil {
		return err
	}
	u.AppendBytes(byte(asdu.FBPTestWord&0xff), byte(asdu.FBPTestWord>>8))
	u.AppendBytes(asdu.CP56Time2a(t, u.InfoObjTimeZone)...)
	return sf.Send(u)
}
