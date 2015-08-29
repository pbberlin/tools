// Package fetch_rss takes http JSON commands;
// downloading html files in parallel from the designated source;
// making them available via static http fileserver.
package fetch_rss

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/osutilpb"
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

	config = addDefaults(w, r, config)

	// Fetching the rssXML takes time.
	// We do it before the timouts of the pipeline stages are set off.
	lg(" ")
	lg(config.Host)
	if config.Host == "test.economist.com" {
		switchTData(w, r)
	}

	// lg(stringspb.IndentedDump(config))
	var tzero time.Time
	dirTree := &DirTree{Name: "/", Dirs: map[string]DirTree{}, LastFound: tzero}

	fnDigest := path.Join(docRoot, config.Host, "digest2.json")
	loadDigest(w, r, fs, fnDigest, dirTree) // previous

	age := time.Now().Sub(dirTree.LastFound)
	lg("DirTree is %5.2v hours old (%v)", age.Hours(), dirTree.LastFound.Format(time.ANSIC))
	if age.Hours() > 0.001 {

		rssUrl := matchingRSSURI(w, r, config)
		if rssUrl == "" {
			err := crawl(w, r, dirTree, fs, config)
			lge(err)
			if err != nil {
				return
			}
		} else {
			rssUrl = path.Join(config.Host, rssUrl)
			rssDoc, rssUrlObj := rssXMLFile(w, r, fs, rssUrl)
			_ = rssUrlObj
			rssDoc2DirTree(w, r, dirTree, rssDoc, config.Host)
		}

		saveDigest(w, r, fs, fnDigest, dirTree)
	}

	// lg(dirTree.String())

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
	uriPrefix := config.SearchPrefix
	found := 0
	uriPrefixExcl := "impossible"
	for i := 0; i < 15; i++ {
		lg("  searching for prefix   %v    - excl %q    - %v of %v", uriPrefix, uriPrefixExcl, found, config.DesiredNumber)
		found += stuffStage1(w, r, config, inn, fin, dirTree,
			uriPrefixExcl, uriPrefix, config.DesiredNumber-found)

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

	//
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
		semanticUri := condenseTrailingDir(u.Path, config.CondenseTrailingDirs)
		dir := path.Dir(semanticUri)
		histoDir[dir]++
	}
	sr := sortmap.SortMapByCount(histoDir)
	_ = sr

	// Create dirs
	for k, _ := range histoDir {
		dir := path.Join(docRoot, k) // config.Host already contained in k
		err := fs.MkdirAll(dir, 0755)
		lge(err)
		err = fs.Chtimes(dir, time.Now(), time.Now())
		lge(err)
	}

	// Saving as files
	for _, a := range fullArticles {
		if len(a.Body) == 0 {
			continue
		}
		u, err := url.Parse(a.Url)
		u.Fragment = ""
		u.RawQuery = ""
		lge(err)
		semanticUri := condenseTrailingDir(u.RequestURI(), config.CondenseTrailingDirs)
		p := path.Join(docRoot, semanticUri)
		err = fs.WriteFile(p, a.Body, 0644)
		lge(err)
		err = fs.Chtimes(p, a.Mod, a.Mod)
		lge(err)
	}

	{
		b, err := json.MarshalIndent(histoDir, "  ", "\t")
		lge(err)
		fnDigest := path.Join(docRoot, config.Host, "fetchDigest.json")
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
	inn chan *FullArticle, fin chan struct{}, dirTree *DirTree,
	uriPrefixExcl, uriPrefixIncl string, nWant int) (nFound int) {

	lg, lge := loghttp.Logger(w, r)
	_ = lge

	depthPrefix := strings.Count(uriPrefixIncl, "/")
	if uriPrefixIncl == "/" {
		depthPrefix = 0
	}

	var subtree *DirTree
	subtree = dirTree
	head, dir, remainder := "", "", uriPrefixIncl
	for {
		dir, remainder, _ = osutilpb.PathDirReverse(remainder)
		head += dir
		// lg("    %10v %10v - dir - rdr", dir, remainder)
		if newSubtr, ok := subtree.Dirs[dir]; ok {
			subtree = &newSubtr
			// lg("    recursion found  %-10v %-10v for %v (%v)", dir, subtree.Name, uriPrefixIncl, subtree.Name)
		} else {
			// lg("    recursion failed %-10v %-10v for %v (%v)", dir, subtree.Name, uriPrefixIncl, subtree.Name)
			subtree = nil // remain on this level
			break
		}

		if remainder == "" {
			break
		}
	}

	if subtree != nil {

		var rec1 func(rump string, d *DirTree)

		rec1 = func(rump string, d *DirTree) {
			keys := make([]string, 0, len(d.Dirs))
			for k, _ := range d.Dirs {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {

				indir := d.Dirs[key]
				lpRump := rump + indir.Name
				// lg("      tryring %v", lpRump)

				isEndpoint := !d.SrcRSS || (d.SrcRSS && d.RSSEndPoint)

				if isEndpoint && checkURL(w, r, lpRump, uriPrefixExcl, uriPrefixIncl, depthPrefix, config) {

					lg("    feed #%02v: %v - %v", nFound, indir.LastFound.Format("15:04:05"), stringspb.ToLen(lpRump, 50))
					art := FullArticle{Url: config.Host + lpRump, Mod: indir.LastFound}

					select {
					case inn <- &art:
						// stage 1 loading
					case <-fin:
						lg("downstream stage has shut down, stop stuffing stage1")
						return
					}

					nFound++
					if nFound >= nWant {
						return
					}
				}

				// articles = append(articles, FullArticle{Url: config.Host + lpRump, Mod: indir.LastFound})
				if len(indir.Dirs) > 0 {
					rec1(lpRump, &indir)
				}
			}
		}

		lg("    starting recursion for %v  (%v)", head, len(subtree.Dirs))
		rec1(head, subtree)
	}

	return

}

func checkURL(w http.ResponseWriter, r *http.Request,
	surl, uriPrefixExcl, uriPrefixIncl string, depthPrefix int, config FetchCommand) bool {

	lg, lge := loghttp.Logger(w, r)
	_ = lg

	u, err := url.Parse(surl)
	lge(err)
	// short := stringspb.Ellipsoider(u.Path, 50)

	if strings.HasPrefix(u.Path, uriPrefixExcl) {
		// lg("\t\tskipping %20v", short)
		return false
	}

	if !strings.HasPrefix(u.Path, uriPrefixIncl) {
		// lg("\t\tskipping %20v", short)
		return false
	}

	semanticUri := condenseTrailingDir(u.Path, config.CondenseTrailingDirs)
	depthUri := strings.Count(semanticUri, "/")
	if depthUri > depthPrefix+1+config.DepthTolerance {
		// lg("\t\tskipping too deep    %v of %v - %20v", depthUri, depthPrefix, semanticUri)
		return false
	}
	if depthUri <= depthPrefix {
		// lg("\t\tskipping too shallow %v of %v - %20v", depthUri, depthPrefix, semanticUri)
		return false
	}

	return true
}
