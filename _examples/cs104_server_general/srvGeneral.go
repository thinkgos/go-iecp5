package main

import (
	"log"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

func main() {
	srv := cs104.NewServer(&mysrv{})
	srv.SetOnConnectionHandler(func(c asdu.Connect) {
		log.Println("on connect")
	})
	srv.SetConnectionLostHandler(func(c asdu.Connect) {
		log.Println("connect lost")
	})
	srv.LogMode(true)
	// go func() {
	// 	time.Sleep(time.Second * 20)
	// 	log.Println("try ooooooo", err)
	// 	err := srv.Close()
	// 	log.Println("ooooooo", err)
	// }()
	if err := srv.ListenAndServer(":2404"); err != nil {
		log.Println(err)
		return
	}
}

type mysrv struct{}

func (sf *mysrv) InterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
	log.Println("qoi", qoi)
	asduPack.SendReplyMirror(c, asdu.ActivationCon)
	err := asdu.Single(c, false, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asdu.GlobalCommonAddr,
		asdu.SinglePointInfo{})
	if err != nil {
		// log.Println("falied")
	} else {
		// log.Println("success")
	}
	// go func() {
	// 	for {
	// 		err := asdu.Single(c, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.GlobalCommonAddr,
	// 			asdu.SinglePointInfo{})
	// 		if err != nil {
	// 			log.Println("falied", err)
	// 		} else {
	// 			log.Println("success", err)
	// 		}

	// 		time.Sleep(time.Second * 1)
	// 	}
	// }()
	asduPack.SendReplyMirror(c, asdu.ActivationTerm)
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
