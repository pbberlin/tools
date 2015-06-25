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

var host, port string

const c_formFetchURL = `

	<style> .ib { display:inline-block; }</style>
    <form action="{{.protocol}}://{{.host}}/{{.path}}" method="post" >
      <div style='margin:8px;'>
      	<span class='ib' style='width:140px'>URL </span>
      	  <input name="url"           size="80"  value="{{.val}}"><br/>
      	<span class='ib' style='width:140px'>Put into pre tags </span>
      	  <input name="renderInPre"    size="4"    value='' ><br/>
      	<span class='ib' style='width:140px'>Elipse Output</span>
      
      	<input type="submit" value="Fetch" accesskey='f'></div>
    </form>

`

const FetchURL = "fetch-url-x"

func handleFetchURL(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c, _ := util_appengine.SafeGaeCheck(r)
	_ = c

	if r.URL.Scheme != "https" && !util_appengine.IsLocalEnviron() {
		r.URL.Scheme = "https"
		r.URL.Host = r.Host
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}

	urlAsPost := ""
	rURL := ""

	/*
		To distinguish between posted and getted value,
		we check the "post-only" slice of values first.
		If nothing's there, but FormValue *has* a value,
		then it was "getted", otherwise "posted"
	*/
	if r.PostFormValue("url") != "" {
		urlAsPost = "url posted"
		rURL = r.PostFormValue("url")
	}

	if r.FormValue("url") != "" {
		if rURL == "" {
			urlAsPost = "url getted"
			rURL = r.FormValue("url")
		}
	}
	_ = urlAsPost

	renderInPre := false
	if len(r.FormValue("renderInPre")) > 0 {
		renderInPre = true
	}

	if len(rURL) == 0 {

		tplAdder, tplExec := tplx.FuncTplBuilder(w, r)
		tplAdder("n_html_title", "Fetch some http data", nil)

		m := map[string]string{"protocol": "https", "host": r.Host, "path": FetchURL, "val": "google.com"}
		if util_appengine.IsLocalEnviron() {
			m["protocol"] = "http"
		}
		tplAdder("n_cont_0", c_formFetchURL, m)
		tplExec(w, r)

	} else {

		w.Header().Set("Content-type", "text/html; charset=utf-8")
		// w.Header().Set("Content-type", "text/html; charset=latin-1")

		bts, err := fetch.UrlGetter(rURL, r, false)
		logif.E(err)

		cntnt := string(bts)
		cntnt = insertNewlines.Replace(cntnt)
		cntnt = undouble.Replace(cntnt)

		if renderInPre {
			fmt.Fprintf(w, "content is: <pre>"+cntnt+"</pre>")
		} else {
			cntnt = domclean1.ModifyHTML(r, cntnt)
			fmt.Fprintf(w, cntnt)
		}

	}

}

func init() {
	http.HandleFunc("/"+FetchURL, loghttp.Adapter(handleFetchURL))
}
