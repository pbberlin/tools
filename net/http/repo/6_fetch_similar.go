package repo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"appengine"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/stringspb"
)

// FetchSimilar is an extended version of Fetch
// It is uses a DirTree of crawled *links*, not actual files.
// As it moves up the DOM, it crawls every document for additional links.
// It first moves up to find similar URLs on the same depth
//                        /\
//          /\           /  \
//    /\   /  \         /    \
// It then moves up the ladder again - to accept higher URLs
//                        /\
//          /\
//    /\
func FetchSimilar(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,

	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	wpf(b, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Find similar HTML URLs"}))
	defer wpf(b, tplx.Foot)

	wpf(b, "<pre>")
	defer wpf(b, "</pre>")

	fs1 := GetFS(appengine.NewContext(r))

	err := r.ParseForm()
	lg(err)

	countSimilar := 3
	sCountSimilar := r.FormValue("cnt")
	if sCountSimilar != "" {
		i, err := strconv.Atoi(strings.TrimSpace(sCountSimilar))
		if err == nil {
			countSimilar = i
		}
	}

	surl := r.FormValue(routes.URLParamKey)
	ourl, err := fetch.URLFromString(surl)
	lg(err)
	if err != nil {
		return
	}
	if ourl.Host == "" {
		lg("host is empty (%v)", surl)
		return
	}

	srcDepth := strings.Count(ourl.Path, "/")

	cmd := FetchCommand{}
	cmd.Host = ourl.Host
	cmd.SearchPrefix = ourl.Path
	cmd = addDefaults(w, r, cmd)

	dirTree := &DirTree{Name: "/", Dirs: map[string]DirTree{}, EndPoint: true}
	fnDigest := path.Join(docRoot, cmd.Host, "digest2.json")
	loadDigest(w, r, lg, fs1, fnDigest, dirTree) // previous
	lg("dirtree 400 chars is %v end of dirtree\n", stringspb.ToLen(dirTree.String(), 400))

	btsSrc, modSrc, err := fetchCrawlSave(w, r, lg, dirTree, fs1, path.Join(cmd.Host, ourl.Path))
	lg(err)
	if err != nil {
		return
	}

	var treePath string
	treePath = "/news"
	treePath = "/blogs"
	treePath = "/blogs/freeexchange"
	treePath = "/news/europe"
	treePath = path.Dir(ourl.Path)

	opt := LevelWiseDeeperOptions{}
	opt.Rump = treePath
	opt.ExcludeDir = "/news/americas"
	opt.ExcludeDir = "/blogs/buttonwood"
	opt.ExcludeDir = "/something-impossible"
	opt.MinDepthDiff = 1
	opt.MaxDepthDiff = 1
	opt.CondenseTrailingDirs = cmd.CondenseTrailingDirs
	opt.MaxNumber = 2266
	opt.MaxNumber = cmd.DesiredNumber + 1  // one more for "self"
	opt.MaxNumber = cmd.DesiredNumber + 10 // collect more, 'cause we filter out those too old later

	var subtree *DirTree
	links := []FullArticle{}

MarkOuter:
	for j := 0; j < srcDepth; j++ {
		treePath = path.Dir(ourl.Path)
	MarkInner:
		for i := 1; i < 22; i++ {

			subtree, treePath = DiveToDeepestMatch(dirTree, treePath)

			lg("\nLooking from height %v to level %v  - %v", srcDepth-i, srcDepth-j, treePath)

			_, _, err = fetchCrawlSave(w, r, lg, dirTree, fs1, path.Join(cmd.Host, treePath))
			lg(err)
			if err != nil {
				return
			}

			if subtree == nil {
				lg("\n#%v treePath %q ; subtree is nil", i, treePath)
			} else {
				// lg("\n#%v treePath %q ; subtree exists", i, treePath)

				opt.Rump = treePath
				opt.MinDepthDiff = i - j
				opt.MaxDepthDiff = i - j
				lvlLinks := LevelWiseDeeper(nil, nil, subtree, opt)
				links = append(links, lvlLinks...)
				for _, art := range lvlLinks {
					lg("#%v     fnd %v", i, stringspb.ToLen(art.Url, 100))
				}

				if len(links) >= opt.MaxNumber {
					lg("found enough")
					break MarkOuter
				}

				pathPrev := treePath
				treePath = path.Dir(treePath)
				// lg("#%v  bef %v - aft %v", i, pathPrev, treePath)

				if pathPrev == "." && treePath == "." ||
					pathPrev == "/" && treePath == "/" ||
					pathPrev == "" && treePath == "." {
					lg("break to innner")
					break MarkInner
				}
			}

		}
	}

	lg("\nNow reading/fetching actual similar files - not just the links")
	//
	tried := 0
	selecteds := []FullArticle{}

	for i, art := range links {

		tried = i + 1

		if art.Url == ourl.Path {
			lg("skipping self")
			continue
		}

		useExisting := false

		semanticUri := condenseTrailingDir(art.Url, cmd.CondenseTrailingDirs)
		p := path.Join(docRoot, cmd.Host, semanticUri)
		lg("reading  %v", p)
		f, err := fs1.Open(p)
		lg(err)
		if err != nil {
			// its no error if file does not exist
		} else {

			// lets put this into a func, so that f.close it called at the end of this func
			// otherwise defer f.close() spans the entire func and prevents
			// overwrites chmods further down
			f := func() {
				defer f.Close()
				fi, err := f.Stat()
				lg(err)
				if err != nil {

				} else {
					age := time.Now().Sub(fi.ModTime())
					if age.Hours() < 10 {
						lg(" using file with age %4.2v hrs", age.Hours())
						art.Mod = fi.ModTime()
						bts, err := ioutil.ReadAll(f)
						lg(err)
						art.Body = bts
						selecteds = append(selecteds, art)
						useExisting = true
					}
				}
			}
			f()

		}

		if !useExisting {
			surl := path.Join(cmd.Host, art.Url)
			lg("fetching %v", surl)
			bts, inf, err := fetch.UrlGetter(r, fetch.Options{URL: surl, RedirectHandling: 1})
			lg(err)

			if inf.Mod.IsZero() {
				inf.Mod = time.Now().Add(-75 * time.Minute)
			}

			lg("saving   %v - with time %v", p, inf.Mod.Format(http.TimeFormat))
			dir := path.Dir(p)
			err = fs1.MkdirAll(dir, 0755)
			lg(err)
			err = fs1.Chtimes(dir, time.Now(), time.Now())
			lg(err)
			err = fs1.WriteFile(p, bts, 0644)
			lg(err)
			err = fs1.Chtimes(p, inf.Mod, inf.Mod)
			lg(err)

			if inf.Mod.After(time.Now().Add(-10 * time.Hour)) {
				lg(" using fetched")
				art.Mod = inf.Mod
				art.Body = bts
				selecteds = append(selecteds, art)
			}

		}

		if len(selecteds) >= countSimilar {
			break
		}

		if tried > countSimilar+4 {
			break
		}

	}

	lg("tried %v to find %v new similars; requested: %v", tried, len(selecteds), countSimilar)

	mp := map[string][]byte{}
	mp["msg"] = b.Bytes()
	mp["url_self"] = []byte(condenseTrailingDir(ourl.Path, cmd.CondenseTrailingDirs))
	mp["mod_self"] = []byte(modSrc.Format(http.TimeFormat))
	mp["bod_self"] = btsSrc

	for i, v := range selecteds {
		mp["url__"+spf("%02v", i)] = []byte(v.Url)
		mp["mod__"+spf("%02v", i)] = []byte(v.Mod.Format(http.TimeFormat))
		mp["bod__"+spf("%02v", i)] = v.Body
	}

	//
	smp, err := json.MarshalIndent(mp, "", "\t")
	if err != nil {
		lg(b, "marshalling mp to []byte failed\n")
		return
	}

	r.Header.Set("X-Custom-Header-Counter", "nocounter")
	w.Header().Set("Content-Type", "application/json")
	w.Write(smp)

	b.Reset()             // this keeps the  buf pointer intact; outgoing defers are still heeded
	b = new(bytes.Buffer) // creates a *new* buf pointer; outgoing defers write into the *old* buf

	return

}
