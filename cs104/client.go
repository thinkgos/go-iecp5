package cs104

import(
	"fmt"
	"net"
	"time"
	"sync"
	"io"
	"strings"
	"context"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/clog"
)

type sendFr struct {
	t 	time.Time
	sn 	uint16
}
type sendBuf struct {
	buf 	[]sendFr						// 已发送的I帧暂存区
	head 	uint16							// 以发送的未确认I帧序号头
	tail 	uint16							// 以发送的未确认I帧序号尾
	mutex   sync.Mutex
}

// Client is an IEC104 master
type Client struct {
	conf 			*Config
	param 			*asdu.Params

	handler 		ClientHandler			// 接口

	conn 			net.Conn

	// channel 
	recvChan		chan []byte				// 接收到的数据包
	sendChan		chan []byte				// 要发送的数据包
	
	// I帧的发送与接收序号
	sendSN 			uint16					// 发送序号
	recvSN 			uint16					// 接收序号
	ackSN			uint16					// 已确认的最大的发送I帧序号

	// I帧发送控制
	*sendBuf								// 已发送的未确认I帧暂存区

	// I帧接收控制
	t2Flag			bool					// 超时时间t2被设置标志
	t2Time			time.Time				// 接收到连续I帧第一帧的时间
	recvCnt			uint16					// 接收到的连续I帧数量

	// IdleTimeout3控制
	t3Time			time.Time

	// u帧接收控制
	uFlag			bool
	uTime			time.Time

	// 服务器是否激活
	isServerActive 	bool

	// 连接是否被关闭(只能通过Disconnect()修改)
	isClosed		bool

	// 其他
	*clog.Clog
	wg              sync.WaitGroup
}

// DefaultConfig is
func DefaultConfig() *Config {
	return &Config{
		SendUnAckTimeout1: 10 * time.Second,
		IdleTimeout3: 20 *time.Second,
	}
}

// DefaultParam is
func DefaultParam() *asdu.Params {
	return asdu.ParamsWide
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
		conf: 			conf,
		param: 			params,
		handler: 		handler,
		Clog: 			clog.NewWithPrefix("IEC104 client =>"),
	}
	err = nil
	c.LogMode(false)

	return 
}

// Connect is 
func (c *Client) Connect(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, c.conf.ConnectTimeout0)
	if err != nil {
		return fmt.Errorf("Failed to dial %s, error: %v", addr, err)
	}
	c.conn = conn
	defer c.conn.Close()

	// initialization
	c.sendSN = 0
	c.recvSN = 0
	c.ackSN = 0
	c.t2Flag = false
	c.recvCnt = 0
	c.isServerActive = false
	c.isClosed = false
	c.recvChan = make(chan []byte, APDUSizeMax)
	c.sendChan = make(chan []byte, APDUSizeMax)
	c.sendBuf = &sendBuf{
		buf:	make([]sendFr, c.conf.SendUnAckLimitK),
		head:	0,
		tail:	0,
	}
	c.t3Time = time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	c.wg.Add(3)
	go c.recvLoop(ctx)
	go c.sendLoop(ctx)
	go c.handleLoop(ctx, cancel)
	c.SendStopDt()							// 发送stopDt激活指令
	time.Sleep(c.conf.SendUnAckTimeout1/2)
	c.SendStartDt()							// 发送startDt激活指令

	defer func() {
		cancel()
		c.wg.Wait()
		c.Debug("Connection to %s Ended!", addr)
		if !c.isClosed {					// 非人为关闭情况下,主动重连
			c.Connect(addr)
		}
	}()


	for {
		if c.isClosed {
			return fmt.Errorf("Connection is closed")
		}

		// TODO: need sleep?
		time.Sleep(time.Second)
		if time.Since(c.t3Time) >= c.conf.IdleTimeout3 {
			c.t3Time = time.Now()
			c.SendTestDt()
		}

		if c.uFlag {
			if time.Since(c.uTime) >= c.conf.SendUnAckTimeout1 {
				return fmt.Errorf("SendUnAckTimeout1 of uFrame expires")
			}
		}
		if c.t2Flag {
			if time.Since(c.t2Time) >= c.conf.RecvUnAckTimeout2 || c.recvCnt >= c.conf.RecvUnAckLimitW {
				c.recvCnt = 0
				c.t2Flag = false
				c.sendSFrame()
			}
		}

		if c.sendBuf.head != c.sendBuf.tail {
			if time.Since(c.sendBuf.buf[c.sendBuf.head].t) >= c.conf.SendUnAckTimeout1 {
				return fmt.Errorf("SendUnAckTimeout1 of iFrame expires")
			}
			c.ackSN = c.sendBuf.buf[c.sendBuf.head].sn - 1
		} else {
			c.ackSN = c.sendSN
		}

		select {
		case <- ctx.Done():
			return fmt.Errorf("ctx done")
		default:
		}
	}
}

