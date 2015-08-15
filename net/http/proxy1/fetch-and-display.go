// Package proxy1 forwards html pages and reduces their size.
package proxy1

import (
	"fmt"
	"net/http"
	"strings"

	_ "html"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/domclean1"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/paths"
	"github.com/pbberlin/tools/net/http/tplx"
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
      	  <input name="url"           size="80"  value="{{.val}}"><br/>
      	<span class='ib' style='width:140px'>Put into pre tags </span>
      	  <input name="renderInPre"    size="4"    value='' ><br/>
      	<span class='ib' style='width:140px'> </span>
      
      	<input type="submit" value="Fetch" accesskey='f'></div>
    </form>

`

// var host, port string

func handleFetchURL(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	// on live server => always use https
	if r.URL.Scheme != "https" && !util_appengine.IsLocalEnviron() {
		r.URL.Scheme = "https"
		r.URL.Host = r.Host
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
	if r.PostFormValue("url") != "" {
		urlAs += "url posted "
		rURL = r.PostFormValue("url")
	}

	if r.FormValue("url") != "" {
		if rURL == "" {
			urlAs += "url getted "
			rURL = r.FormValue("url")
		}
	}
	// lg("received %v:  %q", urlAs, rURL)

	if len(rURL) == 0 {

		tplAdder, tplExec := tplx.FuncTplBuilder(w, r)
		tplAdder("n_html_title", "Fetch some http data", nil)

		m := map[string]string{
			"protocol": "https",
			"host":     r.Host, // not  fetch.HostFromReq(r)
			"path":     paths.FetchUrl,
			"val":      "google.com",
		}
		if util_appengine.IsLocalEnviron() {
			m["protocol"] = "http"
		}
		tplAdder("n_cont_0", c_formFetchUrl, m)
		tplExec(w, r)

	} else {

		w.Header().Set("Content-type", "text/html; charset=utf-8")
		// w.Header().Set("Content-type", "text/html; charset=latin-1")

		bts, u, err := fetch.UrlGetter(r, fetch.Options{URL: rURL})
		lge(err)

		cntnt := string(bts)
		cntnt = insertNewlines.Replace(cntnt)
		cntnt = undouble.Replace(cntnt)

		lg("clean %v - %v - %v - %v", u, rURL, u.Host, fetch.HostFromUrl(u))
		cntnt = domclean1.ModifyHTML(r, u, cntnt)
		fmt.Fprintf(w, cntnt)

	}

}

func init() {
	http.HandleFunc(paths.FetchUrl, loghttp.Adapter(handleFetchURL))
}
