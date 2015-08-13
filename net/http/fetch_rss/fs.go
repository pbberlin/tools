package fetch_rss

import (
	"net/http"
	"os"
	"strconv"

	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/httpfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"

	"appengine"
)

var docRoot = ""                                           // no relative path, 'cause working dir too flippant
var whichType = 2                                          // which type of filesystem, default is memfs
var memMapFileSys = memfs.New(memfs.DirSort("byDateDesc")) // package variable required as "persistence"

//
var httpFSys = &httpfs.HttpFs{SourceFs: fsi.FileSystem(memMapFileSys)} // memMap is always ready
var fileserver1 = http.FileServer(httpFSys.Dir(docRoot))

const uriSetType = "/fetch/set-fs-type"

const mountName = "mntftch"
const uriMountNameY = "/" + mountName + "/serve-file/"

func getFs(c appengine.Context) (fs fsi.FileSystem) {
	switch whichType {
	case 0:
		// must be re-instantiated for each request
		docRoot = ""
		dsFileSys := dsfs.New(dsfs.DirSort("byDateDesc"), dsfs.MountName(mountName), dsfs.AeContext(c))
		fs = fsi.FileSystem(dsFileSys)
	case 1:
		docRoot = "c:/docroot/"
		os.Chdir(docRoot)
		osFileSys := osfs.New(osfs.DirSort("byDateDesc"))
		fs = fsi.FileSystem(osFileSys)
	case 2:
		// re-instantiation would delete contents
		docRoot = ""
		fs = fsi.FileSystem(memMapFileSys)
	default:
		panic("invalid whichType ")
	}

	return
}

func setFSType(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Set fetcher reservoir filesystem type"}))
	defer wpf(w, tplx.Foot)

	stp := r.FormValue("type")
	newTp, err := strconv.Atoi(stp)

	if err == nil && newTp >= 0 && newTp <= 2 {
		whichType = newTp
		wpf(w, "new type: %v<br><br>\n", whichType)
	}

	if whichType != 0 {
		wpf(w, "<a href='%v?type=0' >dsfs</a><br>\n", uriSetType)
	} else {
		wpf(w, "<b>dsfs</b><br>\n")
	}
	if whichType != 1 {
		wpf(w, "<a href='%v?type=1' >osfs</a><br>\n", uriSetType)
	} else {
		wpf(w, "<b>osfs</b><br>\n")
	}
	if whichType != 2 {
		wpf(w, "<a href='%v?type=2' >memfs</a><br>\n", uriSetType)
	} else {
		wpf(w, "<b>memfs</b><br>\n")
	}

}
