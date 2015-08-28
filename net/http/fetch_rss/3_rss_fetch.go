package fetch_rss

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/url"
	"path"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/stringspb"
)

// Fetches the RSS.xml file.
func rssXMLFile(w http.ResponseWriter, r *http.Request, fs fsi.FileSystem, rssUrl string) (rssDoc RSS, rssUrlObj *url.URL) {

	lg, lge := loghttp.Logger(w, r)

	var bts []byte
	var err error

	bts, rssUrlObj, err = fetch.UrlGetter(r, fetch.Options{URL: rssUrl})
	lge(err)

	bts = bytes.Replace(bts, []byte("content:encoded>"), []byte("content-encoded>S"), -1) // hack

	err = xml.Unmarshal(bts, &rssDoc)
	lge(err)

	// save it
	bdmp := stringspb.IndentedDumpBytes(rssDoc)
	err = fs.MkdirAll(path.Join(docRoot, rssUrlObj.Host), 0755)
	lge(err)
	err = fs.WriteFile(path.Join(docRoot, rssUrlObj.Host, "outp_rss.xml"), bdmp, 0755)
	lge(err)
	lg("RSS resp size, outp_rss.xml, : %v", len(bdmp))

	return
}
