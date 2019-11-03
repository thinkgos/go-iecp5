package cs104

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
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
	handler ClientHandler // 接口

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

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	server            *url.URL      // 连接的服务器端
	autoReconnect     bool          // 是否启动重连
	reconnectInterval time.Duration // 重连间隔时间
	TLSConfig         *tls.Config   // tls配置
	onConnect         func(c *Client) error
	onConnectionLost  func(c *Client)
}

// NewClient returns an IEC104 master
func NewClient(conf *Config, params *asdu.Params, handler ClientHandler) (*Client, error) {
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
		Clog:     clog.NewWithPrefix("IEC104 client =>"),

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
	ctx, sf.cancel = context.WithCancel(context.Background())
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
			if err := sf.handleIFrame(asduPack); err != nil {
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
func (sf *Client) handleIFrame(a *asdu.ASDU) error {
	defer func() {
		if err := recover(); err != nil {
			sf.Critical("Client handler %+v", err)
		}
	}()

	sf.Debug("ASDU %+v", a)

	// check common addr
	if a.CommonAddr == asdu.InvalidCommonAddr {
		return a.SendReplyMirror(sf, asdu.UnknownCA)
	}

	if a.Identifier.Coa.Cause == asdu.UnknownTypeID ||
		a.Identifier.Coa.Cause == asdu.UnknownCOT ||
		a.Identifier.Coa.Cause == asdu.UnknownCA ||
		a.Identifier.Coa.Cause == asdu.UnknownIOA {
		return fmt.Errorf("GOT COT %v", a.Identifier.Coa.Cause)
	}

	switch a.Identifier.Type {
	case asdu.M_SP_NA_1, asdu.M_SP_TA_1, asdu.M_SP_TB_1: // 遥信 单点信息 01 02 30
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.Background ||
			a.Identifier.Coa.Cause == asdu.Spontaneous ||
			a.Identifier.Coa.Cause == asdu.Request ||
			a.Identifier.Coa.Cause == asdu.ReturnInfoRemote ||
			a.Identifier.Coa.Cause == asdu.ReturnInfoLocal ||
			a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
			(a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
				a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16)) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		info := a.GetSinglePoint()
		sf.handler.Handle01_02_1e(sf, a, info)
	case asdu.M_DP_NA_1, asdu.M_DP_TA_1, asdu.M_DP_TB_1: // 遥信 双点信息 3,4,31
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.Background ||
			a.Identifier.Coa.Cause == asdu.Spontaneous ||
			a.Identifier.Coa.Cause == asdu.Request ||
			a.Identifier.Coa.Cause == asdu.ReturnInfoRemote ||
			a.Identifier.Coa.Cause == asdu.ReturnInfoLocal ||
			a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
			(a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
				a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16)) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		info := a.GetDoublePoint()
		sf.handler.Handle03_04_1f(sf, a, info)
	case asdu.M_ST_NA_1, asdu.M_ST_TB_1: // 遥信 步调节信息 5,32
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.Background ||
			a.Identifier.Coa.Cause == asdu.Spontaneous ||
			a.Identifier.Coa.Cause == asdu.Request ||
			a.Identifier.Coa.Cause == asdu.ReturnInfoRemote ||
			a.Identifier.Coa.Cause == asdu.ReturnInfoLocal ||
			a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
			(a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
				a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16)) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		info := a.GetStepPosition()
		sf.handler.Handle05_20(sf, a, info)
	case asdu.M_BO_NA_1, asdu.M_BO_TA_1, asdu.M_BO_TB_1: // 遥信 比特串信息 07,08,33								// 比特串,07
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.Background ||
			a.Identifier.Coa.Cause == asdu.Spontaneous ||
			a.Identifier.Coa.Cause == asdu.Request ||
			a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
			(a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
				a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16)) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		info := a.GetBitString32()
		sf.handler.Handle07_08_21(sf, a, info)
	case asdu.M_ME_NA_1, asdu.M_ME_TA_1, asdu.M_ME_TD_1, asdu.M_ME_ND_1: // 遥测 归一化测量值 09,10,21,34
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.Periodic ||
			a.Identifier.Coa.Cause == asdu.Background ||
			a.Identifier.Coa.Cause == asdu.Spontaneous ||
			a.Identifier.Coa.Cause == asdu.Request ||
			a.Identifier.Coa.Cause == asdu.InterrogatedByStation) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		value := a.GetMeasuredValueNormal()
		sf.handler.Handle09_0a_15_22(sf, a, value)
	case asdu.M_ME_NB_1, asdu.M_ME_TB_1, asdu.M_ME_TE_1: //遥测 标度化值 11,12,35
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.Periodic ||
			a.Identifier.Coa.Cause == asdu.Background ||
			a.Identifier.Coa.Cause == asdu.Spontaneous ||
			a.Identifier.Coa.Cause == asdu.Request ||
			a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
			(a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
				a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16)) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		value := a.GetMeasuredValueScaled()
		sf.handler.Handle0b_0c_23(sf, a, value)
	case asdu.M_ME_NC_1, asdu.M_ME_TC_1, asdu.M_ME_TF_1: // 遥信 短浮点数 13,14,16
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.Periodic ||
			a.Identifier.Coa.Cause == asdu.Background ||
			a.Identifier.Coa.Cause == asdu.Spontaneous ||
			a.Identifier.Coa.Cause == asdu.Request ||
			a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
			(a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
				a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16)) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		value := a.GetMeasuredValueFloat()
		sf.handler.Handle0d_0e_10(sf, a, value)
	case asdu.M_EI_NA_1: // 站初始化结束 70
		// check cause of transmission
		if !(a.Identifier.Coa.Cause == asdu.Initialized) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		ioa, coi := a.GetEndOfInitialization()
		if ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		sf.handler.Handle46(sf, coi)
	case asdu.C_IC_NA_1: // 总召唤 100
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.ActivationCon ||
			a.Identifier.Coa.Cause == asdu.DeactivationCon ||
			a.Identifier.Coa.Cause == asdu.ActivationTerm) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		// get ioa and qoi
		ioa, qua := a.GetInterrogationCmd()
		// check ioa
		if ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		sf.handler.Handle64(sf, a, qua)
	case asdu.C_CI_NA_1: // 计数量召唤 101
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.ActivationCon ||
			a.Identifier.Coa.Cause == asdu.DeactivationCon ||
			a.Identifier.Coa.Cause == asdu.ActivationTerm) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		// get ioa and qoi
		ioa, qua := a.GetCounterInterrogationCmd()
		// check ioa
		if ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		sf.handler.Handle65(sf, a, qua)
	case asdu.C_CS_NA_1: // 时钟同步 103
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.ActivationCon ||
			a.Identifier.Coa.Cause == asdu.ActivationTerm) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		ioa, t := a.GetClockSynchronizationCmd()
		// check ioa
		if ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		sf.handler.Handle67(sf, a, t)
	case asdu.C_RP_NA_1: // 复位进程 105
		// check cot
		if !(a.Identifier.Coa.Cause == asdu.ActivationCon) {
			return a.SendReplyMirror(sf, asdu.UnknownCOT)
		}
		ioa, qua := a.GetResetProcessCmd()
		// check ioa
		if ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(sf, asdu.UnknownIOA)
		}
		sf.handler.Handle69(sf, a, qua)
	// case asdu.C_TS_TA_1:										// 测试命令 107
	// 	// check cot
	// 	if !( a.Identifier.Coa.Cause == asdu.ActivationCon) {
	// 	  	return a.SendReplyMirror(sf, asdu.UnknownCOT)
	// 	}
	// 	sf.handler.Handle6b(sf, a, true)
	default:
		return a.SendReplyMirror(sf, asdu.UnknownTypeID)
	}

	// if err := sf.handler.ASDUHandler(sf, a); err != nil {
	// 	return a.SendReplyMirror(sf, asdu.UnknownTypeID)
	// }
	return nil
}

