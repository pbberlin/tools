package aefs

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"appengine"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi/fsc"
	"github.com/pbberlin/tools/util"
)

func init() {
	http.HandleFunc("/fs/aefs/create-objects", loghttp.Adapter(demoSaveRetrieve))
	http.HandleFunc("/fs/aefs/retrieve-by-query", loghttp.Adapter(retrieveByQuery))
	http.HandleFunc("/fs/aefs/delete-all", loghttp.Adapter(deleteAll))
	http.HandleFunc("/fs/aefs/walk", loghttp.Adapter(walkH))
}

func deleteAll(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>")
	defer wpf(w, "</pre>")
	defer wpf(w, tplx.Foot)

	rts := fmt.Sprintf("mnt%02v", util.CounterLast())
	fs := NewAeFs(rts, AeContext(appengine.NewContext(r)))
	msg, err := fs.DeleteAll()
	if err != nil {
		wpf(w, "err during delete %v\n", err)
	}
	wpf(w, msg)

}
func demoSaveRetrieve(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>")
	defer wpf(w, tplx.Foot)
	defer wpf(w, "</pre>") // strange order

	rt := <-util.Counter
	rts := fmt.Sprintf("mnt%02v", rt)
	bb := CreateSys(appengine.NewContext(r), rts)
	w.Write(bb.Bytes())

}

func retrieveByQuery(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	wpf(w, "<pre>")
	defer wpf(w, tplx.Foot)
	defer wpf(w, "</pre>")

	rts := fmt.Sprintf("mnt%02v", util.CounterLast())
	bb := RetrieveDirs(appengine.NewContext(r), rts)
	w.Write(bb.Bytes())

}

func walkH(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	fmt.Fprint(w, tplx.Head)
	defer wpf(w, tplx.Foot)

	rts := fmt.Sprintf("mnt%02v", util.CounterLast())
	fs := NewAeFs(rts, AeContext(appengine.NewContext(r)))
	loghttp.Pf(w, r, "-------filewalk----<br>")

	var bb bytes.Buffer
	walkFunc := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			bb.WriteString(spf("Visting path %s: error %v \n<br>", path, err))
			return err
		} else {
			tp := "file"
			if f != nil {
				if f.IsDir() {
					tp = "dir "
				}
			}
			bb.WriteString(spf("Visited: %s %s \n<br>", tp, path))
		}
		return nil
	}

	var err error

	bb = bytes.Buffer{}
	err = fsc.Walk(fs, fs.RootName(), walkFunc)
	bb.WriteString(spf("fs.Walk() returned %v\n<br><br>", err))
	w.Write(bb.Bytes())

	bb = bytes.Buffer{}
	err = fsc.Walk(fs, "ch1/ch2", walkFunc)
	bb.WriteString(spf("fs.Walk() returned %v\n<br><br>", err))
	w.Write(bb.Bytes())

	bb = bytes.Buffer{}
	err = fsc.Walk(fs, "ch1/ch2/ch3", walkFunc)
	bb.WriteString(spf("fs.Walk() returned %v\n<br><br>", err))
	w.Write(bb.Bytes())

	//
	err = fs.RemoveAll("ch1/ch2/ch3")
	wpf(w, "fs.RemoveAll() returned %v\n<br><br>", err)

	fmt.Fprint(w, tplx.Foot)

}
