package backend

import (
	"bytes"
	"net/http"

	"appengine"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi/aefs"
)

var backendFragFsiAefs = new(bytes.Buffer)

func init() {
	http.HandleFunc("/fs/aefs/delete-all", loghttp.Adapter(deleteAll))

	http.HandleFunc("/fs/aefs/create-objects", loghttp.Adapter(demoSaveRetrieve))
	http.HandleFunc("/fs/aefs/retrieve-by-query", loghttp.Adapter(retrieveByQuery))
	http.HandleFunc("/fs/aefs/walk", loghttp.Adapter(walkH))

	http.HandleFunc("/fs/aefs/reset", loghttp.Adapter(resetMountPoint))
	http.HandleFunc("/fs/aefs/decr", loghttp.Adapter(decrMountPoint))

	htmlfrag.Wb(backendFragFsiAefs, "aefs create", "/fs/aefs/create-objects")

	htmlfrag.Wb(backendFragFsiAefs, "query", "/fs/aefs/retrieve-by-query")
	htmlfrag.Wb(backendFragFsiAefs, "walk", "/fs/aefs/walk")
	htmlfrag.Wb(backendFragFsiAefs, "delete all fs entities", "/fs/aefs/delete-all")

	htmlfrag.Wb(backendFragFsiAefs, "reset", "/fs/aefs/reset")
	htmlfrag.Wb(backendFragFsiAefs, "decr", "/fs/aefs/decr")

}

func deleteAll(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")
	defer wpf(w, tplx.Foot)

	fs := aefs.NewAeFs(aefs.MountPointLast(), aefs.AeContext(appengine.NewContext(r)))
	msg, err := fs.DeleteAll()
	if err != nil {
		wpf(w, "err during delete %v\n", err)
	}
	wpf(w, msg)

}

func decrMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")
	defer wpf(w, tplx.Foot)

	wpf(w, "counted down %v\n", aefs.MountPointDecr())

}

func resetMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")
	defer wpf(w, tplx.Foot)

	wpf(w, "reset %v\n", aefs.MountPointReset())

}

func demoSaveRetrieve(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, tplx.Foot)
	defer wpf(w, "\n</pre>") // strange order

	bb := aefs.CreateSys(appengine.NewContext(r))
	w.Write(bb.Bytes())

}

func retrieveByQuery(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, tplx.Foot)
	defer wpf(w, "\n</pre>")

	bb := aefs.RetrieveDirs(appengine.NewContext(r))
	w.Write(bb.Bytes())

}

func walkH(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, tplx.Foot)
	defer wpf(w, "\n</pre>")

	bb := aefs.WalkDirs(appengine.NewContext(r))
	w.Write(bb.Bytes())

}
