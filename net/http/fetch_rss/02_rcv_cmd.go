package fetch_rss

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"appengine"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/stringspb"
)

// fetchCommandReceiver takes http post requests, extracts the JSON commands
// and submits them to FetchHTML
func fetchCommandReceiver(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	var fcs []FetchCommand

	// The type of resp.body  <io.Reader> lends itself to using decode.
	// http://stackoverflow.com/ - ... using-json-unmarshal-vs-json-newdecoder-decode
	//
	// Nevertheless, we use Unmarshal here, because we want to inspect the bytes of body.
	var Unmarshal_versus_Decode = true

	if Unmarshal_versus_Decode {

		body, err := ioutil.ReadAll(r.Body) // no response write until here !
		lge(err)

		if len(body) == 0 {
			lg("empty body")
			return
		}

		err = json.Unmarshal(body, &fcs)
		if err != nil {
			lge(err)
			lg("body is %s", body)
			return
		}

	} else {

		dec := json.NewDecoder(r.Body)
		for {
			if err := dec.Decode(&fcs); err == io.EOF {
				break
			} else if err != nil {
				lge(err)
				return
			}
			lg("command loop is: %s", stringspb.IndentedDump(fcs))
		}

	}

	FetchHTML(w, r, fcs)

}

// FetchHTML executes the fetch commands.
// It creates the configured filesystem and calls the fetcher.
func FetchHTML(w http.ResponseWriter, r *http.Request, fcs []FetchCommand) {

	lg, lge := loghttp.Logger(w, r)
	var err error

	fs := GetFS(appengine.NewContext(r))
	// fs = fsi.FileSystem(memMapFileSys)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Requesting files"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>")
	defer wpf(w, "</pre>")

	err = fs.WriteFile(path.Join(docRoot, "msg.html"), msg, 0644)
	lge(err)

	// err = fs.WriteFile(path.Join(docRoot, "index.html"), []byte("content of index.html"), 0644)
	// lge(err)

	err = fs.MkdirAll(path.Join(docRoot, "testDirX/testDirY"), 0755)
	lge(err)

	for _, config := range fcs {
		config = addDefaults(config)
		Fetch(w, r, fs, config)
	}

	lg("fetching complete")

}

func addDefaults(in FetchCommand) FetchCommand {

	var preset FetchCommand

	h := in.Host
	if exactPreset, ok := ConfigDefaults[h]; ok {
		preset = exactPreset
	} else {
		preset = ConfigDefaults["unspecified"]
	}

	in.DepthTolerance = preset.DepthTolerance
	in.CondenseTrailingDirs = preset.CondenseTrailingDirs
	if in.DesiredNumber == 0 {
		in.DesiredNumber = preset.DesiredNumber
	}

	if in.RssXMLURI == nil || len(in.RssXMLURI) == 0 {
		in.RssXMLURI = preset.RssXMLURI
	}

	return in
}