func (c *Client) handleLoop(ctx context.Context, cancel context.CancelFunc) {
	c.Debug("HandleLoop Started")
	defer func() {
		c.wg.Done()
		c.Debug("HandleLoop Ended")
	}()

	for {
		select {
		case <- ctx.Done():
			return
		case apdu := <- c.recvChan:
			apci, rawAsdu := parse(apdu)
			c.t3Time = time.Now()
			switch apci := apci.(type) {
			case uAPCI:
				c.Debug(apci.String())
				c.uFlag = false
				switch apci.function {
				// case uStartDtActive:
				case uStartDtConfirm:
					c.isServerActive = true
					Activate67(c)								// 激活之后进行时钟同步?
				case uTestFrActive:
					c.SendTestCon()
				// case uTestFrConfirm:
				// case uStopDtActive:
				case uStopDtConfirm:
				}
			case sAPCI:
				c.Debug(apci.String())
				if err := c.checkRecvSN(apci.rcvSN); err != nil {
					cancel()
					return
				}
			case iAPCI:
				c.Debug(apci.String())

				// 接收到I帧后开始RecvUnAckTimeout2计时
				c.recvCnt++
				if !c.t2Flag {
					c.t2Flag = true
					c.t2Time = time.Now()
				} else {
					c.t2Time = time.Now()
				}

				// 判断接收到的I帧发送序号是否等于客户端的I帧接收序号,第一帧时同为0
				if apci.sendSN != c.recvSN {
					c.Debug("IFrame sequence error, close connection!")
					cancel()
					return
				}

				// 判断接收到的I帧接收序号与客户端的I帧发送情况是否匹配
				if err := c.checkRecvSN(apci.rcvSN); err != nil {
					c.Debug(err.Error())
					cancel()
					return
				}

				c.recvSN = (c.recvSN + 1) % 32768

				asduPack := asdu.NewEmptyASDU(c.param)
				if err := asduPack.UnmarshalBinary(rawAsdu); err != nil {
					c.Error("asdu UnmarshalBinary failed,%+v", err)
					continue
				}
				if err := c.handleIFrame(asduPack); err != nil {
					c.Error("Falied handling I frame, error: %v", err)
				}
			}
		}
	}
}

// 判断接收到的I帧接收序号与客户端的I帧发送情况是否匹配
func (c *Client) checkRecvSN(recvSN uint16) error {
	c.sendBuf.mutex.Lock()
	defer c.sendBuf.mutex.Unlock()
	if c.sendBuf.head == c.sendBuf.tail {				// sendBuf为空,没有未确认的已发送I帧
		if recvSN == c.sendSN {
			return nil
		}
	} else {											// sendBuf不为空,有未被确认的以发送帧
		head, tail := c.sendBuf.buf[c.sendBuf.head].sn, c.sendBuf.buf[c.sendBuf.tail].sn
		if recvSN == tail {								// S帧确认了所有已发送的I帧
			c.sendBuf.head, c.sendBuf.tail = 0, 0
			return nil
		}
		if head < tail {								// 客户端I帧发送序号未溢出
			if recvSN >= head && recvSN <= tail {
				for recvSN >= head {
					c.sendBuf.head++
					if c.sendBuf.head == c.sendBuf.tail {
						c.sendBuf.head, c.sendBuf.tail = 0, 0
						return nil
					}
					head = c.sendBuf.buf[c.sendBuf.head].sn
				}
				return nil
			}
		} else {										//客户端I帧发送序号溢出
			if recvSN >= head && recvSN <= 32767 {		// 发送和接收序号最大为15位,2^15-1
				for recvSN >= head {
					c.sendBuf.head++
					head = c.sendBuf.buf[c.sendBuf.head].sn
					if head == 0 {
						return nil
					}
				}
				return nil
			} else if recvSN <= tail{
				for head != 0 {
					c.sendBuf.head++
					head = c.sendBuf.buf[c.sendBuf.head].sn
				}
				for recvSN >= head {
					c.sendBuf.head++
					if c.sendBuf.head == c.sendBuf.tail {
						c.sendBuf.head, c.sendBuf.tail = 0, 0
						return nil
					}
				}
				return nil
			}
		}
	}

	return fmt.Errorf("wrong sequence number, close connection")
}

