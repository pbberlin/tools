// Package fetch_rss downloads html files in parallel.
package fetch_rss

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/osfs"
	"github.com/pbberlin/tools/sort/sortmap"
	"github.com/pbberlin/tools/stringspb"
)

type FullArticle struct {
	Url  string
	Body []byte
}

var fs fsi.FileSystem

func init() {
	fs = osfs.New()
	os.Chdir(docRoot)
	// fs = memfs.New()
}

func Fetch(rssUrl, uriPrefix string, numberArticles int) {

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
				pf("        fetched          %v \n", stringspb.Ellipsoider(path.Dir(u.RequestURI()), 50))
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
	// stage 1
	rssDoc := rssDir(rssUrl)

	cntr := 0
	for i, lpItem := range rssDoc.Items.ItemList {

		u, err := url.Parse(lpItem.Link)
		logif.E(err)
		short := stringspb.Ellipsoider(path.Dir(u.RequestURI()), 50)

		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		logif.E(err)

		if !strings.HasPrefix(u.RequestURI(), uriPrefix) {
			// pf("\t\t      want %20v - got %v\n", uriPrefix, short)
			pf("\t\tskipping %20v\n", short)
			continue
		}

		pf("    feed #%02v: %v - %v\n", i, t.Format("15:04:05"), short)
		inn <- &FullArticle{Url: lpItem.Link} // stage 1 loading

		cntr++
		if cntr > numberArticles {
			break
		}
	}

	pf("stage3Wait.Wait() before\n")
	stage3Wait.Wait()
	pf("stage3Wait.Wait() after\n")

	time.Sleep(3 * time.Millisecond) // not needed - workers spin down earlier

	// Saving as files
	for idx, a := range fullArticles {
		orig, numbered := fetchFileName(a.Url, idx+len(testDocs))
		err := fs.WriteFile(orig, a.Body, 0755)
		logif.E(err)

		err = fs.WriteFile(numbered, a.Body, 0755)
		logif.E(err)
	}

	// Write out directory statistics
	histoDir := map[string]int{}
	for _, a := range fullArticles {
		u, err := url.Parse(a.Url)
		logif.E(err)
		dir := filepath.Dir(u.RequestURI())
		dir = filepath.Dir(dir)
		histoDir[dir]++
	}
	sr := sortmap.SortMapByCount(histoDir)

	{
		b, err := json.MarshalIndent(sr, "  ", "\t")
		logif.E(err)
		fnDigest := filepath.Join(docRoot, "digest_detailed.json")
		err = fs.WriteFile(fnDigest, b, 0755)
		logif.E(err)
	}

	{
		b, err := json.MarshalIndent(histoDir, "  ", "\t")
		logif.E(err)
		fnDigest := filepath.Join(docRoot, "digest.json")
		err = fs.WriteFile(fnDigest, b, 0755)
		logif.E(err)
	}
}

//
func rssDir(rssUrl string) (rssDoc RSS) {

	bts, urlRSS, err := fetch.UrlGetter(rssUrl, nil, false)
	logif.F(err)

	bts = bytes.Replace(bts, []byte("content:encoded>"), []byte("content-encoded>S"), -1) // hack

	err = xml.Unmarshal(bts, &rssDoc)
	logif.E(err)

	// save it
	bdmp := stringspb.IndentedDumpBytes(rssDoc)
	err = fs.WriteFile(path.Join(docRoot, urlRSS.Host, "outp_rss.xml"), bdmp, 0755)
	logif.E(err)
	pf("RSS resp size, outp_rss.xml, : %v\n", len(bdmp))

	return
}
