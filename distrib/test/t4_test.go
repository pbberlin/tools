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

var opt = distrib.NewDefaultOptions()

type MyWorker struct {
	Inp, Res int
}

func (m *MyWorker) Work() {
	m.Res = 500 + m.Inp
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)))
}

func Test4(t *testing.T) {
	opt.CollectRemainder = true
	do(t, 3)

}

func Test5(t *testing.T) {
	// opt.CollectRemainder = false
	// do(t, 3)
}

func Test6(t *testing.T) {
	opt.CollectRemainder = true
	opt.Want = 4
	do(t, 20)
}

func Test7(t *testing.T) {
	opt.CollectRemainder = false
	opt.Want = 4
	do(t, 20)
}

func do(t *testing.T, precreatedPackets int) {

	rand.Seed(time.Now().UnixNano())

	fmt.Printf("\n--------------------\n")

	jobs := make([]distrib.Worker, 0, precreatedPackets)
	for i := 0; i < precreatedPackets; i++ {
		job := distrib.Worker(&MyWorker{Inp: i})
		jobs = append(jobs, job)
	}

	ret, msg := distrib.Distrib(jobs, opt)
	fmt.Print(msg.String())

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

}
