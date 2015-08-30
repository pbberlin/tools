// +build weed2
// go test -tags=weed2

package weedout

import (
	"testing"

	"github.com/pbberlin/tools/net/http/loghttp"
)

func Test2(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	logdir := prepareLogDir()
	_ = logdir

	least3Files, err := similar("")
	lge(err)
	if err != nil {
		return
	}

	iter := make([]int, numTotal)

	for i, _ := range iter {
		lg("iterating %v", least3Files[i])
	}

}
