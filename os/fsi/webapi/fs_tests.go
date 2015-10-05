// Package webapi contains handlers to manage and test
// fsi filesystems from a html UI.
package webapi

import (
	"bytes"
	"net/http"
	"strconv"

	"appengine"

	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
	"github.com/pbberlin/tools/os/fsi/tests"
)

var memMapFileSys = memfs.New()
var osFileSys = osfs.New()

// var dsFileSys = // cannot be instantiated without ae.context

var whichType = 0

func setFSType(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Set filesystem type"}))
	defer wpf(w, tplx.Foot)

	stp := r.FormValue("type")
	newTp, err := strconv.Atoi(stp)

	if err == nil && newTp >= 0 && newTp <= 2 {
		whichType = newTp
		wpf(w, "new type: %v<br><br>\n", whichType)
	}

	if whichType != 0 {
		wpf(w, "<a href='%v?type=0' >dsfs</a><br>\n", UriSetFSType)
	} else {
		wpf(w, "<b>dsfs</b><br>\n")
	}
	if whichType != 1 {
		wpf(w, "<a href='%v?type=1' >osfs</a><br>\n", UriSetFSType)
	} else {
		wpf(w, "<b>osfs</b><br>\n")
	}
	if whichType != 2 {
		wpf(w, "<a href='%v?type=2' >memfs</a><br>\n", UriSetFSType)
	} else {
		wpf(w, "<b>memfs</b><br>\n")
	}

}

func getFS(c appengine.Context, mnt string) fsi.FileSystem {

	var fs fsi.FileSystem

	switch whichType {
	case 0:
		// must be re-instantiated for each request
		dsFileSys := dsfs.New(dsfs.MountName(mnt), dsfs.AeContext(c))
		fs = fsi.FileSystem(dsFileSys)
	case 1:
		fs = fsi.FileSystem(osFileSys)
	case 2:
		// re-instantiation would delete everything
		fs = fsi.FileSystem(memMapFileSys)
	default:
		panic("invalid whichType ")
	}

	return fs
}

func runTestX(
	w http.ResponseWriter,
	r *http.Request,
	f1 func() string,
	f2 func(fsi.FileSystem) (*bytes.Buffer, string),
) {

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Run a test"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>\n")
	defer wpf(w, "\n</pre>")

	if f1 == nil {
		f1 = dsfs.MountPointLast
	}
	mnt := f1()
	fs := getFS(appengine.NewContext(r), mnt)

	bb := new(bytes.Buffer)
	msg := ""
	wpf(bb, "created fs %v\n\n", mnt)
	bb, msg = f2(fs)
	w.Write([]byte(msg))
	w.Write([]byte("\n\n"))
	w.Write(bb.Bytes())

}

func createSys(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, tests.CreateSys)
	// runTestX(w, r, dsfs.MountPointIncr, tests.CreateSys)
}

func retrieveByQuery(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, tests.RetrieveByQuery)
}

func retrieveByReadDir(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, tests.RetrieveByReadDir)
}

func walkH(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, tests.WalkDirs)
}

func removeSubtree(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	runTestX(w, r, nil, tests.RemoveSubtree)
}
