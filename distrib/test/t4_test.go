// +build p4
// go test -tags=p4

package test

import (
	"fmt"
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
	distrib.DefaultOptions.CollectRemainder = true
	do(t, 3)

}

func Test5(t *testing.T) {
	distrib.DefaultOptions.CollectRemainder = false
	do(t, 3)
}

func ATest6(t *testing.T) {
	distrib.DefaultOptions.CollectRemainder = true
	do(t, 33)
}

func ATest7(t *testing.T) {
	distrib.DefaultOptions.CollectRemainder = false
	do(t, 33)
}

func do(t *testing.T, precreatedPackets int) {

	jobs := make([]distrib.Worker, 0, precreatedPackets)
	for i := 0; i < precreatedPackets; i++ {
		job := distrib.Worker(&MyWorker{Inp: i})
		jobs = append(jobs, job)
	}

	ret := distrib.Distrib(jobs, distrib.DefaultOptions)

	if len(ret) != precreatedPackets {
		// t.Errorf("wnt %v got %v", precreatedPackets, len(ret))
	}

	for k, v := range ret {
		v1, _ := v.Worker.(*MyWorker)
		fmt.Printf("   %2v  %2v  => %v\n", k, v1.Inp, v1.Res)
		if v1.Inp != v1.Res-500 {
			t.Errorf("%v %v", v1.Inp, v1.Res)
		}
	}
	fmt.Printf("\n\n")

}
