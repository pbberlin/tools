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

	"github.com/pbberlin/tools/distrib"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/stringspb"
)

type MyWorker struct {
	SURL     string
	Protocol string

	r *http.Request

	lg loghttp.FuncBufUniv

	fs1 fsi.FileSystem

	err error
	FA  *FullArticle
}

func (m *MyWorker) Work() {

	bts, mod, _, err := fetchSave(m)

	if err != nil {
		m.err = err
		return
	}

	if len(bts) < 200 {
		if bytes.Contains(bts, []byte(fetch.MsgNoRdirects)) {
			return
		}
	}

	m.FA = &FullArticle{}
	m.FA.Mod = mod
	m.FA.Body = bts
	m.FA.Url = m.SURL

}

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

	knownProtocol := ""
	if r.FormValue("prot") != "" {
		knownProtocol = r.FormValue("prot")
	}

	srcDepth := strings.Count(ourl.Path, "/")

	cmd := FetchCommand{}
	cmd.Host = ourl.Host
	cmd.SearchPrefix = ourl.Path
	cmd = addDefaults(cmd)

	dirTree := &DirTree{Name: "/", Dirs: map[string]DirTree{}, EndPoint: true}
	fnDigest := path.Join(docRoot, cmd.Host, "digest2.json")
	loadDigest(w, r, lg, fs1, fnDigest, dirTree) // previous
	lg("dirtree 400 chars is %v end of dirtree\n", stringspb.ToLen(dirTree.String(), 400))

	m1 := new(MyWorker)
	m1.r = r
	m1.lg = lg
	m1.fs1 = fs1
	m1.SURL = path.Join(cmd.Host, ourl.Path)
	m1.Protocol = knownProtocol
	btsSrc, modSrc, usedExisting, err := fetchSave(m1)
	if usedExisting {
		addAnchors(lg, cmd.Host, btsSrc, dirTree)
	}
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
	opt.MaxNumber = cmd.DesiredNumber + 1  // one more for "self"
	opt.MaxNumber = cmd.DesiredNumber + 40 // collect more, 'cause we filter out those too old later

	var subtree *DirTree
	links := []FullArticle{}

MarkOuter:
	for j := 0; j < srcDepth; j++ {
		treePath = path.Dir(ourl.Path)
	MarkInner:
		for i := 1; i < srcDepth; i++ {

			subtree, treePath = DiveToDeepestMatch(dirTree, treePath)

			lg("\nLooking from height %v to level %v  - %v", srcDepth-i, srcDepth-j, treePath)

			m2 := new(MyWorker)
			m2.r = r
			m2.lg = lg
			m2.fs1 = fs1
			m2.SURL = path.Join(cmd.Host, treePath)
			m2.Protocol = knownProtocol

			btsPar, _, usedExisting, err := fetchSave(m2)
			if usedExisting {
				addAnchors(lg, cmd.Host, btsPar, dirTree)
			}
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
					lg("#%v fnd    %v", i, stringspb.ToLen(art.Url, 100))
				}

				if len(links) >= opt.MaxNumber {
					lg("found enough links")
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

	//
	//
	//
	//

	lg("\nNow reading/fetching actual similar files - not just the links\n")
	//
	tried := 0
	selecteds := []FullArticle{}

	nonExisting := []FullArticle{}
	nonExistFetched := []FullArticle{}

	for _, art := range links {

		if art.Url == ourl.Path {
			lg("skipping self\t%v", art.Url)
			continue
		}

		tried++

		useExisting := false

		semanticUri := condenseTrailingDir(art.Url, cmd.CondenseTrailingDirs)
		p := path.Join(docRoot, cmd.Host, semanticUri)
		lg("reading %q", p)
		f, err := fs1.Open(p)
		// lg(err) // its no error if file does not exist
		if err != nil {

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
						lg("\t\tusing existing file with age %4.2v hrs", age.Hours())
						art.Mod = fi.ModTime()
						bts, err := ioutil.ReadAll(f)
						lg(err)
						art.Body = bts
						if len(bts) < 200 {
							if bytes.Contains(bts, []byte(fetch.MsgNoRdirects)) {
								return
							}
						}
						selecteds = append(selecteds, art)
						useExisting = true
					}
				}
			}
			f()

		}

		if !useExisting {
			nonExisting = append(nonExisting, art)
		}

		if len(selecteds) >= countSimilar {
			break
		}

	}
	lg("tried %v links - yielding %v existing similars; %v were requested.\n",
		tried, len(selecteds), countSimilar)

	if len(selecteds) < countSimilar {
		jobs := make([]distrib.Worker, 0, len(nonExisting))
		for _, art := range nonExisting {
			surl := path.Join(cmd.Host, art.Url)
			wrkr := MyWorker{SURL: surl}
			wrkr.Protocol = knownProtocol
			wrkr.r = r
			wrkr.lg = lg
			wrkr.fs1 = fs1
			job := distrib.Worker(&wrkr)
			jobs = append(jobs, job)
		}

		opt := distrib.NewDefaultOptions()
		opt.TimeOutDur = 3500 * time.Millisecond
		opt.Want = int32(countSimilar - len(selecteds) + 4) // get some more, in case we have "redirected" bodies
		opt.NumWorkers = int(opt.Want)                      // 5s query limit; => hurry; spawn as many as we want
		opt.CollectRemainder = false                        // 5s query limit; => hurry; dont wait for stragglers

		ret, msg := distrib.Distrib(jobs, opt)
		lg("\n" + msg.String())
		for _, v := range ret {
			v1, _ := v.Worker.(*MyWorker)
			if v1.FA != nil {
				age := time.Now().Sub(v1.FA.Mod)
				if age.Hours() < 10 {
					lg("\t\tusing fetched file with age %4.2v hrs", age.Hours())
					nonExistFetched = append(nonExistFetched, *v1.FA)
					if len(nonExistFetched) > (countSimilar - len(selecteds)) {
						break
					}
				}
			}
			if v1.err != nil {
				lg(err)
			}
		}

		lg("tried %v links - yielding %v fetched - jobs %v\n", len(nonExisting), len(nonExistFetched), len(jobs))
		selecteds = append(selecteds, nonExistFetched...)

		//
		//
		// Extract links
		for _, v := range nonExistFetched {
			// lg("links -> memory dirtree for %q", v.Url)
			addAnchors(lg, cmd.Host, v.Body, dirTree)
		}

	}

	//
	if time.Now().Sub(dirTree.LastFound).Seconds() < 10 {
		lg("saving accumulated (new) links to digest")
		saveDigest(w, r, fs1, fnDigest, dirTree)
	}

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

	mp["lensimilar"] = []byte(spf("%02v", len(selecteds)))

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
