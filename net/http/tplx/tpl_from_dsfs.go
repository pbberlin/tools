package tplx

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"

	"appengine"
)

const TplPrefix = "/mnt02"
const pathToTmpl = "/page/programmatic_content/index.html"

var fs1 = memfs.New(
	memfs.Ident(TplPrefix[1:]), // a closured variable in init() did not survive map-pointer reallocation
)

// Implements fsi.File
type SyncedMap struct {
	sync.Mutex
	mp map[string]string
}

var mp = SyncedMap{mp: map[string]string{}}

func BootstrapTemplate(w http.ResponseWriter, r *http.Request) string {

	if _, ok := mp.mp[pathToTmpl]; ok {
		return mp.mp[pathToTmpl]
	}

	lg, _ := loghttp.BuffLoggerUniversal(w, r)

	fs2 := dsfs.New(
		dsfs.MountName(TplPrefix[1:]),
		dsfs.AeContext(appengine.NewContext(r)),
	)
	fs1.SetOption(
		memfs.ShadowFS(fs2),
	)

	bts, err := fs1.ReadFile(pathToTmpl)
	lg(err)
	if err != nil {
		return err.Error()
	}

	bts = bytes.Replace(bts, []byte("[REPLACE_TITLE]"), []byte("{{ .HtmlTitle }}"), -1)
	bts = bytes.Replace(bts, []byte("[REPLACE_DESC]"), []byte("{{ .HtmlDescription }}"), -1)
	bts = bytes.Replace(bts, []byte("<p>[REPLACE_CONTENT]</p>"), []byte("{{ .HtmlContent }}"), -1)
	bts = bytes.Replace(bts, []byte("[REPLACE_CONTENT]"), []byte("{{ .HtmlContent }}"), -1)

	mp.Lock()

	mp.mp[pathToTmpl] = string(bts)
	mp.Unlock()

	return mp.mp[pathToTmpl]

}
