package upload

import (
	"net/http"
	"strings"

	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/os/fsi/dsfs"

	"appengine"
)

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
