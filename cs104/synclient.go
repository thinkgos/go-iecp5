package cs104

import (
	"context"
	// "errors"
	"fmt"
	"io"

	// "math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

const (
	// Timeout for syncSendIFrame
	syncSendTimeout = time.Second
)

type AsduInfo struct {
	asdu.Identifier
	Ioa       asdu.InfoObjAddr
	Value     interface{}
	Quality   asdu.QualityDescriptor
	Timestamp time.Time
}

// Client is an IEC104 master
type HighLevelClient struct {
	option           ClientOption
	conn             net.Conn
	responseHandler  map[uint64]chan *asdu.ASDU
	subscriptionChan chan AsduInfo

	// channel
	rcvAsdu chan *asdu.ASDU
	rcvRaw  chan []byte // for recvLoop raw cs104 frame
	sendRaw chan []byte // for sendLoop raw cs104 frame

	// I Frame send and receive sequence
	seqNoSend uint16 // sequence number of next outbound I-frame
	ackNoSend uint16 // outbound sequence number yet to be confirmed
	seqNoRcv  uint16 // sequence number of next inbound I-frame
	ackNoRcv  uint16 // inbound sequence number yet to be confirmed
	// maps sendTime I-frames to their respective sequence number
	pending []seqPending

	// 连接状态
	status   uint32
	rwMux    sync.RWMutex
	isActive uint32

	// 其他

	clog.Clog

	wg sync.WaitGroup

	onConnect        func()
	onConnectionLost func()
}

// NewClient returns an IEC104 master,default config and default asdu.ParamsWide params
func NewHighLevelClient(o *ClientOption) *HighLevelClient {
	return &HighLevelClient{
		option:           *o,
		responseHandler:  make(map[uint64]chan *asdu.ASDU),
		rcvAsdu:          make(chan *asdu.ASDU, o.config.RecvUnAckLimitW<<5),
		rcvRaw:           make(chan []byte, o.config.RecvUnAckLimitW<<5),
		sendRaw:          make(chan []byte, o.config.SendUnAckLimitK<<5), // may not block!
		Clog:             clog.NewLogger("cs104 client => "),
		onConnect:        func() {},
		onConnectionLost: func() {},
	}
}

// SetOption set the client option
func (sf *HighLevelClient) SetOption(o *ClientOption) {
	sf.option = *o
}

// SetOnConnectHandler set on connect handler
func (sf *HighLevelClient) SetOnConnectHandler(f func()) {
	if f != nil {
		sf.onConnect = f
	}
}

// SetConnectionLostHandler set connection lost handler
func (sf *HighLevelClient) SetConnectionLostHandler(f func()) {
	if f != nil {
		sf.onConnectionLost = f
	}
}

// IsConnected get server session connected state
func (sf *HighLevelClient) IsConnected() bool {
	return sf.connectStatus() == connected
}

// Params returns params of client
func (sf *HighLevelClient) Params() *asdu.Params {
	return &sf.option.param
}

// UnderlyingConn returns underlying conn of client
func (sf *HighLevelClient) UnderlyingConn() net.Conn {
	return sf.conn
}

// // Close close all
// func (sf *HighLevelClient) Close() error {
// 	sf.rwMux.Lock()
// 	if sf.closeCancel != nil {
// 		sf.closeCancel()
// 	}
// 	sf.rwMux.Unlock()
// 	return nil
// }

func (sf *HighLevelClient) Subscribe(sub chan AsduInfo) {
	sf.rwMux.Lock()
	sf.subscriptionChan = sub
	sf.rwMux.Unlock()
}

func (sf *HighLevelClient) Write(id asdu.TypeID, ca asdu.CommonAddr, ioa asdu.InfoObjAddr, v interface{}) error {

	// TODO: QualifierOfCommand Ignored
	var qoc byte
	// TODO: QualifierOfSetpointCmd Ignored
	var qos byte

	asduPack := asdu.NewASDU(sf.Params(), asdu.Identifier{
		Type:       id,
		Variable:   asdu.VariableStruct{IsSequence: false, Number: 1},
		Coa:        asdu.ParseCauseOfTransmission(byte(asdu.Activation)),
		OrigAddr:   0,
		CommonAddr: ca,
	})
	if err := asduPack.AppendInfoObjAddr(ioa); err != nil {
		return err
	}

	switch id {
	case asdu.C_SC_NA_1, asdu.C_SC_TA_1:
		if vv, ok := v.(bool); ok {
			if vv {
				asduPack.AppendBytes(qoc | 0x01)
			} else {
				asduPack.AppendBytes(qoc | 0x00)
			}
		} else {
			return fmt.Errorf("Should provide value in boolean type")
		}

		if id == asdu.C_SC_TA_1 {
			asduPack.AppendBytes(asdu.CP56Time2a(time.Now(), asduPack.InfoObjTimeZone)...)
		}
	case asdu.C_DC_NA_1, asdu.C_DC_TA_1:
		if vv, ok := v.(byte); ok {
			asduPack.AppendBytes(qoc | byte(vv&0x03))
		} else {
			return fmt.Errorf("Should provide value in byte type ")
		}

		if id == asdu.C_DC_TA_1 {
			asduPack.AppendBytes(asdu.CP56Time2a(time.Now(), asduPack.InfoObjTimeZone)...)
		}
	case asdu.C_RC_NA_1, asdu.C_RC_TA_1:
		if vv, ok := v.(byte); ok {
			asduPack.AppendBytes(qoc | byte(vv&0x03))
		} else {
			return fmt.Errorf("Should provide value in byte type ")
		}
		if id == asdu.C_RC_TA_1 {
			asduPack.AppendBytes(asdu.CP56Time2a(time.Now(), asduPack.InfoObjTimeZone)...)
		}
	case asdu.C_SE_NA_1, asdu.C_SE_TA_1:
		if vv, ok := v.(int16); ok {
			asduPack.AppendNormalize(asdu.Normalize(vv)).AppendBytes(qos)
		} else {
			return fmt.Errorf("Should provide value in int16 type ")
		}
		if id == asdu.C_SE_TA_1 {
			asduPack.AppendBytes(asdu.CP56Time2a(time.Now(), asduPack.InfoObjTimeZone)...)
		}
	case asdu.C_SE_NB_1, asdu.C_SE_TB_1:
		if vv, ok := v.(int16); ok {
			asduPack.AppendScaled(vv).AppendBytes(qos)
		} else {
			return fmt.Errorf("Should provide value in int16 type ")
		}
		if id == asdu.C_SE_TB_1 {
			asduPack.AppendBytes(asdu.CP56Time2a(time.Now(), asduPack.InfoObjTimeZone)...)
		}
	case asdu.C_SE_NC_1, asdu.C_SE_TC_1:
		if vv, ok := v.(float32); ok {
			asduPack.AppendFloat32(vv).AppendBytes(qos)
		} else {
			return fmt.Errorf("Should provide value in float32 type ")
		}
		if id == asdu.C_SE_TC_1 {
			asduPack.AppendBytes(asdu.CP56Time2a(time.Now(), asduPack.InfoObjTimeZone)...)
		}
	case asdu.C_BO_NA_1, asdu.C_BO_TA_1:
		if vv, ok := v.(uint32); ok {
			asduPack.AppendBitsString32(vv)
		} else {
			return fmt.Errorf("Should provide value in float32 type ")
		}
		if id == asdu.C_BO_TA_1 {
			asduPack.AppendBytes(asdu.CP56Time2a(time.Now(), asduPack.InfoObjTimeZone)...)
		}
	}

	// this uID is used to matching response packet to request packet
	uID := uint64(ioa) + (uint64(ca) << (8 * sf.option.param.InfoObjAddrSize))

	// Request target on ioa which is already in operating state will be rejected
	if _, ok := sf.responseHandler[uID]; ok {
		return fmt.Errorf("Last Read Request has not completed")
	}
	ch := make(chan *asdu.ASDU)
	sf.responseHandler[uID] = ch
	defer func() {
		sf.rwMux.Lock()
		delete(sf.responseHandler, uID)
		sf.rwMux.Unlock()
	}()
	resp, err := sf.syncSendIFrame(asduPack, ch)
	if err != nil {
		return err
	}

	if resp.Coa.Cause == asdu.ActivationCon {
		return nil
	}
	return fmt.Errorf("Unexpected CauseOfTransmission: %v", resp.Coa.Cause)
}

// Read executes a synchronous read request
func (sf *HighLevelClient) Read(ca asdu.CommonAddr, ioa asdu.InfoObjAddr) (*AsduInfo, error) {
	asduPack := asdu.NewASDU(sf.Params(), asdu.Identifier{
		Type:       asdu.C_RD_NA_1,
		Variable:   asdu.VariableStruct{IsSequence: false, Number: 1},
		Coa:        asdu.ParseCauseOfTransmission(byte(asdu.Request)),
		OrigAddr:   0,
		CommonAddr: ca,
	})
	if err := asduPack.AppendInfoObjAddr(ioa); err != nil {
		return nil, err
	}

	// this id is used to matching response packet to request packet
	uID := uint64(ioa) + (uint64(ca) << (8 * sf.option.param.InfoObjAddrSize))

	// Request target on ioa which is already in operating state will be rejected
	if _, ok := sf.responseHandler[uID]; ok {
		return nil, fmt.Errorf("Last Read on ioa: %v has not completed", ioa)
	}

	ch := make(chan *asdu.ASDU)
	sf.responseHandler[uID] = ch
	defer func() {
		sf.rwMux.Lock()
		delete(sf.responseHandler, uID)
		sf.rwMux.Unlock()
	}()
	resp, err := sf.syncSendIFrame(asduPack, ch)
	if err != nil {
		return nil, err
	}
	data := &AsduInfo{}
	switch resp.Type {
	case asdu.M_SP_NA_1:
		// Single Info
		v := resp.GetSinglePoint()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_SP_TB_1:
		// Single Info with time
		v := resp.GetSinglePoint()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
		data.Timestamp = v.Time
	case asdu.M_DP_NA_1:
		// Double Info
		v := resp.GetDoublePoint()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_DP_TB_1:
		// Double Info with time
		v := resp.GetDoublePoint()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
		data.Timestamp = v.Time
	case asdu.M_ST_NA_1:
		// Step Position Info
		v := resp.GetStepPosition()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_ST_TB_1:
		// Step Position Info with time
		v := resp.GetStepPosition()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
		data.Timestamp = v.Time
	case asdu.M_BO_NA_1:
		// 32 Bit string
		v := resp.GetBitString32()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_BO_TB_1:
		// 32 Bit string with time
		v := resp.GetBitString32()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
		data.Timestamp = v.Time
	case asdu.M_ME_NA_1:
		// Normalized Measured Value
		v := resp.GetMeasuredValueNormal()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_ME_TD_1:
		// Normalized Measured Value with time
		v := resp.GetMeasuredValueNormal()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
		data.Timestamp = v.Time
	case asdu.M_ME_ND_1:
		// Normalized Measured Value without quality description
		v := resp.GetMeasuredValueNormal()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_ME_NB_1:
		// Scaled Measured Value
		v := resp.GetMeasuredValueScaled()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_ME_TE_1:
		// Scaled Measured Value with time
		v := resp.GetMeasuredValueScaled()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
		data.Timestamp = v.Time
	case asdu.M_ME_NC_1:
		// Short Float Measured Value
		v := resp.GetMeasuredValueFloat()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
	case asdu.M_ME_TF_1:
		// Short Float Measured Value with time
		v := resp.GetMeasuredValueFloat()[0]
		if v.Qds != asdu.QDSGood {
			return nil, fmt.Errorf("Quality not Good: %v", v.Qds)
		}
		data.Value = v.Value()
		data.Timestamp = v.Time
	default:
		return nil, fmt.Errorf("TypeID: %v Not Supported", resp.Type)
	}
	return data, nil
}

//InterrogationCmd wrap asdu.InterrogationCmd
func (sf *HighLevelClient) InterrogationCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, qoi asdu.QualifierOfInterrogation) error {
	return asdu.InterrogationCmd(sf, coa, ca, qoi)
}