func (c *Client) recvLoop(ctx context.Context) {
	c.Debug("RecvLoop Started")
	defer func() {
		c.wg.Done()
		c.Debug("RecvLoop Ended")
	}()

	apdu := make([]byte, APDUSizeMax)
	for head, tail := 0, 1; head < tail; {
		rdCnt, err := io.ReadFull(c.conn, apdu[head:tail])
		if err != nil {
			// See: https://github.com/golang/go/issues/4373
			if err != io.EOF && err != io.ErrClosedPipe ||
				strings.Contains(err.Error(), "use of closed network connection") {
				c.Error("receive failed, %v", err)
				return
			}
			if e, ok := err.(net.Error); ok && !e.Temporary() {
				c.Error("receive failed, %v", err)
				return
			}
			if rdCnt == 0 && err == io.EOF {
				c.Error("remote connect closed, %v", err)
				return
			}
		}
			
		switch head {
		case 0:
			if apdu[head] == startFrame {
				head += rdCnt
				tail += rdCnt
			}
		case 1:
			tail += int(apdu[head])
			head += rdCnt
			if tail < APCICtlFiledSize + 2 || tail > APDUSizeMax {
				head = 0
				tail = 1
				apdu = make([]byte, APDUSizeMax)
			}
		default:
			head += rdCnt
			if tail == head {
				c.Debug("RX [% x]", apdu[:tail])
				c.recvChan <- apdu[:tail]
				head = 0
				tail = 1
				apdu = make([]byte, APDUSizeMax)
			}
		}

		select {
		case <- ctx.Done():
			return
		default:
		}
	}
}

func (c *Client) sendLoop(ctx context.Context) {
	c.Debug("SendLoop Started")
	defer func() {
		c.Debug("SendLoop Ended")
		c.wg.Done()
	}()
	for {
		select {
		case <- ctx.Done():
			return
		case apdu := <- c.sendChan:
			c.Debug("TX [% X]", apdu)
			for wrCnt := 0; len(apdu) > wrCnt; {
				byteCount, err := c.conn.Write(apdu[wrCnt:])
				if err != nil {
					// See: https://github.com/golang/go/issues/4373
					if err != io.EOF && err != io.ErrClosedPipe ||
						strings.Contains(err.Error(), "use of closed network connection") {
						c.Error("send failed, %v", err)
						return
					}
					if e, ok := err.(net.Error); !ok || !e.Temporary() {
						c.Error("send failed, %v", err)
						return
					}
					// temporary error may be recoverable
				}
				wrCnt += byteCount
			}
		}
	}
}

func (c *Client) sendIFrame(asdu []byte) {
	iFrame, err := newIFrame(c.sendSN, c.recvSN, asdu)
	if err != nil {
		c.Debug(err.Error())
	}
	c.sendChan <- iFrame
	c.Debug("TX iFrame %v", iAPCI{c.sendSN, c.recvSN})
	c.sendSN = (c.sendSN + 1) % 32768

	c.sendBuf.mutex.Lock()
	defer c.sendBuf.mutex.Unlock()
	c.sendBuf.tail++
	c.sendBuf.buf[c.sendBuf.tail].t = time.Now()
	c.sendBuf.buf[c.sendBuf.tail].sn = c.sendSN
}
func (c *Client) sendSFrame() {
	c.Debug("TX sFrame %v", sAPCI{c.recvSN})
	c.sendChan <- newSFrame(c.recvSN)
}
func (c *Client) sendUFrame(b byte) {
	c.Debug("TX uFrame %v", uAPCI{b})
	c.sendChan <- newUFrame(b)
}

