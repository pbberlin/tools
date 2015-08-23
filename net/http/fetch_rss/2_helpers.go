package fetch_rss

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/runtimepb"
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

func condenseTrailingDir(uri string, n int) (ret string) {

	switch n {
	case 0:
		return uri
	case 1:
		return uri
	case 2:
		base1 := path.Base(uri)
		rdir1 := path.Dir(uri) // rightest Dir

		base2 := path.Base(rdir1)

		rdir2 := path.Dir(rdir1)

		ret = path.Join(rdir2, base2+"-"+base1)
	default:
		runtimepb.StackTrace(4)
		log.Fatalf("not implemented n > 2 (%v)", n)
	}

	return

}
