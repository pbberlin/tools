// +build weed2
// go test -tags=weed2

package weedout

import (
	"testing"

	"appengine/aetest"

	"github.com/pbberlin/tools/net/http/domclean2"
	"github.com/pbberlin/tools/net/http/loghttp"
)

func Test2(t *testing.T) {

	lg, b := loghttp.BuffLoggerUniversal(nil, nil)
	_ = b
	// closureOverBuf := func(bUnused *bytes.Buffer) {
	// 	loghttp.Pf(nil, nil, b.String())
	// }
	// defer closureOverBuf(b) // the argument is ignored,

	c, err := aetest.NewContext(nil)
	lg(err)
	if err != nil {
		return
	}
	defer c.Close()
	fs := GetFS(c, 2)

	least3Files := DecodeJSON(URLs[0], lg, fs)
	doc := WeedOut(least3Files, lg, fs)

	fNamer := domclean2.FileNamer(logDir, 0)
	fNamer() // first call yields key
	fsPerm := GetFS(c, 2)
	fileDump(lg, fsPerm, doc, fNamer, "_fin.html")

	pf("MapSimiliarCompares: %v SimpleCompares: %v LevenstheinComp: %v\n", breakMapsTooDistinct, appliedLevenshtein, appliedCompare)
	pf("Finish\n")

}
