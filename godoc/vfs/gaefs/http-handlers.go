package gaefs

import (
	"fmt"
	"net/http"
	pth "path"

	"appengine"
	"appengine/datastore"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/util"
)

var nestedOrRooted = true

func deleteAll(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	{
		q := datastore.NewQuery(tfil).KeysOnly()
		var files []AeFile
		keys, err := q.GetAll(c, &files)
		loghttp.E(w, r, err, false)
		err = datastore.DeleteMulti(c, keys)
		loghttp.E(w, r, err, false)
		loghttp.Pf(w, r, "%v files deleted", len(keys))
	}

	{
		q := datastore.NewQuery(tdir).KeysOnly()
		var dirs []AeDir
		keys, err := q.GetAll(c, &dirs)
		loghttp.E(w, r, err, false)
		err = datastore.DeleteMulti(c, keys)
		loghttp.E(w, r, err, false)
		loghttp.Pf(w, r, "%v directories deleted", len(keys))
	}

}
func demoSaveRetrieve(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	fmt.Fprint(w, tplx.Head)

	nestedOrRooted = !nestedOrRooted

	rt := <-util.Counter
	rts := fmt.Sprintf("mount%03v", rt)
	fs := NewFs(rts, appengine.NewContext(r), nestedOrRooted)
	loghttp.Pf(w, r, "created fs %v<br>\n", rts)

	fc1 := func(p []string) {
		path := pth.Join(p...)
		path = cleanseLeadingSlash(path)

		dir, err := fs.SaveDirByPath(path)
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child created %v - %v <br>", dir.Name, dir.Key)
	}

	loghttp.Pf(w, r, "-----------------<br>\n")

	fc1([]string{"ch1"})
	fc1([]string{"ch1", "ch2"})
	fc1([]string{"ch1", "ch2", "ch3"})
	fc1([]string{"ch1A"})

	loghttp.Pf(w, r, "-----------------<br>\n")

	// retrieval
	fc2 := func(p []string) {
		path := pth.Join(p...)
		path = cleanseLeadingSlash(path)

		loghttp.Pf(w, r, "searching  %v<br>", path)
		f, err := fs.GetDirByPath(path)
		if err != nil {
			loghttp.Pf(w, r, "  nothing retrieved - err %v<br>", err)
		} else {
			loghttp.Pf(w, r, "  child retrieved %v, %v<br>", f.Name, f.dir)
		}
	}
	fc2([]string{"ch1"})
	fc2([]string{"ch1", "ch2"})
	fc2([]string{"ch1", "cha2"})
	fc2([]string{"ch1", "ch2", "ch3"})
	fc2([]string{"fsd,mount000", "fsd,ch1", "ch2", "ch3"})
	fc2([]string{"ch1A"})

	loghttp.Pf(w, r, "-----------------<br>\n")

	fc3 := func(path string) {
		loghttp.Pf(w, r, "searching  %v<br>", path)
		f, err := fs.GetDirByPathQuery(path)
		if err != nil {
			loghttp.Pf(w, r, "  nothing retrieved - err %v<br>", err)
		} else {
			loghttp.Pf(w, r, "  child retrieved %v, %v<br>", f.Name, f.dir)
		}
	}
	fc3(spf(`/fsd,%v/fsd,ch1/fsd,ch2/fsd,ch3`, rts))
	fc3(spf(`/fsd,%v/fsd,ch1/ch2/ch3`, rts))

	loghttp.Pf(w, r, "-----------------<br>\n")

	fc4 := func(name, content string) {
		f := AeFile{}
		dir, base := pth.Split(name)
		f.name = base
		f.data = []byte(content)

		err := fs.SaveFile(&f, dir)
		logif.E(err)
	}

	fc4("ch1/ch2/file1", "content 1")
	fc4("ch1/ch2/file2", "content 2")
	fc4("ch1/ch2/file3", "another content")

	fc4("file4", "root content")

	//____________________

	{
		files, err := fs.GetFiles("ch1/ch2")
		logif.E(err)
		for k, v := range files {
			loghttp.Pf(w, r, "%v  -  %v %s<br>\n", k, v.Name, v.data)
		}
	}

	{
		files, err := fs.GetFiles("")
		logif.E(err)
		for k, v := range files {
			loghttp.Pf(w, r, "%v  -  %v %s<br>\n", k, v.Name, v.data)
		}
	}

	fmt.Fprint(w, tplx.Foot)
	//
}

func init() {
	http.HandleFunc("/fs/vfs-gae-demo", loghttp.Adapter(demoSaveRetrieve))
	http.HandleFunc("/fs/delete-all", loghttp.Adapter(deleteAll))
}
