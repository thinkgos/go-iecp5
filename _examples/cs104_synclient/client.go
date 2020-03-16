// This examble uses FreyrSCADA IEC 60870-5-104 Simulator as Server
// Information Object should be setted according to
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

type ReadTable struct {
	ca  asdu.CommonAddr
	ioa asdu.InfoObjAddr
}
type WriteTable struct {
	ca    asdu.CommonAddr
	ioa   asdu.InfoObjAddr
	id    asdu.TypeID
	value interface{}
	quali interface{}
}

var readTable []*ReadTable
var writeTable []*WriteTable
var opt *cs104.ClientOption
var synclientInstance *cs104.Synclient

func init() {
	opt = cs104.NewOption()
	opt.AddRemoteServer("192.168.137.2:2404")
	opt.SetReconnectInterval(10 * time.Second)
	synclientInstance = cs104.NewSynclient(opt)
	// synclientInstance.LogMode(false)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	readTable = []*ReadTable{
		{1, 1},
		{1, 3},
		{1, 5},
		{1, 7},
		{1, 9},
		{1, 11},
		{1, 13},

		{1, 30},
		{1, 31},
		{1, 32},
		{1, 33},
		{1, 34},
		{1, 35},
		{1, 36},
	}

	writeTable = []*WriteTable{
		{1, 45, 45, asdu.SingleCommand(r.Intn(2) == 1), asdu.ParseQualifierOfCommand(0)},
		{1, 46, 46, asdu.DoubleCommand(r.Intn(4)), asdu.ParseQualifierOfCommand(0)},
		// According to 7.2.4.17 StepCommand value could be 0, 1, 2, 3
		// but the server returns cause: unknown info obj address
		{1, 47, 47, asdu.StepCommand(r.Intn(2) + 1), asdu.ParseQualifierOfCommand(0)},
		{1, 48, 48, asdu.NormalizedMeasurement(r.Intn(65536) - 32768), asdu.ParseQualifierOfCommand(0)},
		{1, 49, 49, asdu.ScaledMeasurement(r.Intn(65536) - 32768), asdu.ParseQualifierOfCommand(0)},
		{1, 50, 50, asdu.ShortFloatMeasurement(r.Float32()), asdu.ParseQualifierOfCommand(0)},
		{1, 51, 51, asdu.BitString(r.Uint32()), asdu.ParseQualifierOfCommand(0)},

		{1, 58, 58, asdu.SingleCommand(r.Intn(2) == 1), asdu.ParseQualifierOfCommand(0)},
		{1, 59, 59, asdu.DoubleCommand(r.Intn(4)), asdu.ParseQualifierOfCommand(0)},
		{1, 60, 60, asdu.StepCommand(r.Intn(2) + 1), asdu.ParseQualifierOfCommand(0)},
		{1, 61, 61, asdu.NormalizedMeasurement(r.Intn(65536) - 32768), asdu.ParseQualifierOfCommand(0)},
		{1, 62, 62, asdu.ScaledMeasurement(r.Intn(65536) - 32768), asdu.ParseQualifierOfCommand(0)},
		{1, 63, 63, asdu.ShortFloatMeasurement(r.Float32()), asdu.ParseQualifierOfCommand(0)},
		{1, 64, 64, asdu.BitString(r.Uint32()), asdu.ParseQualifierOfCommand(0)},
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	synclientInstance.Connect(ctx)
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return
			}
			ReadAndWrite()
		}
	}()
	Subscribing(ctx)
}

func ReadAndWrite() {
	for _, v := range writeTable {
		err := synclientInstance.Write(v.ca, v.ioa, v.id, v.value, v.quali)
		if err != nil {
			log.Fatalf("WriteRead: %v, %v, %v, %v, %+v Failed: %v", v.ca, v.ioa, v.id, v.value, v.quali, err)
		} else {
			log.Printf("WriteRead: %v, %v, %v, %v, %+v Succeed", v.ca, v.ioa, v.id, v.value, v.quali)
		}
	}

	for _, v := range readTable {
		resp, err := synclientInstance.Read(v.ca, v.ioa)
		if err != nil {
			log.Fatalf("TestRead: %v, %v Failed: %v", v.ca, v.ioa, err)
		} else {
			log.Printf("TestRead: %v, %v Succeed: %v", v.ca, v.ioa, resp.Value)
		}
	}
}

func Subscribing(ctx context.Context) {
	sub := make(chan *cs104.AsduInfo, 1)
	synclientInstance.Subscribe(sub)
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-sub:
			if uint8(data.Quality) != 0 {
				log.Printf("Wrong data quality: %v", data.Quality)
				continue
			}

			log.Printf("%v", String(data))
		}
	}
}

func String(m *cs104.AsduInfo) string {
	b, err := json.Marshal(*m)
	if err != nil {
		return fmt.Sprintf("%+v", *m)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *m)
	}
	return out.String()
}
