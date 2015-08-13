// Package fetch_rss downloads html files in parallel.
package fetch_rss

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/sort/sortmap"
	"github.com/pbberlin/tools/stringspb"
)

var hosts = map[string]map[string]interface{}{
	"www.handelsblatt.com": map[string]interface{}{
		"host":                   "www.handelsblatt.com",
		"rss-xml-uri":            "/contentexport/feed/schlagzeilen",
		"condense-trailing-dirs": 2,
		"depth-tolerance":        1, // The last one or two directories might be article titles or ids
	},
}

type FullArticle struct {
	Url  string
	Mod  time.Time
	Body []byte
}

const numWorkers = 3

func Fetch(w http.ResponseWriter, r *http.Request, fs fsi.FileSystem, config map[string]interface{}, uriPrefix string, numberArticles int) {

	lg, lge := loghttp.Logger(w, r)

	rssUrl := path.Join(config["host"].(string), config["rss-xml-uri"].(string))

	// Fetching the rssXML takes time.
	// We do it before the timouts of the pipeline stages are set off.
	rssDoc, rssUrlObj := rssXMLFile(w, r, fs, rssUrl)

	//
	//
	// setting up a 3 staged pipeline from bottom up
	//
	var fullArticles []FullArticle

	var inn chan *FullArticle = make(chan *FullArticle) // jobs are stuffed in here
	var out chan *FullArticle = make(chan *FullArticle) // completed jobs are delivered here
	var fin chan struct{} = make(chan struct{})         // downstream signals end to upstream
	var stage3Wait sync.WaitGroup

	// stage 3
	// fire up the "collector", a fan-in
	go func() {
		stage3Wait.Add(1)
		const delayInitial = 1200
		const delayRefresh = 800 // 400 good value; critical point at 35
		cout := time.After(time.Millisecond * delayInitial)
		for {
			select {

			case fa := <-out:
				fullArticles = append(fullArticles, *fa)
				u, _ := url.Parse(fa.Url)
				lg("    fetched              %v ", stringspb.Ellipsoider(path.Dir(u.RequestURI()), 50))
				cout = time.After(time.Millisecond * delayRefresh) // refresh timeout
			case <-cout:
				lg("timeout after %v articles", len(fullArticles))
				// we are using channel == nil - channel closed combinations
				// inspired by http://dave.cheney.net/2013/04/30/curious-channels
				out = nil // not close(out) => case above is now blocked
				close(fin)
				lg("fin closed; out nilled")
				stage3Wait.Done()
				return
			}
		}
	}()

	//
	// stage 2
	for i := 0; i < numWorkers; i++ {
		// fire up a dedicated fetcher routine, a worker
		// we are using channel == nil - channel closed combinations
		// inspired by http://dave.cheney.net/2013/04/30/curious-channels
		go func() {
			var a *FullArticle
			for {
				select {
				case a = <-inn:
					var err error
					a.Body, _, err = fetch.UrlGetter(a.Url, r, false)
					lge(err)
					out <- a
					a = new(FullArticle)
				case <-fin:
					if a != nil && a.Url != "" {
						u, _ := url.Parse(a.Url)
						lg("    abandoned %v", u.RequestURI())
					} else {
						lg("    worker spinning down")
					}
					return
				}
			}
		}()
	}

	//
	//
	//
	// loading stage 1
	found := 0
	uriPrefixExcl := "impossible"
	for i := 0; i < 15; i++ {
		lg("  searching for prefix %q excl %q   - %v of %v", uriPrefix, uriPrefixExcl, found, numberArticles)
		found += stuffStage1(w, r, config, inn, &rssDoc, uriPrefix, uriPrefixExcl, numberArticles-found)
		if found >= numberArticles {
			break
		}

		if uriPrefix == "/" {
			lg("  root exhausted")
			break
		}

		newPrefix := path.Dir(uriPrefix)
		uriPrefixExcl = uriPrefix
		uriPrefix = newPrefix
	}
	lg("  found %v of %v", found, numberArticles)

	lg("stage3Wait.Wait() before")
	stage3Wait.Wait()
	lg("stage3Wait.Wait() after")

	time.Sleep(3 * time.Millisecond) // not needed - workers spin down earlier

	// compile out directory statistics
	histoDir := map[string]int{}
	condenseTrailingDirs := config["condense-trailing-dirs"].(int)
	for _, a := range fullArticles {
		u, err := url.Parse(a.Url)
		lge(err)
		semanticUri := condenseTrailingDir(u.RequestURI(), condenseTrailingDirs)
		dir := path.Dir(semanticUri)
		histoDir[dir]++
	}
	sr := sortmap.SortMapByCount(histoDir)

	// Create dirs
	for k, _ := range histoDir {
		dir := path.Join(docRoot, rssUrlObj.Host, k)
		err := fs.MkdirAll(dir, 0755)
		lge(err)
	}

	// Saving as files
	for _, a := range fullArticles {
		u, err := url.Parse(a.Url)
		lge(err)
		semanticUri := condenseTrailingDir(u.RequestURI(), condenseTrailingDirs)
		p := path.Join(docRoot, u.Host, semanticUri)
		err = fs.WriteFile(p, a.Body, 0644)
		lge(err)
		err = fs.Chtimes(p, a.Mod, a.Mod)
		lge(err)
	}

	// Save digests
	{
		b, err := json.MarshalIndent(sr, "  ", "\t")
		lge(err)
		fnDigest := path.Join(docRoot, "digest_detailed.json")
		err = fs.WriteFile(fnDigest, b, 0755)
		lge(err)
	}

	{
		b, err := json.MarshalIndent(histoDir, "  ", "\t")
		lge(err)
		fnDigest := path.Join(docRoot, "digest.json")
		err = fs.WriteFile(fnDigest, b, 0755)
		lge(err)
	}

	// fsm, ok := memfs.Unwrap(fs)
	// if ok {
	// 	fsm.Dump()
	// }

}

