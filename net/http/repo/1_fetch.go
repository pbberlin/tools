// Package repo takes http JSON commands;
// downloading html files in parallel from the designated source;
// making them available via quasi-static http fileserver.
package repo

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
	dirTree := &DirTree{Name: "/", Dirs: map[string]DirTree{}, EndPoint: true}

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
				pth := fetch.PathFromStringUrl(fa.Url)
				lg("    fetched   %v - %v ", fa.Mod.Format("15:04:05"), stringspb.Ellipsoider(pth, 50))
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
					var inf fetch.Info
					a.Body, inf, err = fetch.UrlGetter(r, fetch.Options{URL: a.Url})
					lge(err)
					if a.Mod.IsZero() {
						a.Mod = inf.Mod
					}
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

	subtree, head := DiveToDeepestMatch(dirTree, uriPrefixIncl)

	if subtree == nil {
		lg("      does not exist in dirtree: %q", uriPrefixIncl)
	} else {

		articles := LevelWiseDeeper(w, r, head, subtree, uriPrefixExcl, config.DepthTolerance, config.CondenseTrailingDirs, nWant)
		// lg("      levelwise deeper found %v articles", len(articles))

		for _, art := range articles {

			lg("    feed #%02v: %v - %v", nFound, art.Mod.Format("15:04:05"), stringspb.Ellipsoider(art.Url, 50))

			art.Url = config.Host + art.Url

			select {
			case inn <- &art:
				// stage 1 loading
			case <-fin:
				lg("downstream stage has shut down, stop stuffing stage1")
				return
			}

			nFound++
			if nFound > nWant-1 {
				return
			}

		}

	}

	return

}

func DiveToDeepestMatch(dirTree *DirTree, uriPrefixIncl string) (*DirTree, string) {

	var subtree *DirTree
	subtree = dirTree
	head, dir, remainder := "", "", uriPrefixIncl

	if subtree.Name == uriPrefixIncl {
		// exception for root
		head = "" // not "/"
	} else {
		// recur deeper
		for {
			dir, remainder, _ = osutilpb.PathDirReverse(remainder)
			head += dir
			// lg("    %-10q %-10q %-10q - head - dir - remainder", head, dir, remainder)
			if newSubtr, ok := subtree.Dirs[dir]; ok {
				subtree = &newSubtr
				// lg("    recursion found  %-10v %-10v for %v (%v)", dir, subtree.Name, uriPrefixIncl, subtree.Name)
			} else {
				// lg("    recursion failed %-10v %-10v for %v (%v)", dir, subtree.Name, uriPrefixIncl, subtree.Name)

				// Calling off searching on this level
				// StuffStage() itsself can step up if it wants to
				//
				// *not* setting subtree = nil
				// would keep us one level higher than this level.
				subtree = nil

				break
			}

			if remainder == "" {
				break
			}
		}
	}

	return subtree, head
}

func LevelWiseDeeper(w http.ResponseWriter, r *http.Request, rump string, dtree *DirTree, excludeDir string,
	maxDepthDiff, CondenseTrailingDirs, maxNumber int) []FullArticle {

	lg, lge := loghttp.Logger(w, r)
	_ = lge

	depthRump := strings.Count(rump, "/")

	arts := []FullArticle{}

	var fc func(rmp1 string, dr1 *DirTree, lvl int)
	fc = func(rmp1 string, dr1 *DirTree, lvl int) {

		// lg("      lvl %2v %v", lvl, dr1.Name)
		if lvl == 1 && excludeDir == dr1.Name {
			lg("    excl %v", excludeDir)
			return
		}

		keys := make([]string, 0, len(dr1.Dirs))
		for k, _ := range dr1.Dirs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			dr2 := dr1.Dirs[key]
			rmp2 := rmp1 + dr2.Name

			// lg("        %v", rmp2)

			//
			// rmp2 a candidate?
			if dr2.EndPoint {
				semanticUri := condenseTrailingDir(rmp2, CondenseTrailingDirs)
				depthUri := strings.Count(semanticUri, "/")

				if depthUri-depthRump <= maxDepthDiff {

					// lg("        %v %v ; %v     %v", depthRump, depthUri, maxDepthDiff, semanticUri)

					art := FullArticle{Url: rmp2}
					if dr2.SrcRSS {
						art.Mod = dr2.LastFound
					}
					if len(arts) > maxNumber-1 {
						return
					}
					arts = append(arts, art)

				}
			}

			//
			// recurse deeper?
			if len(dr2.Dirs) == 0 {
				// lg("      LevelWiseDeeper - no children")
				continue
			}
			fc(rmp2, &dr2, lvl+1)

		}
	}

	fc(rump, dtree, 0)

	return arts

}
