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

const UriSetFSType = "/fs/dsfs/set-fs-type"
const UriDeleteSubtree = "/fs/dsfs/delete-subtree"

func InitHandlers() {
	http.HandleFunc(UriSetFSType, loghttp.Adapter(setFSType))
	http.HandleFunc("/fs/dsfs/create-objects", loghttp.Adapter(createSys))
	http.HandleFunc("/fs/dsfs/retrieve-by-query", loghttp.Adapter(retrieveByQuery))
	http.HandleFunc("/fs/dsfs/retrieve-by-read-dir", loghttp.Adapter(retrieveByReadDir))
	http.HandleFunc("/fs/dsfs/walk", loghttp.Adapter(walkH))
	http.HandleFunc("/fs/dsfs/remove", loghttp.Adapter(removeSubtree))

	http.HandleFunc("/fs/dsfs/delete-all", loghttp.Adapter(deleteAll))
	http.HandleFunc(UriDeleteSubtree, loghttp.Adapter(deleteSubtree))

	http.HandleFunc("/fs/dsfs/reset", loghttp.Adapter(resetMountPoint))
	http.HandleFunc("/fs/dsfs/incr", loghttp.Adapter(incrMountPoint))
	http.HandleFunc("/fs/dsfs/decr", loghttp.Adapter(decrMountPoint))
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {

	var b1 = new(bytes.Buffer)

	htmlfrag.Wb(b1, "filesystem interface", "")

	htmlfrag.Wb(b1, "set type", UriSetFSType)
	htmlfrag.Wb(b1, "create", "/fs/dsfs/create-objects")

	htmlfrag.Wb(b1, "query", "/fs/dsfs/retrieve-by-query")
	htmlfrag.Wb(b1, "readdir", "/fs/dsfs/retrieve-by-read-dir")
	htmlfrag.Wb(b1, "walk", "/fs/dsfs/walk")
	htmlfrag.Wb(b1, "remove", "/fs/dsfs/remove")

	htmlfrag.Wb(b1, "delete all", "/fs/dsfs/delete-all", "all fs types")
	htmlfrag.Wb(b1, "delete tree", "/fs/dsfs/delete-subtree", "of selected fs")

	// htmlfrag.Wb(b1, , "")
	htmlfrag.Wb(b1, "dsfs mount", "nobr")
	htmlfrag.Wb(b1, "decr", "/fs/dsfs/decr")
	htmlfrag.Wb(b1, "incr", "/fs/dsfs/incr")
	htmlfrag.Wb(b1, "reset", "/fs/dsfs/reset")

	return b1

}

func resetMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Mountpoint reset"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	wpf(w, "reset %v\n", dsfs.MountPointReset())

}

func incrMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Mountpoint increment"}))
	defer wpf(w, "\n</pre>")

	wpf(w, "<pre>\n")
	defer wpf(w, tplx.Foot)

	xx := r.Header.Get("adapter_01")
	wpf(w, "adapter set %q\n", xx)

	wpf(w, "counted up %v\n", dsfs.MountPointIncr())

}

func decrMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Mountpoint decrement"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	wpf(w, "counted down %v\n", dsfs.MountPointDecr())

}
