package main

import (
	"fmt"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

type myClient struct{}

func main() {
	mycli := &myClient{}

	client, err := cs104.NewClient(&cs104.Config{}, asdu.ParamsWide, mycli)
	if err != nil {
		fmt.Printf("Failed to creat cs104 client. error:%v\n", err)
	}
	client.LogMode(true)
	if err = client.AddRemoteServer("127.0.0.1:2404"); err != nil {
		panic(err)
	}
	err = client.Start()
	if err != nil {
		fmt.Printf("Failed to connect. error:%v\n", err)
	}

	for {
		time.Sleep(time.Second * 100)
	}

}
func (myClient) InterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfInterrogation) error {
	return nil
}

func (myClient) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error {
	return nil
}
func (myClient) ReadHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr) error {
	return nil
}
func (myClient) ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error {
	return nil
}
func (myClient) ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error {
	return nil
}
func (myClient) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error {
	return nil
}
func (myClient) ASDUHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
