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
	err = client.Connect("192.168.9.33:2404")
	if err != nil {
		fmt.Printf("Failed to connect. error:%v\n", err)
	}
}

// Handle01_02_1e ...
// 01:[M_SP_NA_1] 不带时标单点信息
// 02:[M_SP_TA_1] 带时标CP24Time2a的单点信息,只有(SQ = 0)单个信息元素集合
// 1e:[M_SP_TB_1] 带时标CP56Time2a的单点信息,只有(SQ = 0)单个信息元素集合
func (m *myClient) Handle01_02_1e(c *cs104.Client, spi []asdu.SinglePointInfo) error {
	fmt.Printf("Receivced SinglePointInfo %v\n", spi)
	return nil
}

// Handle03_04_1f ...
// 03:[M_DP_NA_1].双点信息
// 04:[M_DP_TA_1] .带CP24Time2a双点信息,只有(SQ = 0)单个信息元素集合
// 1f:[M_DP_TB_1].带CP56Time2a的双点信息,只有(SQ = 0)单个信息元素集合
func (m *myClient) Handle03_04_1f([]asdu.DoublePointInfo) error {
	return nil
}

// Handle07_08_21 ...
// 07:[M_BO_NA_1] 比特位串
// 08:[M_BO_TA_1] 带时标CP24Time2a比特位串，只有(SQ = 0)单个信息元素集合
// 21:[M_BO_TB_1] 带时标CP56Time2a比特位串，只有(SQ = 0)单个信息元素集
func (m *myClient) Handle07_08_21([]asdu.BitString32Info) error {
	return nil
}

// Handle09_0a_15_22 ...
// 09:[M_ME_NA_1] 测量值,规一化值
// 0a:[M_ME_TA_1] 带时标CP24Time2a的测量值,规一化值,只有(SQ = 0)单个信息元素集合
// 15:[M_ME_ND_1] 不带品质的测量值,规一化值
// 22:[M_ME_TD_1] 带时标CP57Time2a的测量值,规一化值,只有(SQ = 0)单个信息元素集合
func (m *myClient) Handle09_0a_15_22(c *cs104.Client, mvni []asdu.MeasuredValueNormalInfo) error {
	fmt.Printf("Receivced MeasuredValueNormal %v\n", mvni)
	return nil
}

// Handle0b_0c_23 ...
// 0b:[M_ME_NB_1].测量值,标度化值
// 0c:[M_ME_TB_1].带时标CP24Time2a的测量值,标度化值,只有(SQ = 0)单个信息元素集合
// 23:[M_ME_TE_1].带时标CP56Time2a的测量值,标度化值,只有(SQ = 0)单个信息元素集合
func (m *myClient) Handle0b_0c_23([]asdu.MeasuredValueScaledInfo) error{
	return nil
}

// Handle0d_0e_10 ...
// 0d:[M_ME_TF_1] 测量值,短浮点数
// 0e:[M_ME_TC_1].带时标CP24Time2a的测量值,短浮点数,只有(SQ = 0)单个信息元素集合
// 10:[M_ME_TF_1].带时标CP56Time2a的测量值,短浮点数,只有(SQ = 0)单个信息元素集合
func (m *myClient) Handle0d_0e_10([]asdu.MeasuredValueFloatInfo) error{
	return nil
}

// Handle46 46:[M_EI_NA_1], 站初始化结束
func (m *myClient) Handle46(c asdu.Connect, coi asdu.CauseOfInitial) error {
	fmt.Printf("Receivced IFrame typeID 46, cause of init: %v, islocalchange: %v", coi.Cause, coi.IsLocalChange)
	if err := cs104.Activate64(c); err != nil {
		fmt.Printf(err.Error())
	}
	return nil

}

// Handle64 64:[C_IC_NA_1], 总召唤
func (m *myClient) Handle64(c asdu.Connect, a *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
	fmt.Printf("Receivced IFrame typeID 64, qualifier of interrogation: %v\n", qoi)
	if a.Identifier.Coa.Cause == asdu.ActivationCon {
		fmt.Printf("Activation of interrogation is confirmed\n")
	} else if a.Identifier.Coa.Cause == asdu.Deactivation {
		fmt.Printf("Interrogation is deactivated\n")
	}
	return nil
}

// Handle65 65:[C_CI_NA_1], 计数量召唤
// func (m *myclient) Handle65(asdu.Connect, *asdu.ASDU, asdu.QualifierOfInterrogation) error {
// 	return nil
// }

// Handle67 67:[C_CS_NA_1], 时钟同步
func (m *myClient) Handle67(c asdu.Connect, a *asdu.ASDU, t time.Time) error {
	fmt.Printf("Receivced IFrame typeID 67, clock sync: %v\n", t)
	return nil
}