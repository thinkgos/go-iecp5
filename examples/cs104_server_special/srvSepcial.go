package main

import (
	"log"
	"net/http"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
	_ "net/http/pprof"
)

func main() {
	srv, _ := cs104.NewServerSpecial(&cs104.Config{}, asdu.ParamsWide, &mysrv{})

	srv.LogMode(true)
	err := srv.AddRemoteServer("192.168.199.214:2404")
	if err != nil {
		log.Println(err)
	}

	srv.SetOnConnectHandler(func(c cs104.ServerSpecial) error {
		_, err := c.UnderlyingConn().Write([]byte{0x68, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x46, 0x01, 0x04, 0x00, 0xa0, 0xaf, 0xbd, 0xd8, 0x0a, 0xf4})
		log.Println("connected")
		return err

	})
	srv.SetConnectionLostHandler(func(cs104.ServerSpecial) {
		log.Println("disconnected")
	})
	err = srv.Start()
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":6060", nil); err != nil {
		panic(err)
	}
}

type mysrv struct{}

func (this *mysrv) InterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
	log.Println("qoi", qoi)
	// asduPack.SendReplyMirror(c, asdu.ActivationCon)
	// err := asdu.Single(c, false, asdu.CauseOfTransmission{Cause: asdu.Inrogen}, asdu.GlobalCommonAddr,
	// 	asdu.SinglePointInfo{})
	// if err != nil {
	// 	// log.Println("falied")
	// } else {
	// 	// log.Println("success")
	// }
	// asduPack.SendReplyMirror(c, asdu.ActivationTerm)
	return nil
}
func (this *mysrv) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error {
	return nil
}
func (this *mysrv) ReadHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr) error { return nil }
func (this *mysrv) ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error   { return nil }
func (this *mysrv) ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error {
	return nil
}
func (this *mysrv) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error { return nil }
func (this *mysrv) ASDUHandler(asdu.Connect, *asdu.ASDU) error                     { return nil }
