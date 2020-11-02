package main

import (
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

func main() {
	option := cs104.NewOption()
	err := option.AddRemoteServer("127.0.0.1:2404")
	if err != nil {
		panic(err)
	}

	srv := cs104.NewServerSpecial(&mysrv{}, option)

	srv.LogMode(true)

	srv.SetOnConnectHandler(func(c asdu.Connect) {
		_, _ = c.UnderlyingConn().Write([]byte{0x68, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x46, 0x01, 0x04, 0x00, 0xa0, 0xaf, 0xbd, 0xd8, 0x0a, 0xf4})
		log.Println("connected")
	})
	srv.SetConnectionLostHandler(func(c asdu.Connect) {
		log.Println("disconnected")
	})
	if err = srv.Start(); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":6060", nil); err != nil {
		panic(err)
	}
}

type mysrv struct{}

func (sf *mysrv) InterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
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
func (sf *mysrv) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error {
	return nil
}
func (sf *mysrv) ReadHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr) error { return nil }
func (sf *mysrv) ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error   { return nil }
func (sf *mysrv) ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error {
	return nil
}
func (sf *mysrv) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error { return nil }
func (sf *mysrv) ASDUHandler(asdu.Connect, *asdu.ASDU) error                     { return nil }
