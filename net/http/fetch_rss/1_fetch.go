// Package fetch_rss takes http JSON commands;
// downloading html files in parallel from the designated source;
// making them available via static http fileserver.
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

// FullArticle is the main struct passed
// between the pipeline stages
type FullArticle struct {
	Url  string
	Mod  time.Time
	Body []byte
}

// parallel fetchers routines
const numWorkers = 3

// Fetch takes a RSS XML uri and fetches some of its documents.
// It uses a three staged pipeline for parallel fetching.
// Results are stored into the given filesystem fs.
// Config points to the source of RSS XML,
// and has some rules for conflating URI directories.
// uriPrefix and config.DesiredNumber tell the func
// which subdirs of the RSS dir should be fetched - and how many at max.
func Fetch(w http.ResponseWriter, r *http.Request,
	fs fsi.FileSystem, config FetchCommand,
) {

	lg, lge := loghttp.Logger(w, r)

	if config.Host == "" {
		lg(" empty host; returning")
		return
	}

	// Fetching the rssXML takes time.
	// We do it before the timouts of the pipeline stages are set off.
	lg(" ")
	lg(config.Host)
	rssUrl := path.Join(config.Host, config.RssXMLURI)
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
		// 400 good value; critical point at 35
		// economist.com required 800 ms
		const delayInitial = 1200
		const delayRefresh = 800
		cout := time.After(time.Millisecond * delayInitial)
		for {
			select {

			case fa := <-out:
				fullArticles = append(fullArticles, *fa)
				u, _ := url.Parse(fa.Url)
				lg("    fetched              %v ", stringspb.Ellipsoider(u.Path, 50))
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
					a.Body, _, err = fetch.UrlGetter(r, fetch.Options{URL: a.Url})
					lge(err)
					select {
					case out <- a:
					case <-fin:
						lg("    worker spinning down; branch 1; abandoning %v", a.Url)
						return
					}
					a = new(FullArticle)
				case <-fin:
					if a != nil && a.Url != "" {
						u, _ := url.Parse(a.Url)
						lg("    abandoned %v", u.RequestURI())
					} else {
						lg("    worker spinning down; branch 2")
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
	for _, uriPrefix := range config.SearchPrefixs {
		found := 0
		uriPrefixExcl := "impossible"
		for i := 0; i < 15; i++ {
			lg("  searching for prefix   %v    - excl %q    - %v of %v", uriPrefix, uriPrefixExcl, found, config.DesiredNumber)
			found += stuffStage1(w, r, config, inn, fin, &rssDoc,
				uriPrefix, uriPrefixExcl, config.DesiredNumber-found)

			if found >= config.DesiredNumber {
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
		lg("  found %v of %v", found, config.DesiredNumber)
	}

	lg("stage3Wait.Wait() before")
	stage3Wait.Wait()
	lg("stage3Wait.Wait() after")

	// workers spin down earlier -
	// but ae log writer and response writer need some time
	// to record the spin-down messages
	time.Sleep(120 * time.Millisecond)

	// compile out directory statistics
	histoDir := map[string]int{}
	for _, a := range fullArticles {
		u, err := url.Parse(a.Url)
		lge(err)
		semanticUri := condenseTrailingDir(u.RequestURI(), config.CondenseTrailingDirs)
		dir := path.Dir(semanticUri)
		histoDir[dir]++
	}
	sr := sortmap.SortMapByCount(histoDir)
	_ = sr

	// Create dirs
	for k, _ := range histoDir {
		dir := path.Join(docRoot, rssUrlObj.Host, k)
		err := fs.MkdirAll(dir, 0755)
		lge(err)
		err = fs.Chtimes(dir, time.Now(), time.Now())
		lge(err)
	}

	// Saving as files
	for _, a := range fullArticles {
		u, err := url.Parse(a.Url)
		u.Fragment = ""
		u.RawQuery = ""
		lge(err)
		semanticUri := condenseTrailingDir(u.RequestURI(), config.CondenseTrailingDirs)
		p := path.Join(docRoot, u.Host, semanticUri)
		err = fs.WriteFile(p, a.Body, 0644)
		lge(err)
		err = fs.Chtimes(p, a.Mod, a.Mod)
		lge(err)
	}

	// Save digests
	// {
	// 	b, err := json.MarshalIndent(sr, "  ", "\t")
	// 	lge(err)
	// 	fnDigest := path.Join(docRoot, rssUrlObj.Host, "digest_detailed.json")
	// 	err = fs.WriteFile(fnDigest, b, 0755)
	// 	lge(err)
	// }

	{
		b, err := json.MarshalIndent(histoDir, "  ", "\t")
		lge(err)
		fnDigest := path.Join(docRoot, rssUrlObj.Host, "digest.json")
		err = fs.WriteFile(fnDigest, b, 0755)
		lge(err)
	}

	// fsm, ok := memfs.Unwrap(fs)
	// if ok {
	// 	fsm.Dump()
	// }

}

// stuffStage1 ranges over the RSS entries and filters out unwanted directories.
// Wanted urls are sent to the stage one channel.
func stuffStage1(w http.ResponseWriter, r *http.Request, config FetchCommand,
	inn chan *FullArticle, fin chan struct{}, rssDoc *RSS,
	uriPrefixIncl, uriPrefixExcl string, nWant int) (nFound int) {

	lg, lge := loghttp.Logger(w, r)

	depthPrefix := strings.Count(uriPrefixIncl, "/")
	if uriPrefixIncl == "/" {
		depthPrefix = 0
	}

	for i, lpItem := range rssDoc.Items.ItemList {

		u, err := url.Parse(lpItem.Link)
		lge(err)
		short := stringspb.Ellipsoider(u.Path, 50)

		t, err := time.Parse(time.RFC1123Z, lpItem.Published)
		//     := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		lge(err)

		if strings.HasPrefix(u.RequestURI(), uriPrefixExcl) {
			// lg("\t\tskipping %20v", short)
			continue
		}

		if !strings.HasPrefix(u.RequestURI(), uriPrefixIncl) {
			// lg("\t\tskipping %20v", short)
			continue
		}

		semanticUri := condenseTrailingDir(u.RequestURI(), config.CondenseTrailingDirs)
		depthUri := strings.Count(semanticUri, "/")
		if depthUri > depthPrefix+1+config.DepthTolerance {
			lg("\t\tskipping too deep    %v of %v - %20v", depthUri, depthPrefix, semanticUri)
			continue
		}
		if depthUri <= depthPrefix {
			lg("\t\tskipping too shallow %v of %v - %20v", depthUri, depthPrefix, semanticUri)
			continue
		}

		lg("    feed #%02v: %v - %v", i, t.Format("15:04:05"), short)

		select {
		case inn <- &FullArticle{Url: lpItem.Link, Mod: t}:
			// stage 1 loading
		case <-fin:
			lg("downstream stage has shut down, stop stuffing stage1")
			return
		}

		nFound++
		if nFound >= nWant {
			break
		}
	}

	return
}