// ClientHandler is
type ClientHandler interface {
	// 01:[M_SP_NA_1] 不带时标单点信息
	// 02:[M_SP_TA_1] 带时标CP24Time2a的单点信息,只有(SQ = 0)单个信息元素集合
	// 1e:[M_SP_TB_1] 带时标CP56Time2a的单点信息,只有(SQ = 0)单个信息元素集合
	Handle01_02_1e(asdu.Connect, *asdu.ASDU, []asdu.SinglePointInfo)
	// 03:[M_DP_NA_1].双点信息
	// 04:[M_DP_TA_1] .带CP24Time2a双点信息,只有(SQ = 0)单个信息元素集合
	// 1f:[M_DP_TB_1].带CP56Time2a的双点信息,只有(SQ = 0)单个信息元素集合
	Handle03_04_1f(asdu.Connect, *asdu.ASDU, []asdu.DoublePointInfo)
	// 07:[M_BO_NA_1] 比特位串
	// 08:[M_BO_TA_1] 带时标CP24Time2a比特位串，只有(SQ = 0)单个信息元素集合
	// 21:[M_BO_TB_1] 带时标CP56Time2a比特位串，只有(SQ = 0)单个信息元素集
	Handle05_20(asdu.Connect, *asdu.ASDU, []asdu.StepPositionInfo)
	// 07:[M_BO_NA_1] 比特位串
	// 08:[M_BO_TA_1] 带时标CP24Time2a比特位串，只有(SQ = 0)单个信息元素集合
	// 21:[M_BO_TB_1] 带时标CP56Time2a比特位串，只有(SQ = 0)单个信息元素集
	Handle07_08_21(asdu.Connect, *asdu.ASDU, []asdu.BitString32Info)
	// 09:[M_ME_NA_1] 测量值,规一化值
	// 0a:[M_ME_TA_1] 带时标CP24Time2a的测量值,规一化值,只有(SQ = 0)单个信息元素集合
	// 15:[M_ME_ND_1] 不带品质的测量值,规一化值
	// 22:[M_ME_TD_1] 带时标CP57Time2a的测量值,规一化值,只有(SQ = 0)单个信息元素集合
	Handle09_0a_15_22(asdu.Connect, *asdu.ASDU, []asdu.MeasuredValueNormalInfo)
	// 0b:[M_ME_NB_1].测量值,标度化值
	// 0c:[M_ME_TB_1].带时标CP24Time2a的测量值,标度化值,只有(SQ = 0)单个信息元素集合
	// 23:[M_ME_TE_1].带时标CP56Time2a的测量值,标度化值,只有(SQ = 0)单个信息元素集合
	Handle0b_0c_23(asdu.Connect, *asdu.ASDU, []asdu.MeasuredValueScaledInfo)
	// 0d:[M_ME_TF_1] 测量值,短浮点数
	// 0e:[M_ME_TC_1].带时标CP24Time2a的测量值,短浮点数,只有(SQ = 0)单个信息元素集合
	// 10:[M_ME_TF_1].带时标CP56Time2a的测量值,短浮点数,只有(SQ = 0)单个信息元素集合
	Handle0d_0e_10(asdu.Connect, *asdu.ASDU, []asdu.MeasuredValueFloatInfo)
	// 46:[M_EI_NA_1] 站初始化结束
	Handle46(asdu.Connect, asdu.CauseOfInitial)
	// 64:[C_IC_NA_1] 总召唤
	Handle64(asdu.Connect, *asdu.ASDU, asdu.QualifierOfInterrogation)
	// 65:[C_CI_NA_1] 计数量召唤
	Handle65(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall)
	// 67:[C_CS_NA_1] 时钟同步
	Handle67(asdu.Connect, *asdu.ASDU, time.Time)
	// 69:[C_RP_NA_1] 复位进程
	Handle69(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd)
	// 6B:[C_TS_TA_1] 测试命令
	Handle6b(asdu.Connect, *asdu.ASDU, bool)
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
	if sf.cancel != nil {
		sf.cancel()
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
