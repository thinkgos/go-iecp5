package cs104

import (
	"time"
)

type seqPending struct {
	seq      uint16
	sendTime time.Time
}
