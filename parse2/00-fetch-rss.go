package parse2

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"sync"
	"time"

	"github.com/pbberlin/tools/pbfetch"
	"github.com/pbberlin/tools/pblog"
	"github.com/pbberlin/tools/pbstrings"
)

var hosts = []string{"www.handelsblatt.com"}

type FullArticle struct {
	Url  string
	Body []byte
}

func Fetch(rssUrl string, numberArticles int) {

	var fullArticles []FullArticle

	var c chan *FullArticle = make(chan *FullArticle)
	var finish chan struct{} = make(chan struct{})
	var wg sync.WaitGroup

	// fire up the "collector"
	go func() {
		wg.Add(1)
		const delay1 = 800
		const delay2 = 400 // 400 good value; critical point at 35
		cout := time.After(time.Millisecond * delay1)
		for {
			select {

			case pfa := <-c:
				fa := *pfa
				fullArticles = append(fullArticles, fa)
				u, _ := url.Parse(fa.Url)
				pf("    fetched %v \n", u.RequestURI())
				cout = time.After(time.Millisecond * delay2) // refresh timeout
			case <-cout:
				pf("timeout after %v articles\n", len(fullArticles))
				close(finish)
				c = nil
				// close(c)
				pf("finish closed; c nilled\n")
				wg.Done()
				return
			}

		}
	}()

	bts, err := pbfetch.UrlGetter(rssUrl, nil, false)
	pblog.Fatal(err)

	bts = bytes.Replace(bts, []byte("content:encoded>"), []byte("content-encoded>S"), -1)

	var rssDoc RSS
	err = xml.Unmarshal(bts, &rssDoc)
	pblog.LogE(err)

	bdmp := pbstrings.IndentedDumpBytes(rssDoc)
	bytes2File("outp_rss.xml", bdmp)
	pf("RSS resp size, outp_rss.xml, : %v\n", len(bdmp))

	for i, lpItem := range rssDoc.Items.ItemList {

		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", lpItem.Published)
		pblog.LogE(err)
		u, err := url.Parse(lpItem.Link)
		pblog.LogE(err)
		pf("    feed #%02v: %v - %v\n", i, t.Format("15:04:05"), u.RequestURI())

		// fire up a dedicated fetcher routine
		go func(argURL string) {

			// defer func() {
			// 	panicSignal := recover()
			// 	if panicSignal != nil {
			// 		log.Printf("panic caught:\n\n%s", panicSignal)
			// 		util_err.StackTrace(15)
			// 	}
			// }()

			bs, err := pbfetch.UrlGetter(argURL, nil, false)
			pblog.Fatal(err)
			select {
			case c <- &FullArticle{argURL, bs}:
				return
			case <-finish:
				u, _ := url.Parse(argURL)
				pf("    abandoned %v\n", u.RequestURI())
				return
			}
		}(lpItem.Link)

		if i+1 >= numberArticles {
			break
		}
	}

	// time.Sleep(4 * time.Second)
	pf("wait() before\n")
	wg.Wait()
	pf("wait() after\n")

	time.Sleep(333 * time.Millisecond)

	// just saving them
	for idx, a := range fullArticles {
		orig, numbered := fetchFileName(a.Url, idx+len(testDocs))
		bytes2File(orig, a.Body)
		bytes2File(numbered, a.Body)
	}

}
