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

func path2DirTree(w http.ResponseWriter, r *http.Request, treeX *DirTree, articles []FullArticle, domain string) {

	lg, lge := loghttp.Logger(w, r)
	_ = lg

	if treeX == nil {
		treeX = &DirTree{Name: "root1", Dirs: map[string]DirTree{}, LastFound: time.Now()}
	}
	var trLp *DirTree
	trLp = treeX

	pfx1 := "http://" + domain
	pfx2 := "https://" + domain

	for _, art := range articles {
		href := art.Url
		if art.Mod.IsZero() {
			art.Mod = time.Now()
		}
		href = strings.TrimPrefix(href, pfx1)
		href = strings.TrimPrefix(href, pfx2)
		if strings.HasPrefix(href, "/") { // ignore other domains
			parsed, err := url.Parse(href)
			lge(err)
			href = parsed.Path
			// lg("%v", href)
			trLp = treeX
			// lg("trLp is %v", trLp.String())
			dir, remainder := "", href
			lvl := 0
			for {

				dir, remainder = osutilpb.PathDirReverse(remainder)

				if dir == "/" && remainder == "" {
					// skip root
					break
				}

				if lvl > 0 {
					trLp.Name = dir // lvl==0 => root
				}
				trLp.LastFound = art.Mod

				// lg("   %v, %v", dir, remainder)

				if _, ok := trLp.Dirs[dir]; !ok {
					trLp.Dirs[dir] = DirTree{Name: dir, Dirs: map[string]DirTree{}}
				}

				// We "cannot assign" to map struct directly:
				// trLp.Dirs[dir].LastFound = art.Mod   // fails
				addressable := trLp.Dirs[dir]
				addressable.LastFound = art.Mod
				trLp.Dirs[dir] = addressable
				trLp = &addressable

				if remainder == "" {
					// lg("break\n")
					break
				}

				lvl++
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

	}

	path2DirTree(w, r, treeX, articleList, domain)

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
