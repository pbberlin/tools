package filesys

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/util"
)

func demoSaveRetrieve(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	NestedOrRooted = !NestedOrRooted

	rt := <-util.Counter
	rts := fmt.Sprintf("mount%03v", rt)
	fs := NewFileSys(w, r, rts)
	loghttp.Pf(w, r, "created fs %v<br>\n", rts)

	fc1 := func(p []string) {
		path := filepath.Join(p...)
		path = filepath.ToSlash(path)
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
		path := filepath.Join(p...)
		path = filepath.ToSlash(path)

		loghttp.Pf(w, r, "searching  %v<br>", path)
		f, err := fs.GetDirByPath(path)
		if err != nil {
			loghttp.Pf(w, r, "  nothing retrieved - err %v<br>", err)
		} else {
			loghttp.Pf(w, r, "  child retrieved %v, %v<br>", f.Name, f.Dir)
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
			loghttp.Pf(w, r, "  child retrieved %v, %v<br>", f.Name, f.Dir)
		}
	}
	fc3(spf(`/fsd,%v/fsd,ch1/fsd,ch2/fsd,ch3`, rts))
	fc3(spf(`/fsd,%v/fsd,ch1/ch2/ch3`, rts))

	loghttp.Pf(w, r, "-----------------<br>\n")

	f := File{}
	f.Name = "file1"
	f.Content = []byte("file content")

	err := fs.SaveFile(&f, "ch1/ch2/ch3")
	logif.E(err)

	f.Name = "file2"
	err = fs.SaveFile(&f, "ch1/ch2/ch3")
	logif.E(err)

	files, err := fs.GetFiles("ch1/ch2/ch3")
	logif.E(err)

	for k, v := range files {
		loghttp.Pf(w, r, "%v  -  %v", k, v)
	}
	//
}

func init() {
	http.HandleFunc("/fs/demo-fullpath", loghttp.Adapter(demoSaveRetrieve))
}
