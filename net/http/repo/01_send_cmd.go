package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/pbberlin/tools/appengine/instance_mgt"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
)

// Submit test commands internally, without http request.
func staticFetchDirect(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	FetchHTML(w, r, testCommands)
}

// Submit test commands by http posting them.
func staticFetchViaPosting2Receiver(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "JSON Post"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>")
	defer wpf(w, "</pre>")

	b, err := Post2Receiver(r, testCommands)

	lge(err)
	lg("msg from Post2Receiver:")
	lg(b.String())

}

// Post2Receiver takes commands and http posts them to
// the command receiver
func Post2Receiver(r *http.Request, commands []FetchCommand) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)

	if commands == nil || len(commands) == 0 {
		return b, fmt.Errorf("Slice of commands nil or empty %v", commands)
	}

	ii := instance_mgt.Get(r)
	fullURL := fmt.Sprintf("https://%s%s", ii.PureHostname, uriFetchCommandReceiver)
	wpf(b, "sending to URL:    %v\n", fullURL)

	bcommands, err := json.MarshalIndent(commands, "", "\t")
	if err != nil {
		wpf(b, "marshalling to []byte failed\n")
		return b, err
	}

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(bcommands))
	if err != nil {
		wpf(b, "creation of POST request failed\n")
		return b, err
	}
	req.Header.Set("X-Custom-Header-Counter", "nocounter")
	req.Header.Set("Content-Type", "application/json")

	bts, reqUrl, err := fetch.UrlGetter(r, fetch.Options{Req: req})
	_, _ = bts, reqUrl
	if err != nil {
		wpf(b, "Sending the POST request failed\n")
		return b, err
	}

	wpf(b, "effective req url: %v\n", reqUrl)
	wpf(b, "response body:\n")
	wpf(b, "%s\n", html.EscapeString(string(bts)))

	return b, nil
}
