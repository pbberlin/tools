// +build weed2
// go test -tags=weed2

package weedout

import (
	"testing"
	"time"

	"appengine/aetest"

	"github.com/pbberlin/tools/net/http/domclean2"
	"github.com/pbberlin/tools/net/http/loghttp"
)

func Test2(t *testing.T) {

	start := time.Now()

	lg, b := loghttp.BuffLoggerUniversal(nil, nil)
	_ = b
	// closureOverBuf := func(bUnused *bytes.Buffer) {
	// 	loghttp.Pf(nil, nil, b.String())
	// }
	// defer closureOverBuf(b) // the argument is ignored,

	var c aetest.Context
	if false {
		var err error
		c, err = aetest.NewContext(nil)
		lg(err)
		if err != nil {
			return
		}
		defer c.Close()
	}
	fs := GetFS(c, 2)

	lg("took1 %4.2v secs", time.Now().Sub(start).Seconds())

	least3Files := FetchAndDecodeJSON(nil, URLs[0], lg, fs)

	lg("took2 %4.2v secs", time.Now().Sub(start).Seconds())

	doc := WeedOut(least3Files, lg, fs)

	fNamer := domclean2.FileNamer(logDir, 0)
	fNamer() // first call yields key
	fsPerm := GetFS(c, 2)
	fileDump(lg, fsPerm, doc, fNamer, "_fin.html")

	pf("MapSimiliarCompares: %v SimpleCompares: %v LevenstheinComp: %v\n", breakMapsTooDistinct, appliedLevenshtein, appliedCompare)
	pf("Finish\n")

}
