package upload

import (
	"net/http"

	"appengine"

	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/httpfs"
)

var cx appengine.Context

func howIsContext(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	if cx == nil {
		wpf(w, "nil\n")
	} else {
		wpf(w, "context is well: %v\n", cx)

	}
}

func serveAefs(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	// Examples
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("c:\\temp"))))

	c := appengine.NewContext(r)
	cx = c

	mountPoint := aefs.MountPointLast()
	fs1 := aefs.New(
		aefs.MountName(mountPoint),
		aefs.AeContext(c),
	)
	httpFSys := &httpfs.HttpFs{SourceFs: fs1}

	http.Handle("/tmp1/", http.StripPrefix("/tmp1/", http.FileServer(httpFSys.Dir("./"))))

	wpf(w, "serving %v", mountPoint)

	// time.Sleep(5 * time.Second)

}
