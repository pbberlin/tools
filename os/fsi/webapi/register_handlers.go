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
	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/fstest"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
)

var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

var memMapFileSys = memfs.New()
var osFileSys = osfs.New()

// var aeFileSys = // cannot be instantiated without ae.context

var whichType = 0

func InitHandlers() {
	http.HandleFunc("/fs/aefs/set-fs-type", loghttp.Adapter(setFSType))
	http.HandleFunc("/fs/aefs/create-objects", loghttp.Adapter(createSys))
	http.HandleFunc("/fs/aefs/retrieve-by-query", loghttp.Adapter(retrieveByQuery))
	http.HandleFunc("/fs/aefs/retrieve-by-read-dir", loghttp.Adapter(retrieveByReadDir))
	http.HandleFunc("/fs/aefs/walk", loghttp.Adapter(walkH))
	http.HandleFunc("/fs/aefs/remove", loghttp.Adapter(removeSubtree))

	http.HandleFunc("/fs/aefs/delete-all", loghttp.Adapter(deleteAll))

	http.HandleFunc("/fs/aefs/reset", loghttp.Adapter(resetMountPoint))
	http.HandleFunc("/fs/aefs/incr", loghttp.Adapter(incrMountPoint))
	http.HandleFunc("/fs/aefs/decr", loghttp.Adapter(decrMountPoint))
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {

	var b1 = new(bytes.Buffer)

	htmlfrag.Wb(b1, "filesystem interface", "")

	htmlfrag.Wb(b1, "set type", "/fs/aefs/set-fs-type")
	htmlfrag.Wb(b1, "create", "/fs/aefs/create-objects")

	htmlfrag.Wb(b1, "query", "/fs/aefs/retrieve-by-query")
	htmlfrag.Wb(b1, "readdir", "/fs/aefs/retrieve-by-read-dir")
	htmlfrag.Wb(b1, "walk", "/fs/aefs/walk")
	htmlfrag.Wb(b1, "remove", "/fs/aefs/remove")

	htmlfrag.Wb(b1, "delete all fs entities", "/fs/aefs/delete-all")

	// htmlfrag.Wb(b1, , "")
	htmlfrag.Wb(b1, "aefs mount", "nobr")
	htmlfrag.Wb(b1, "decr", "/fs/aefs/decr")
	htmlfrag.Wb(b1, "incr", "/fs/aefs/incr")
	htmlfrag.Wb(b1, "reset", "/fs/aefs/reset")

	return b1

}

func setFSType(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	defer wpf(w, tplx.Foot)

	stp := r.FormValue("type")
	newTp, err := strconv.Atoi(stp)

	if err == nil && newTp >= 0 && newTp <= 2 {
		whichType = newTp
		wpf(w, "new type: %v<br><br>\n", whichType)
	}

	if whichType != 0 {
		wpf(w, "<a href='/fs/aefs/set-fs-type?type=0' >aefs</a><br>\n")
	} else {
		wpf(w, "<b>aefs</b><br>\n")
	}
	if whichType != 1 {
		wpf(w, "<a href='/fs/aefs/set-fs-type?type=1' >osfs</a><br>\n")
	} else {
		wpf(w, "<b>osfs</b><br>\n")
	}
	if whichType != 2 {
		wpf(w, "<a href='/fs/aefs/set-fs-type?type=2' >memfs</a><br>\n")
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

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, tplx.Foot)
	defer wpf(w, "\n</pre>")

	var fs fsi.FileSystem

	switch whichType {
	case 0:
		// must be re-instantiated for each request
		if f1 == nil {
			f1 = aefs.MountPointLast
		}
		aeFileSys := aefs.New(aefs.MountName(f1()), aefs.AeContext(appengine.NewContext(r)))
		fs = fsi.FileSystem(aeFileSys)
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
	wpf(bb, "created fs %v\n\n", aefs.MountPointLast())
	bb, msg = f2(fs)
	w.Write([]byte(msg))
	w.Write([]byte("\n\n"))
	w.Write(bb.Bytes())

}

func createSys(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, aefs.MountPointIncr, fstest.CreateSys)
}

func retrieveByQuery(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, fstest.RetrieveByQuery)
}

func retrieveByReadDir(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, fstest.RetrieveByReadDir)
}

func walkH(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, fstest.WalkDirs)
}

func removeSubtree(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, fstest.RemoveSubtree)
}

//
// aefs specific
func deleteAll(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")
	defer wpf(w, tplx.Foot)

	fs := aefs.New(aefs.AeContext(appengine.NewContext(r)))
	wpf(w, "aefs:\n")
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

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")
	defer wpf(w, tplx.Foot)

	wpf(w, "reset %v\n", aefs.MountPointReset())

}

func incrMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")
	defer wpf(w, tplx.Foot)

	wpf(w, "counted up %v\n", aefs.MountPointIncr())

}

func decrMountPoint(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")
	defer wpf(w, tplx.Foot)

	wpf(w, "counted down %v\n", aefs.MountPointDecr())

}
