// +build fetch1
// go test -tags=fetch1

package repo

import (
	"testing"

	"github.com/pbberlin/tools/net/http/loghttp"

	"appengine/aetest"
)

func Test1(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	c, err := aetest.NewContext(nil)
	if err != nil {
		lge(err)
		t.Fatal(err)
	}
	defer c.Close()

	whichType = 2
	fs := GetFS(c)
	lg(fs.Name() + "-" + fs.String())

	for _, config := range testCommands {
		Fetch(nil, nil, fs, config)
	}

}