// CounterInterrogationCmd wrap asdu.CounterInterrogationCmd
func (sf *HighLevelClient) CounterInterrogationCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, qcc asdu.QualifierCountCall) error {
	return asdu.CounterInterrogationCmd(sf, coa, ca, qcc)
}

// ClockSynchronizationCmd wrap asdu.ClockSynchronizationCmd
func (sf *HighLevelClient) ClockSynchronizationCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, t time.Time) error {
	return asdu.ClockSynchronizationCmd(sf, coa, ca, t)
}

// ResetProcessCmd wrap asdu.ResetProcessCmd
func (sf *HighLevelClient) ResetProcessCmd(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, qrp asdu.QualifierOfResetProcessCmd) error {
	return asdu.ResetProcessCmd(sf, coa, ca, qrp)
}

// DelayAcquireCommand wrap asdu.DelayAcquireCommand
func (sf *HighLevelClient) DelayAcquireCommand(coa asdu.CauseOfTransmission, ca asdu.CommonAddr, msec uint16) error {
	return asdu.DelayAcquireCommand(sf, coa, ca, msec)
}

// TestCommand  wrap asdu.TestCommand
func (sf *HighLevelClient) TestCommand(coa asdu.CauseOfTransmission, ca asdu.CommonAddr) error {
	return asdu.TestCommand(sf, coa, ca)
}

