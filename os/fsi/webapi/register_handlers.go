package webapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"appengine"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
	"github.com/pbberlin/tools/os/fsi/test"
)

var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

var memMapFileSys = memfs.New()
var osFileSys = osfs.New()

// var dsFileSys = // cannot be instantiated without ae.context

var whichType = 0

func InitHandlers() {
	http.HandleFunc("/fs/dsfs/set-fs-type", loghttp.Adapter(setFSType))
	http.HandleFunc("/fs/dsfs/create-objects", loghttp.Adapter(createSys))
	http.HandleFunc("/fs/dsfs/retrieve-by-query", loghttp.Adapter(retrieveByQuery))
	http.HandleFunc("/fs/dsfs/retrieve-by-read-dir", loghttp.Adapter(retrieveByReadDir))
	http.HandleFunc("/fs/dsfs/walk", loghttp.Adapter(walkH))
	http.HandleFunc("/fs/dsfs/remove", loghttp.Adapter(removeSubtree))

	http.HandleFunc("/fs/dsfs/delete-all", loghttp.Adapter(deleteAll))

	http.HandleFunc("/fs/dsfs/reset", loghttp.Adapter(resetMountPoint))
	http.HandleFunc("/fs/dsfs/incr", loghttp.Adapter(incrMountPoint))
	http.HandleFunc("/fs/dsfs/decr", loghttp.Adapter(decrMountPoint))
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {

	var b1 = new(bytes.Buffer)

	htmlfrag.Wb(b1, "filesystem interface", "")

	htmlfrag.Wb(b1, "set type", "/fs/dsfs/set-fs-type")
	htmlfrag.Wb(b1, "create", "/fs/dsfs/create-objects")

	htmlfrag.Wb(b1, "query", "/fs/dsfs/retrieve-by-query")
	htmlfrag.Wb(b1, "readdir", "/fs/dsfs/retrieve-by-read-dir")
	htmlfrag.Wb(b1, "walk", "/fs/dsfs/walk")
	htmlfrag.Wb(b1, "remove", "/fs/dsfs/remove")

	htmlfrag.Wb(b1, "delete all fs entities", "/fs/dsfs/delete-all")

	// htmlfrag.Wb(b1, , "")
	htmlfrag.Wb(b1, "dsfs mount", "nobr")
	htmlfrag.Wb(b1, "decr", "/fs/dsfs/decr")
	htmlfrag.Wb(b1, "incr", "/fs/dsfs/incr")
	htmlfrag.Wb(b1, "reset", "/fs/dsfs/reset")

	return b1

}

func setFSType(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Set filesystem type"}))
	defer wpf(w, tplx.Foot)

	stp := r.FormValue("type")
	newTp, err := strconv.Atoi(stp)

	if err == nil && newTp >= 0 && newTp <= 2 {
		whichType = newTp
		wpf(w, "new type: %v<br><br>\n", whichType)
	}

	if whichType != 0 {
		wpf(w, "<a href='/fs/dsfs/set-fs-type?type=0' >dsfs</a><br>\n")
	} else {
		wpf(w, "<b>dsfs</b><br>\n")
	}
	if whichType != 1 {
		wpf(w, "<a href='/fs/dsfs/set-fs-type?type=1' >osfs</a><br>\n")
	} else {
		wpf(w, "<b>osfs</b><br>\n")
	}
	if whichType != 2 {
		wpf(w, "<a href='/fs/dsfs/set-fs-type?type=2' >memfs</a><br>\n")
	} else {
		wpf(w, "<b>memfs</b><br>\n")
	}

}

func runTestX(
	w http.ResponseWriter,
	r *http.Request,
	f1 func() string,
	f2 func(fsi.FileSystem) (*bytes.Buffer, string),
) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Run a test"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	var fs fsi.FileSystem

	switch whichType {
	case 0:
		// must be re-instantiated for each request
		if f1 == nil {
			f1 = dsfs.MountPointLast
		}
		dsFileSys := dsfs.New(dsfs.MountName(f1()), dsfs.AeContext(appengine.NewContext(r)))
		fs = fsi.FileSystem(dsFileSys)
	case 1:
		fs = fsi.FileSystem(osFileSys)
	case 2:
		// re-instantiation would delete everything
		fs = fsi.FileSystem(memMapFileSys)
	default:
		panic("invalid whichType ")
	}

	bb := new(bytes.Buffer)
	msg := ""
	wpf(bb, "created fs %v\n\n", dsfs.MountPointLast())
	bb, msg = f2(fs)
	w.Write([]byte(msg))
	w.Write([]byte("\n\n"))
	w.Write(bb.Bytes())

}

func createSys(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, dsfs.MountPointIncr, test.CreateSys)
}

func retrieveByQuery(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, test.RetrieveByQuery)
}

func retrieveByReadDir(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, test.RetrieveByReadDir)
}

func walkH(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, test.WalkDirs)
}

func removeSubtree(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, test.RemoveSubtree)
}

//
// dsfs specific
func deleteAll(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Delete all filesystem data"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	fs := dsfs.New(dsfs.AeContext(appengine.NewContext(r)))
	wpf(w, "dsfs:\n")
	msg, err := fs.DeleteAll()
	if err != nil {
		wpf(w, "err during delete %v\n", err)
	}
	wpf(w, msg)

	memMapFileSys = memfs.New()
	wpf(w, "\n")
	wpf(w, "memMapFs new")

	// cleanup must be manual
	osFileSys = osfs.New()

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
