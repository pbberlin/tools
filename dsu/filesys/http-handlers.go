package filesys

import (
	"fmt"
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/util"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf

func createFileSystem(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	rt := <-util.Counter
	rts := spf("mount%03v", rt)
	fs := NewFileSys(w, r, rts)
	w.Write([]byte("root  created " + rts + "<br>\n"))

	fso1 := fs.newFso("ch1", fs.RootDir.Key, true)
	w.Write([]byte("child created " + fso1.Name + "<br>\n"))

	fso2 := fs.newFso("ch2a", fso1.Key, true)
	w.Write([]byte("child created " + fso2.Name + "<br>\n"))

	fso2b := fs.newFso("ch2b", fso1.Key, true)
	w.Write([]byte("child created " + fso2b.Name + "<br>\n"))

	fso3 := fs.newFso("ch3", fso2.Key, true)
	w.Write([]byte("child created " + fso3.Name + "<br>\n"))

	path := `/fso,mount000/fso,ch1/fso,ch2a`
	path = `/fso,mount000/fso,ch1/fso,ch2a/fso,ch3`

	fso4, err := fs.GetFso(path)
	w.Write([]byte("child retrieved " + fso4.Name + "<br>\n"))
	w.Write([]byte("child retrieved " + fso4.SKey + "<br>\n"))
	w.Write([]byte(" &nbsp; &nbsp; err " + spf("%v", err) + "<br>\n"))

}

func init() {
	http.HandleFunc("/fs/new", loghttp.Adapter(createFileSystem))
}
