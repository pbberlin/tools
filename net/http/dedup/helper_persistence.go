package dedup

import (
	"bytes"
	"os"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/common"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
	"golang.org/x/net/html"

	"appengine"
)

var logDir = "c:/tmp/dedup/"

var memMapFileSys = memfs.New(memfs.DirSort("byDateDesc")) // package variable required as "persistence"

func GetFS(c appengine.Context, whichType int) (fs fsi.FileSystem) {

	switch whichType {
	case 0:
		// re-instantiation would delete contents
		fs = fsi.FileSystem(memMapFileSys)
	case 1:
		// must be re-instantiated for each request
		dsFileSys := dsfs.New(dsfs.DirSort("byDateDesc"), dsfs.MountName("mntTest"), dsfs.AeContext(c))
		fs = fsi.FileSystem(dsFileSys)
	case 2:

		osFileSys := osfs.New(osfs.DirSort("byDateDesc"))
		fs = fsi.FileSystem(osFileSys)
		os.Chdir(logDir)
	default:
		panic("invalid whichType ")
	}

	return
}

func fileDump(lg loghttp.FuncBufUniv, fs fsi.FileSystem,
	content interface{}, fNamer func() string, secondPart string) {

	if fNamer != nil {
		fn := fNamer() + secondPart
		switch casted := content.(type) {
		case *html.Node:
			var b bytes.Buffer
			err := html.Render(&b, casted)
			lg(err)
			if err != nil {
				return
			}
			err = common.WriteFile(fs, fn, b.Bytes())
			lg(err)
		case []byte:
			err := common.WriteFile(fs, fn, casted)
			lg(err)
		}

	}

}
