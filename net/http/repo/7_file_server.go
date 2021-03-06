package repo

import (
	"net/http"
	"os"
	"strconv"

	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// GetFS instantiates a filesystem, depending on whichtype
func GetFS(c context.Context) (fs fsi.FileSystem) {
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

// setFSType sets an internal variable, determining what FileSystems
// should be used. Default is dsfs.
func setFSType(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Set fetcher reservoir filesystem type"}))
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

// UNUSED, since appengine context is required for our filesystems
func serveSingleRootFile(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename) // filename refers to local path; unusable for fsi
	})
}

// serveFile makes the previously fetched files available like
// a static fileserver.
func serveFile(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	fs1 := GetFS(appengine.NewContext(r))
	fileserver.FsiFileServer(w, r, fileserver.Options{FS: fs1, Prefix: UriMountNameY})
}
