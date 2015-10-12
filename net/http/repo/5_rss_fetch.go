package repo

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/stringspb"
)

func rssDoc2DirTree(w http.ResponseWriter, r *http.Request, treeX *DirTree, rssDoc RSS, domain string) {

	lg, lge := loghttp.Logger(w, r)
	_ = lg

	if treeX == nil {
		treeX = &DirTree{Name: "root1", Dirs: map[string]DirTree{}, LastFound: time.Now().Truncate(time.Minute)}
	}

	articleList := []FullArticle{}

	for _, lpItem := range rssDoc.Items.ItemList {

		t, err := time.Parse(time.RFC1123Z, lpItem.Published)
		//     := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		lge(err)

		articleList = append(articleList, FullArticle{Url: lpItem.Link, Mod: t})

	}

	lg1, _ := loghttp.BuffLoggerUniversal(w, r)

	path2DirTree(lg1, treeX, articleList, domain, true)

}

//
//
// Fetches the RSS.xml file.
func rssXMLFile(w http.ResponseWriter, r *http.Request, fs fsi.FileSystem, rssUrl string) (RSS, *url.URL) {

	lg, lge := loghttp.Logger(w, r)

	bts, respInf, err := fetch.UrlGetter(r, fetch.Options{URL: rssUrl})
	lge(err)

	bts = bytes.Replace(bts, []byte("content:encoded>"), []byte("content-encoded>S"), -1) // hack

	rssDoc := RSS{}
	err = xml.Unmarshal(bts, &rssDoc)
	lge(err)

	// save it
	bdmp := stringspb.IndentedDumpBytes(rssDoc)
	err = fs.MkdirAll(path.Join(docRoot, respInf.URL.Host), 0755)
	lge(err)
	err = fs.WriteFile(path.Join(docRoot, respInf.URL.Host, "outp_rss.xml"), bdmp, 0755)
	lge(err)
	lg("RSS resp size %5.2vkB, saved to %v", len(bdmp)/1024, respInf.URL.Host+"/outp_rss.xml")

	return rssDoc, respInf.URL
}
