package fetch_rss

import (
	"testing"

	"github.com/pbberlin/tools/net/http/loghttp"

	"appengine/aetest"
)

func Test2(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	c, err := aetest.NewContext(nil)
	if err != nil {
		lge(err)
		t.Fatal(err)
	}
	defer c.Close()

	whichType = 2
	fs := getFs(c)
	lg(fs.Name() + "-" + fs.String())

	for _, config := range hosts {
		Fetch(nil, nil, fs, config, "/politik/international/aa/bb", 12)
	}

}
