package fetch_rss

import (
	"bytes"
	"net/http"
	"path"

	"appengine"

	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
)

func InitHandlers() {
	http.HandleFunc(uriSetType, loghttp.Adapter(setFSType))

	http.HandleFunc("/fetch/request", loghttp.Adapter(requestFetch))

	// working only for memfs
	http.Handle("/fetch/reservoire/static/", http.StripPrefix("/fetch/reservoire/static/", fileserver1))

	http.HandleFunc(uriMountNameY, loghttp.Adapter(serveFile))

}

func serveFile(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	fs1 := getFs(appengine.NewContext(r))
	fileserver.FsiFileServer(fs1, uriMountNameY, w, r)
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Fetcher", "")
	htmlfrag.Wb(b1, "type", uriSetType, "storage type")

	htmlfrag.Wb(b1, "add", "/fetch/request", "add")

	htmlfrag.Wb(b1, "reservoire static", "/fetch/reservoire/static/", "browse - memfs only")

	htmlfrag.Wb(b1, "reservoire BOTH", uriMountNameY, "browse ANY fsi.FileSystem")
	return b1
}

func requestFetch(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)
	var err error

	fs := getFs(appengine.NewContext(r))
	// fs = fsi.FileSystem(memMapFileSys)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Requesting files"}))
	defer wpf(w, tplx.Foot)
	wpf(w, "<pre>")
	defer wpf(w, "</pre>")

	err = fs.WriteFile(path.Join(docRoot, "msg.html"), msg, 0644)
	lge(err)

	// err = fs.WriteFile(path.Join(docRoot, "index.html"), []byte("content of index.html"), 0644)
	// lge(err)

	err = fs.Mkdir(path.Join(docRoot, "dirX"), 0755)
	lge(err)

	lg("static fileserver created")

	for _, config := range hosts {
		Fetch(w, r, fs, config, "/politik/international/aa/bb", 12)
	}

	lg("fetching complete")

}
