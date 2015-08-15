package fetch_rss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"appengine"

	"github.com/pbberlin/tools/appengine/instance_mgt"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/stringspb"
)

type FetchCommand struct {
	Host                 string   // www.handelsblatt.com,
	RssXMLURI            string   // /contentexport/feed/schlagzeilen,
	SearchPrefixs        []string // /politik/international/aa/bb,
	CondenseTrailingDirs int      // The last one or two directories might be article titles or ids
	DepthTolerance       int      // The last one or two directories might be article titles or ids
}

var fcs = []FetchCommand{
	FetchCommand{
		Host:                 "www.handelsblatt.com",
		RssXMLURI:            "/contentexport/feed/schlagzeilen",
		SearchPrefixs:        []string{"/politik/international/aa/bb", "/politik/deutschland/aa/bb"},
		CondenseTrailingDirs: 2,
		DepthTolerance:       1,
	},
	FetchCommand{
		"www.economist.com",
		"/sections/international/rss.xml",
		[]string{"/news/international"},
		1,
		2,
	},
}

func fetchCommandReceiver(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	var fc []FetchCommand

	// Decode from io.Reader seems better:
	// http://stackoverflow.com/questions/21197239/decoding-json-in-golang-using-json-unmarshal-vs-json-newdecoder-decode
	//
	// We use Unmarshal here, because we want to inspect the bytes of body.
	var Unmarshal_versus_Decode = true

	if Unmarshal_versus_Decode {

		body, err := ioutil.ReadAll(r.Body) // no response write until here !
		lge(err)

		lg("body is %s", body)

		err = json.Unmarshal(body, &fc)
		lge(err)
		if err == nil {
			lg("command is: %s", *stringspb.IndentedDump(fc))
		}

	} else {

		//
		dec := json.NewDecoder(r.Body)
		for {
			if err := dec.Decode(&fc); err == io.EOF {
				break
			} else if err != nil {
				lge(err)
			}
			lg("command loop is: %s", *stringspb.IndentedDump(fc))
		}

	}

}

func fetchCommandSender(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	//  curl -X POST -d "{\"tes\": \"that\"}" localhost:8082/test

	lg, lge := loghttp.Logger(w, r)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "JSON Post"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>")
	defer wpf(w, "</pre>")

	ii := instance_mgt.Get(appengine.NewContext(r))
	fullURL := fmt.Sprintf("https://%s%s", ii.PureHostname, uriFetchCommandReceiver)
	lg("URL: %v", fullURL)

	command, err := json.MarshalIndent(fcs, "", "\t")
	lge(err)

	// This should cause an error upon unmarshal into []fetch.Command:
	// command = []byte(`{"title":"Buy cheese and bread for breakfast."}`)

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(command))
	lge(err)

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	bts, reqUrl, err := fetch.UrlGetter(r, fetch.Options{Req: req})
	lge(err)
	lg("response from: %v", reqUrl)

	lg("response Body: %s", bts)

}
