package pbfetch

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"

	_ "html"
	"io/ioutil"
	"strings"

	"github.com/pbberlin/tools/tpl_html"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_appengine"
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
      	 <input name="elipseOutput"  size="4"    value='' ><br/>
      
      	<input type="submit" value="Fetch" accesskey='f'></div>
    </form>

`

const FetchURL = "fetch-url-x"

func handleFetchURL(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	if r.URL.Scheme != "https" && !util_appengine.IsLocalEnviron() {
		r.URL.Scheme = "https"
		r.URL.Host = r.Host
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}

	rURL := ""
	urlAsPost := ""

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

	renderInPre := false
	if len(r.FormValue("renderInPre")) > 0 {
		renderInPre = true
	}

	elipseOutput := r.FormValue("elipseOutput")

	var msg, cntnt string

	if len(rURL) == 0 {

		tplAdder, tplExec := tpl_html.FuncTplBuilder(w, r)
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

		if !strings.HasPrefix(rURL, "http://") && !strings.HasPrefix(rURL, "https://") {
			rURL = "https://" + rURL
		}

		u, err := url.Parse(rURL)
		if err != nil {
			panic(err)
		}
		host, port, err = net.SplitHostPort(u.Host)
		if err != nil {
			host = u.Host
		}
		log.Println("host and port: ", host, port, "of", rURL, "standalone:", u.Host)

		client := urlfetch.Client(c)
		resp, err := client.Get(rURL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(w, "HTTP GET returned status %v<br>\n\n", resp.Status)
			return
		}

		defer resp.Body.Close()
		byteContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Errorf("%s", err)
			fmt.Fprintf(w, "Error %v<br>\n\n", err.Error())
			return
		} else {
			msg = fmt.Sprintf("%v bytes read<br>", len(byteContent))
			if urlAsPost != "" {
				msg += fmt.Sprintf("get overwritten by post: %v <br>", urlAsPost)
			}
			cntnt = string(byteContent)
		}

		cntnt = insertNewlines.Replace(cntnt)
		cntnt = undouble.Replace(cntnt)

		if len(elipseOutput) > 0 {
			cutoff := util.Min(100, len(cntnt))
			fmt.Fprintf(w, "content is: <pre>"+cntnt[:cutoff]+" ... "+cntnt[len(cntnt)-cutoff:]+"</pre>")
		} else {
			if renderInPre {
				fmt.Fprintf(w, "content is: <pre>"+cntnt+"</pre>")
			} else {
				cntnt = ModifyHTML(r, cntnt)
				fmt.Fprintf(w, cntnt)
			}
		}

	}

	// cntnt = html.EscapeString(cntnt)

	fmt.Fprintf(w, " %s \n\n", msg)

}

func init() {
	http.HandleFunc("/"+FetchURL, handleFetchURL)
}