// stuffStage1 ranges of the RSS entries and filters out unwanted directories.
// Wanted urls are sent to the stage one channel.
func stuffStage1(w http.ResponseWriter, r *http.Request, config map[string]interface{}, inn chan *FullArticle, rssDoc *RSS,
	uriPrefixIncl, uriPrefixExcl string, nWant int) (nFound int) {

	lg, lge := loghttp.Logger(w, r)

	depthPrefix := strings.Count(uriPrefixIncl, "/")
	if uriPrefixIncl == "/" {
		depthPrefix = 0
	}
	condenseTrailingDirs := config["condense-trailing-dirs"].(int)
	depthTolerance := config["depth-tolerance"].(int)

	for i, lpItem := range rssDoc.Items.ItemList {

		u, err := url.Parse(lpItem.Link)
		lge(err)
		short := stringspb.Ellipsoider(path.Dir(u.RequestURI()), 50)

		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		lge(err)

		if strings.HasPrefix(u.RequestURI(), uriPrefixExcl) {
			// lg("\t\tskipping %20v", short)
			continue
		}

		if !strings.HasPrefix(u.RequestURI(), uriPrefixIncl) {
			// lg("\t\tskipping %20v", short)
			continue
		}

		semanticUri := condenseTrailingDir(u.RequestURI(), condenseTrailingDirs)
		depthUri := strings.Count(semanticUri, "/")
		if depthUri > depthPrefix+1+depthTolerance {
			// lg("\t\tskipping %20v - too deep (%v - %v)", semanticUri, depthPrefix, depthUri)
			continue
		}

		lg("    feed #%02v: %v - %v", i, t.Format("15:04:05"), short)
		inn <- &FullArticle{Url: lpItem.Link, Mod: t} // stage 1 loading

		nFound++
		if nFound >= nWant {
			break
		}
	}

	return
}
