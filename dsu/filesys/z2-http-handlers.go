package filesys

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/util"
)

func demoSaveRetrieve(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	rt := <-util.Counter
	rts := fmt.Sprintf("mount%03v", rt)
	fs := NewFileSys(w, r, rts)
	loghttp.Pf(w, r, "created fs %v<br>\n", rts)

	f := func(p []string) {
		fso1, err := fs.saveDirByPath(filepath.Join(p...))
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child created %v - %v <br>\n", fso1.Name, fso1.Key)
	}

	f([]string{"ch1"})
	f([]string{"ch1", "ch2"})
	f([]string{"ch1", "ch2", "ch3"})
	f([]string{"ch1A"})

	//
	path := filepath.Join("ch1", "ch2", "ch3")
	loghttp.Pf(w, r, "<br>\nget by path for  %q<br>\n", path)

	{
		f, err := fs.getDirByPath(path)
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child retrieved %v, %v<br>\n", f.Name, f.Dir)
	}

	{
		f, err := fs.getDirByPath(path + "something")
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child retrieved %v, %v<br>\n", f.Name, f.Dir)
	}
}

func init() {
	http.HandleFunc("/fs/demo-fullpath", loghttp.Adapter(demoSaveRetrieve))
}
