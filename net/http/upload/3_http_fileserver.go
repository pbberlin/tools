package upload

import (
	"net/http"
	"strings"

	"appengine"

	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/httpfs"
)

var cx appengine.Context
var mountPoint string

func howIsContext(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	if cx == nil {
		wpf(w, "nil\n")
	} else {
		wpf(w, "context for mp %q is well: %v\n", mountPoint, cx)
	}
}

//
// A static fileserver is NOT working
// Since we need an appengine.context
//
// UNUSED, NOT WORKING
func serveDsFs(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	// Examples
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("c:\\temp"))))

	c := appengine.NewContext(r)
	cx = c

	mountPoint = dsfs.MountPointLast()
	fs1 := dsfs.New(
		dsfs.MountName(mountPoint),
		dsfs.AeContext(c),
	)
	httpFSys := &httpfs.HttpFs{SourceFs: fs1}

	// HERE is the trouble!
	http.Handle("/tmp1/", http.StripPrefix("/tmp1/", http.FileServer(httpFSys.Dir("./"))))

	wpf(w, "serving %v", mountPoint)

	// time.Sleep(5 * time.Second)

}

// We cannot use http.FileServer(http.Dir("./css/") to dispatch our dsfs files.
// We need the appengine context to initialize dsfs.
//
// Thus we re-implement a serveFile method:
func ServeDsFsFile(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	urlPath := m["dir"].(string)
	if len(urlPath) > 0 {
		urlPath = urlPath[1:]
	}

	prefix := "/mnt00"

	pos := strings.Index(urlPath, "/")
	if pos > 0 {
		prefix = "/" + urlPath[:pos]
	}
	if pos == -1 {
		prefix = "/" + urlPath
	}

	fs2 := dsfs.New(
		dsfs.MountName(prefix[1:]),
		dsfs.AeContext(appengine.NewContext(r)),
	)

	fileserver.FsiFileServer(fs2, prefix+"/", w, r)

}
