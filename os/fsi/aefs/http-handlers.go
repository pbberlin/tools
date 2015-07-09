package aefs

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	pth "path"

	"appengine"

	"github.com/pbberlin/tools/logif"
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

	rts := fmt.Sprintf("mnt%02v", util.CounterLast())
	fs := NewAeFs(rts, AeContext(appengine.NewContext(r)))
	msg, err := fs.DeleteAll()

	loghttp.E(w, r, err, true)

	wpf(w, "<pre>%v</pre>", msg)

	wpf(w, tplx.Foot)

}
func demoSaveRetrieve(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	wpf(w, tplx.Head)
	defer wpf(w, tplx.Foot)

	rt := <-util.Counter
	rts := fmt.Sprintf("mnt%02v", rt)
	fs := NewAeFs(rts, AeContext(appengine.NewContext(r)))
	loghttp.Pf(w, r, "created fs %v<br>", rts)

	fc1 := func(p []string) {
		path := pth.Join(p...)
		path = cleanseLeadingSlash(path)

		dir, err := fs.saveDirByPath(path)
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child created %3v - %v ", dir.Name(), dir.Key)
	}

	wpf(w, "<pre>")

	loghttp.Pf(w, r, "--------create-dirs---------")

	fc1([]string{"ch1"})
	fc1([]string{"ch1", "ch2"})
	fc1([]string{"ch1", "ch2a"})
	fc1([]string{"ch1", "ch2", "ch3"})
	fc1([]string{"ch1", "ch2", "ch3", "ch4"})
	fc1([]string{"ch1A"})
	fc1([]string{"ch1B"})
	fc1([]string{"ch1", "d2", "d3", "d4"})
	fc1([]string{"d1", "d2", "d3", "d4"})

	loghttp.Pf(w, r, "\n--------retrieve-dirs---------")

	// retrieval
	fc2 := func(p []string) {
		path := pth.Join(p...)
		loghttp.Pf(w, r, "searching... %v", path)
		f, err := fs.dirByPath(path)
		if err != nil {
			loghttp.Pf(w, r, "   nothing retrieved - err %v", err)
		} else {
			loghttp.Pf(w, r, "   fnd %v ", f.Dir+f.Name())
		}
	}
	fc2([]string{"ch1"})
	fc2([]string{"ch1", "ch2"})
	fc2([]string{"ch1", "cha2"})
	fc2([]string{"ch1", "ch2", "ch3"})
	fc2([]string{"fsd,mount000", "fsd,ch1", "ch2", "ch3"})
	fc2([]string{"ch1A"})
	fc2([]string{fs.RootDir()})

	loghttp.Pf(w, r, "\n-------create and save some files----")

	fc4a := func(name, content string) {
		err := fs.WriteFile(name, []byte(content), os.ModePerm)
		logif.E(err)
	}
	fc4b := func(name, content string) {
		f, err := fs.Create(name)
		logif.E(err)
		if err != nil {
			return
		}
		_, err = f.WriteString(content)
		logif.E(err)
		err = f.Sync()
		logif.E(err)
	}

	fc4a("ch1/ch2/file1", "content 1")
	fc4b("ch1/ch2/file2", "content 2")
	fc4a("ch1/ch2/ch3/file3", "another content")
	fc4b(fs.RootDir()+"file4", "chq content 2")

	loghttp.Pf(w, r, "\n-------retrieve files again----")

	fc5 := func(path string) {
		files, err := fs.filesByPath(fs.RootDir() + path)
		logif.E(err)
		loghttp.Pf(w, r, " srch %v  ", fs.RootDir()+path)
		for k, v := range files {
			loghttp.Pf(w, r, "     %v  -  %v %s", k, v.Name(), v.Data)
		}
	}

	fc5("ch1/ch2")
	fc5("ch1/ch2/ch3")
	fc5("")

	wpf(w, "</pre>")

}

func retrieveByQuery(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	fmt.Fprint(w, tplx.Head)
	defer wpf(w, tplx.Foot)

	rts := fmt.Sprintf("mnt%02v", util.CounterLast())
	fs := NewAeFs(rts, AeContext(appengine.NewContext(r)))
	loghttp.Pf(w, r, "created fs %v<br>", rts)

	loghttp.Pf(w, r, "--------retrieve by query---------<br>")

	fc3 := func(path string, direct bool) {
		loghttp.Pf(w, r, "searching ---  %v", path)
		children, err := fs.subdirsByPath(path, direct)
		if err != nil {
			loghttp.Pf(w, r, "   nothing retrieved - err %v", err)
		} else {
			for k, v := range children {
				loghttp.Pf(w, r, "child #%2v %-24v", k, v.Dir+v.Name())
			}
		}
	}
	wpf(w, "<pre>")
	fc3(spf(`ch1/ch2/ch3`), false)
	fc3(spf(`ch1/ch2/ch3`), true)
	fc3(spf(`ch1`), false)
	fc3(spf(`ch1`), true)
	fc3(spf(``), true)
	wpf(w, "</pre>")

}

func walkH(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	fmt.Fprint(w, tplx.Head)
	defer wpf(w, tplx.Foot)

	rts := fmt.Sprintf("mnt%02v", util.CounterLast())
	fs := NewAeFs(rts, AeContext(appengine.NewContext(r)))
	loghttp.Pf(w, r, "-------filewalk----<br>")

	var bb bytes.Buffer
	walkFunc := func(path string, f os.FileInfo, err error) error {
		tp := "file"
		if f != nil {
			if f.IsDir() {
				tp = "dir "
			}
		}
		bb.WriteString(spf("Visited: %s %s \n<br>", tp, path))
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

	fmt.Fprint(w, tplx.Foot)

}
