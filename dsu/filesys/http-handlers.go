package filesys

import (
	"fmt"
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/util"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf

func counterReset(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
}

func createFileSystem(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	rt := <-util.Counter
	rts := spf("mount%03v", rt)
	fs := NewFileSys(w, r, rts)
	w.Write([]byte("root  created " + rts + "<br>\n"))

	fso1 := fs.newFsoByParentKey("ch1", fs.RootDir.Key, true)
	w.Write([]byte("child created " + fso1.Name + "<br>\n"))

	fso2 := fs.newFsoByParentKey("ch2a", fso1.Key, true)
	w.Write([]byte("child created " + fso2.Name + "<br>\n"))

	fso2b := fs.newFsoByParentKey("ch2b", fso1.Key, true)
	w.Write([]byte("child created " + fso2b.Name + "<br>\n"))

	fso2c := fs.newFsoByParentKey("ch2c", fso1.Key, true)
	w.Write([]byte("child created " + fso2c.Name + "<br>\n"))

	fso3 := fs.newFsoByParentKey("ch3", fso2.Key, true)
	w.Write([]byte("child created " + fso3.Name + "<br>\n"))

	path := ""
	path = spf("/fso,%s/fso,ch1/fso,ch2a/fso,ch3", rts)

	fso4, err := fs.GetFsoByQuery(path)
	loghttp.E(w, r, err, true)
	loghttp.Pf(w, r, "child retrieved %v, %v<br>\n", fso4.Name, fso4.SKey)

}

func init() {
	http.HandleFunc("/fs/new", loghttp.Adapter(createFileSystem))
	http.HandleFunc("/fs/reset", loghttp.Adapter(counterReset))
}
