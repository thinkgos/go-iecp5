package cs104

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

// TODO: should check all
// TODO: check when remote connect closed
// TODO: make simple

type sendFr struct {
	t  time.Time
	sn uint16
}
type sendBuf struct {
	buf   []sendFr // 已发送的I帧暂存区
	head  uint16   // 以发送的未确认I帧序号头
	tail  uint16   // 以发送的未确认I帧序号尾
	mutex sync.Mutex
}

type states uint8

const (
	Waiting states = iota
	Running
	Closing
)

// Client is an IEC104 master
type Client struct {
	conf    *Config
	param   *asdu.Params
	handler ClientHandler // 接口

	conn       net.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc

	// channel
	in      chan []byte // for received asdu
	out     chan []byte // for send asdu
	rawRcv  chan []byte // for recvLoop raw cs104 frame
	rawSend chan []byte // for sendLoop raw cs104 frame

	// I帧的发送与接收序号
	sendSN uint16 // 发送序号
	recvSN uint16 // 接收序号
	ackSN  uint16 // 已确认的最大的发送I帧序号

	// I帧发送控制
	*sendBuf // 已发送的未确认I帧暂存区

	// I帧接收控制
	t2Flag  bool      // 超时时间t2被设置标志
	t2Time  time.Time // 接收到连续I帧第一帧的时间
	recvCnt uint16    // 接收到的连续I帧数量

	// IdleTimeout3控制
	t3Time time.Time

	// u帧接收控制
	uFlag bool
	uTime time.Time

	// 服务器是否激活
	isServerActive bool

	// 连接是否被关闭(只能通过Disconnect()修改)
	state states

	// 测试命令计数
	testCnt uint16

	// 其他
	*clog.Clog
	wg sync.WaitGroup
}

// NewClient returns an IEC104 master
func NewClient(conf *Config, params *asdu.Params, handler ClientHandler) (c *Client, err error) {
	if err := conf.Valid(); err != nil {
		return nil, err
	}
	if err := params.Valid(); err != nil {
		return nil, err
	}
	c = &Client{
		conf:    conf,
		param:   params,
		handler: handler,
		state:   Waiting,
		Clog:    clog.NewWithPrefix("IEC104 client =>"),
	}
	err = nil

	return
}

// Connect is
func (sf *Client) Connect(addr string) error {
	sf.state = Running
	conn, err := net.DialTimeout("tcp", addr, sf.conf.ConnectTimeout0)
	if err != nil {
		sf.state = Waiting
		return fmt.Errorf("Failed to dial %s, error: %v", addr, err)
	}
	sf.conn = conn
	sf.Debug("Connected to %s!", addr)

	// initialization
	sf.sendSN = 0
	sf.recvSN = 0
	sf.ackSN = 0
	sf.t2Flag = false
	sf.recvCnt = 0
	sf.testCnt = 0
	sf.isServerActive = false
	sf.rawRcv = make(chan []byte, APDUSizeMax)
	sf.rawSend = make(chan []byte, APDUSizeMax)
	sf.sendBuf = &sendBuf{
		buf:  make([]sendFr, sf.conf.SendUnAckLimitK),
		head: 0,
		tail: 0,
	}
	sf.t3Time = time.Now()

	sf.ctx, sf.cancelFunc = context.WithCancel(context.Background())
	sf.wg.Add(3)
	go sf.recvLoop()
	go sf.sendLoop()
	go sf.handleLoop()
	sf.SendStopDt() // 发送stopDt激活指令
	time.Sleep(sf.conf.SendUnAckTimeout1 / 2)
	sf.SendStartDt() // 发送startDt激活指令

	defer func() {
		sf.cancelFunc()
		sf.wg.Wait()
		sf.conn.Close()
		sf.Debug("Connection to %s Ended!", addr)
		if sf.state != Closing { // 非人为关闭情况下,主动重连
			sf.Connect(addr)
		}
		sf.state = Waiting
	}()

	for {
		// TODO: need sleep?
		time.Sleep(time.Second)
		if time.Since(sf.t3Time) >= sf.conf.IdleTimeout3 {
			sf.t3Time = time.Now()
			sf.SendTestDt()
		}

		if sf.uFlag {
			if time.Since(sf.uTime) >= sf.conf.SendUnAckTimeout1 {
				sf.Error("SendUnAckTimeout1 of uFrame expires, timeout:%v", time.Since(sf.uTime))
				return nil
			}
		}
		if sf.t2Flag {
			if time.Since(sf.t2Time) >= sf.conf.RecvUnAckTimeout2 || sf.recvCnt >= sf.conf.RecvUnAckLimitW {
				sf.recvCnt = 0
				sf.t2Flag = false
				sf.sendSFrame()
			}
		}

		if sf.sendBuf.head != sf.sendBuf.tail {
			if time.Since(sf.sendBuf.buf[sf.sendBuf.head].t) >= sf.conf.SendUnAckTimeout1 {
				sf.Error("SendUnAckTimeout1 of iFrame expires, timeout:%v", time.Since(sf.sendBuf.buf[sf.sendBuf.head].t))
				return nil
			}
			sf.ackSN = sf.sendBuf.buf[sf.sendBuf.head].sn - 1
		} else {
			sf.ackSN = sf.sendSN
		}

		select {
		case <-sf.ctx.Done():
			return fmt.Errorf("ctx done")
		default:
		}
	}
}

