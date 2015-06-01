package fetch

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"github.com/pbberlin/tools/util_appengine"
	"github.com/pbberlin/tools/util_err"

	"github.com/mjibson/appstats"
)

func formRedirector(c appengine.Context, w http.ResponseWriter, r *http.Request) {

	var msg, cntnt, rURL string

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	// w.Header().Set("Content-type", "text/html; charset=latin-1")

	rURL = r.FormValue("redirect-to")
	//util_err.LogAndShow(w, r, "url: %q <br>\n", rURL)

	u, err := url.Parse(rURL)
	util_err.Err_http(w, r, err, false)

	host, port, err = net.SplitHostPort(u.Host)
	util_err.Err_http(w, r, err, true)
	if err != nil {
		host = u.Host
	}
	//util_err.LogAndShow(w, r, "host and port: %q : %q of %q<br>\n", host, port, rURL)
	//util_err.LogAndShow(w, r, " &nbsp;  &nbsp;  &nbsp; standalone %q <br>\n", u.Host)

	client := urlfetch.Client(c)

	if len(r.PostForm) > 0 {
		// util_err.LogAndShow(w, r, "post unimplemented:<br> %#v <br>\n", r.PostForm)
		// return
		msg += fmt.Sprintf("post converted to get<br>")
	}

	rURL = fmt.Sprintf("%v?1=2&", rURL)
	for key, vals := range r.Form {
		if key == "redirect-to" {
			continue
		}
		val := vals[0]
		if !util_appengine.IsLocalEnviron() {
			val = strings.Replace(val, " ", "%20", -1)
		}
		rURL = fmt.Sprintf("%v&%v=%v", rURL, key, val)
	}

	resp, err := client.Get(rURL)
	util_err.Err_http(w, r, err, false)

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "HTTP GET returned status %v<br>\n\n%v<br>\n\n", resp.Status, rURL)
		return
	}

	defer resp.Body.Close()
	byteContent, err := ioutil.ReadAll(resp.Body)
	util_err.Err_http(w, r, err, false)
	if err != nil {
		return
	} else {
		msg += fmt.Sprintf("%v bytes read<br>", len(byteContent))
		cntnt = string(byteContent)
	}

	cntnt = insertNewlines.Replace(cntnt)
	cntnt = undouble.Replace(cntnt)

	cntnt = ModifyHTML(r, cntnt)

	fmt.Fprintf(w, "%s \n\n", cntnt)
	fmt.Fprintf(w, "%s \n\n", msg)

}

func init() {
	http.Handle("/blob2/form-redirector", appstats.NewHandler(formRedirector))
}