func (sf *HighLevelClient) Connect(ctx context.Context) {
	defer sf.setConnectStatus(initial)
	if !atomic.CompareAndSwapUint32(&sf.status, initial, disconnected) {
		return
	}

	go func() {
		waitChan := make(chan struct{}, 1)
		for {
			sf.Debug("connecting server %+v", sf.option.server)
			conn, err := openConnection(sf.option.server, sf.option.TLSConfig, sf.option.config.ConnectTimeout0)
			if err != nil {
				sf.Error("connect failed, %v", err)
				if !sf.option.autoReconnect {
					return
				}
				go func() {
					time.Sleep(sf.option.reconnectInterval)
					waitChan <- struct{}{}
				}()
				select {
				case <-ctx.Done():
					return
				case <-waitChan:
					continue
				}
			}
			sf.Debug("connect success")
			sf.conn = conn
			sf.run(ctx)
			sf.conn.Close()

			if !sf.option.autoReconnect {
				return
			}
			go func() {
				time.Sleep(sf.option.reconnectInterval)
				waitChan <- struct{}{}
			}()
			select {
			case <-ctx.Done():
				return
			case <-waitChan:
				continue
			}
		}
	}()
}

// run is the big fat state machine.
func (sf *HighLevelClient) run(ctx context.Context) {
	// before any thing make sure init
	sf.cleanUp()
	sf.Debug("Connected server %+v", sf.option.server)
	defer sf.Debug("disconnected server %+v", sf.option.server)

	runCtx, runCancel := context.WithCancel(ctx)
	defer runCancel()

	sf.wg.Add(2)
	defer sf.wg.Wait()
	go sf.recvLoop(runCtx, runCancel)
	go sf.subscribeLoop(runCtx)

	var checkTicker = time.NewTicker(timeoutResolution)
	defer checkTicker.Stop()

	// Will not timeout until 100 years
	var willNotTimeout = time.Now().AddDate(100, 0, 0)

	var unAckRcvSince = willNotTimeout
	var idleTimeout3Sine = time.Now()         // 空闲间隔发起testFrAlive
	var testFrAliveSendSince = willNotTimeout // 当发起testFrAlive时,等待确认回复的超时间隔

	var startDtActiveSendSince atomic.Value
	var stopDtActiveSendSince atomic.Value
	startDtActiveSendSince.Store(willNotTimeout)
	stopDtActiveSendSince.Store(willNotTimeout)

	sendStartDt := func() {
		startDtActiveSendSince.Store(time.Now())
		sf.sendUFrame(uStartDtActive)
	}

	sendStopDt := func() {
		stopDtActiveSendSince.Store(time.Now())
		sf.sendUFrame(uStopDtActive)
		// the data Transfer is inactived when client try to send a stopDt
		// but it is actived only when client received a startDt confirm
		atomic.StoreUint32(&sf.isActive, inactive)
	}

	sf.setConnectStatus(connected)
	defer sf.setConnectStatus(disconnected)
	sf.onConnect()
	defer sf.onConnectionLost()
	sendStartDt()
	defer sendStopDt()

	for {
		select {
		case <-runCtx.Done():
			return
		case now := <-checkTicker.C:
			// check all timeouts
			if now.Sub(testFrAliveSendSince) >= sf.option.config.SendUnAckTimeout1 ||
				now.Sub(startDtActiveSendSince.Load().(time.Time)) >= sf.option.config.SendUnAckTimeout1 ||
				now.Sub(stopDtActiveSendSince.Load().(time.Time)) >= sf.option.config.SendUnAckTimeout1 {
				sf.Error("test frame alive confirm timeout t₁")
				return
			}
			// check oldest unacknowledged outbound
			if sf.ackNoSend != sf.seqNoSend &&
				//now.Sub(sf.peek()) >= sf.SendUnAckTimeout1 {
				now.Sub(sf.pending[0].sendTime) >= sf.option.config.SendUnAckTimeout1 {
				sf.ackNoSend++
				sf.Error("fatal transmission timeout t₁")
				return
			}

			// // 确定最早发送的i-Frame是否超时,超时则回复sFrame
			// if sf.ackNoRcv != sf.seqNoRcv &&
			// 	(now.Sub(unAckRcvSince) >= sf.option.config.RecvUnAckTimeout2 ||
			// 		now.Sub(idleTimeout3Sine) >= timeoutResolution) {
			// 	sf.sendSFrame()
			// 	sf.ackNoRcv = sf.seqNoRcv
			// }

			// 确定最早发送的i-Frame是否超时,超时则回复sFrame
			if sf.ackNoRcv != sf.seqNoRcv &&
				now.Sub(unAckRcvSince) >= sf.option.config.RecvUnAckTimeout2 {
				sf.sendSFrame()
				sf.ackNoRcv = sf.seqNoRcv
			}

			// 空闲时间到，发送TestFrActive帧,保活
			if now.Sub(idleTimeout3Sine) >= sf.option.config.IdleTimeout3 {
				sf.sendUFrame(uTestFrActive)
				testFrAliveSendSince = time.Now()
				idleTimeout3Sine = testFrAliveSendSince
			}
		case apdu := <-sf.sendRaw:
			idleTimeout3Sine = time.Now()
			sf.send(apdu)
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
				if atomic.LoadUint32(&sf.isActive) == inactive {
					sf.Warn("station not active")
					break // not active, discard apdu
				}
				if !sf.updateAckNoOut(head.rcvSN) || head.sendSN != sf.seqNoRcv {
					sf.Error("fatal incoming acknowledge either earlier than previous or later than sendTime")
					return
				}

				if sf.ackNoRcv == sf.seqNoRcv { // first unacked
					unAckRcvSince = time.Now()
				}
				sf.seqNoRcv = (sf.seqNoRcv + 1) & 32767
				if seqNoCount(sf.ackNoRcv, sf.seqNoRcv) >= sf.option.config.RecvUnAckLimitW {
					sf.sendSFrame()
					sf.ackNoRcv = sf.seqNoRcv
				}

				asduPack := asdu.NewEmptyASDU(&sf.option.param)
				if err := asduPack.UnmarshalBinary(asduVal); err != nil {
					sf.Error("Error unmarshaling asdu: %v", err)
				} else {
					// this uID is used to matching response packet to request packet
					ca := asduPack.CommonAddr
					ioa := asduPack.Clone().DecodeInfoObjAddr()
					uID := uint64(ioa) + (uint64(ca) << (8 * sf.option.param.InfoObjAddrSize))

					if resp, ok := sf.responseHandler[uID]; ok {
						resp <- asduPack
					} else {
						sf.rcvAsdu <- asduPack
					}
				}

			case uAPCI:
				sf.Debug("RX uFrame %v", head)
				switch head.function {
				case uStartDtConfirm:
					// the data Transfer is inactived when client try to send a stopDt
					// but it is actived only when client received a startDt confirm
					atomic.StoreUint32(&sf.isActive, active)
					startDtActiveSendSince.Store(willNotTimeout)
				case uStopDtConfirm:
					// atomic.StoreUint32(&sf.isActive, inactive)
					stopDtActiveSendSince.Store(willNotTimeout)
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

func (sf *HighLevelClient) subscribeLoop(ctx context.Context) {
	sf.Debug("SubscribeLoop started")
	defer sf.Debug("SubscribeLoop stopped")
	defer sf.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-sf.rcvAsdu:
			if sf.subscriptionChan == nil {
				sf.Debug("subscriptionChan nil")
				continue
			}
			ioa := data.Clone().DecodeInfoObjAddr()
			switch data.Type {
			case asdu.M_SP_NA_1:
				// Single Info
				value := data.GetSinglePoint()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_SP_TB_1:
				// Single Info with time
				value := data.GetSinglePoint()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_DP_NA_1:
				// Double Info
				value := data.GetDoublePoint()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_DP_TB_1:
				// Double Info with time
				value := data.GetDoublePoint()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ST_NA_1:
				// Step Position Info
				value := data.GetStepPosition()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ST_TB_1:
				// Step Position Info with time
				value := data.GetStepPosition()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_BO_NA_1:
				// 32 Bit string
				value := data.GetBitString32()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_BO_TB_1:
				// 32 Bit string with time
				value := data.GetBitString32()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ME_NA_1:
				// Normalized Measured Value
				value := data.GetMeasuredValueNormal()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ME_TD_1:
				// Normalized Measured Value with time
				value := data.GetMeasuredValueNormal()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ME_ND_1:
				// Normalized Measured Value without quality description
				value := data.GetMeasuredValueNormal()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ME_NB_1:
				// Scaled Measured Value
				value := data.GetMeasuredValueScaled()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ME_TE_1:
				// Scaled Measured Value with time
				value := data.GetMeasuredValueScaled()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ME_NC_1:
				// Short Float Measured Value
				value := data.GetMeasuredValueFloat()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			case asdu.M_ME_TF_1:
				// Short Float Measured Value with time
				value := data.GetMeasuredValueFloat()[0]
				v := AsduInfo{
					Identifier: data.Identifier,
					Ioa:        ioa,
					Value:      value.Value(),
					Quality:    value.Qds,
					Timestamp:  value.Time,
				}
				sf.subscriptionChan <- v
			}
		}
	}

}

func (sf *HighLevelClient) recvLoop(ctx context.Context, cancel context.CancelFunc) {
	sf.Debug("recvLoop started")
	defer sf.Debug("recvLoop stopped")
	defer func() {
		cancel()
		sf.wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
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
					rdCnt, length = 0, 2
					continue
				}
				length = int(rawData[1]) + 2
				if length < APCICtlFiledSize+2 || length > APDUSizeMax {
					rdCnt, length = 0, 2
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

func (sf *HighLevelClient) setConnectStatus(status uint32) {
	sf.rwMux.Lock()
	atomic.StoreUint32(&sf.status, status)
	sf.rwMux.Unlock()
}

func (sf *HighLevelClient) connectStatus() uint32 {
	sf.rwMux.RLock()
	status := atomic.LoadUint32(&sf.status)
	sf.rwMux.RUnlock()
	return status
}

func (sf *HighLevelClient) cleanUp() {
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
		// case <-sf.rcvASDU:
		// case <-sf.sendASDU:
		default:
			break loop
		}
	}
}

func (sf *HighLevelClient) sendSFrame() {
	sf.Debug("TX sFrame %v", sAPCI{sf.seqNoRcv})
	sf.sendRaw <- newSFrame(sf.seqNoRcv)
}

func (sf *HighLevelClient) sendUFrame(which byte) {
	sf.Debug("TX uFrame %v", uAPCI{which})
	sf.sendRaw <- newUFrame(which)
}

func (sf *HighLevelClient) updateAckNoOut(ackNo uint16) (ok bool) {
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

// Send send asdu
func (sf *HighLevelClient) Send(a *asdu.ASDU) error {
	if !sf.IsConnected() {
		return ErrUseClosedConnection
	}
	if atomic.LoadUint32(&sf.isActive) == inactive {
		return ErrNotActive
	}

	data, err := a.MarshalBinary()
	if err != nil {
		return err
	}

	seqNo := sf.seqNoSend

	iframe, err := newIFrame(seqNo, sf.seqNoRcv, data)
	if err != nil {
		return err
	}
	sf.ackNoRcv = sf.seqNoRcv
	sf.seqNoSend = (seqNo + 1) & 32767
	sf.pending = append(sf.pending, seqPending{seqNo & 32767, time.Now()})

	sf.Debug("TX iFrame %v", iAPCI{seqNo, sf.seqNoRcv})
	sf.sendRaw <- iframe

	return nil
}

func (sf *HighLevelClient) syncSendIFrame(asduPack *asdu.ASDU, resp chan *asdu.ASDU) (*asdu.ASDU, error) {
	data, err := asduPack.MarshalBinary()
	if err != nil {
		return nil, err
	}
	seqNo := sf.seqNoSend
	iframe, err := newIFrame(seqNo, sf.seqNoRcv, data)
	if err != nil {
		return nil, err
	}
	sf.ackNoRcv = sf.seqNoRcv
	sf.seqNoSend = (seqNo + 1) & 32767
	sf.pending = append(sf.pending, seqPending{seqNo & 32767, time.Now()})

	if seqNoCount(sf.ackNoSend, sf.seqNoSend) > sf.option.config.SendUnAckLimitK {
		// TODO: Wait some time? RecvUnAckTimeout2?
		return nil, ErrBufferFulled
	}

	sf.Debug("TX iFrame %v", iAPCI{seqNo, sf.seqNoRcv})

	sf.sendRaw <- iframe

	// timer := time.NewTimer(sf.option.config.SendUnAckTimeout1)
	timer := time.NewTimer(syncSendTimeout)
	defer timer.Stop()

	select {
	case ch := <-resp:
		return ch, nil
	case <-timer.C:
		return nil, fmt.Errorf("ErrorBadTimeOut")
	}
}

func (sf *HighLevelClient) send(apdu []byte) {
	sf.Debug("TX Raw[% x]", apdu)
	for wrCnt := 0; len(apdu) > wrCnt; {
		byteCount, err := sf.conn.Write(apdu[wrCnt:])
		if err != nil {
			// See: https://github.com/golang/go/issues/4373
			if err != io.EOF && err != io.ErrClosedPipe ||
				strings.Contains(err.Error(), "use of closed network connection") {
				// sf.Error("sendRaw failed, %v", err)
				return
			}
			if e, ok := err.(net.Error); !ok || !e.Temporary() {
				// sf.Error("sendRaw failed, %v", err)
				return
			}
			// temporary error may be recoverable
		}
		wrCnt += byteCount
	}
}
