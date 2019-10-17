package main

import(
	"fmt"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

type myClient struct{}

func main() {
	mycli := &myClient{}

	client, err := cs104.NewClient(cs104.DefaultConfig(), cs104.DefaultParam(), mycli)
	if err != nil {
		fmt.Printf("Failed to creat cs104 client. error:%v\n", err)
	}
	client.LogMode(true)
	err = client.Connect("127.0.0.1:2404")
	if err != nil {
		fmt.Printf("Failed to connect. error:%v\n", err)
	}
}

// Handle01_02_1e ...
// 01:[M_SP_NA_1] 不带时标单点信息
// 02:[M_SP_TA_1] 带时标CP24Time2a的单点信息,只有(SQ = 0)单个信息元素集合
// 1e:[M_SP_TB_1] 带时标CP56Time2a的单点信息,只有(SQ = 0)单个信息元素集合
func (c *myClient) Handle01_02_1e(conn asdu.Connect, a *asdu.ASDU, infos []asdu.SinglePointInfo) {
	for i := range infos {
		fmt.Println(bool(infos[i].Value), uint8(a.Identifier.Type), uint(infos[i].Ioa))
	}
}

// Handle03_04_1f ...
// 03:[M_DP_NA_1].双点信息
// 04:[M_DP_TA_1] .带CP24Time2a双点信息,只有(SQ = 0)单个信息元素集合
// 1f:[M_DP_TB_1].带CP56Time2a的双点信息,只有(SQ = 0)单个信息元素集合
func (c *myClient) Handle03_04_1f(conn asdu.Connect, a *asdu.ASDU, infos []asdu.DoublePointInfo) {
	for i := range infos {
		fmt.Println(byte(infos[i].Value), uint8(a.Identifier.Type), uint(infos[i].Ioa))
	}
}

// Handle05_20 ...
// 05:[M_ST_NA_1].步位置信息
// 20:[M_ST_TB_1].带时标CP56Time2a的双点信息,只有(SQ = 0)单个信息元素集合
func (c *myClient) Handle05_20(conn asdu.Connect, a *asdu.ASDU, infos []asdu.StepPositionInfo) {
	for i := range infos {
		fmt.Println(byte(infos[i].Value.Value()), a.Identifier.Type, infos[i].Ioa)
	}
}

// Handle07_08_21 ...
// 07:[M_BO_NA_1] 比特位串
// 08:[M_BO_TA_1] 带时标CP24Time2a比特位串，只有(SQ = 0)单个信息元素集合
// 21:[M_BO_TB_1] 带时标CP56Time2a比特位串，只有(SQ = 0)单个信息元素集
func (c *myClient) Handle07_08_21(conn asdu.Connect, a *asdu.ASDU, infos []asdu.BitString32Info) {
	for i := range infos {
		fmt.Println(infos[i].Value, uint8(a.Identifier.Type), uint(infos[i].Ioa))
	}
}

// Handle09_0a_15_22 ...
// 09:[M_ME_NA_1] 测量值,规一化值
// 0a:[M_ME_TA_1] 带时标CP24Time2a的测量值,规一化值,只有(SQ = 0)单个信息元素集合
// 15:[M_ME_ND_1] 不带品质的测量值,规一化值
// 22:[M_ME_TD_1] 带时标CP57Time2a的测量值,规一化值,只有(SQ = 0)单个信息元素集合
func (c *myClient) Handle09_0a_15_22(conn asdu.Connect, a *asdu.ASDU, infos []asdu.MeasuredValueNormalInfo) {
	for i := range infos {
		fmt.Println(int16(infos[i].Value), uint8(a.Identifier.Type), uint(infos[i].Ioa))
	}
}

// Handle0b_0c_23 ...
// 0b:[M_ME_NB_1].测量值,标度化值
// 0c:[M_ME_TB_1].带时标CP24Time2a的测量值,标度化值,只有(SQ = 0)单个信息元素集合
// 23:[M_ME_TE_1].带时标CP56Time2a的测量值,标度化值,只有(SQ = 0)单个信息元素集合
func (c *myClient) Handle0b_0c_23(conn asdu.Connect, a *asdu.ASDU, infos []asdu.MeasuredValueScaledInfo) {
	for i := range infos {
		fmt.Println(int16(infos[i].Value), uint8(a.Identifier.Type), uint(infos[i].Ioa))
	}
}

// Handle0d_0e_10 ...
// 0d:[M_ME_TF_1] 测量值,短浮点数
// 0e:[M_ME_TC_1].带时标CP24Time2a的测量值,短浮点数,只有(SQ = 0)单个信息元素集合
// 10:[M_ME_TF_1].带时标CP56Time2a的测量值,短浮点数,只有(SQ = 0)单个信息元素集合
func (c *myClient) Handle0d_0e_10(conn asdu.Connect, a *asdu.ASDU, infos []asdu.MeasuredValueFloatInfo) {
	for i := range infos {
		fmt.Println(float32(infos[i].Value), uint8(a.Identifier.Type), uint(infos[i].Ioa))
	}
}

// Handle46 46:[M_EI_NA_1], 站初始化结束
func (c *myClient) Handle46(conn asdu.Connect, coi asdu.CauseOfInitial) {
	fmt.Printf("Receivced IFrame typeID 46, cause of init: %v, islocalchange: %v", coi.Cause, coi.IsLocalChange)
}

// Handle64 64:[C_IC_NA_1], 总召唤
func (c *myClient) Handle64(conn asdu.Connect, a *asdu.ASDU, qoi asdu.QualifierOfInterrogation) {
	fmt.Println(fmt.Sprintf("Receivced IFrame typeID 64, qualifier of interrogation: %v", qoi))
}

// Handle65 65:[C_CI_NA_1], 计数量召唤
// func (m *myclient) Handle65(asdu.Connect, *asdu.ASDU, asdu.QualifierOfInterrogation) {
// 	
// }

// Handle67 67:[C_CS_NA_1], 时钟同步
func (c *myClient) Handle67(conn asdu.Connect, a *asdu.ASDU, t time.Time) {
	fmt.Println(fmt.Sprintf("Receivced IFrame typeID 67, clock sync: %v", t))
	
}