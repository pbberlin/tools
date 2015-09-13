// +build p4
// go test -tags=p4

package test

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/pbberlin/tools/distrib"
)

type MyWorker struct {
	Inp, Res int
}

func (m *MyWorker) Work() {
	m.Res = 500 + m.Inp
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)))
}

func Test4(t *testing.T) {

	// Precreate some distrib.WLoad packets.
	// Create more than we want.

	precreatedPackets := 33
	packets := make([]distrib.WLoad, 0, precreatedPackets)
	for i := 0; i < precreatedPackets; i++ {
		pack := distrib.WLoad{}
		pack.Workload = &MyWorker{Inp: i}

		packets = append(packets, pack)
	}

	ret := distrib.Distrib(packets, distrib.DefaultOptions)

	for k, v := range ret {
		log.Printf("%2v  %v ", k, v.Workload)
	}

}