func (sf *Client) handleLoop() {
	sf.Debug("HandleLoop Started")
	defer func() {
		sf.wg.Done()
		sf.Debug("HandleLoop Ended")
	}()

	for {
		select {
		case <-sf.ctx.Done():
			return
		case apdu := <-sf.rawRcv:
			apci, rawAsdu := parse(apdu)
			sf.t3Time = time.Now()
			switch apci := apci.(type) {
			case uAPCI:
				sf.Debug(apci.String())
				sf.uFlag = false
				switch apci.function {
				// case uStartDtActive:
				case uStartDtConfirm:
					sf.isServerActive = true
					// 激活之后进行时钟同步?
					_ = sf.ClockSynchronizationCmd(asdu.CauseOfTransmission{Cause: asdu.Activation}, 0xffff, time.Now())
				case uTestFrActive:
					sf.SendTestCon()
				// case uTestFrConfirm:
				// case uStopDtActive:
				case uStopDtConfirm:
				}
			case sAPCI:
				sf.Debug(apci.String())
				if err := sf.checkRecvSN(apci.rcvSN); err != nil {
					sf.Error("Check SFrame recvSN error:%v", err)
					sf.cancelFunc()
					return
				}
			case iAPCI:
				sf.Debug(apci.String())

				// 接收到I帧后开始RecvUnAckTimeout2计时
				sf.recvCnt++
				if !sf.t2Flag {
					sf.t2Flag = true
					sf.t2Time = time.Now()
				} else {
					sf.t2Time = time.Now()
				}

				// 判断接收到的I帧发送序号是否等于客户端的I帧接收序号,第一帧时同为0
				if apci.sendSN != sf.recvSN {
					sf.Error("IFrame sendSN error, close connection!")
					sf.cancelFunc()
					return
				}

				// 判断接收到的I帧接收序号与客户端的I帧发送情况是否匹配
				if err := sf.checkRecvSN(apci.rcvSN); err != nil {
					sf.Error("Check IFrame recvSN error:%v", err)
					sf.cancelFunc()
					return
				}

				sf.recvSN = (sf.recvSN + 1) % 32768

				asduPack := asdu.NewEmptyASDU(sf.param)
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
}

// 判断接收到的I帧接收序号与客户端的I帧发送情况是否匹配
func (sf *Client) checkRecvSN(recvSN uint16) error {
	sf.sendBuf.mutex.Lock()
	defer sf.sendBuf.mutex.Unlock()
	if sf.sendBuf.head == sf.sendBuf.tail { // sendBuf为空,没有未确认的已发送I帧
		if recvSN == sf.sendSN {
			return nil
		}
	} else { // sendBuf不为空,有未被确认的以发送帧
		head, tail := sf.sendBuf.buf[sf.sendBuf.head].sn, sf.sendBuf.buf[sf.sendBuf.tail].sn
		if recvSN == tail { // S帧确认了所有已发送的I帧
			sf.sendBuf.head, sf.sendBuf.tail = 0, 0
			return nil
		}
		if head < tail { // 客户端I帧发送序号未溢出
			if recvSN >= head && recvSN <= tail {
				for recvSN >= head {
					sf.sendBuf.head++
					if sf.sendBuf.head == sf.sendBuf.tail {
						sf.sendBuf.head, sf.sendBuf.tail = 0, 0
						return nil
					}
					head = sf.sendBuf.buf[sf.sendBuf.head].sn
				}
				return nil
			}
		} else { //客户端I帧发送序号溢出
			if recvSN >= head && recvSN <= 32767 { // 发送和接收序号最大为15位,2^15-1
				for recvSN >= head {
					sf.sendBuf.head++
					head = sf.sendBuf.buf[sf.sendBuf.head].sn
					if head == 0 {
						return nil
					}
				}
				return nil
			} else if recvSN <= tail {
				for head != 0 {
					sf.sendBuf.head++
					head = sf.sendBuf.buf[sf.sendBuf.head].sn
				}
				for recvSN >= head {
					sf.sendBuf.head++
					if sf.sendBuf.head == sf.sendBuf.tail {
						sf.sendBuf.head, sf.sendBuf.tail = 0, 0
						return nil
					}
				}
				return nil
			}
		}
	}

	return fmt.Errorf("wrong sequence number, close connection")
}

func (sf *Client) recvLoop() {
	sf.Debug("recvLoop started")
	defer func() {
		sf.cancelFunc()
		sf.wg.Done()
		sf.Debug("recvLoop stopped")
	}()

	for {
		rawData := make([]byte, APDUSizeMax)
		for rdCnt, length := 0, 1; rdCnt < length; {
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
					sf.rawRcv <- apdu
				}
			}
		}
	}
}

func (sf *Client) sendLoop() {
	sf.Debug("sendLoop started")
	defer func() {
		sf.cancelFunc()
		sf.wg.Done()
		sf.Debug("sendLoop stopped")
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

func (sf *Client) sendSFrame() {
	sf.Debug("TX sFrame %v", sAPCI{sf.recvSN})
	sf.rawSend <- newSFrame(sf.recvSN)
}
func (sf *Client) sendUFrame(which byte) {
	sf.Debug("TX uFrame %v", uAPCI{which})
	sf.rawSend <- newUFrame(which)
}

func (sf *Client) sendIFrame(asdu []byte) error {
	iFrame, err := newIFrame(sf.sendSN, sf.recvSN, asdu)
	if err != nil {
		return err
	}
	sf.rawSend <- iFrame
	sf.Debug("TX iFrame %v", iAPCI{sf.sendSN, sf.recvSN})
	sf.sendSN = (sf.sendSN + 1) % 32768

	sf.sendBuf.mutex.Lock()
	defer sf.sendBuf.mutex.Unlock()
	sf.sendBuf.buf[sf.sendBuf.tail].t = time.Now()
	sf.sendBuf.buf[sf.sendBuf.tail].sn = sf.sendSN
	sf.sendBuf.tail++
	return nil
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

func (sf *Client) SendStartDt() {
	if !sf.uFlag {
		sf.uFlag = true
		sf.uTime = time.Now()
	}
	sf.sendUFrame(uStartDtActive)
}
func (sf *Client) SendStopDt() {
	if !sf.uFlag {
		sf.uFlag = true
		sf.uTime = time.Now()
	}
	sf.sendUFrame(uStopDtActive)
}
func (sf *Client) SendTestDt() {
	if !sf.uFlag {
		sf.uFlag = true
		sf.uTime = time.Now()
	}
	sf.sendUFrame(uTestFrActive)
}
func (sf *Client) SendTestCon() {
	sf.sendUFrame(uTestFrConfirm)
}

// Params returns params of client
func (sf *Client) Params() *asdu.Params {
	return sf.param
}

// Send send asdu
func (sf *Client) Send(a *asdu.ASDU) error {
	if !sf.isServerActive {
		return fmt.Errorf("ErrorUnactive")
	}
	data, err := a.MarshalBinary()
	if err != nil {
		return err
	}
	//select {
	//case sf.out <- data:
	//default:
	//	return ErrBufferFulled
	//}
	//return nil

	return sf.sendIFrame(data)
}

// UnderlyingConn returns underlying conn of client
func (sf *Client) UnderlyingConn() net.Conn {
	return sf.conn
}

func (sf *Client) GetState() states {
	return sf.state
}

func (sf *Client) Disconnect() {
	sf.state = Closing
	sf.cancelFunc()
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
	u.AppendBytes(byte(sf.testCnt&0xff), byte(sf.testCnt>>8))
	u.AppendBytes(asdu.CP56Time2a(t, u.InfoObjTimeZone)...)
	return sf.Send(u)
}

// // GetTestCommand ...
// func GetTestCommand(a *asdu.ASDU) (ioa asdu.InfoObjAddr, tsc uint16, t time.Time) {
// 	ioa = a.DecodeInfoObjAddr()
// 	tsc = a.infoObj[0] + (a.infoObj[1] << 8)
// 	t = asdu.ParseCP56Time2a(a.infoObj, a.InfoObjTimeZone)
// 	return
// }
