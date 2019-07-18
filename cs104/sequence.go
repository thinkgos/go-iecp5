package cs104

import (
	"time"
)

type seqPending struct {
	seq      uint16
	sendTime time.Time
}

//
//type seqManage struct {
//	list []seqPending
//}
//
//func (this *seqManage) peek() time.Time {
//	return this.list[0].sendTime
//}
//
//func (this *seqManage) push(pending seqPending) {
//	this.list = append(this.list, pending)
//
//}
//
//func (this *seqManage) confirmReception(ackNo uint16) {
//	for i, v := range this.list {
//		if v.seq == (ackNo - 1) {
//			this.list = this.list[i+1:]
//			return
//		}
//	}
//}
