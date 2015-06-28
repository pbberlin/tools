package nested

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/util"
)

func demoFullpath(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	rt := <-util.Counter
	rts := fmt.Sprintf("mount%03v", rt)
	fs := NewNestedFileSys(w, r, rts)

	loghttp.Pf(w, r, "created fs %v<br>\n", rts)

	//
	fso1 := fs.saveDirByPath(filepath.Join("ch1"))
	loghttp.Pf(w, r, "child created %v - %v <br>\n", fso1.Name, fso1.Key)

	fso2 := fs.saveDirByPath(filepath.Join("ch1", "ch2"))
	loghttp.Pf(w, r, "child created %v - %v <br>\n", fso2.Name, fso2.Key)

	fso3 := fs.saveDirByPath(filepath.Join("ch1", "ch2", "ch3"))
	loghttp.Pf(w, r, "child created %v - %v <br>\n", fso3.Name, fso3.Key)

	//
	path := filepath.Join("ch1", "ch2", "ch3")
	loghttp.Pf(w, r, "<br>\nget by path for  %q<br>\n", path)

	{
		f, err := fs.dirByPath(path)
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child retrieved %v, %v<br>\n", f.Name, f.Dir)
	}

	{
		f, err := fs.dirByPath(path + "something")
		loghttp.E(w, r, err, true)
		loghttp.Pf(w, r, "child retrieved %v, %v<br>\n", f.Name, f.Dir)
	}
}

func init() {
	http.HandleFunc("/fs/demo-fullpath", loghttp.Adapter(demoFullpath))
}
