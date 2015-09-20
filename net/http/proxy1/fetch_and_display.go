// Package proxy1 forwards html pages and reduces their size.
package proxy1

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"

	_ "html"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/domclean2"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/net/http/tplx"
	"golang.org/x/net/html"
)

var insertNewlines = strings.NewReplacer(
	"<head", "\n<head",
	"</head>", "</head>\n",
	"<meta", "\n<meta",
	"</script>", "</script>\n",
	"</style>", "</style>\n",
	"</div>", "</div>\n",
	"<style", "\n<style",
	"<script", "\n<script")

var replTabsNewline = strings.NewReplacer("\n", " ", "\t", " ")
var undouble = strings.NewReplacer("\n\n\n", "\n", "\n\n", "\n")

const c_formFetchUrl = `

	<style> .ib { display:inline-block; }</style>
    <form action="{{.protocol}}://{{.host}}{{.path}}" method="post" >
      <div style='margin:8px;'>
      	<span class='ib' style='width:140px'>URL </span>
      	  <input name="{{.name}}"           size="80"  value="{{.val}}"><br/>
      	<span class='ib' style='width:140px'>Put into pre tags </span>
      	  <input name="renderInPre"    size="4"    value='' ><br/>
      	<span class='ib' style='width:140px'> </span>
      
      	<input type="submit" value="Fetch" accesskey='f'></div>
    </form>

`

// handleFetchURL either displays a form for requesting an url
// or it returns the URL´s contents.
func handleFetchURL(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	// on live server => always use https
	if r.URL.Scheme != "https" && !util_appengine.IsLocalEnviron() {
		r.URL.Scheme = "https"
		r.URL.Host = r.Host
		lg("lo - redirect %v", r.URL.String())
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}

	/*
		To distinguish between posted and getted value,
		we check the "post-only" slice of values first.
		If nothing's there, but FormValue *has* a value,
		then it was "getted", otherwise "posted"
	*/
	rURL := ""
	urlAs := ""
	err := r.ParseForm()
	logif.E(err)
	if r.PostFormValue(routes.URLParamKey) != "" {
		urlAs += "url posted "
		rURL = r.PostFormValue(routes.URLParamKey)
	}

	if r.FormValue(routes.URLParamKey) != "" {
		if rURL == "" {
			urlAs += "url getted "
			rURL = r.FormValue(routes.URLParamKey)
		}
	}
	// lg("received %v:  %q", urlAs, rURL)

	if len(rURL) == 0 {

		tplAdder, tplExec := tplx.FuncTplBuilder(w, r)
		tplAdder("n_html_title", "Fetch some http data", nil)

		m := map[string]string{
			"protocol": "https",
			"host":     r.Host, // not  fetch.HostFromReq(r)
			"path":     routes.ProxifyURI,
			"name":     routes.URLParamKey,
			"val":      "google.com",
		}
		if util_appengine.IsLocalEnviron() {
			m["protocol"] = "http"
		}
		tplAdder("n_cont_0", c_formFetchUrl, m)
		tplExec(w, r)

	} else {

		r.Header.Set("X-Custom-Header-Counter", "nocounter")

		bts, inf, err := fetch.UrlGetter(r, fetch.Options{URL: rURL})
		lge(err)

		tp := mime.TypeByExtension(path.Ext(inf.URL.Path))
		if false {
			ext := path.Ext(rURL)
			ext = strings.ToLower(ext)
			tp = mime.TypeByExtension(ext)
		}
		w.Header().Set("Content-Type", tp)
		// w.Header().Set("Content-type", "text/html; charset=latin-1")

		if r.FormValue("dbg") != "" {
			w.Header().Set("Content-type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "%s<br>\n  %s<br>\n %v", inf.URL.Path, tp, inf.URL.String())
			return
		}

		opts := domclean2.CleaningOptions{Proxify: true}
		opts.Beautify = true // "<a> Linktext without trailing space"
		opts.RemoteHost = fetch.HostFromStringUrl(rURL)

		// opts.ProxyHost = routes.AppHost01
		opts.ProxyHost = fetch.HostFromReq(r)
		if !util_appengine.IsLocalEnviron() {
			opts.ProxyHost = fetch.HostFromReq(r)
		}

		doc, err := domclean2.DomClean(bts, opts)

		var bufRend bytes.Buffer
		err = html.Render(&bufRend, doc)
		lge(err)
		w.Write(bufRend.Bytes())

	}

}

func init() {
	http.HandleFunc(routes.ProxifyURI, loghttp.Adapter(handleFetchURL))
}
