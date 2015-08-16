package proxy1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/domclean1"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/paths"
)

func formRedirector(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	var msg, cntnt, rURL string

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	// w.Header().Set("Content-type", "text/html; charset=latin-1")

	rURL = r.FormValue("redirect-to")
	lg("lo redirect to: %v", rURL)

	if len(r.PostForm) > 0 {
		// loghttp.Pf(w, r, "post unimplemented:<br> %#v <br>\n", r.PostForm)
		// return
		msg += fmt.Sprintf("post converted to get<br>")
	}

	rURL = fmt.Sprintf("%v?1=2&", rURL)
	for key, vals := range r.Form {
		if key == "redirect-to" {
			continue
		}
		val := vals[0]
		if util_appengine.IsLocalEnviron() {
			val = strings.Replace(val, " ", "%20", -1)
		}
		rURL = fmt.Sprintf("%v&%v=%v", rURL, key, val)
	}

	bts, u, err := fetch.UrlGetter(r, fetch.Options{URL: rURL})
	lge(err)

	cntnt = string(bts)

	cntnt = insertNewlines.Replace(cntnt)
	cntnt = undouble.Replace(cntnt)

	cntnt = domclean1.ModifyHTML(r, u, cntnt)

	fmt.Fprintf(w, "%s \n\n", cntnt)
	fmt.Fprintf(w, "%s \n\n", msg)

}

func init() {
	http.Handle(paths.FormRedirector, loghttp.Adapter(formRedirector))
}
