package fetch_rss

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/osutilpb"
	"github.com/pbberlin/tools/stringspb"
)

func path2DirTree(w http.ResponseWriter, r *http.Request, treeX *DirTree, hrefs []FullArticle, domain string) {

	lg, lge := loghttp.Logger(w, r)
	_ = lg

	if treeX == nil {
		treeX = &DirTree{Name: "root1", Dirs: map[string]DirTree{}, LastFound: time.Now()}
	}
	var trLp *DirTree
	trLp = treeX

	pfx1 := "http://" + domain
	pfx2 := "https://" + domain

	for _, art := range hrefs {
		href := art.Url
		href = strings.TrimPrefix(href, pfx1)
		href = strings.TrimPrefix(href, pfx2)
		if strings.HasPrefix(href, "/") { // ignore other domains
			parsed, err := url.Parse(href)
			lge(err)
			href = parsed.Path
			// lg("href is %v", href)
			trLp = treeX
			// lg("trLp is %v", trLp.String())
			dir, remainder := "", href
			for {

				dir, remainder = osutilpb.PathDirReverse(remainder)
				trLp.Name = dir
				trLp.LastFound = art.Mod

				if _, ok := trLp.Dirs[dir]; !ok {
					trLp.Dirs[dir] = DirTree{Name: dir, Dirs: map[string]DirTree{}}
				}

				addressTaker := trLp.Dirs[dir]
				trLp = &addressTaker
				// Since we "cannot assign" to map struct directly:
				// trLp.Dirs[dir].LastFound = art.Mod   // fails
				// trLp.LastFound = art.Mod

				if remainder == "" {
					break
				}
			}

		}
	}

}

func rssDoc2DirTree(w http.ResponseWriter, r *http.Request, treeX *DirTree, rssDoc RSS, domain string) {

	lg, lge := loghttp.Logger(w, r)
	_ = lg

	if treeX == nil {
		treeX = &DirTree{Name: "root1", Dirs: map[string]DirTree{}, LastFound: time.Now()}
	}

	articleList := []FullArticle{}

	for _, lpItem := range rssDoc.Items.ItemList {

		t, err := time.Parse(time.RFC1123Z, lpItem.Published)
		//     := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		lge(err)

		articleList = append(articleList, FullArticle{Url: lpItem.Link, Mod: t})

		path2DirTree(w, r, treeX, articleList, domain)

	}

}

//
//
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
	lg("RSS resp size %5.2vkB, saved to %v", len(bdmp)/1024, rssUrlObj.Host+"/outp_rss.xml")

	return
}