// 接收到I帧后根据TYPEID进行不同的处理,分别调用对应的接口函数
func (c *Client) handleIFrame(a *asdu.ASDU) error {
	
	defer func() {
		if err := recover(); err != nil {
			c.Critical("Client handler %+v", err)
		}
	}()

	c.Debug("ASDU %+v", a)

	// check common addr
	if 	a.CommonAddr == asdu.InvalidCommonAddr {
		return a.SendReplyMirror(c, asdu.UnknownCA)
	}

	if 	a.Identifier.Coa.Cause == asdu.UnknownTypeID ||
		a.Identifier.Coa.Cause == asdu.UnknownCOT ||
		a.Identifier.Coa.Cause == asdu.UnknownCA ||
		a.Identifier.Coa.Cause == asdu.UnknownIOA {
		return fmt.Errorf("GOT COT %v", a.Identifier.Coa.Cause)
	}

	switch a.Identifier.Type {
	case asdu.M_SP_NA_1, asdu.M_SP_TA_1, asdu.M_SP_TB_1:		// 遥信 单点信息 01 02 30
		// check cot
		if !( a.Identifier.Coa.Cause == asdu.Background ||
			  a.Identifier.Coa.Cause == asdu.Spontaneous ||
			  a.Identifier.Coa.Cause == asdu.Request ||
			  a.Identifier.Coa.Cause == asdu.ReturnInfoRemote ||
			  a.Identifier.Coa.Cause == asdu.ReturnInfoLocal ||
			  a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
		    ( a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
			  a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16 ) ) {
			return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		info := a.GetSinglePoint()
		c.handler.Handle01_02_1e(c, a, info)
	case asdu.M_DP_NA_1, asdu.M_DP_TA_1, asdu.M_DP_TB_1: 		// 遥信 双点信息 3,4,31
		// check cot
		if !( a.Identifier.Coa.Cause == asdu.Background ||
			  a.Identifier.Coa.Cause == asdu.Spontaneous ||
			  a.Identifier.Coa.Cause == asdu.Request ||
			  a.Identifier.Coa.Cause == asdu.ReturnInfoRemote ||
			  a.Identifier.Coa.Cause == asdu.ReturnInfoLocal ||
			  a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
		    ( a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
			  a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16 ) ) {
			return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		info := a.GetDoublePoint()
		c.handler.Handle03_04_1f(c, a, info)
		case asdu.M_ST_NA_1, asdu.M_ST_TB_1: 					// 遥信 步调节信息 5,32
			// check cot
			if !( a.Identifier.Coa.Cause == asdu.Background ||
				  a.Identifier.Coa.Cause == asdu.Spontaneous ||
				  a.Identifier.Coa.Cause == asdu.Request ||
				  a.Identifier.Coa.Cause == asdu.ReturnInfoRemote ||
				  a.Identifier.Coa.Cause == asdu.ReturnInfoLocal ||
				  a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
				( a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
				  a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16 ) ) {
				return a.SendReplyMirror(c, asdu.UnknownCOT)
			}
			info := a.GetStepPosition()
			c.handler.Handle05_20(c, a, info)
	case asdu.M_BO_NA_1, asdu.M_BO_TA_1, asdu.M_BO_TB_1:		// 遥信 比特串信息 07,08,33								// 比特串,07
		// check cot
		if !( a.Identifier.Coa.Cause == asdu.Background ||
			  a.Identifier.Coa.Cause == asdu.Spontaneous ||
			  a.Identifier.Coa.Cause == asdu.Request ||
			  a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
			( a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
			  a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16 ) ) {
		  	return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		info := a.GetBitString32()
		c.handler.Handle07_08_21(c, a, info)
	case asdu.M_ME_NA_1, asdu.M_ME_TA_1, asdu.M_ME_TD_1, asdu.M_ME_ND_1:	// 遥测 归一化测量值 09,10,21,34
		// check cot
		if !( a.Identifier.Coa.Cause == asdu.Periodic ||
			  a.Identifier.Coa.Cause == asdu.Background ||
			  a.Identifier.Coa.Cause == asdu.Spontaneous ||
			  a.Identifier.Coa.Cause == asdu.Request ||
			  a.Identifier.Coa.Cause == asdu.InterrogatedByStation ) {
			return a.SendReplyMirror(c, asdu.UnknownCOT)
		  }
		  value := a.GetMeasuredValueNormal()
		  c.handler.Handle09_0a_15_22(c, a, value)
	case asdu.M_ME_NB_1, asdu.M_ME_TB_1, asdu.M_ME_TE_1:		//遥测 标度化值 11,12,35
		// check cot
		if !( a.Identifier.Coa.Cause == asdu.Periodic ||
			  a.Identifier.Coa.Cause == asdu.Background ||
			  a.Identifier.Coa.Cause == asdu.Spontaneous ||
			  a.Identifier.Coa.Cause == asdu.Request ||
			  a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
		    ( a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
			  a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16 ) ) {
			return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		value := a.GetMeasuredValueScaled()
		c.handler.Handle0b_0c_23(c, a, value)
	case asdu.M_ME_NC_1, asdu.M_ME_TC_1, asdu.M_ME_TF_1:		// 遥信 短浮点数 13,14,16
		// check cot
		if !( a.Identifier.Coa.Cause == asdu.Periodic ||
			  a.Identifier.Coa.Cause == asdu.Background ||
			  a.Identifier.Coa.Cause == asdu.Spontaneous ||
			  a.Identifier.Coa.Cause == asdu.Request ||
			  a.Identifier.Coa.Cause == asdu.InterrogatedByStation ||
		    ( a.Identifier.Coa.Cause >= asdu.InterrogatedByGroup1 &&
			  a.Identifier.Coa.Cause <= asdu.InterrogatedByGroup16 ) ) {
			return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		value := a.GetMeasuredValueFloat()
		c.handler.Handle0d_0e_10(c, a, value)
	case asdu.M_EI_NA_1:										// 站初始化结束 70
		// check cause of transmission
		if 	!( a.Identifier.Coa.Cause == asdu.Initialized ) {
			return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		ioa, coi := a.GetEndOfInitialization()
		if 	ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(c, asdu.UnknownIOA)
		}
		c.handler.Handle46(c, coi)
	case asdu.C_IC_NA_1: 										// 总召唤 100
		// check cot
		if 	!( a.Identifier.Coa.Cause == asdu.ActivationCon ||
			  a.Identifier.Coa.Cause == asdu.Deactivation ||
			  a.Identifier.Coa.Cause == asdu.ActivationTerm ) {
			return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		// get ioa and qoi
		ioa, qoi := a.GetInterrogationCmd()
		// check ioa
		if 	ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(c, asdu.UnknownIOA)
		}
		c.handler.Handle64(c, a, qoi)
	case asdu.C_CS_NA_1:										// 时钟同步 103
		// check cot
		if !( a.Identifier.Coa.Cause == asdu.ActivationCon ||
			  a.Identifier.Coa.Cause == asdu.ActivationTerm ||
			  a.Identifier.Coa.Cause == asdu.UnknownTypeID ) {
		  	return a.SendReplyMirror(c, asdu.UnknownCOT)
		}
		ioa, t := a.GetClockSynchronizationCmd()
		// check ioa
		if 	ioa != asdu.InfoObjAddrIrrelevant {
			return a.SendReplyMirror(c, asdu.UnknownIOA)
		}
		c.handler.Handle67(c, a, t)
	default:
		return a.SendReplyMirror(c, asdu.UnknownTypeID)
	}

	// if err := c.handler.ASDUHandler(c, a); err != nil {
	// 	return a.SendReplyMirror(c, asdu.UnknownTypeID)
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
	// Handle65(asdu.Connect, *asdu.ASDU, asdu.QualifierOfInterrogation)
	// 67:[C_CS_NA_1] 时钟同步
	Handle67(asdu.Connect, *asdu.ASDU, time.Time)
}

func (c *Client) SendStartDt() {
	if !c.uFlag {
		c.uFlag = true
		c.uTime = time.Now()
	}
	c.sendUFrame(uStartDtActive)
}
func (c *Client) SendStopDt() {
	if !c.uFlag {
		c.uFlag = true
		c.uTime = time.Now()
	}
	c.sendUFrame(uStopDtActive)
}
func (c *Client) SendTestDt() {
	if !c.uFlag {
		c.uFlag = true
		c.uTime = time.Now()
	}
	c.sendUFrame(uTestFrActive)
}
func (c *Client) SendTestCon() {
	c.sendUFrame(uTestFrConfirm)
}

// Send sends
func (c *Client) Send(a *asdu.ASDU) error {
	if !c.isServerActive {
		return fmt.Errorf("ErrorUnactive")
	}
	data, err := a.MarshalBinary()
	if err != nil {
		return err
	}
	c.sendIFrame(data)
	return nil
}
//Params returns params of client
func (c *Client) Params() *asdu.Params {
	return c.param
}
//UnderlyingConn returns underlying conn of client
func (c *Client) UnderlyingConn() net.Conn {
	return c.conn
}

func (c *Client) EnableLogging(b bool) {
	c.LogMode(b)
}

func (c *Client) Disconnect() {
	c.isClosed = true
}

// Activate64 wraps InterrogationCmd for easy use
func Activate64(c asdu.Connect) error {
	if err := asdu.InterrogationCmd(c, asdu.ParseCauseOfTransmission(0x06), 0xFFFF, 0x14); err != nil {
		return err
	}
	return nil
}
// Activate67 wraps ClockSynchronizationCmd for easy use
func Activate67(c asdu.Connect) error {
	if err := asdu.ClockSynchronizationCmd(c, asdu.ParseCauseOfTransmission(0x06), 0xFFFF, time.Now()); err != nil {
		return err
	}
	return nil
}

// Read66 wraps ReadCmd for easy use
func Read66(c asdu.Connect, i asdu.InfoObjAddr) error {
	if err := asdu.ReadCmd(c, asdu.ParseCauseOfTransmission(0x05), 0xFFFF, i); err != nil {
		return err
	}
	return nil
}
