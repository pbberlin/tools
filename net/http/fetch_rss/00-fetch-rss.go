// Package fetch_rss downloads html files in parallel.
package fetch_rss

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
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
	Body []byte
}

var fs fsi.FileSystem

var docRoot = "c:/docroot/" // no relative path, 'cause working dir too flippant

func init() {
	fs = osfs.New()
	os.Chdir(docRoot)

	fs = memfs.New(memfs.DirSort("byDateDesc"))

}

func Fetch(config map[string]interface{}, uriPrefix string, numberArticles int) {

	rssUrl := path.Join(config["host"].(string), config["rss-xml-uri"].(string))

	//
	// setting up a 3 staged pipeline from bottom up
	//
	var fullArticles []FullArticle
	const numWorkers = 3

	var inn chan *FullArticle = make(chan *FullArticle) // jobs are stuffed in here
	var out chan *FullArticle = make(chan *FullArticle) // completed jobs are delivered here
	var fin chan struct{} = make(chan struct{})         // downstream signals end to upstream
	var stage3Wait sync.WaitGroup

	// stage 3
	// fire up the "collector", a fan-in
	go func() {
		stage3Wait.Add(1)
		const delay1 = 800
		const delay2 = 400 // 400 good value; critical point at 35
		cout := time.After(time.Millisecond * delay1)
		for {
			select {

			case fa := <-out:
				fullArticles = append(fullArticles, *fa)
				u, _ := url.Parse(fa.Url)
				pf("    fetched              %v \n", stringspb.Ellipsoider(path.Dir(u.RequestURI()), 50))
				cout = time.After(time.Millisecond * delay2) // refresh timeout
			case <-cout:
				pf("timeout after %v articles\n", len(fullArticles))
				// we are using channel == nil - channel closed combinations
				// inspired by http://dave.cheney.net/2013/04/30/curious-channels
				out = nil // not close(c)
				close(fin)
				pf("fin closed; out nilled\n")
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
					a.Body, _, err = fetch.UrlGetter(a.Url, nil, false)
					logif.F(err)
					out <- a
					a = new(FullArticle)
				case <-fin:
					if a != nil && a.Url != "" {
						u, _ := url.Parse(a.Url)
						pf("    abandoned %v\n", u.RequestURI())
					} else {
						pf("    worker spinning down\n")
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
	rssDoc, rssUrlObj := rssXMLFile(rssUrl)
	found := 0
	uriPrefixExcl := "impossible"
	for i := 0; i < 15; i++ {
		pf("  searching for prefix %q excl %q   - %v of %v\n", uriPrefix, uriPrefixExcl, found, numberArticles)
		found += stuffStage1(config, inn, &rssDoc, uriPrefix, uriPrefixExcl, numberArticles-found)
		if found >= numberArticles {
			break
		}

		if uriPrefix == "/" {
			pf("  root exhausted\n")
			break
		}

		newPrefix := path.Dir(uriPrefix)
		uriPrefixExcl = uriPrefix
		uriPrefix = newPrefix
	}
	pf("  found %v of %v\n", found, numberArticles)

	pf("stage3Wait.Wait() before\n")
	stage3Wait.Wait()
	pf("stage3Wait.Wait() after\n")

	time.Sleep(3 * time.Millisecond) // not needed - workers spin down earlier

	// compile out directory statistics
	histoDir := map[string]int{}
	condenseTrailingDirs := config["condense-trailing-dirs"].(int)
	for _, a := range fullArticles {
		u, err := url.Parse(a.Url)
		logif.E(err)
		semanticUri := condenseTrailingDir(u.RequestURI(), condenseTrailingDirs)
		dir := path.Dir(semanticUri)
		histoDir[dir]++
	}
	sr := sortmap.SortMapByCount(histoDir)

	// Create dirs
	for k, _ := range histoDir {
		dir := path.Join(docRoot, rssUrlObj.Host, k)
		err := fs.MkdirAll(dir, 0755)
		logif.E(err)
	}

	// Saving as files
	for _, a := range fullArticles {
		u, err := url.Parse(a.Url)
		logif.E(err)
		semanticUri := condenseTrailingDir(u.RequestURI(), condenseTrailingDirs)
		p := path.Join(docRoot, u.Host, semanticUri)
		err = fs.WriteFile(p, a.Body, 0644)
		logif.E(err)
	}

	// digests
	{
		b, err := json.MarshalIndent(sr, "  ", "\t")
		logif.E(err)
		fnDigest := path.Join(docRoot, "digest_detailed.json")
		err = fs.WriteFile(fnDigest, b, 0755)
		logif.E(err)
	}

	{
		b, err := json.MarshalIndent(histoDir, "  ", "\t")
		logif.E(err)
		fnDigest := path.Join(docRoot, "digest.json")
		err = fs.WriteFile(fnDigest, b, 0755)
		logif.E(err)
	}

	// fsm, ok := memfs.Unwrap(fs)
	// if ok {
	// 	fsm.Dump()
	// }

}

func stuffStage1(config map[string]interface{}, inn chan *FullArticle, rssDoc *RSS, uriPrefixIncl, uriPrefixExcl string, nWant int) (nFound int) {

	depthPrefix := strings.Count(uriPrefixIncl, "/")
	if uriPrefixIncl == "/" {
		depthPrefix = 0
	}
	condenseTrailingDirs := config["condense-trailing-dirs"].(int)
	depthTolerance := config["depth-tolerance"].(int)

	for i, lpItem := range rssDoc.Items.ItemList {

		u, err := url.Parse(lpItem.Link)
		logif.E(err)
		short := stringspb.Ellipsoider(path.Dir(u.RequestURI()), 50)

		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		logif.E(err)

		if strings.HasPrefix(u.RequestURI(), uriPrefixExcl) {
			// pf("\t\tskipping %20v\n", short)
			continue
		}

		if !strings.HasPrefix(u.RequestURI(), uriPrefixIncl) {
			// pf("\t\tskipping %20v\n", short)
			continue
		}

		semanticUri := condenseTrailingDir(u.RequestURI(), condenseTrailingDirs)
		depthUri := strings.Count(semanticUri, "/")
		if depthUri > depthPrefix+1+depthTolerance {
			pf("\t\tskipping %20v - too deep (%v - %v)\n", semanticUri, depthPrefix, depthUri)
			continue
		}

		pf("    feed #%02v: %v - %v\n", i, t.Format("15:04:05"), short)
		inn <- &FullArticle{Url: lpItem.Link} // stage 1 loading

		nFound++
		if nFound >= nWant {
			break
		}
	}

	return
}
