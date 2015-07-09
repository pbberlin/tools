package aefs

import (
	"fmt"
	"net/http"
	"os"
	pth "path"

	"appengine"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/util"
)

func deleteAll(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	fs := NewAeFs("anyRoot", AeContext(appengine.NewContext(r)))
	err, msg := fs.DeleteAll()

	loghttp.E(w, r, err, true)

	wpf(w, "<pre>%v</pre>", msg)

}
func demoSaveRetrieve(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	fmt.Fprint(w, tplx.Head)

	rt := <-util.Counter
	rts := fmt.Sprintf("mnt%02v", rt)
	fs := NewAeFs(rts, AeContext(appengine.NewContext(r)))
	loghttp.Pf(w, r, "created fs %v<br>", rts)

	fc1 := func(p []string) {
		path := pth.Join(p...)
		path = cleanseLeadingSlash(path)

		dir, err := fs.saveDirByPath(path)
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child created %v - %v <br>", dir.Name(), dir.Key)
	}

	loghttp.Pf(w, r, "--------create-dirs---------<br>")

	fc1([]string{"ch1"})
	fc1([]string{"ch1", "ch2"})
	fc1([]string{"ch1", "ch2a"})
	fc1([]string{"ch1", "ch2", "ch3"})
	fc1([]string{"ch1", "ch2", "ch3", "ch4"})
	fc1([]string{"ch1A"})
	fc1([]string{"ch1B"})

	loghttp.Pf(w, r, "--------retrieve-dirs---------<br>")

	// retrieval
	fc2 := func(p []string) {
		path := pth.Join(p...)
		path = fs.RootDir() + path
		loghttp.Pf(w, r, "%v - searching...<br>", path)
		f, err := fs.dirByPath(path)
		if err != nil {
			loghttp.Pf(w, r, "  nothing retrieved - err %v<br>", err)
		} else {
			loghttp.Pf(w, r, " - %v retrieved under %v <br>", f.Name(), f.Dir)
		}
	}
	fc2([]string{"ch1"})
	fc2([]string{"ch1", "ch2"})
	fc2([]string{"ch1", "cha2"})
	fc2([]string{"ch1", "ch2", "ch3"})
	fc2([]string{"fsd,mount000", "fsd,ch1", "ch2", "ch3"})
	fc2([]string{"ch1A"})

	loghttp.Pf(w, r, "--------retrieve by query---------<br>")

	fc3 := func(path string) {
		loghttp.Pf(w, r, "searching  %v<br>", path)
		children, err := fs.subdirsByPath(path, true)
		if err != nil {
			loghttp.Pf(w, r, "  nothing retrieved - err %v<br>", err)
		} else {
			for k, v := range children {
				loghttp.Pf(w, r, "  %v child retrieved %v, %v<br>", k, v.Name(), v.Dir)
			}
		}
	}
	fc3(spf(`/fsd,%v/fsd,ch1/fsd,ch2/fsd,ch3`, rts))
	fc3(spf(`/fsd,%v/fsd,ch1/ch2/ch3`, rts))

	loghttp.Pf(w, r, "-------create and save some files----<br>")

	fc4 := func(name, content string) {
		err := fs.WriteFile(name, []byte(content), os.ModePerm)
		logif.E(err)
	}

	fc4("ch1/ch2/file1", "content 1")
	fc4("ch1/ch2/file2", "content 2")
	fc4("ch1/ch2/ch3/file3", "another content")
	fc4(fs.RootDir()+"file4", "root content")

	loghttp.Pf(w, r, "-------retrieve files again----<br>")

	fc5 := func(path string) {
		files, err := fs.filesByPath(fs.RootDir() + path)
		logif.E(err)
		loghttp.Pf(w, r, " srch %v  <br>", fs.RootDir()+path)
		for k, v := range files {
			loghttp.Pf(w, r, "%v  -  %v %s<br>", k, v.Name(), v.Data)
		}
	}

	fc5("ch1/ch2")
	fc5("ch1/ch2/ch3")
	fc5("")

}

func init() {
	http.HandleFunc("/fs/aefs/vfs-gae-demo", loghttp.Adapter(demoSaveRetrieve))
	http.HandleFunc("/fs/aefs/delete-all", loghttp.Adapter(deleteAll))
}
