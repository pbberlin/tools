// +build p4
// go test -tags=p4

package distrib

import (
	"math/rand"
	"testing"
	"time"
)

func Test4(t *testing.T) {

	// Precreate some WLoad packets.
	// Create more than we want.

	precreatedPackets := 33
	packets := make([]WLoad, 0, precreatedPackets)
	for i := 0; i < precreatedPackets; i++ {
		pack := WLoad{}
		pack.TaskID = i
		pack.Func = func(px WLoad) WLoad {
			px.Result = px.TaskID + 100
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)))
			return px
		}
		packets = append(packets, pack)
	}

	Distrib(packets, DefaultOptions)
}
