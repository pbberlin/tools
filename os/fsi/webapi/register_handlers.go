package webapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi/dsfs"
)

var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

const UriSetFSType = "/fsi/set-fs-type"
const UriDeleteSubtree = "/fsi/delete-subtree"

func InitHandlers() {
	http.HandleFunc(UriSetFSType, loghttp.Adapter(setFSType))
	http.HandleFunc("/fsi/create-objects", loghttp.Adapter(createSys))
	http.HandleFunc("/fsi/retrieve-by-query", loghttp.Adapter(retrieveByQuery))
	http.HandleFunc("/fsi/retrieve-by-read-dir", loghttp.Adapter(retrieveByReadDir))
	http.HandleFunc("/fsi/walk", loghttp.Adapter(walkH))
	http.HandleFunc("/fsi/remove", loghttp.Adapter(removeSubtree))

	// http.HandleFunc("/fsi/delete-all", loghttp.Adapter(deleteAll))
	http.HandleFunc(UriDeleteSubtree, loghttp.Adapter(DeleteSubtree))

	http.HandleFunc("/fsi/cntr/last", loghttp.Adapter(lastMountPoint))
	http.HandleFunc("/fsi/cntr/reset", loghttp.Adapter(resetMountPoint))
	http.HandleFunc("/fsi/cntr/incr", loghttp.Adapter(incrMountPoint))
	http.HandleFunc("/fsi/cntr/decr", loghttp.Adapter(decrMountPoint))
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {

	var b1 = new(bytes.Buffer)

	htmlfrag.Wb(b1, "filesystem interface", "")

	htmlfrag.Wb(b1, "set type", UriSetFSType)
	htmlfrag.Wb(b1, "create", "/fsi/create-objects")

	htmlfrag.Wb(b1, "query", "/fsi/retrieve-by-query")
	htmlfrag.Wb(b1, "readdir", "/fsi/retrieve-by-read-dir")
	htmlfrag.Wb(b1, "walk", "/fsi/walk")
	htmlfrag.Wb(b1, "remove subset", "/fsi/remove")

	// htmlfrag.Wb(b1, "delete all", "/fsi/delete-all", "all fs types")
	htmlfrag.Wb(b1, "delete tree", UriDeleteSubtree, "of selected fs")

	// htmlfrag.Wb(b1, , "")
	htmlfrag.Wb(b1, "dsfs mount", "nobr")
	htmlfrag.Wb(b1, "last", "/fsi/cntr/last")
	htmlfrag.Wb(b1, "decr", "/fsi/cntr/decr")
	htmlfrag.Wb(b1, "incr", "/fsi/cntr/incr")
	htmlfrag.Wb(b1, "reset", "/fsi/cntr/reset")

	return b1

}

func lastMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "last Mountpoint"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	wpf(w, "reset %v\n", dsfs.MountPointLast())

}
func resetMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Mountpoint reset"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	wpf(w, "reset %v\n", dsfs.MountPointReset())

}

func incrMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Mountpoint increment"}))
	defer wpf(w, "\n</pre>")

	wpf(w, "<pre>\n")
	defer wpf(w, tplx.Foot)

	xx := r.Header.Get("adapter_01")
	wpf(w, "adapter set %q\n", xx)

	wpf(w, "counted up %v\n", dsfs.MountPointIncr())

}

func decrMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Mountpoint decrement"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	wpf(w, "counted down %v\n", dsfs.MountPointDecr())

}
