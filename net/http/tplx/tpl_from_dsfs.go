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
	memfs.Ident(TplPrefix[1:]), // a closured variable in init( ) did not survive map-pointer reallocation
)

// Implements fsi.File
type SyncedMap struct {
	sync.Mutex
	mp map[string]string
}

// We need this, since we have that additional step of putting the {{ .Variable }} strings into the hugo html
var mp = SyncedMap{mp: map[string]string{}}

func TemplateFromHugoReset(w http.ResponseWriter, r *http.Request, mapOnerous map[string]interface{}) {
	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	mp.Lock()
	if _, ok := mp.mp[pathToTmpl]; ok {
		mp = SyncedMap{mp: map[string]string{}}
	}
	mp.Unlock()
	w.Write([]byte("reset successful"))
}

func TemplateFromHugoPage(w http.ResponseWriter, r *http.Request) string {

	mp.Lock()
	if _, ok := mp.mp[pathToTmpl]; ok {
		mp.Unlock()
		return mp.mp[pathToTmpl]
	}
	mp.Unlock()

	lg, _ := loghttp.BuffLoggerUniversal(w, r)

	//
	fs2 := dsfs.New(
		dsfs.MountName(TplPrefix[1:]),
		dsfs.AeContext(appengine.NewContext(r)),
	)
	fs1.SetOption(
		memfs.ShadowFS(fs2),
	)

	bts, err := fs1.ReadFile(pathToTmpl)
	if err != nil {
		lg(err)
		bts = hugoTplFallback
	}

	bts = bytes.Replace(bts, []byte("[REPLACE_TITLE]"), []byte("{{ .HtmlTitle }}"), -1)
	bts = bytes.Replace(bts, []byte("[REPLACE_DESC]"), []byte("{{ .HtmlDescription }}"), -1)
	bts = bytes.Replace(bts, []byte("</head>"), []byte("{{ .HtmlHeaders }}\n</head>"), -1)
	bts = bytes.Replace(bts, []byte("<p>[REPLACE_CONTENT]</p>"), []byte("{{ .HtmlContent }}"), -1)
	bts = bytes.Replace(bts, []byte("[REPLACE_CONTENT]"), []byte("{{ .HtmlContent }}"), -1)
	bts = bytes.Replace(bts, []byte("[REPLACE_FOOTER]"), []byte("{{ .HtmlFooter }}"), -1)

	mp.Lock()
	mp.mp[pathToTmpl] = string(bts)
	mp.Unlock()

	return mp.mp[pathToTmpl]

}
